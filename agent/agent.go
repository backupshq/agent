package agent

import (
	"fmt"
	"time"

	"github.com/backupshq/agent/actions"
	"github.com/backupshq/agent/api"
	"github.com/backupshq/agent/config"
	"github.com/backupshq/agent/log"
	"github.com/robfig/cron/v3"
)

type Agent struct {
	logger        *log.Logger
	apiClient     *api.ApiClient
	config        *config.Config
	syncFrequency string
	backups       map[string]api.Backup
	crons         map[string]*cron.Cron
}

func Create(c *config.Config, syncFrequency string) *Agent {
	return &Agent{
		logger:        log.CreateStdoutLogger(c.LogLevel.Level),
		apiClient:     api.NewClient(c),
		config:        c,
		syncFrequency: syncFrequency,
		backups:       make(map[string]api.Backup),
		crons:         make(map[string]*cron.Cron),
	}
}

func (a *Agent) update() {
	a.logger.Debug("Checking for changes to backups...")
	backups := a.apiClient.ListBackups(api.BACKUP_TYPE_SCHEDULED)
	a.logger.Debug(fmt.Sprintf("Scheduled backups pulled from the API: %d", len(backups)))

	updatedCount := 0
	for id := range backups {
		if a.backups[id] != backups[id] {
			a.backups[id] = backups[id]
			if cron, ok := a.crons[id]; ok { // checks if there's an existing cron job for this backup
				cron.Stop()
			}
			a.crons[id] = actions.Schedule(a.apiClient, a.backups[id], a.logger)
			updatedCount++
			a.logger.Debug("Updated backup: " + backups[id].Name)
		}
	}

	if updatedCount > 0 {
		a.logger.Info(fmt.Sprintf("Updated %d backup definitions", updatedCount))
	} else {
		a.logger.Debug("No changes detected")
	}
}

func (a *Agent) Start() {
	a.logger.Info(`
========================
Starting BackupsHQ agent
========================
`)
	a.apiClient.Authenticate()

	a.logger.Info("Sync frequency: " + a.syncFrequency)
	a.update()
	cr := cron.New()
	cr.AddFunc(a.syncFrequency, func() {
		a.update()
	})
	cr.Start()

	time.Sleep(time.Minute * 100)
}
