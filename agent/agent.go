package agent

import (
	"fmt"
	"sync"
	"time"

	"github.com/backupshq/agent/api"
	"github.com/backupshq/agent/config"
	"github.com/backupshq/agent/log"
	"github.com/backupshq/agent/utils"
	"github.com/robfig/cron/v3"
)

type Agent struct {
	logger            *log.Logger
	apiClient         *api.ApiClient
	config            *config.Config
	principal         api.Principal
	account           api.Account
	token             api.AgentToken
	backups           map[string]api.Backup
	crons             map[string]*cron.Cron
	workerQueue       chan api.Job
	jobCancelChannels *utils.SignalMap
}

func Create(c *config.Config) *Agent {
	return &Agent{
		logger:            log.CreateStdoutLogger(c.LogLevel.Level),
		apiClient:         api.NewClient(c),
		config:            c,
		backups:           make(map[string]api.Backup),
		crons:             make(map[string]*cron.Cron),
		workerQueue:       make(chan api.Job, 10),
		jobCancelChannels: &utils.SignalMap{},
	}
}

func (a *Agent) ping() {
	pingResponse := a.apiClient.Ping(a.token)
	a.logger.Debug("Ping")

	if len(pingResponse.AssignedJobs) > 0 {
		for _, job := range pingResponse.AssignedJobs {
			if job.Status == "pending" {
				a.workerQueue <- job
			}
			if job.Status == "cancelled" {
				if cancel, ok := a.jobCancelChannels.Load(job.ID); ok {
					cancel <- true
				}
			}
		}
	}

	if pingResponse.UpdatedBackupCount > 0 {
		a.logger.Debug(fmt.Sprintf("Ping response returned %d updated backups", pingResponse.UpdatedBackupCount))
		a.update()
	}
}

func (a *Agent) update() {
	backups := a.apiClient.ListBackups(a.principal.ID)
	a.logger.Debug(fmt.Sprintf("Fetched %d managed backups assigned to this agent", len(backups)))

	for id := range backups {
		fullBackup := a.apiClient.GetBackup(id)
		if a.backups[id].UpdatedAt != fullBackup.UpdatedAt {
			a.logger.Debug(fmt.Sprintf("Updated definition of %s", fullBackup.Name))
			a.backups[id] = fullBackup

			a.configureSchedule(fullBackup)
		}
	}
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
			a,
		)

		go worker.work(a.workerQueue)
	}

	wg.Wait()
}
