package provider

import (
	"errors"
	"fmt"
	"github.com/mmcdole/gofeed"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type rss struct {
	name string
	url  string
}

type V2EX struct {
	RSS     []rss
	Sources []Source
}

func NewV2EX() DataProvider {
	return &V2EX{
		RSS: []rss{
			{
				name: "V2EX全站",
				url:  "https://www.v2ex.com/index.xml",
			},
			{
				name: "V2EX技术",
				url:  "https://www.v2ex.com/feed/tab/tech.xml",
			},
			{
				name: "V2EX创意",
				url:  "https://www.v2ex.com/feed/tab/creative.xml",
			},
			{
				name: "V2EX好玩",
				url:  "https://www.v2ex.com/feed/tab/play.xml",
			},
			{
				name: "V2EX问答",
				url:  "https://www.v2ex.com/feed/tab/qna.xml",
			},
		},
	}
}

func (v *V2EX) Fetch() error {
	var returnErrs []error
	for _, rss := range v.RSS {
		request := gorequest.New()
		resp, _, errs := request.Get(rss.url).End()
		if len(errs) > 0 {
			returnErrs = append(returnErrs, errs...)
			continue
		}
		feed, err := gofeed.NewParser().Parse(resp.Body)
		if err != nil {
			returnErrs = append(returnErrs, err)
			continue
		}

		var news []News
		for idx, item := range feed.Items {
			news = append(news, News{
				Index: idx + 1,
				Title: item.Title,
				Hot:   fmt.Sprintf("%d 回复", getReplyNumber(item.Link)),
				URL:   item.Link,
			})
		}
		v.Sources = append(v.Sources, Source{
			Name:       rss.name,
			SubTitle:   "",
			UpdateTime: time.Now(),
			News:       news,
		})
	}
	return errors.Join(returnErrs...)
}

func (v *V2EX) Get() []Source {
	return v.Sources
}

// https://www.v2ex.com/t/1058380#reply2 -> 2
func getReplyNumber(v2exURL string) int {
	URL, err := url.Parse(v2exURL)
	if err != nil {
		log.Println(err)
		return 0
	}
	fragment := URL.Fragment
	countString := strings.TrimLeft(fragment, "reply")
	count, err := strconv.Atoi(countString)
	if err != nil {
		log.Println(err)
		return 0
	}
	return count
}
