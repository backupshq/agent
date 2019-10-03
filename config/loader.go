package config

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/urfave/cli"
)

type Config struct {
	Auth struct {
		ClientId     string `toml:"client_id"`
		ClientSecret string `toml:"client_secret"`
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
		"env": func(val string) string {
			return l.env[val]
		},
	}

	var templateReader bytes.Buffer
	tpl, err := template.New("config").Funcs(funcMap).Parse(tomlText)
	if err != nil {
		log.Fatal("Error templating config file: ", err)
	}
	err = tpl.Execute(&templateReader, map[string]string{})
	if err != nil {
		log.Fatal("Error templating config file: ", err)
	}
	tomlText = templateReader.String()

	var config Config
	if _, err := toml.Decode(tomlText, &config); err != nil {
		return nil, err
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

func (l *ConfigLoader) LoadCli(c *cli.Context) *Config {
	filePath := c.GlobalString("config")

	if filePath == "" {
		log.Fatal("Error: configuration file required. Use the --config flag: `backupshq --config config.toml`.")
	}
	config, err := l.LoadFile(filePath)

	if err != nil {
		log.Fatal("Error load configuration file: " + err.Error())
	}

	return config
}
