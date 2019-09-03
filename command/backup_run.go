package command

import "fmt"
import "github.com/urfave/cli"
import "time"
import "log"
import "net/http"
import "../config"
import "../auth"
import "../actions"

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
		job := actions.StartJob(client, tokenResponse, c.Args().Get(0), nil)
		fmt.Printf("Started a new job with id %q.\n", job.ID)

		fmt.Println("Running backup...")
		time.Sleep(3 * time.Second)

		fmt.Println("Publishing results to API.")
		actions.FinishJob(client, tokenResponse, job)

		return nil
	},
}
