package command

import "fmt"
import "github.com/urfave/cli"
import "time"
import "io/ioutil"
import "log"
import "net/http"
import "encoding/json"
import "../config"
import "../auth"

var BackupRun = cli.Command{
	Name:  "run",
	Usage: "Run a one-off backup and send the results to the API",
	Action: func(c *cli.Context) error {
		config := config.LoadCli(c)

		tokenResponse, err := auth.GetAccessToken(config)
		if err != nil {
			log.Fatal(err)
		}

		client := &http.Client{
			Timeout: time.Second * 3,
		}
		job := startJob(client, tokenResponse, c.Args().Get(0))
		fmt.Printf("Started a new job with id %q.\n", job.ID)

		fmt.Println("Running backup...")
		time.Sleep(3 * time.Second)

		fmt.Println("Publishing results to API.")
		finishJob(client, tokenResponse, job)

		return nil
	},
}

type job struct {
	ID string
}

func startJob(client *http.Client, tokenResponse auth.AccessTokenResponse, backupId string) job {
	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:8000/backups/%s/start", backupId), nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	auth.AddAuthHeader(req, tokenResponse)

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

	var startedJob job
	err = json.Unmarshal(body, &startedJob)

	return startedJob
}

func finishJob(client *http.Client, tokenResponse auth.AccessTokenResponse, job job) {
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
