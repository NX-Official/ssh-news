package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"regexp"
	"ssh-news/provider"
)

type sourcesItem struct {
	name string
}

func (i sourcesItem) Title() string       { return i.name }
func (i sourcesItem) Description() string { return "" }
func (i sourcesItem) FilterValue() string { return "" }

type newsItem struct {
	title string
	desc  string
	url   string
}

func (i newsItem) Title() string       { return i.title }
func (i newsItem) Description() string { return i.desc }
func (i newsItem) URL() string         { return i.url }
func (i newsItem) FilterValue() string { return "" }

type model struct {
	sources  list.Model
	news     map[string]list.Model
	selected int
}

const (
	selectedSources = iota
	selectedNews
)

var m model

func UpdateModel(sources []provider.Source) {
	var newModel model
	var sourcesItems []list.Item
	var newsItems = map[string][]list.Item{}

	for _, source := range sources {
		sourcesItems = append(sourcesItems, sourcesItem{name: fmt.Sprintf("%s", source.Name)})
		newsItems[source.Name] = []list.Item{}
		for _, item := range source.News {
			newsItems[source.Name] = append(newsItems[source.Name], newsItem{title: fmt.Sprintf("%d.%s", item.Index, removeControlChars(item.Title)), desc: item.Hot, url: item.URL})
		}
	}

	newModel.sources = list.New(sourcesItems, list.DefaultDelegate{
		ShowDescription: false,
		Styles:          list.NewDefaultItemStyles(),
	}, 20, 20)

	newModel.news = map[string]list.Model{}
	for k, v := range newsItems {
		newModel.news[k] = list.New(v, list.NewDefaultDelegate(), 0, 0)
	}
	m = newModel

}

func removeControlChars(input string) string {
	// 定义一个正则表达式，匹配所有控制字符（包括换行符）
	re := regexp.MustCompile(`[\x00-\x1F\x7F]`)
	// 使用正则表达式替换掉所有匹配的字符
	return re.ReplaceAllString(input, "")
}
