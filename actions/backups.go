package actions

import (
	"fmt"

	"github.com/backupshq/agent/api"
	"github.com/backupshq/agent/log"
	"github.com/backupshq/agent/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

func RunBackup(client *api.ApiClient, backup api.Backup, logger *log.Logger) {
	job := client.StartJob(backup.ID)
	logger.Info(fmt.Sprintf("Starting a new job with id %q.", job.ID))
	status := "succeeded"

	tmpDir, err := ioutil.TempDir(os.TempDir(), "backupshq-")
	if err != nil {
		logger.Error(fmt.Sprintf("Cannot create temporary file ", err))
	}
	defer os.RemoveAll(tmpDir)
	for _, definition := range backup.StepDefinitions {
		scriptPath := filepath.Join(tmpDir, fmt.Sprintf("%d.sh", definition.SortOrder))
		if err := ioutil.WriteFile(scriptPath, []byte(definition.Script.Script), 0700); err != nil {
			logger.Error(fmt.Sprintf("Cannot create temporary file ", err))
		}
		step := client.CreateStep(job.ID, definition.Name, definition.SortOrder)

		logger.Debug(fmt.Sprintf(`Running backup command: "%s"`, scriptPath))
		out, err := utils.ExecuteCommand(scriptPath, []string{})
		if err != nil {
			// In the future we can use this block to determine status code, but for now just send the error to the API
			logger.Warn(err.Error())
			out = err.Error()
			status = "failed"
		}
		logger.Debug("\n" + out)
		client.SendLogs(step, out)
		client.FinishStep(step.ID, status)
		logger.Debug(step.ID)
		if status == "failed" {
			logger.Warn(fmt.Sprintf("Job step %d failed", step.SortOrder))
			break
		}
	}

	logger.Debug("Publishing job result to the API.")
	client.FinishJob(job, status)
}
