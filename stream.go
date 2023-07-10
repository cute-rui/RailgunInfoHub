package railgun

import (
	"errors"
	"strconv"
	"strings"
)

type VideoStream struct {
	Quality  string
	VideoURL string
	AudioURL string
	Suffix   string
}

type AudioStream struct {
	Index    int
	AudioURL string
	Dolby    bool
	Flac     bool
}

func (r *fetchVideoResp) toDownloadArgs(request *request) ([]*DownloadArgs, error) {
	var argsList []*DownloadArgs

	for i := range r.Data.Pages {
		//Todo: timer
		d, err := request.fetchVideoDownloadStreamURL(r.Data.Bvid, strconv.Itoa(r.Data.Pages[i].Cid))
		if err != nil {
			return nil, err
		}

		s, err := request.fetchSubtitle(r.Data.Bvid, strconv.Itoa(r.Data.Pages[i].Cid))
		if err != nil {
			return nil, err
		}

		subtitles := s.getSubtitles()

		args := DownloadArgs{
			Title:     r.Data.Title,
			SubTitle:  r.Data.Pages[i].Part,
			CID:       strconv.Itoa(r.Data.Pages[i].Cid),
			BVID:      s.Data.Bvid,
			Author:    r.Data.Owner.Name,
			Subtitles: subtitles,
		}

		err = d.getStreamURL(&args)
		if err != nil {
			return nil, err
		}

		args.Subtitles = subtitles

		argsList = append(argsList, &args)
	}

	if len(argsList) == 0 {
		return nil, errors.New(`no video found`)
	}

	return argsList, nil
}

func (r *fetchSeasonResp) toDownloadArgs(epid string, region Region, request *request) ([]*DownloadArgs, error) {
	var argsList []*DownloadArgs

	for i := range r.Result.Episodes {
		if strconv.Itoa(r.Result.Episodes[i].Id) != epid && epid != `` {
			continue
		}

		d, err := request.fetchEpisodeDownloadStreamURL(
			strconv.Itoa(r.Result.Episodes[i].Id),
			region,
		)
		if err != nil {
			return nil, err
		}

		s, err := request.fetchSubtitle(r.Result.Episodes[i].Bvid, strconv.Itoa(r.Result.Episodes[i].Cid))
		if err != nil {
			return nil, err
		}

		subtitles := s.getSubtitles()

		args := DownloadArgs{
			Title:     r.Result.SeasonTitle,
			SubTitle:  stringJoiner(` - `, r.Result.Episodes[i].LongTitle, r.Result.Episodes[i].Title),
			EpisodeID: strconv.Itoa(r.Result.Episodes[i].Id),
			SeasonID:  strconv.Itoa(r.Result.SeasonId),
			CID:       strconv.Itoa(r.Result.Episodes[i].Cid),
			BVID:      s.Data.Bvid,
			Subtitles: subtitles,
		}

		err = d.getStreamURL(&args)
		if err != nil {
			return nil, err
		}

		args.Subtitles = subtitles

		argsList = append(argsList, &args)
	}

	if len(argsList) == 0 {
		return nil, errors.New(`no episode found`)
	}

	return argsList, nil
}

