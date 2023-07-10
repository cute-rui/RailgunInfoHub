package railgun

import (
	"errors"
	"github.com/go-resty/resty/v2"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/multierr"
	"net/http"
	"strconv"
	"time"
)

type request struct {
	client   *resty.Client
	wbiCache *WBICache
}

func newRequest() *request {
	return &request{client: newRestyClient()}
}

func newRestyClient() *resty.Client {
	return resty.New().SetRetryCount(3).
		SetRetryWaitTime(5*time.Second).
		SetRetryMaxWaitTime(20*time.Second).
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, nil
		}).
		SetCookie(&http.Cookie{Name: `CURRENT_FNVAL`, Value: DEFAULT_FNVAL, Path: `/`, Expires: time.Now().AddDate(0, 0, 1)}).
		SetCookie(&http.Cookie{Name: `CURRENT_QUALITY`, Value: DEFAULT_QUALITY, Path: `/`, Expires: time.Now().AddDate(0, 0, 1)}).
		SetHeader(`Referer`, DEFAULT_REFERER).
		SetHeader(`User-Agent`, DEFAULT_USER_AGENT)
}

func (r *request) setSessData(sessdata string) {
	r.client.SetCookie(&http.Cookie{Name: `SESSDATA`, Value: sessdata, Path: `/`, Expires: time.Now().AddDate(0, 0, 1)})
}

