package config

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/backupshq/agent/log"
)

type Config struct {
	Auth struct {
		ClientId     string `toml:"client_id"`
		ClientSecret string `toml:"client_secret"`
	}
	LogLevel LogLevel `toml:"log_level"`
}

type LogLevel struct {
	Level int
	Label string
}

var logLevels = map[string]int{
	"debug": log.Debug,
	"info":  log.Info,
	"warn":  log.Warn,
	"error": log.Error,
}

func (l *LogLevel) UnmarshalText(text []byte) error {
	value := string(text)
	if level, ok := logLevels[value]; ok {
		l.Level = level
		l.Label = value
		return nil
	}
	return errors.New(fmt.Sprintf("Unknown log level '%s'. Valid levels are 'debug', 'info', 'warn', and 'error'.", value))
}

type ConfigLoader struct {
	env map[string]string
}

func NewConfigLoader(env map[string]string) ConfigLoader {
	return ConfigLoader{env: env}
}

func (l *ConfigLoader) ApplyEnvironmentVariables(tomlText string) (string, error) {
	funcMap := template.FuncMap{
		"env": func(key string) (string, error) {
			val := l.env[key]
			if val == "" {
				return "", errors.New("Missing environment variable: " + key)
			}
			return l.env[key], nil
		},
	}

	var buffer bytes.Buffer
	tpl, err := template.New("config").Funcs(funcMap).Parse(tomlText)
	if err != nil {
		return "", errors.New("Template syntax error: " + err.Error())
	}
	err = tpl.Execute(&buffer, map[string]string{})
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (l *ConfigLoader) LoadString(tomlText string) (*Config, error) {
	tomlText, err := l.ApplyEnvironmentVariables(tomlText)
	if err != nil {
		return nil, err
	}

	config := Config{
		LogLevel: LogLevel{
			Level: log.Info,
			Label: "info",
		},
	}
	metadata, err := toml.Decode(tomlText, &config)
	if err != nil {
		return nil, errors.New("TOML syntax error: " + err.Error())
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
