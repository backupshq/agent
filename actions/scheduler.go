package actions

import (
	"github.com/backupshq/agent/api"
	"github.com/backupshq/agent/config"
	"github.com/backupshq/agent/log"
	"github.com/robfig/cron/v3"
)

func Schedule(client *api.ApiClient, backup api.Backup, logger *log.Logger, config *config.Config) *cron.Cron {
	logger.Info("Schedule " + backup.Name + " for " + backup.Schedule)
	c := cron.New()
	var cancelChannel = make(chan bool)
	c.AddFunc(backup.Schedule, func() { RunBackup(client, backup, logger, config, cancelChannel) })
	c.Start()
	return c
}