func (d *respVideoDetail) getAudioStreams() ([]*AudioStream, error) {
	var streamList []*AudioStream

	// Dolby Atmos
	if d.Data.Dash.Dolby.Type != 0 {
		var (
			index       int
			bestQuality int
		)
		for j := range d.Data.Dash.Dolby.Audio {
			if d.Data.Dash.Dolby.Audio[j].Id > bestQuality {
				index = j
				bestQuality = d.Data.Dash.Dolby.Audio[j].Id
			}
		}

		streamList = append(streamList, &AudioStream{
			Index:    len(streamList),
			AudioURL: d.Data.Dash.Dolby.Audio[index].BaseUrl,
			Dolby:    true,
			Flac:     false,
		})
	}

	// Flac
	if d.Data.Dash.Flac.Audio.Id == 30251 {
		streamList = append(streamList, &AudioStream{
			Index:    len(streamList),
			AudioURL: d.Data.Dash.Flac.Audio.BaseUrl,
			Dolby:    false,
			Flac:     true,
		})

		return streamList, nil
	}

	// Normal audio track
	tmp, audioIndex := 0, 0
	for j := range d.Data.Dash.Audio {
		if d.Data.Dash.Audio[j].Id > tmp {
			tmp = d.Data.Dash.Audio[j].Id
			audioIndex = j
		}
	}

	streamList = append(streamList, &AudioStream{
		Index:    len(streamList),
		AudioURL: d.Data.Dash.Audio[audioIndex].BaseUrl,
		Dolby:    false,
		Flac:     false,
	})

	if len(streamList) == 0 {
		return nil, errors.New(`no audio found`)
	}

	return streamList, nil
}

// Should only provide dolby atmos track by default
func (d *respVideoDetail) getDolbyVisionVideo(audioList []*AudioStream) ([]*VideoStream, error) {
	var DolbyList []*VideoStream
	DolbyVision := VideoStream{
		Quality: getQualityString(126),
		Suffix:  ` Dolby Vision `,
	}

	for i := range audioList {
		if audioList[i].Dolby {
			DolbyVision.AudioURL = audioList[i].AudioURL
			break
		}
	}

	if DolbyVision.AudioURL == `` {
		return nil, errors.New(`no dolby vision audio found`)
	}

	var avc VideoStream
	for i := range d.Data.Dash.Video {

		if d.Data.Dash.Video[i].Id == 126 {
			suffix := DolbyVision.Suffix
			if strings.Contains(d.Data.Dash.Video[i].Codecs, `avc1`) {
				avc = VideoStream{
					Quality:  DolbyVision.Quality,
					VideoURL: d.Data.Dash.Video[i].BaseUrl,
					AudioURL: DolbyVision.AudioURL,
					Suffix:   stringJoiner(`-`, suffix, ` AVC `),
				}

				continue
			} else if strings.Contains(d.Data.Dash.Video[i].Codecs, `av01`) {
				suffix = stringJoiner(`-`, suffix, ` AV1 `)
			}

			DolbyList = append(DolbyList, &VideoStream{
				Quality:  DolbyVision.Quality,
				VideoURL: d.Data.Dash.Video[i].BaseUrl,
				AudioURL: DolbyVision.AudioURL,
				Suffix:   suffix,
			})
		}
	}

	if len(DolbyList) == 0 {
		if avc.VideoURL == `` {
			return nil, errors.New(`no dolby vision video found`)
		}

		DolbyList = append(DolbyList, &avc)
	}

	return DolbyList, nil
}

func (d *respVideoDetail) getHDRVideo(audioList []*AudioStream) ([]*VideoStream, error) {
	var VideoList []*VideoStream
	HDR := VideoStream{
		Quality: getQualityString(125),
		Suffix:  ` HDR `,
	}

	avcList, err := d.getAVCVideo(audioList, 125, HDR.Suffix)
	if err != nil {
		return nil, err
	}

	av1List, err := d.getAV1Video(audioList, 125, HDR.Suffix)
	if err != nil {
		return nil, err
	}

	hevcList, err := d.getHEVCVideo(audioList, 125, HDR.Suffix)
	if err != nil {
		return nil, err
	}

	VideoList = append(VideoList, av1List...)
	VideoList = append(VideoList, hevcList...)

	if len(VideoList) == 0 {
		if len(avcList) == 0 {
			return nil, errors.New(`no hdr video found`)
		}

		VideoList = append(VideoList, avcList...)
	}

	return VideoList, nil
}

func (d *respVideoDetail) getAV1Video(audioList []*AudioStream, quality int, suffix string) ([]*VideoStream, error) {
	return d.getVideoStreams(audioList, quality, stringJoiner(`-`, suffix, AV1_TITLE_SUFFIX), AV1_CODEC_PREFIX)
}

