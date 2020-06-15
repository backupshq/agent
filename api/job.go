package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Job struct {
	ID string
}

func (c *ApiClient) StartJob(backupId string) Job {
	var requestBody = []byte(fmt.Sprintf(`{"backup_id":"%s", "start": true}`, backupId))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/jobs", c.server), bytes.NewBuffer(requestBody))
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
		log.Fatal("Unable to start job: HTTP " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	var startedJob Job
	err = json.Unmarshal(body, &startedJob)
	return startedJob
}

func (c *ApiClient) FinishJob(job Job, status string) {
	var requestBody = []byte(fmt.Sprintf(`{"status":"%s"}`, status))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/jobs/%s/finish", c.server, job.ID), bytes.NewBuffer(requestBody))
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
		log.Fatal("Failed to finish job: " + resp.Status)
	}

	return
}

func (c *ApiClient) SendLogs(job Job, logStr string) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/jobs/%s/logs", c.server, job.ID), strings.NewReader(logStr))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	c.AddAuthHeader(req)

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	return
}
