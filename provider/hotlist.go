package provider

import (
	"errors"
	"github.com/parnurzeal/gorequest"
	"time"
)

type HotList struct {
	URL     string
	Sources []Source
}

func NewHotList() DataProvider {
	return &HotList{
		URL: "https://api.vvhan.com/api/hotlist/all",
	}
}

type HotListResp struct {
	Success bool `json:"success"`
	Data    []struct {
		Name       string `json:"name"`
		Subtitle   string `json:"subtitle"`
		UpdateTime string `json:"update_time"`
		Data       []struct {
			Index    int    `json:"index"`
			Title    string `json:"title"`
			Hot      string `json:"hot"`
			Url      string `json:"url"`
			MobilUrl string `json:"mobilUrl"`
		} `json:"data"`
	} `json:"data"`
}

func (l *HotList) Fetch() error {
	hotListResp := &HotListResp{}

	request := gorequest.New()
	_, _, errs := request.Get(l.URL).EndStruct(hotListResp)
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	var tmpSources []Source
	for _, source := range hotListResp.Data {
		var news []News
		for _, item := range source.Data {
			news = append(news, News{
				Index: item.Index,
				Title: item.Title,
				Hot:   item.Hot,
				URL:   item.Url,
			})
		}
		tmpSources = append(tmpSources, Source{
			Name:       source.Name,
			SubTitle:   source.Subtitle,
			UpdateTime: parseTime(source.UpdateTime),
			News:       news,
		})
	}

	l.Sources = tmpSources

	return nil
}

func (l *HotList) Get() []Source {
	return l.Sources
}

func parseTime(timeStr string) time.Time {
	t, err := time.Parse(time.DateTime, timeStr)
	if err != nil {
		// log.Println(err)
		return time.Time{}
	}
	return t
}
