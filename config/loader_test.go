package config

import (
	"strings"
	"testing"

	"github.com/backupshq/agent/log"
)

func TestString(t *testing.T) {
	t.Run("load empty configuration", func(t *testing.T) {
		loader := NewConfigLoader(map[string]string{})
		config, err := loader.LoadString("")

		if err != nil {
			t.Errorf("expected a Config struct to be created without error, got %q", err.Error())
			return
		}
		expected := "https://api.backupshq.com"
		if config.ApiServer != expected {
			t.Errorf("got %q want %q", config.ApiServer, expected)
		}
	})

	t.Run("load simple config", func(t *testing.T) {
		loader := NewConfigLoader(map[string]string{})
		config, err := loader.LoadString(`
[auth]
client_id = 'id'
client_secret = 'secret'
`)
		if err != nil {
			t.Errorf("expected a Config struct to be created without error, got %q", err.Error())
			return
		}
		if config.Auth.ClientId != "id" {
			t.Errorf("got %q want %q", config.Auth.ClientSecret, "id")
		}
		if config.Auth.ClientSecret != "secret" {
			t.Errorf("got %q want %q", config.Auth.ClientSecret, "secret")
		}
	})

	t.Run("load simple config with odd casing", func(t *testing.T) {
		loader := NewConfigLoader(map[string]string{})
		config, err := loader.LoadString(`
[Auth]
CLIENT_ID = 'id'
Client_Secret = 'secret'
`)
		if err != nil {
			t.Errorf("expected a Config struct to be created without error, got %q", err.Error())
			return
		}
		if config.Auth.ClientId != "id" {
			t.Errorf("got %q want %q", config.Auth.ClientSecret, "id")
		}
		if config.Auth.ClientSecret != "secret" {
			t.Errorf("got %q want %q", config.Auth.ClientSecret, "secret")
		}
	})

	t.Run("load log level", func(t *testing.T) {
		loader := NewConfigLoader(map[string]string{})
		config, err := loader.LoadString(`
log_level = 'info'
`)
		if err != nil {
			t.Errorf("expected a Config struct to be created without error, got %q", err.Error())
			return
		}
		if config.LogLevel.Level != log.Info {
			t.Errorf("got %d want %d", config.LogLevel.Level, log.Info)
		}
	})

	t.Run("load invalid log level", func(t *testing.T) {
		loader := NewConfigLoader(map[string]string{})
		_, err := loader.LoadString(`
log_level = 'critical'
`)
		if err == nil {
			t.Errorf("expected an error loading config")
			return
		}
		if err.Error() != "TOML syntax error: Unknown log level 'critical'. Valid levels are 'debug', 'info', 'warn', and 'error'." {
			t.Errorf("unexpected error message: %q", err.Error())
		}
	})

	t.Run("invalid toml returns error", func(t *testing.T) {
		loader := NewConfigLoader(map[string]string{})
		config, err := loader.LoadString(`
[]Not valid
this doesn't work
= "
`)
		if config != nil || err == nil {
			t.Errorf("error should be returned from bad toml")
			return
		}
		if !strings.Contains(err.Error(), "TOML syntax error: ") {
			t.Errorf("error message should mention a template error, got %q", err.Error())
		}
	})

	t.Run("replace env key with val", func(t *testing.T) {
		loader := NewConfigLoader(map[string]string{"TESTVAR": "test"})
		config, err := loader.LoadString(`
[auth]
client_id = '{{ env "TESTVAR"}}'
client_secret = 'secret'
`)
		if err != nil {
			t.Errorf("expected a Config struct to be created without error, got %q", err.Error())
			return
		}
		if config.Auth.ClientId != "test" {
			t.Errorf("got %q want %q", config.Auth.ClientId, "test")
		}
	})

	t.Run("handle invalid env var", func(t *testing.T) {
		loader := NewConfigLoader(map[string]string{})
		config, err := loader.LoadString(`
[auth]
client_id = '{{ env "TESTVAR"}}'
client_secret = 'secret'
`)
		if config != nil || err == nil {
			t.Errorf("unknown env variable should return an error")
			return
		}
		if !strings.Contains(err.Error(), "Missing environment variable: TESTVAR") {
			t.Errorf("error message should reference the missing environment variable, got %q", err.Error())
		}
	})

	t.Run("handle invalid template syntax", func(t *testing.T) {
		loader := NewConfigLoader(map[string]string{})
		config, err := loader.LoadString(`
[auth]
client_id = '{{'
client_secret = 'secret'
`)
		if config != nil || err == nil {
			t.Errorf("bad template syntax should return an error")
			return
		}
		if !strings.Contains(err.Error(), "Template syntax error: ") {
			t.Errorf("error message should mention a template error, got %q", err)
		}
	})

	t.Run("nested double quotes work", func(t *testing.T) {
		loader := NewConfigLoader(map[string]string{"CLIENT_ID": "100"})
		config, err := loader.LoadString(`
[auth]
client_id = "{{ env "CLIENT_ID" }}"
client_secret = "secret"
`)
		if err != nil {
			t.Errorf("expected a Config struct to be created without error, got %q", err.Error())
			return
		}
		if config.Auth.ClientId != "100" {
			t.Errorf("got %q want %q", config.Auth.ClientId, "100")
		}
	})

	t.Run("unknown TOML key returns error", func(t *testing.T) {
		loader := NewConfigLoader(map[string]string{})
		config, err := loader.LoadString(`
[auth]
not_a_real_key = "id"
client_secret = "secret"
`)
		if config != nil || err == nil {
			t.Errorf("unknown TOML key should return an error")
			return
		}
		if !strings.Contains(err.Error(), "Unrecognized TOML key(s) given: ") {
			t.Errorf("error message should mention unrecognized key, got %q", err)
		}
	})

	t.Run("custom api server", func(t *testing.T) {
		loader := NewConfigLoader(map[string]string{})
		config, err := loader.LoadString(`
api_server = "https://api.example.com"
`)

		if err != nil {
			t.Errorf("expected a Config struct to be created without error, got %q", err.Error())
			return
		}
		expected := "https://api.example.com"
		if config.ApiServer != expected {
			t.Errorf("got %q want %q", config.ApiServer, expected)
		}
	})
}

func TestFile(t *testing.T) {
	t.Run("load simple file", func(t *testing.T) {
		loader := NewConfigLoader(map[string]string{})
		config, err := loader.LoadFile("example.toml")

		if err != nil {
			t.Errorf("file should be loaded successfully, got error %q", err)
			return
		}

		if config.Auth.ClientId != "id" {
			t.Errorf("got %q want %q", config.Auth.ClientSecret, "id")
		}
		if config.Auth.ClientSecret != "secret" {
			t.Errorf("got %q want %q", config.Auth.ClientSecret, "secret")
		}
	})
}
