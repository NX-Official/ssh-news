package provider

import (
	"github.com/charmbracelet/log"
	"sync"
	"time"
)

type DataProvider interface {
	Fetch() error
	Get() []Source
}

var sources []Source
var mutex = sync.RWMutex{}

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
			log.Error(err)
		}
		newSources = append(newSources, p.Get()...)
	}
	mutex.Lock()
	defer mutex.Unlock()
	sources = newSources
}

func Get() []Source {
	mutex.RLock()
	defer mutex.RUnlock()
	return sources
}
