package config

import (
	"bytes"
	"errors"
	"html/template"
	"io/ioutil"
	"strings"

	"../utils"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

type Config struct {
	Auth struct {
		ClientId     string `toml:"client_id"`
		ClientSecret string `toml:"client_secret"`
		Next         struct {
			ClientId     string `toml:"red"`
			ClientSecret string `toml:"blue"`
		}
	}
	High struct {
		Test string `toml:"test"`
	}
}

type ConfigLoader struct {
	env map[string]string
}

func NewConfigLoader(envrionment map[string]string) ConfigLoader {
	return ConfigLoader{env: envrionment}
}

func (l *ConfigLoader) LoadString(tomlText string) (*Config, error) {

	funcMap := template.FuncMap{
		"env": func(key string) (string, error) {
			val := l.env[key]
			if val == "" {
				return "", errors.New("Cannot find envrionment variable: " + key)
			}
			return l.env[key], nil
		},
	}

	var templateReader bytes.Buffer
	tpl, err := template.New("config").Funcs(funcMap).Parse(tomlText)
	if err != nil {
		return nil, err
	}
	err = tpl.Execute(&templateReader, map[string]string{})
	if err != nil {
		return nil, err
	}
	tomlText = templateReader.String()

	var config Config
	metadata, err := toml.Decode(tomlText, &config)
	if err != nil {
		return nil, err
	}

	invalidKeys := metadata.Undecoded()
	if len(invalidKeys) != 0 {
		var keysAsString []string
		for _, key := range invalidKeys {
			keysAsString = append(keysAsString, key.String())
		}
		return nil, errors.New("Unrecognized TOML key(s) given: " + strings.Join(keysAsString, ", "))
	}

	return &config, nil
}

func (l *ConfigLoader) LoadFile(filePath string) (*Config, error) {
	tomlText, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return l.LoadString(string(tomlText))
}

func LoadCli(c *cli.Context) (*Config, error) {
	loader := NewConfigLoader(utils.GetEvnVariables())
	filePath := c.GlobalString("config")

	if filePath == "" {
		return nil, errors.New("Error: configuration file required. Use the --config flag: `backupshq --config config.toml`.")
	}
	config, err := loader.LoadFile(filePath)

	if err != nil {
		return nil, errors.New("Error loading configuration file: " + err.Error())
	}

	return config, nil
}
