package railgun

const (
	DEFAULT_HOSTNAME_WITH_SCHEME = `https://api.bilibili.com`
	DEFAULT_HOSTNAME             = `api.bilibili.com`
	DEFAULT_FNVAL                = `4048`
	DEFAULT_QUALITY              = `127`
	DEFAULT_REFERER              = `https://www.bilibili.com`
	DEFAULT_USER_AGENT           = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/97.0.4692.99 Safari/537.36"
)

const (
	VIDEO_URI      = `/x/web-interface/view`
	SEASON_URI     = `/pgc/view/web/season`
	EPISODE_URI    = `/pgc/view/web/season`
	COLLECTION_URI = `/x/polymer/space/seasons_archives_list`
	MEDIA_URI      = `/pgc/review/user`
	PGC_URI        = `/pgc/player/web/playurl`
	PLAY_URI       = `/x/player/playurl`
	PLAYER_URI     = `/x/player/v2`
)

const (
	CHECK_COOKIE_URL   = `https://passport.bilibili.com/x/passport-login/web/cookie/info`
	CORESPONDING_URL   = `https://www.bilibili.com/correspond/1/`
	REFRESH_COOKIE_URL = `https://passport.bilibili.com/x/passport-login/web/cookie/refresh`
)

const (
	DEFAULT_REFRESH_COOKIE_SOURCE = `main_web`
)

const (
	NAV_INFO_URL = `https://api.bilibili.com/x/web-interface/nav`
)

const (
	AV_REGEX         = `(?:av|AV)\d+`
	BV_REGEX         = `((bv|BV)[0-9A-z]{10})`
	DYNAMIC_REGEX    = `(t\.bilibili\.com/(h5/dynamic/detail/)?)([0-9]{18})`
	ROOM_REGEX       = `(live\.bilibili\.com/)(\d+)`
	SHORT_LINK_REGEX = `((b23\.tv\/)[0-9A-z]+)`
	SPACE_REGEX      = `(space\.bilibili\.com/)(\d+)`
	SEASON_REGEX     = `(?:bilibili\.com/bangumi/play/ss)\d+`
	EPISODE_REGEX    = `(?:bilibili\.com/bangumi/play/ep)\d+`
	MEDIA_REGEX      = `(?:bilibili\.com/bangumi/media/md)\d+`
	COLLECTION_REGEX = `(?:channel/collectiondetail\?sid=)\d+`
)

const (
	NOTMATCHED DataType = -1
	SHORTLINK  DataType = 0
	AVVIDEO    DataType = 1
	BVVIDEO    DataType = 2
	SEASON     DataType = 3
	EPISODE    DataType = 4
	MEDIA      DataType = 5
	AUDIO      DataType = 6
	ARTICLE    DataType = 7
	COLLECTION DataType = 8
)

const (
	AVC_TITLE_SUFFIX  = ` AVC `
	AVC_CODEC_PREFIX  = `avc1`
	AV1_TITLE_SUFFIX  = ` AV1 `
	AV1_CODEC_PREFIX  = `av01`
	HEVC_CODEC_PREFIX = `hev1`

	FLAC_TITLE_SUFFIX     = ` Hi-Res `
	FLAC_FILENAME_SUFFIX  = `.flac`
	DOLBY_TITLE_SUFFIX    = ` Dolby Atmos `
	DOLBY_FILENAME_SUFFIX = `.ac3`
)

const (
	REGION_CN Region = iota
	REGION_HK
	REGION_TW
	REGION_TH
	REGION_UNLOCATED
)
