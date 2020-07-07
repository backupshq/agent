package command

import (
	"github.com/backupshq/agent/agent"
	"github.com/backupshq/agent/config"
	"github.com/urfave/cli"
)

var Agent = cli.Command{
	Name:  "agent",
	Usage: "Run the BackupsHQ agent",
	Action: func(c *cli.Context) error {
		config := config.LoadCli(c)
		agent := agent.Create(config)
		agent.Start()

		return nil
	},
}
