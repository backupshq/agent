package actions


import "encoding/json"
import "fmt"
import "io/ioutil"
import "log"
import "net/http"
import "net/url"
import "strings"
import "../auth"


type Job struct {
	ID string
}

func StartJob(client *http.Client, tokenResponse auth.AccessTokenResponse, backupId string, context url.Values) Job {
	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8000/backups/%s/start", backupId), strings.NewReader(context.Encode()))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	auth.AddAuthHeader(req, tokenResponse)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	
	resp, err := client.Do(req)
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

func FinishJob(client *http.Client, tokenResponse auth.AccessTokenResponse, job Job) {
	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8000/jobs/%s/finish", job.ID), nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	auth.AddAuthHeader(req, tokenResponse)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	return
}
