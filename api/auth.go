package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/backupshq/agent/utils"
)

func (c *ApiClient) GetAccessToken() (string, error) {
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {c.credentials.clientId},
		"client_secret": {c.credentials.clientSecret},
		"scope":         {"agent"},
	}

	req, err := http.NewRequest("POST", c.server+"/auth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", errors.New("Error reading request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", errors.New("Error reading response: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", errors.New("Unable to retrieve access token: HTTP " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("Error reading JSON response body: " + string(body))
	}

	response := struct {
		Token string `json:"access_token"`
	}{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Token, nil
}

func (c *ApiClient) Authenticate() {
	token, err := c.GetAccessToken()
	if err != nil {
		log.Fatal(err)
	}

	c.accessToken = token
	c.accessTokenExpiry, _ = utils.GetAccessTokenExpiry(token)
}

func (c *ApiClient) AddAuthHeader(req *http.Request) {
	req.Header.Set("Authorization", "Bearer "+c.accessToken)
	req.Header.Set("Accept", fmt.Sprintf("application/vnd.backupshq.v%d+json", c.version))
}
