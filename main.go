package main

import (
	"os"

	"./command"
	"github.com/urfave/cli"
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
			Name:  "backup",
			Usage: "Interact with backups",
			Subcommands: []cli.Command{
				command.BackupRun,
				command.BackupStartUnmanaged,
				command.BackupFinishUnmanaged,
			},
		},
		{
			Name:  "config",
			Usage: "Interact with TOML configurations",
			Subcommands: []cli.Command{
				command.ConfigValidate,
			},
		},
	}

	return app
}

func main() {
	app := createApp()

	app.Run(os.Args)
}
