package config

import (
	"log"

	"../utils"

	"github.com/urfave/cli"
)

func LoadCli(c *cli.Context) *Config {
	loader := NewConfigLoader(utils.EnvMap())
	filePath := c.GlobalString("config")

	if filePath == "" {
		log.Fatal("Error: configuration file required. Use the --config flag: `backupshq --config config.toml`.")
	}
	config, err := loader.LoadFile(filePath)

	if err != nil {
		log.Fatal("Error loading configuration file: " + err.Error())
	}

	return config
}
