package api

import (
	"log"
	"net/http"
	"time"

	"github.com/backupshq/agent/auth"
	"github.com/backupshq/agent/config"
	"strings"
)

type credentials struct {
	clientId     string
	clientSecret string
}

type ApiClient struct {
	client        *http.Client
	credentials   credentials
	tokenResponse auth.AccessTokenResponse
	server        string
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
		server: config.ApiServer,
	}
}

func (c *ApiClient) Authenticate() {
	tokenResponse, err := auth.GetAccessToken(c.credentials.clientId, c.credentials.clientSecret, c.server+"/auth/token")
	if err != nil {
		log.Fatal(err)
	}

	c.tokenResponse = tokenResponse
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
