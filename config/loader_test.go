package config

import "testing"

func TestString(t *testing.T) {
	t.Run("load empty configuration", func(t *testing.T) {
		config, _ := LoadString("")

		if config.Auth.ClientId != "" {
			t.Errorf("got %q want %q", config.Auth.ClientSecret, "")
		}
	})

	t.Run("load simple config", func(t *testing.T) {
		config, _ := LoadString(`
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
		config, _ := LoadString(`
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
		config, err := LoadString(`
[]Not valid
this doesn't work
= "
`)
		if config != nil || err == nil {
			t.Errorf("error should be returned from bad toml")
		}
	})
}

func TestFile(t *testing.T) {
	t.Run("load simple file", func(t *testing.T) {
		config, err := LoadFile("example.toml")

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
