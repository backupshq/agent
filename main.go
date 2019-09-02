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
	app.Commands = []cli.Command{
		command.Agent,
	}

	return app
}

func main() {
	app := createApp()

	app.Run(os.Args)
}
