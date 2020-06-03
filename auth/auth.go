package auth

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func GetAccessToken(clientId string, clientSecret string, endpoint string) (AccessTokenResponse, error) {
	client := &http.Client{
		Timeout: time.Second * 3,
	}
	tokenResponse := AccessTokenResponse{}

	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"scope":         {"agent"},
	}

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return tokenResponse, errors.New("Error reading request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return tokenResponse, errors.New("Error reading response: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return tokenResponse, errors.New("Unable to retrieve access token: HTTP " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tokenResponse, errors.New("Error reading JSON response body: " + string(body))
	}

	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		return tokenResponse, err
	}

	return tokenResponse, nil
}

func AddAuthHeader(req *http.Request, token AccessTokenResponse) {
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("Accept", "application/vnd.backupshq.v1+json")
}
