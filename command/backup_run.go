package command

import (
	"log"

	"../actions"
	"../api"
	"../config"
	"github.com/urfave/cli"
)

var BackupRun = cli.Command{
	Name:  "run",
	Usage: "Run a one-off backup and send the results to the API",
	Action: func(c *cli.Context) error {
		config := config.LoadCli(c)

		client := api.NewClient(config)

		backup := client.GetBackup(c.Args().Get(0))
		if backup.Type == api.BACKUP_TYPE_UNMANAGED {
			log.Fatal("Cannot start unmanaged backup using `run` command, try `start-unmanaged`")
		}

		actions.RunBackup(client, backup)

		return nil
	},
}
