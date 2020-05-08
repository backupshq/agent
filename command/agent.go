package command

import (
	"github.com/backupshq/agent/agent"
	"github.com/backupshq/agent/config"
	"github.com/urfave/cli"
)

var Agent = cli.Command{
	Name:  "agent",
	Usage: "Run the BackupsHQ agent",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "sync",
			Value: "* * * * *",
			Usage: "Cron expression representing how frequently the agent will sync with the API.",
		},
	},
	Action: func(c *cli.Context) error {
		config := config.LoadCli(c)
		agent := agent.Create(config, c.String("sync"))
		agent.Start()

		return nil
	},
}
