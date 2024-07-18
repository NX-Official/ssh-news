package provider

import (
	"time"
)

type DataProvider interface {
	Fetch() error
	Get() []Source
}

var sources []Source

type Source struct {
	Name       string
	SubTitle   string
	UpdateTime time.Time
	News       []News
}

type News struct {
	Index int
	Title string
	Hot   string
	URL   string
}

var providers = []DataProvider{
	NewV2EX(),
	NewHotList(),
}

func Fetch() {
	var newSources []Source
	for _, p := range providers {
		err := p.Fetch()
		if err != nil {
			// log.Println(err)
		}
		newSources = append(newSources, p.Get()...)
	}
	sources = newSources
}

func Get() []Source {
	return sources
}
