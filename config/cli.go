package config

import (
	"log"

	"github.com/backupshq/agent/utils"

	"github.com/urfave/cli"
)

func LoadCli(c *cli.Context) *Config {
	filePath := CliFilePath(c)
	loader := NewConfigLoader(utils.EnvMap())

	config, err := loader.LoadFile(filePath)

	if err != nil {
		log.Fatal("Error loading configuration file: " + err.Error())
	}

	return config
}

func CliFilePath(c *cli.Context) string {
	filePath := c.GlobalString("config")

	if filePath == "" {
		log.Fatal("Error: configuration file required. Use the --config flag: `backupshq --config config.toml`.")
	}

	return filePath
}
