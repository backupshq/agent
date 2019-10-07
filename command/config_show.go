package command

import "fmt"
import "github.com/urfave/cli"
import "../utils"
import "../config"
import "io/ioutil"
import "log"

var ConfigShow = cli.Command{
	Name:  "show",
	Usage: "Print the configuration file with environment variables applied",
	Action: func(c *cli.Context) error {
		filePath := config.CliFilePath(c)

		tomlText, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		loader := config.NewConfigLoader(utils.EnvMap())

		parsedToml, err := loader.ApplyEnvironmentVariables(string(tomlText))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(parsedToml)

		return nil
	},
}
