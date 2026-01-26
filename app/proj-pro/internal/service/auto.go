package service

import (
	"store/pkg/sdk/helper/crond"

	"github.com/robfig/cron/v3"
)

func (t ProjService) Start() {

	c := cron.New(cron.WithSeconds())

	_, _ = c.AddJob("@every 1s", crond.NewJobWrapper(t.create))

}

func (t ProjService) create() {

}
