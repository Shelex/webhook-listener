package scheduler

import "github.com/mileusna/crontab"

type Scheduler struct {
	tab *crontab.Crontab
}

func NewScheduler() *Scheduler {
	ctab := crontab.New()
	return &Scheduler{
		tab: ctab,
	}
}

func (cron *Scheduler) Schedule(cronSyntax string, fn interface{}) {
	cron.tab.MustAddJob(cronSyntax, fn)
}
