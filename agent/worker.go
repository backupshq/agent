package agent

import (
	"fmt"
	"github.com/backupshq/agent/api"
	"github.com/backupshq/agent/log"
)

type AgentWorker struct {
	number    int
	logger    *log.Logger
	apiClient *api.ApiClient
}

func CreateWorker(n int, l *log.Logger, c *api.ApiClient) *AgentWorker {
	return &AgentWorker{
		number:    n,
		logger:    l,
		apiClient: c,
	}
}

func (w *AgentWorker) work(c <-chan int) {
	for i := range c {
		w.handleMessage(i)
	}
}

func (w *AgentWorker) handleMessage(i int) {
	w.logger.Debug(fmt.Sprintf("%d Handling message %d", w.number, i))
}
