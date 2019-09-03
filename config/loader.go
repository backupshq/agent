package config

import "github.com/BurntSushi/toml"
import "github.com/urfave/cli"
import "io/ioutil"
import "log"

type Config struct {
	Auth struct {
		ClientId     string `toml:"client_id"`
		ClientSecret string `toml:"client_secret"`
	}
}

func LoadString(tomlText string) (*Config, error) {
	var config Config
	if _, err := toml.Decode(tomlText, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func LoadFile(filePath string) (*Config, error) {
	tomlText, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return LoadString(string(tomlText))
}

func LoadCli(c *cli.Context) *Config {
	filePath := c.GlobalString("config")

	if filePath == "" {
		log.Fatal("Error: configuration file required. Use the --config flag: `backupshq --config config.toml`.")
	}
	config, err := LoadFile(filePath)

	if err != nil {
		log.Fatal("Error load configuration file: " + err.Error())
	}

	return config
}
