package agent

import (
	"fmt"
	"sync"
	"time"

	"github.com/backupshq/agent/actions"
	"github.com/backupshq/agent/api"
	"github.com/backupshq/agent/config"
	"github.com/backupshq/agent/log"
	"github.com/robfig/cron/v3"
)

type Agent struct {
	logger      *log.Logger
	apiClient   *api.ApiClient
	config      *config.Config
	principal   api.Principal
	account     api.Account
	token       api.AgentToken
	backups     map[string]api.Backup
	crons       map[string]*cron.Cron
	workerQueue chan int
}

func Create(c *config.Config) *Agent {
	return &Agent{
		logger:      log.CreateStdoutLogger(c.LogLevel.Level),
		apiClient:   api.NewClient(c),
		config:      c,
		backups:     make(map[string]api.Backup),
		crons:       make(map[string]*cron.Cron),
		workerQueue: make(chan int, 10),
	}
}

func (a *Agent) ping() {
	a.logger.Debug("Checking for changes to backups...")
	pingResponse := a.apiClient.Ping(a.token)

	if len(pingResponse.AssignedPendingJobs) > 0 {
		for i, _ := range pingResponse.AssignedPendingJobs {
			a.workerQueue <- i
		}
	}

	if pingResponse.UpdatedBackupCount > 0 {
		a.logger.Debug("Changes to backups found... Syncing...")
		a.update()
		return
	}
	a.logger.Debug("No changes found")
}

func (a *Agent) update() {
	backups := a.apiClient.ListBackups(a.principal.ID)
	a.logger.Debug(fmt.Sprintf("Scheduled backups pulled from the API: %d", len(backups)))

	updatedCount := 0
	for id := range backups {
		fullBackup := a.apiClient.GetBackup(id)
		if a.backups[id].UpdatedAt != fullBackup.UpdatedAt {
			a.backups[id] = fullBackup

			if cron, ok := a.crons[id]; ok { // checks if there's an existing cron job for this backup
				cron.Stop()
			}
			a.crons[id] = actions.Schedule(a.apiClient, a.backups[id], a.logger, a.config)

			updatedCount++
			a.logger.Debug("Updated backup: " + fullBackup.Name)
		}
	}

	a.logger.Info(fmt.Sprintf("Updated %d backup definitions", updatedCount))
}

func (a *Agent) Start() {
	a.logger.Info(`
========================
Starting BackupsHQ agent
========================
`)
	a.apiClient.Authenticate()
	tokenInfo := a.apiClient.GetCurrentToken()
	a.principal = a.apiClient.GetPrincipal(tokenInfo.PrincipalId)
	a.logger.Info(fmt.Sprintf(`Authenticated as principal %s "%s"`, a.principal.ID, a.principal.Name))
	a.account = a.apiClient.GetAccount(tokenInfo.AccountId)
	a.logger.Info(fmt.Sprintf(`This agent belongs to account %s "%s"`, a.account.ID, a.account.Name))
	a.token = a.apiClient.Register()

	a.update()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			a.ping()
			time.Sleep(time.Second * 10)
		}
	}()

	for i := 1; i < 5; i++ {
		worker := CreateWorker(
			i,
			a.logger,
			a.apiClient,
		)

		go worker.work(a.workerQueue)
	}

	wg.Wait()
}
