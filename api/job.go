package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type Job struct {
	ID         string
	BackupID   string `json:"backup_id"`
	BackupName string `json:"backup_name"`
	JobNumber  int    `json:"job_number"`
	Status     string
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

func (c *ApiClient) StartExistingJob(job Job) Job {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/jobs/%s/start", c.server, job.ID), nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	c.AddAuthHeader(req)

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	if resp.StatusCode != 200 {
		log.Fatal(fmt.Sprintf("Unable to start job: HTTP "+resp.Status+" %v", string(body)))
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

// TODO: Replace all PATCH job endpoints with this
func (c *ApiClient) UpdateJob(job Job, finishedAt time.Time) {
	var requestBody = []byte(fmt.Sprintf(`{"finished_at":"%s"}`, finishedAt.UTC().Format(time.RFC3339)))
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/jobs/%s", c.server, job.ID), bytes.NewBuffer(requestBody))
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
		log.Fatal("Failed to update job: " + resp.Status)
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
