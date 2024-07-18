package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/ssh"
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
	session  ssh.Session

	isCopiedNews    bool
	copiedNewsName  string
	copiedNewsValue list.Model
	copiedNewsIdx   int
}

const (
	selectedSources = iota
	selectedNews
)

func newModel(sources []provider.Source, session ssh.Session) model {
	var m model
	var sourcesItems []list.Item
	var newsItems = map[string][]list.Item{}

	for _, source := range sources {
		sourcesItems = append(sourcesItems, sourcesItem{name: fmt.Sprintf("%s", source.Name)})
		newsItems[source.Name] = []list.Item{}
		for _, item := range source.News {
			newsItems[source.Name] = append(newsItems[source.Name], newsItem{title: fmt.Sprintf("%d.%s", item.Index, removeControlChars(item.Title)), desc: item.Hot, url: item.URL})
		}
	}

	m.sources = list.New(sourcesItems, list.DefaultDelegate{
		ShowDescription: false,
		Styles:          list.NewDefaultItemStyles(),
	}, 20, 20)

	m.news = map[string]list.Model{}
	for k, v := range newsItems {
		m.news[k] = list.New(v, list.NewDefaultDelegate(), 0, 0)
	}

	m.sources.SetFilteringEnabled(false)
	m.sources.SetShowTitle(false)
	m.sources.SetShowStatusBar(false)
	m.sources.SetShowHelp(false)

	for k := range m.news {
		v := m.news[k]
		v.SetShowTitle(false)
		v.SetShowHelp(false)
		v.SetFilteringEnabled(false)
		v.SetShowStatusBar(false)
		m.news[k] = v
	}

	m.session = session
	return m
}

func removeControlChars(input string) string {
	// 定义一个正则表达式，匹配所有控制字符（包括换行符）
	re := regexp.MustCompile(`[\x00-\x1F\x7F]`)
	// 使用正则表达式替换掉所有匹配的字符
	return re.ReplaceAllString(input, "")
}
