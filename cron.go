package main

import (
	"github.com/robfig/cron/v3"
	"ssh-news/provider"
)

func cronFetch() {
	provider.Fetch()
	c := cron.New(cron.WithSeconds())
	c.AddFunc("0 */10 * * * *", func() {
		provider.Fetch()
	})
	c.Start()
}
