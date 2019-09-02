package config

import "github.com/BurntSushi/toml"
import "io/ioutil"

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
