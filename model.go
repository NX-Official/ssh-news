package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
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

func UpdateModel() {
	var newModel model
	var sourcesItems []list.Item
	var newsItems = map[string][]list.Item{}

	for _, source := range news.Data {
		sourcesItems = append(sourcesItems, sourcesItem{name: source.Name})
		newsItems[source.Name] = []list.Item{}
		for _, item := range source.Data {
			newsItems[source.Name] = append(newsItems[source.Name], newsItem{title: fmt.Sprintf("%d.%s", item.Index, item.Title), desc: item.Hot, url: item.Url})
		}
	}

	newModel.sources = list.New(sourcesItems, list.DefaultDelegate{
		ShowDescription: false,
		Styles:          list.NewDefaultItemStyles(),
	}, 40, 21)

	newModel.news = map[string]list.Model{}
	for k, v := range newsItems {
		newModel.news[k] = list.New(v, list.NewDefaultDelegate(), 80, 21)
	}
	m = newModel
}