func (d *respVideoDetail) getHEVCVideo(audioList []*AudioStream, quality int, suffix string) ([]*VideoStream, error) {
	return d.getVideoStreams(audioList, quality, suffix, HEVC_CODEC_PREFIX)
}

func (d *respVideoDetail) getAVCVideo(audioList []*AudioStream, quality int, suffix string) ([]*VideoStream, error) {
	return d.getVideoStreams(audioList, quality, stringJoiner(`-`, suffix, AVC_TITLE_SUFFIX), AVC_CODEC_PREFIX)
}

func (d *respVideoDetail) getVideoStreams(audioList []*AudioStream, quality int, suffix, codec string) ([]*VideoStream, error) {
	template := VideoStream{
		Quality: getQualityString(quality),
		Suffix:  suffix,
	}

	for i := range d.Data.Dash.Video {
		if d.Data.Dash.Video[i].Id != quality {
			continue
		}

		if !strings.Contains(d.Data.Dash.Video[i].Codecs, codec) {
			continue
		}

		template.VideoURL = d.Data.Dash.Video[i].BaseUrl
	}

	if template.VideoURL == `` {
		return nil, nil
	}

	var Result []*VideoStream
	for i := range audioList {
		if audioList[i].Dolby {
			Result = append(Result, &VideoStream{
				Quality:  template.Quality,
				VideoURL: template.VideoURL,
				AudioURL: audioList[i].AudioURL,
				Suffix:   stringJoiner(`-`, template.Suffix, ` Dolby Atmos `),
			})
			continue
		}

		if audioList[i].Flac {
			Result = append(Result, &VideoStream{
				Quality:  template.Quality,
				VideoURL: template.VideoURL,
				AudioURL: audioList[i].AudioURL,
				Suffix:   stringJoiner(`-`, template.Suffix, ` Hi-Res `),
			})
			continue
		}

		Result = append(Result, &VideoStream{
			Quality:  template.Quality,
			VideoURL: template.VideoURL,
			AudioURL: audioList[i].AudioURL,
			Suffix:   template.Suffix,
		})
	}

	if len(Result) == 0 {
		return nil, errors.New(`no audio found`)
	}

	return Result, nil
}

func (d *respVideoDetail) getStreamURL(downloadArg *DownloadArgs) (err error) {
	audioList, err := d.getAudioStreams()
	if err != nil {
		return
	}

	var videoList []*VideoStream
	quality, qualitySP := getBestQualityVideo(d.Data.AcceptQuality)
	if qualitySP != QualitySP_None {
		switch qualitySP {
		case QualitySP_DolbyAndHDR:
			hdrList, err := d.getHDRVideo(audioList)
			if err != nil {
				return err
			}
			videoList = append(videoList, hdrList...)

			fallthrough
		case QualitySP_DolbyVision:
			dolby, err := d.getDolbyVisionVideo(audioList)
			if err != nil {
				return err
			}
			videoList = append(videoList, dolby...)
		case QualitySP_HDR:
			hdrList, err := d.getHDRVideo(audioList)
			if err != nil {
				return err
			}
			videoList = append(videoList, hdrList...)
		}
	}

	avcList, err := d.getAVCVideo(audioList, quality, ``)
	if err != nil {
		return
	}

	av1List, err := d.getAV1Video(audioList, quality, ``)
	if err != nil {
		return
	}

	hevcList, err := d.getHEVCVideo(audioList, quality, ``)
	if err != nil {
		return
	}

	videoList = append(videoList, av1List...)
	videoList = append(videoList, hevcList...)

	if len(videoList) == 0 {
		if len(avcList) != 0 {
			downloadArg.AudioStreamList = audioList
			downloadArg.VideoStreamList = avcList
			return
		}

		return errors.New(`no video found`)
	}

	downloadArg.AudioStreamList = audioList
	downloadArg.VideoStreamList = videoList
	return
}
