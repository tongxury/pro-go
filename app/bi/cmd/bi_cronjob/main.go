package main

import (
	"github.com/robfig/cron/v3"
	"store/pkg/sdk/helper/crond"
)

func main() {

	c := cron.New(cron.WithSeconds())

	_, _ = c.AddJob("@every 1m", crond.NewJobWrapper(func() {

	}))

	c.Run()
}
