package config

import (
	"testing"
)

func TestString(t *testing.T) {
	loader := NewConfigLoader(map[string]string{"TESTVAR": "test"})
	t.Run("load empty configuration", func(t *testing.T) {
		config, _ := loader.LoadString("")

		if config.Auth.ClientId != "" {
			t.Errorf("got %q want %q", config.Auth.ClientSecret, "")
		}
	})

	t.Run("load simple config", func(t *testing.T) {
		config, _ := loader.LoadString(`
[auth]
client_id = "id"
client_secret = "secret"
`)
		if config.Auth.ClientId != "id" {
			t.Errorf("got %q want %q", config.Auth.ClientSecret, "id")
		}
		if config.Auth.ClientSecret != "secret" {
			t.Errorf("got %q want %q", config.Auth.ClientSecret, "secret")
		}
	})

	t.Run("load simple config with odd casing", func(t *testing.T) {
		config, _ := loader.LoadString(`
[Auth]
CLIENT_ID = "id"
Client_Secret = "secret"
`)
		if config.Auth.ClientId != "id" {
			t.Errorf("got %q want %q", config.Auth.ClientSecret, "id")
		}
		if config.Auth.ClientSecret != "secret" {
			t.Errorf("got %q want %q", config.Auth.ClientSecret, "secret")
		}
	})

	t.Run("invalid toml returns error", func(t *testing.T) {
		config, err := loader.LoadString(`
[]Not valid
this doesn't work
= "
`)
		if config != nil || err == nil {
			t.Errorf("error should be returned from bad toml")
		}
	})

	t.Run("replace env key with val", func(t *testing.T) {
		config, _ := loader.LoadString(`
[auth]
client_id = "{{ env "TESTVAR"}}"
client_secret = "secret"
`)

		if config.Auth.ClientId != "test" {
			t.Errorf("got %q want %q", config.Auth.ClientId, "test")
		}
	})
}

func TestFile(t *testing.T) {
	loader := NewConfigLoader(map[string]string{"TESTVAR": "test"})
	t.Run("load simple file", func(t *testing.T) {
		config, err := loader.LoadFile("example.toml")

		if err != nil {
			t.Errorf("file should be loaded successfully")
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
