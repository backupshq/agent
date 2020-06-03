package api

import (
	"log"
	"net/http"
	"time"

	"github.com/backupshq/agent/auth"
	"github.com/backupshq/agent/config"
	"strings"
)

type ApiClient struct {
	client        *http.Client
	tokenResponse auth.AccessTokenResponse
	server        string
}

func NewClient(config *config.Config) *ApiClient {
	tokenResponse, err := auth.GetAccessToken(config)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: time.Second * 3,
	}

	return &ApiClient{
		client:        client,
		tokenResponse: tokenResponse,
		server:        config.ApiServer,
	}
}

func (c *ApiClient) get(path string) (*http.Request, error) {
	path = "/" + strings.TrimLeft(path, "/")
	req, err := http.NewRequest("GET", c.server+path, nil)
	if err != nil {
		return nil, err
	}
	auth.AddAuthHeader(req, c.tokenResponse)

	return req, nil
}
