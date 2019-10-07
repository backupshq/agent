package command

import (
	"fmt"
	"os"

	"../config"
	"../utils"
	"github.com/urfave/cli"
)

var ConfigValidate = cli.Command{
	Name:  "validate",
	Usage: "Validate a BackupsHQ TOML configuration file",
	Action: func(c *cli.Context) {
		filePath := config.CliFilePath(c)
		loader := config.NewConfigLoader(utils.EnvMap())

		_, err := loader.LoadFile(filePath)
		if err != nil {
			fmt.Println("Error found in config file:")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		fmt.Println("Config is valid!")
		os.Exit(0)
	},
}
