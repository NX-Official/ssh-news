package provider

import (
	"log"
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
	NewHotList(),
}

func Fetch() {
	for _, p := range providers {
		err := p.Fetch()
		if err != nil {
			log.Println(err)
		}
		sources = append(sources, p.Get()...)
	}
}

func Get() []Source {
	return sources
}
