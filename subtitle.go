package railgun

import "strings"

func (s *fetchPlayerInfoResp) getSubtitles() (subtitles []*Subtitle) {
	for i := range s.Data.Subtitle.Subtitles {
		if strings.HasPrefix(s.Data.Subtitle.Subtitles[i].Lan, `ai`) {
			continue
		}
		url := s.Data.Subtitle.Subtitles[i].SubtitleUrl
		if strings.HasPrefix(url, `//`) {
			url = stringBuilder(`https:`, url)
		}
		subtitles = append(subtitles, &Subtitle{
			Index:       int(int32(i)),
			Locale:      s.Data.Subtitle.Subtitles[i].Lan,
			LocaleText:  s.Data.Subtitle.Subtitles[i].LanDoc,
			SubtitleURL: url,
		})
	}

	return
}
