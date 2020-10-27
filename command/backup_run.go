package command

import (
	"github.com/backupshq/agent/actions"
	"github.com/backupshq/agent/api"
	"github.com/backupshq/agent/config"
	"github.com/backupshq/agent/log"
	"github.com/urfave/cli"
)

var BackupRun = cli.Command{
	Name:  "run",
	Usage: "Run a one-off backup and send the results to the API",
	Action: func(c *cli.Context) error {
		config := config.LoadCli(c)

		logger := log.CreateStdoutLogger(config.LogLevel.Level)
		client := api.NewClient(config)
		client.Authenticate()

		backup := client.GetBackup(c.Args().Get(0))
		if !backup.Managed {
			return cli.NewExitError("Cannot start an unmanaged backup using `run` command, try `start-unmanaged`.", 1)
		}

		actions.RunBackup(client, backup, logger, config)

		return nil
	},
}
