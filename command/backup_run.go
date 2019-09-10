package command

import (
	"fmt"
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

		job := actions.StartJob(client, tokenResponse, c.Args().Get(0))
		fmt.Printf("Started a new job with id %q.\n", job.ID)

		fmt.Println("Running backup...")
		fmt.Printf("Command: %s\n\n", backup.Command)
		out, err := utils.ExecuteCommand(backup.Command)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(out)

		fmt.Println("Publishing results to API.")
		actions.FinishJob(client, tokenResponse, job)

		return nil
	},
}
