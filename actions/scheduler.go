package actions

import (
	"log"

	"../api"
	"github.com/robfig/cron"
)

func Schedule(client *api.ApiClient, backup api.Backup) *cron.Cron {
	log.Println("Schedule " + backup.Name + " for " + backup.Schedule)
	c := cron.New()
	c.AddFunc(backup.Schedule, func() { RunBackup(client, backup) })
	c.Start()
	return c
}
