package agent

import (
	"github.com/backupshq/agent/actions"
	"github.com/backupshq/agent/api"
	"github.com/robfig/cron/v3"
)

func (a *Agent) configureSchedule(backup api.Backup) {
	if cron, ok := a.crons[backup.ID]; ok {
		cron.Stop()
	}
	if len(backup.Schedule) < 1 {
		return
	}
	c := cron.New()
	var cancelChannel = make(chan bool)
	c.AddFunc(backup.Schedule, func() {
		job := a.apiClient.StartJob(backup.ID)
		a.jobCancelChannels.Store(job.ID, cancelChannel)
		actions.RunJob(a.apiClient, backup, job, a.logger, a.config, cancelChannel)
		a.jobCancelChannels.Delete(job.ID)
	})
	a.crons[backup.ID] = c
	c.Start()
}
