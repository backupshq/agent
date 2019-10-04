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
		filePath := c.GlobalString("config")

		if filePath == "" {
			log.Fatal("Error: configuration file required. Use the --config flag: `backupshq --config config.toml`.")
		}
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
