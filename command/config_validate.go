package command

import (
	"log"

	"github.com/urfave/cli"
)

var ConfigValidate = cli.Command{
	Name:  "validate",
	Usage: "Validate a BackupsHQ TOML configuration file",
	Action: func(c *cli.Context) error {
		log.Println("works")

		return nil
	},
}
