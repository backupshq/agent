package command

import (
	"bufio"
	"log"
	"os"

	"../api"
	"../config"
	"../utils"
	"github.com/urfave/cli"
)

var BackupFinishUnmanaged = cli.Command{
	Name:  "finish-unmanaged",
	Usage: "Finish an unmanaged backup",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "log-stdin",
			Usage: "Log the stdin channel",
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

		env := utils.GetEvnVariables()
		loader := config.NewConfigLoader(env)
		config := loader.LoadCli(c)

		client := api.NewClient(config)

		job := api.Job{ID: c.Args().Get(0)}
		client.FinishJob(job)
		log.Printf("Finished Job: %q.\n", job.ID)

		return nil
	},
}
