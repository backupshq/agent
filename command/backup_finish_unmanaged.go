package command

import "github.com/urfave/cli"
import "time"
import "os"
import "bufio"
import "log"
import "net/http"
import "../config"
import "../auth"
import "../actions"

var BackupFinishUnmanaged = cli.Command{
	Name:  "finish-unmanaged",
	Usage: "Finish an unmanaged backup",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:   "log-stdin",
			Usage:  "Log the stdin channel",
		},
	},
	Action: func(c *cli.Context) error {
		scanner := bufio.NewScanner(os.Stdin)
		stdin := ""

		if c.Bool("log-stdin") {
			for scanner.Scan() {
				text := scanner.Text()
				log.Println(text)
				stdin += text
			}
		}

		config := config.LoadCli(c)

		tokenResponse, err := auth.GetAccessToken(config)
		if err != nil {
			log.Fatal(err)
		}

		client := &http.Client{
			Timeout: time.Second * 3,
		}

		job := actions.Job{ID: c.Args().Get(0)}
		actions.FinishJob(client, tokenResponse, job)
		log.Printf("Finished Job: %q.\n", job.ID)

		return nil
	},
}
