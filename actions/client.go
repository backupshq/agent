package actions

import (
	"log"
	"net/http"
	"time"

	"../auth"
	"../config"
)

type ApiClient struct {
	client        *http.Client
	tokenResponse auth.AccessTokenResponse
}

func NewClient(config *config.Config) *ApiClient {
	tokenResponse, err := auth.GetAccessToken(config)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: time.Second * 3,
	}

	return &ApiClient{client: client, tokenResponse: tokenResponse}
}
