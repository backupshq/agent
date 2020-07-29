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

type JobStep struct {
	ID        string
	Name      string
	SortOrder int
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

func (c *ApiClient) CreateStep(jobId string, name string, sortOrder int) JobStep {
	var requestBody = []byte(fmt.Sprintf(`{"job_id":"%s", "name": "%s", "sort_order": %d}`, jobId, name, sortOrder))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/job-steps", c.server), bytes.NewBuffer(requestBody))
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
		log.Fatal("Unable to create job step: HTTP " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	var step JobStep
	err = json.Unmarshal(body, &step)
	return step
}

func (c *ApiClient) SendLogs(step JobStep, logStr string) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/job-steps/%s/logs", c.server, step.ID), strings.NewReader(logStr))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	req.Header.Set("Content-Type", "text/plain")

	c.AddAuthHeader(req)

	_, err = c.client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
}

func (c *ApiClient) FinishStep(stepId string, status string) {
	var requestBody = []byte(fmt.Sprintf(`{"status": "%s"}`, status))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/job-steps/%s/finish", c.server, stepId), bytes.NewBuffer(requestBody))
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
		log.Fatal("Unable to finish job step: HTTP " + resp.Status)
	}
}
