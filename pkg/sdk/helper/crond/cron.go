package crond

import (
	"github.com/robfig/cron/v3"
)

type CronLogger struct {
}

func (c CronLogger) Info(msg string, keysAndValues ...interface{}) {
	//slf.WithField("msg", msg).Debugln(keysAndValues...)
}

// Error logs an error condition.
func (c CronLogger) Error(err error, msg string, keysAndValues ...interface{}) {
}

func NewJob(job cron.Job) cron.Job {
	return cron.NewChain(cron.SkipIfStillRunning(CronLogger{})).Then(job)
}

func NewJobWrapper(f func()) cron.Job {
	return cron.NewChain(cron.SkipIfStillRunning(CronLogger{})).Then(&jobWrapper{f: f})
}

type jobWrapper struct {
	f func()
}

func (j *jobWrapper) Run() {
	j.f()
}
