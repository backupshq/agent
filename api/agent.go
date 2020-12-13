package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type AgentToken struct {
	Token string
}

func (c *ApiClient) Register() AgentToken {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("Unable to get hostname. ", err)
	}
	var requestBody = []byte(fmt.Sprintf(`{"hostname":"%s"}`, hostname))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/agents/register", c.server), bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	c.AddAuthHeader(req)

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("Unable to register agent: HTTP " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	var token AgentToken
	err = json.Unmarshal(body, &token)
	return token
}

type PingResponse struct {
	UpdatedBackupCount int   `json:"updated_backup_count"`
	AssignedJobs       []Job `json:"assigned_jobs"`
}

func (c *ApiClient) Ping(token AgentToken) PingResponse {
	req, err := c.get(fmt.Sprintf("/agents/ping/%s", token.Token))
	if err != nil {
		log.Fatal("Error creating request. ", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 && resp.StatusCode != 304 {
		log.Fatal("Unable to ping API: HTTP " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	var response PingResponse
	err = json.Unmarshal(body, &response)

	return response
}
