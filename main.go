package main

import (
	"./command"
	"github.com/urfave/cli"
	"os"
)

func createApp() *cli.App {
	app := cli.NewApp()
	app.Name = "backupshq"
	app.Usage = "The backupshq.com command line agent"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Load configuration from `FILE`",
		},
	}
	app.Commands = []cli.Command{
		command.Agent,
		command.ExampleConfig,
		{
			Name:  "job",
			Usage: "Interact with backup jobs",
			Subcommands: []cli.Command{
				command.JobRun,
				command.JobStartUnmanaged,
				command.JobFinishUnmanaged,
			},
		},
	}

	return app
}

func main() {
	app := createApp()

	app.Run(os.Args)
}
