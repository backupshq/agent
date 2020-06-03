package actions

import (
	"fmt"

	"github.com/backupshq/agent/api"
	"github.com/backupshq/agent/log"
	"github.com/backupshq/agent/utils"
)

func RunBackup(client *api.ApiClient, backup api.Backup, logger *log.Logger) {
	job := client.StartJob(backup.ID)
	logger.Info(fmt.Sprintf("Starting a new job with id %q.", job.ID))

	logger.Debug(fmt.Sprintf(`Running backup command: "%s"`, backup.Command))
	status := "succeeded"
	out, err := utils.ExecuteCommand(backup.Command)
	if err != nil {
		// In the future we can use this block to determine status code, but for now just send the error the the API
		logger.Warn(err.Error())
		out = err.Error()
		status = "failed"
	}
	logger.Debug(out)

	logger.Debug("Publishing job result to the API.")
	client.FinishJob(job, status)
	client.SendLogs(job, out)
}
