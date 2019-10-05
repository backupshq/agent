package command

import (
	"fmt"
	"os"

	"../config"
	"github.com/urfave/cli"
)

var ConfigValidate = cli.Command{
	Name:  "validate",
	Usage: "Validate a BackupsHQ TOML configuration file",
	Action: func(c *cli.Context) {
		_, err := config.LoadCli(c)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("Config is valid!")
		os.Exit(0)
	},
}
