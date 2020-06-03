package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/backupshq/agent/config"
)

type credentials struct {
	clientId     string
	clientSecret string
}

type ApiClient struct {
	client            *http.Client
	credentials       credentials
	server            string
	version           int
	accessToken       string
	accessTokenExpiry time.Time
	PrincipalId       string
}

func NewClient(config *config.Config) *ApiClient {
	client := &http.Client{
		Timeout: time.Second * 3,
	}

	return &ApiClient{
		client: client,
		credentials: credentials{
			clientId:     config.Auth.ClientId,
			clientSecret: config.Auth.ClientSecret,
		},
		server:  config.ApiServer,
		version: 1,
	}
}

func (c *ApiClient) get(path string) (*http.Request, error) {
	path = "/" + strings.TrimLeft(path, "/")
	req, err := http.NewRequest("GET", c.server+path, nil)
	if err != nil {
		return nil, err
	}
	c.AddAuthHeader(req)

	return req, nil
}
