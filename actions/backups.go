package actions

import (
	"fmt"
	"strings"

	"github.com/backupshq/agent/api"
	"github.com/backupshq/agent/config"
	"github.com/backupshq/agent/expression"
	"github.com/backupshq/agent/log"
	"github.com/backupshq/agent/utils"
	"io/ioutil"
	"os"
	"path/filepath"
)

func RunBackup(client *api.ApiClient, backup api.Backup, logger *log.Logger, config *config.Config) {
	job := client.StartJob(backup.ID)
	logger.Info(fmt.Sprintf("Starting a new job with id %q.", job.ID))
	status := "succeeded"

	tmpDir, err := ioutil.TempDir(os.TempDir(), "backupshq-")
	if err != nil {
		logger.Error(fmt.Sprintf("Cannot create temporary file: %s", err))
	}
	defer os.RemoveAll(tmpDir)
	for _, definition := range backup.StepDefinitions {
		scriptPath := filepath.Join(tmpDir, fmt.Sprintf("%d.sh", definition.SortOrder))
		if err := ioutil.WriteFile(scriptPath, []byte(definition.Script.Script), 0700); err != nil {
			logger.Error(fmt.Sprintf("Cannot create temporary file: %s", err))
		}
		step := client.CreateStep(job.ID, definition.Name, definition.SortOrder)

		env, err := evaluateExpressions(client, definition, logger, config)
		var out string
		if err == nil {
			logger.Debug(fmt.Sprintf(`Running backup command: "%s"`, scriptPath))
			out, err = utils.ExecuteCommand(scriptPath, env)
		}

		if err != nil {
			// In the future we can use this block to determine status code, but for now just send the error to the API
			logger.Warn(err.Error())
			out = err.Error()
			status = "failed"
		}
		client.SendLogs(step, out)
		client.FinishStep(step.ID, status)
		if status == "failed" {
			logger.Warn(fmt.Sprintf("Job step %d failed", step.SortOrder))
			break
		}
	}

	logger.Debug("Publishing job result to the API.")
	client.FinishJob(job, status)
}

func evaluateExpressions(client *api.ApiClient, definition api.StepDefinition, logger *log.Logger, config *config.Config) ([]string, error) {
	var env []string
	expressionManager := expression.CreateExpressionManager()
	context := expression.Context{
		map[string]string{},
		map[string]func(args ...string) string{
			"server_secret": func(args ...string) string {
				secret := client.GetSecretByName(args[0])

				return secret.Value
			},
			"client_secret": func(args ...string) string {
				fmt.Printf("%v", config.Secrets)
				val, ok := config.Secrets[strings.ToUpper(args[0])]
				if !ok {
					return "not found"
				}
				return val
			},
		},
	}

	for hash, expression := range definition.Expressions {
		logger.Debug(fmt.Sprintf("Evaluating expression %s: %s", hash, expression))
		result, err := expressionManager.Evaluate(expression, context)
		if err != nil {
			return env, err
		}
		env = append(env, fmt.Sprintf("EXPR_%s=%s", hash, result))
	}

	return env, nil
}
