package agent

import (
	"github.com/backupshq/agent/actions"
	"github.com/backupshq/agent/api"
)

type AgentWorker struct {
	number int
	agent  *Agent
}

func CreateWorker(number int, agent *Agent) *AgentWorker {
	return &AgentWorker{
		number: number,
		agent:  agent,
	}
}

func (w *AgentWorker) work(c <-chan api.Job) {
	for i := range c {
		w.handleMessage(i)
	}
}

func (w *AgentWorker) handleMessage(job api.Job) {
	w.agent.apiClient.StartExistingJob(job)
	var cancelChannel = make(chan bool)
	w.agent.jobCancelChannels.Store(job.ID, cancelChannel)
	actions.RunJob(w.agent.apiClient, w.agent.backups[job.BackupID], job, w.agent.logger, w.agent.config, cancelChannel)
	w.agent.jobCancelChannels.Delete(job.ID)
}
