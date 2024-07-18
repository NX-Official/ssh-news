package main

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var selectedTitleStyle = lipgloss.NewStyle().Inherit(list.NewDefaultItemStyles().
	SelectedTitle).Background(lipgloss.Color("250"))
var selectedDescStyle = lipgloss.NewStyle().Inherit(list.NewDefaultItemStyles().
	SelectedDesc).Background(lipgloss.Color("250"))

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

var m model

func UpdateModel() {
	var newModel model
	var sourcesItems []list.Item
	var newsItems = map[string][]list.Item{}

	for _, source := range news.Data {
		sourcesItems = append(sourcesItems, sourcesItem{name: source.Name})
		newsItems[source.Name] = []list.Item{}
		for _, item := range source.Data {
			newsItems[source.Name] = append(newsItems[source.Name], newsItem{title: item.Title, desc: item.Hot, url: item.Url})
		}
	}

	newModel.sources = list.New(sourcesItems, list.DefaultDelegate{
		ShowDescription: false,
		Styles:          list.NewDefaultItemStyles(),
	}, 40, 20)

	newModel.news = map[string]list.Model{}
	for k, v := range newsItems {
		newModel.news[k] = list.New(v, list.NewDefaultDelegate(), 100, 50)
	}
	m = newModel
}
