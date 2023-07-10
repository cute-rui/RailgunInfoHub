package railgun

import "strconv"

type Video struct {
	avid            string
	bvid            string
	checkCollection bool
	collectionID    int
	videoDetail     *fetchVideoResp

	client *request
}

func (v *Video) GetDetail() (*Detail, error) {
	err := v.fetch()
	if err != nil {
		return nil, err
	}

	return &Detail{
		Type:        BVVIDEO,
		Picture:     v.videoDetail.Data.Pic,
		BV:          v.videoDetail.Data.Bvid,
		AV:          v.videoDetail.Data.Aid,
		Title:       v.videoDetail.Data.Title,
		Author:      v.videoDetail.Data.Owner.Name,
		CreateTime:  v.videoDetail.Data.Ctime,
		PublicTime:  v.videoDetail.Data.Pubdate,
		Duration:    v.videoDetail.Data.Duration,
		Description: v.videoDetail.Data.Desc,
		Dynamic:     v.videoDetail.Data.Dynamic,
	}, nil
}

func (v *Video) GetDownloadArgs() ([]*DownloadArgs, error) {
	err := v.fetch()
	if err != nil {
		return nil, err
	}

	if v.collectionID != 0 {
		collectionParser := newCollection(strconv.Itoa(v.collectionID))
		collectionParser.setRequestClient(v.client)
		return collectionParser.GetDownloadArgs()
	}

	return v.videoDetail.toDownloadArgs(v.client)
}

func (v *Video) setRequestClient(client *request) {
	v.client = client
}

func (v *Video) fetch() error {
	resp, err := v.client.fetchVideo(v.avid, v.bvid)

	if err != nil {
		return err
	}

	v.videoDetail = resp
	if v.checkCollection {
		if resp.Data.UgcSeason.Id != 0 {
			v.collectionID = resp.Data.UgcSeason.Id
		}
	}
	return nil
}

func (v *Video) GetType() DataType {
	if v.avid != "" {
		return AVVIDEO
	}

	return BVVIDEO
}

func newAVVideo(avid string, collection bool) *Video {
	return &Video{
		avid:            avid,
		checkCollection: collection,
	}
}

func newBVVideo(bvid string, collection bool) *Video {
	return &Video{
		bvid:            bvid,
		checkCollection: collection,
	}
}
