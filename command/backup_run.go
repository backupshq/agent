package command

import (
	"log"
	"net/http"
	"time"

	"../actions"
	"../auth"
	"../config"
	"../utils"
	"github.com/urfave/cli"
)

var BackupRun = cli.Command{
	Name:  "run",
	Usage: "Run a one-off backup and send the results to the API",
	Action: func(c *cli.Context) error {
		config := config.LoadCli(c)

		tokenResponse, err := auth.GetAccessToken(config)
		if err != nil {
			log.Fatal(err)
		}

		client := &http.Client{
			Timeout: time.Second * 3,
		}

		backup := actions.GetBackup(client, tokenResponse, c.Args().Get(0))
		if backup.Type == actions.BACKUP_TYPE_UNMANAGED {
			log.Fatal("Cannot start unmanaged backup using `run` command, try `start-unmanaged`")
		}

		utils.RunBackup(client, tokenResponse, backup)

		return nil
	},
}
