package railgun

import (
	"errors"
	"math"
	"strconv"
	"time"
)

type Collection struct {
	collectionID     string
	collectionDetail *fetchCollectionResp

	client *request
}

func (c *Collection) GetType() DataType {
	return COLLECTION
}

func (c *Collection) setRequestClient(client *request) {
	c.client = client
}

func (c *Collection) GetDownloadArgs() ([]*DownloadArgs, error) {
	err := c.fetch()
	if err != nil {
		return nil, err
	}

	var results []*DownloadArgs

	for i := range c.collectionDetail.Data.Archives {
		time.Sleep(1 * time.Second)
		video, err := c.client.fetchVideo(``, c.collectionDetail.Data.Archives[i].Bvid)
		if err != nil {
			return nil, err
		}

		args, err := video.toDownloadArgs(c.client)
		if err != nil {
			return nil, err
		}
		results = append(results, args...)
	}

	return results, nil
}

func (c *Collection) fetch() error {
	resp, err := c.client.fetchCollection(c.collectionID, 1)
	if err != nil {
		return err
	}

	if resp.Code != 0 {
		return errors.New(resp.Message)
	}

	pages := math.Ceil(float64(resp.Data.Page.Total) / 100)
	for i := 2; float64(i) <= pages; i++ {
		part, err := c.client.fetchCollection(c.collectionID, i)
		if err != nil {
			return err
		}

		if part.Code != 0 {
			return errors.New(part.Message)
		}

		resp.Data.Aids = append(resp.Data.Aids, part.Data.Aids...)
		resp.Data.Archives = append(resp.Data.Archives, part.Data.Archives...)
	}

	c.collectionDetail = resp

	return nil
}

func (c *Collection) GetDetail() (*Detail, error) {
	err := c.fetch()
	if err != nil {
		return nil, err
	}

	return &Detail{
		Type:        COLLECTION,
		Picture:     c.collectionDetail.Data.Meta.Cover,
		Title:       c.collectionDetail.Data.Meta.Name,
		Description: c.collectionDetail.Data.Meta.Description,
		ShareURL: stringBuilder(
			`https://space.bilibili.com/`,
			strconv.Itoa(c.collectionDetail.Data.Meta.Mid),
			`/channel/collectiondetail?sid=`,
			strconv.Itoa(c.collectionDetail.Data.Meta.SeasonId),
		),
	}, nil
}

func newCollection(cid string) *Collection {
	return &Collection{
		collectionID: cid,
	}
}
