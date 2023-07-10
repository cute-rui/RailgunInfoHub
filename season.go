package railgun

type Season struct {
	seasonID  string
	episodeID string

	episodeQueryAll bool
	region          Region
	seasonDetail    *fetchSeasonResp

	client *request
}

func (s *Season) GetType() DataType {
	if s.seasonID != "" {
		return SEASON
	} else {
		return EPISODE
	}
}

func (s *Season) fetch() error {
	resp, region, err := s.client.fetchSeasonOrEpisode(
		s.seasonID,
		s.episodeID,
		Config.getDefaultRegion())
	if err != nil {
		return err
	}

	s.seasonDetail = resp
	s.region = region
	return nil
}

func (s *Season) GetDetail() (*Detail, error) {
	err := s.fetch()
	if err != nil {
		return nil, err
	}

	var area string
	for i := range s.seasonDetail.Result.Areas {
		stringJoiner(`，`, area, s.seasonDetail.Result.Areas[i].Name)
	}

	return &Detail{
		Type:                    SEASON,
		Picture:                 s.seasonDetail.Result.Cover,
		Title:                   s.seasonDetail.Result.Title,
		Evaluate:                s.seasonDetail.Result.Evaluate,
		Area:                    area,
		SerialStatusDescription: s.seasonDetail.Result.NewEp.Desc,
		ShareURL:                s.seasonDetail.Result.ShareUrl,
	}, nil
}

func (s *Season) GetDownloadArgs() ([]*DownloadArgs, error) {
	err := s.fetch()
	if err != nil {
		return nil, err
	}

	if s.episodeQueryAll {
		s.episodeID = ""
	}
	return s.seasonDetail.toDownloadArgs(s.episodeID, s.region, s.client)
}

func (s *Season) setRequestClient(client *request) {
	s.client = client
}

func newSeason(sid string) *Season {
	return &Season{
		seasonID: sid,
	}
}

func newEpisode(eid string, queryAll bool) *Season {
	return &Season{
		episodeID:       eid,
		episodeQueryAll: queryAll,
	}
}
