package railgun

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type DataType int
type Region int

type Detail struct {
	Type        DataType
	Picture     string
	BV          string
	AV          int64
	Title       string
	Author      string
	CreateTime  int64
	PublicTime  int64
	Duration    int64
	Description string
	Dynamic     string
	//Season and Episode
	Evaluate                string
	Area                    string
	SerialStatusDescription string
	ShareURL                string
	MediaType               string
}

type DownloadArgs struct {
	BVID            string
	EpisodeID       string
	SeasonID        string
	CID             string
	Author          string
	Title           string
	SubTitle        string
	Subtitles       []*Subtitle
	VideoStreamList []*VideoStream
	AudioStreamList []*AudioStream
}

type Subtitle struct {
	Index       int
	Locale      string
	LocaleText  string
	SubtitleURL string
}

type Parser interface {
	GetType() DataType
	GetDetail() (*Detail, error)
	GetDownloadArgs() ([]*DownloadArgs, error)
	setRequestClient(client *request)
}

func (r *RailGun) getMediaTokenType(token string, checkCollection bool) (Parser, error) {
	token, mediaType, err := ParseBilibiliURL(token)
	if err != nil {
		return nil, err
	}

	if mediaType == SHORTLINK {
		token, err = r.client.getShortLinkLocation(stringBuilder(`https://`, token))
		if err != nil {
			return nil, err
		}

		token, mediaType, err = ParseBilibiliURL(token)
	}

	if mediaType == NOTMATCHED {
		return nil, errors.New(`invalid bilibili url`)
	}

	if checkCollection && (mediaType == AVVIDEO || mediaType == BVVIDEO) {
		collection, err := r.checkCollection(mediaType, token)
		if err != nil {
			return nil, err
		}

		if collection != 0 {
			mediaType = COLLECTION
			token = strconv.Itoa(collection)
		}
	}

	switch mediaType {
	case EPISODE:
		return newEpisode(token, checkCollection), nil
	case SEASON:
		return newSeason(token), nil
	case MEDIA:
		return newMedia(token), nil
	case AVVIDEO:
		return newAVVideo(token, checkCollection), nil
	case BVVIDEO:
		return newBVVideo(token, checkCollection), nil
	case COLLECTION:
		return newCollection(token), nil
	}

	return nil, errors.New(`invalid bilibili url`)
}

func (r *RailGun) checkCollection(Type DataType, token string) (int, error) {
	if Type == AVVIDEO {
		resp, err := r.client.fetchVideo(token, ``)
		if err != nil {
			return 0, err
		}

		return resp.Data.UgcSeason.Id, nil
	}

	if Type == BVVIDEO {
		resp, err := r.client.fetchVideo(``, token)
		if err != nil {
			return 0, err
		}

		return resp.Data.UgcSeason.Id, nil
	}

	return 0, nil
}

func ParseBilibiliURL(str string) (string, DataType, error) {
	switch {
	case strings.Contains(str, "b23.tv"):
		r, err := regexp.Compile(SHORT_LINK_REGEX)
		if err != nil {
			return "", NOTMATCHED, err
		}
		return r.FindString(str), SHORTLINK, err
	case strings.Contains(str, `bilibili.com/bangumi/play/ss`):
		r, err := regexp.Compile(SEASON_REGEX)
		if err != nil {
			return "", NOTMATCHED, err
		}
		return strings.TrimPrefix(r.FindString(str), `bilibili.com/bangumi/play/ss`), SEASON, err
	case strings.Contains(str, `bilibili.com/bangumi/play/ep`):
		r, err := regexp.Compile(EPISODE_REGEX)
		if err != nil {
			return "", NOTMATCHED, err
		}
		return strings.TrimPrefix(r.FindString(str), `bilibili.com/bangumi/play/ep`), EPISODE, err
	case strings.Contains(str, `bilibili.com/bangumi/media/md`):
		r, err := regexp.Compile(MEDIA_REGEX)
		if err != nil {
			return "", NOTMATCHED, err
		}
		return strings.TrimPrefix(r.FindString(str), `bilibili.com/bangumi/media/md`), MEDIA, err
	case strings.Contains(str, `av`):
		r, err := regexp.Compile(AV_REGEX)
		if err != nil {
			return "", NOTMATCHED, err
		}
		return strings.TrimPrefix(r.FindString(str), `av`), AVVIDEO, err
	case strings.Contains(str, `AV`):
		r, err := regexp.Compile(AV_REGEX)
		if err != nil {
			return "", NOTMATCHED, err
		}
		return strings.TrimPrefix(r.FindString(str), `AV`), AVVIDEO, err
	case strings.Contains(str, `BV`) || strings.Contains(str, `bv`):
		r, err := regexp.Compile(BV_REGEX)
		if err != nil {
			return "", NOTMATCHED, err
		}
		return r.FindString(str), BVVIDEO, err
	case strings.Contains(str, `channel/collectiondetail?sid=`):
		r, err := regexp.Compile(COLLECTION_REGEX)
		if err != nil {
			return "", NOTMATCHED, err
		}
		return strings.TrimPrefix(r.FindString(str), `channel/collectiondetail?sid=`), COLLECTION, err
	default:
		return "", NOTMATCHED, nil
	}
}