func (r *request) checkRefresh() (*needRenewCookieResp, error) {
	respRaw, err := r.client.R().
		Get(CHECK_COOKIE_URL)
	if err != nil {
		return nil, err
	}

	var resp = needRenewCookieResp{}
	err = jsoniter.Unmarshal(respRaw.Body(), &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r *request) getRefreshCsrf(path string) (string, error) {
	respRaw, err := r.client.R().
		Get(stringBuilder(CORESPONDING_URL, path))
	if err != nil {
		return ``, err
	}

	return ByteToString(respRaw.Body()), nil
}

func (r *request) refreshCookie(data map[string]string) (bili_jct string, sessData string, refresh_token string, err error) {
	respRaw, err := r.client.R().
		SetHeader(`content-type`, `application/x-www-form-urlencoded`).
		SetFormData(data).
		Post(REFRESH_COOKIE_URL)

	if err != nil {
		return ``, ``, ``, err
	}

	for i := range respRaw.Cookies() {
		if respRaw.Cookies()[i].Name == `bili_jct` {
			bili_jct = respRaw.Cookies()[i].Value
		}
		if respRaw.Cookies()[i].Name == `SESSDATA` {
			sessData = respRaw.Cookies()[i].Value
		}
	}

	var resp = refreshCookieResp{}
	err = jsoniter.Unmarshal(respRaw.Body(), &resp)
	if err != nil {
		return ``, ``, ``, err
	}

	return bili_jct, sessData, resp.Data.RefreshToken, nil

}

func (r *request) getNavInfo() (*navResp, error) {
	respRaw, err := r.client.R().
		Get(NAV_INFO_URL)
	if err != nil {
		return nil, err
	}

	var resp = navResp{}
	err = jsoniter.Unmarshal(respRaw.Body(), &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (r *request) getShortLinkLocation(url string) (string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := client.Get(url)
	if err != nil {
		return ``, err
	}

	if !(res.StatusCode == 301 || res.StatusCode == 302) {
		return ``, errors.New("no redirection")
	}

	return res.Header.Get(`Location`), nil
}

func (r *request) fetchCollection(seasonID string, pagenum int) (*fetchCollectionResp, error) {
	respRaw, err := r.client.R().
		SetQueryParam(`mid`, `1`).
		SetQueryParam(`sort_reverse`, `false`).
		SetQueryParam(`page_num`, strconv.Itoa(pagenum)).
		SetQueryParam(`page_size`, `100`).
		SetQueryParam(`season_id`, seasonID).
		Get(stringBuilder(DEFAULT_HOSTNAME_WITH_SCHEME, COLLECTION_URI))
	if err != nil {
		return nil, err
	}

	var resp = fetchCollectionResp{}
	err = jsoniter.Unmarshal(respRaw.Body(), &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(resp.Message)
	}

	return &resp, nil
}

// Todo: proxy
func (r *request) fetchVideo(avidStr, bvid string) (*fetchVideoResp, error) {
	var option restyOption
	if avidStr != `` {
		option = setAVIDQueryParam(avidStr)
	} else {
		option = setBVIDQueryParam(bvid)
	}

	respRaw, err := r.fetchRawVideo(option)

	var resp = fetchVideoResp{}
	err = jsoniter.Unmarshal(respRaw.Body(), &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(resp.Message)
	}

	return &resp, nil
}

func (r *request) fetchRawVideo(options ...restyOption) (*resty.Response, error) {
	if len(options) == 0 {
		return nil, errors.New(`no video id parameter present`)
	}

	req := r.client.R()

	for i := range options {
		options[i](req)
	}

	return req.Get(stringBuilder(DEFAULT_HOSTNAME_WITH_SCHEME, VIDEO_URI))
}

func (r *request) fetchSeasonOrEpisode(seasonID, episodeID string, region Region) (*fetchSeasonResp, Region, error) {
	var option restyOption
	if episodeID != `` {
		option = setEPIDQueryParam(episodeID)
	} else {
		option = setSeasonIDQueryParam(seasonID)
	}

	test := region
	resp, err := r.fetchRawSeasonOrEpisode(&test, option)

	if err != nil {
		return nil, REGION_UNLOCATED, err
	}

	if resp.Code == -404 {
		return r.locateSeasonOrEpisode(region, setSeasonIDQueryParam(seasonID))
	}

	if resp.Code != 0 {
		return nil, REGION_UNLOCATED, errors.New(resp.Message)
	}

	return resp, region, nil
}

func (r *request) locateSeasonOrEpisode(region Region, options ...restyOption) (*fetchSeasonResp, Region, error) {
	var err error

	for i := REGION_CN; i <= REGION_UNLOCATED; i++ {
		if i == region {
			continue
		}

		test := i
		resp, e := r.fetchRawSeasonOrEpisode(&test, options...)
		if e != nil {
			err = multierr.Append(err, e)
			continue
		}

		if test == REGION_UNLOCATED {
			continue
		}

		return resp, test, nil
	}

	return nil, REGION_UNLOCATED, err
}

type restyOption func(req *resty.Request)

func setAVIDQueryParam(avid string) restyOption {
	return func(req *resty.Request) {
		req.SetQueryParam(`aid`, avid)
	}
}

func setBVIDQueryParam(bvid string) restyOption {
	return func(req *resty.Request) {
		req.SetQueryParam(`bvid`, bvid)
	}
}

func setEPIDQueryParam(epID string) restyOption {
	return func(req *resty.Request) {
		req.SetQueryParam(`ep_id`, epID)
	}
}

func setSeasonIDQueryParam(seasonID string) restyOption {
	return func(req *resty.Request) {
		req.SetQueryParam(`season_id`, seasonID)
	}
}

func (r *request) fetchRawSeasonOrEpisode(region *Region, options ...restyOption) (*fetchSeasonResp, error) {
	u := r.getProxyURL(*region)

	if len(options) == 0 {
		return nil, errors.New(`no param present`)
	}

	req := r.client.R()
	for i := range options {
		options[i](req)
	}

	respRaw, err := req.
		Get(stringBuilder(u, SEASON_URI))
	if err != nil {
		return nil, err
	}

	var resp = fetchSeasonResp{}
	err = jsoniter.Unmarshal(respRaw.Body(), &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 && resp.Code != -404 {
		return nil, errors.New(resp.Message)
	}

	if resp.Code == -404 {
		*region = REGION_UNLOCATED
	}

	return &resp, nil
}

func (r *request) getProxyURL(region Region) string {
	var u string
	if region == REGION_UNLOCATED {
		u = DEFAULT_HOSTNAME_WITH_SCHEME
	} else {
		u = Config.getProxy(region).GetProxyURL()
		if u == `` {
			u = DEFAULT_HOSTNAME_WITH_SCHEME
		}
	}

	return u
}

func (r *request) fetchMedia(mediaID string) (*fetchMediaResp, error) {
	respRaw, err := r.client.R().
		SetQueryParam(`media_id`, mediaID).
		Get(stringBuilder(DEFAULT_HOSTNAME_WITH_SCHEME, MEDIA_URI))

	if err != nil {
		return nil, err
	}

	var resp = fetchMediaResp{}
	err = jsoniter.Unmarshal(respRaw.Body(), &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(resp.Message)
	}

	return &resp, nil
}

func (r *request) fetchVideoDownloadStreamURL(bvidStr, cidStr string) (*respVideoDetail, error) {
	respRaw, err := r.client.R().
		SetQueryParam(`bvid`, bvidStr).
		SetQueryParam(`cid`, cidStr).
		SetQueryParam(`fourk`, `1`).
		SetQueryParam(`fnval`, `4048`).
		SetQueryParam(`qn`, `127`).
		Get(stringBuilder(DEFAULT_HOSTNAME_WITH_SCHEME, PLAY_URI))

	if err != nil {
		return nil, err
	}

	var resp = respVideoDetail{}
	err = jsoniter.Unmarshal(respRaw.Body(), &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(`invalid result on requesting video download stream url`)
	}

	return &resp, nil
}

func (r *request) fetchEpisodeDownloadStreamURL(epid string, region Region) (*respVideoDetail, error) {
	u := r.getProxyURL(region)
	respRaw, err := r.client.R().
		SetQueryParam(`ep_id`, epid).
		SetQueryParam(`fourk`, `1`).
		SetQueryParam(`fnval`, `4048`).
		SetQueryParam(`qn`, `127`).
		Get(stringBuilder(u, PGC_URI))

	if err != nil {
		return nil, err
	}

	var resp = respSeasonDetail{}
	err = jsoniter.Unmarshal(respRaw.Body(), &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(`invalid result on requesting episode download stream url`)
	}

	return &respVideoDetail{
		Code:    0,
		Message: resp.Message,
		Data:    resp.Data,
	}, nil
}

func (r *request) fetchSubtitle(bvidStr, cidStr string) (*fetchPlayerInfoResp, error) {
	respRaw, err := r.client.R().
		SetQueryParam(`bvid`, bvidStr).
		SetQueryParam(`cid`, cidStr).
		Get(stringBuilder(DEFAULT_HOSTNAME_WITH_SCHEME, PLAYER_URI))

	if err != nil {
		return nil, err
	}

	var resp = fetchPlayerInfoResp{}
	err = jsoniter.Unmarshal(respRaw.Body(), &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(resp.Message)
	}

	return &resp, nil
}
