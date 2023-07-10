package railgun

import "strconv"

type Media struct {
	mediaID     string
	mediaDetail *fetchMediaResp

	client *request
}

func (m *Media) GetType() DataType {
	return MEDIA
}

func (m *Media) GetDetail() (*Detail, error) {
	err := m.fetch()
	if err != nil {
		return nil, err
	}

	var area string
	for i := range m.mediaDetail.Result.Media.Areas {
		stringJoiner(`，`, area, m.mediaDetail.Result.Media.Areas[i].Name)
	}

	return &Detail{
		Type:                    MEDIA,
		Picture:                 m.mediaDetail.Result.Media.Cover,
		Title:                   m.mediaDetail.Result.Media.Title,
		Area:                    area,
		SerialStatusDescription: m.mediaDetail.Result.Media.NewEp.IndexShow,
		ShareURL:                m.mediaDetail.Result.Media.ShareUrl,
		MediaType:               m.mediaDetail.Result.Media.TypeName,
	}, nil
}

func (m *Media) GetDownloadArgs() ([]*DownloadArgs, error) {
	err := m.fetch()
	if err != nil {
		return nil, err
	}

	parser := newSeason(strconv.Itoa(m.mediaDetail.Result.Media.SeasonId))
	parser.setRequestClient(m.client)
	return parser.GetDownloadArgs()
}

func (m *Media) fetch() error {
	resp, err := m.client.fetchMedia(m.mediaID)

	if err != nil {
		return err
	}

	m.mediaDetail = resp
	return nil
}

func (m *Media) setRequestClient(client *request) {
	m.client = client
}

func newMedia(mid string) *Media {
	return &Media{
		mediaID: mid,
	}
}
