package utils

import (
	"fmt"
	"log"
	"net/http"

	"../actions"
	"../auth"
	"github.com/robfig/cron"
)

func Schedule(client *http.Client, tokenResponse auth.AccessTokenResponse, backup actions.Backup) *cron.Cron {
	log.Println("Schedule " + backup.Name + " for " + backup.Cron)
	c := cron.New()
	c.AddFunc(backup.Cron, func() { RunBackup(client, tokenResponse, backup) })
	c.Start()
	return c
}

func RunBackup(client *http.Client, tokenResponse auth.AccessTokenResponse, backup actions.Backup) {
	job := actions.StartJob(client, tokenResponse, backup.ID)
	log.Printf("Started a new job with id %q.\n", job.ID)

	log.Println("Running backup...")
	log.Printf("Command: %s\n\n", backup.Command)
	out, err := ExecuteCommand(backup.Command)
	if err != nil {
		// In the future we can use this block to determine status code, but for now just send the error the the API
		log.Println(err)
		out = err.Error()
	}
	fmt.Println(out)

	log.Println("Publishing results to API.")
	actions.FinishJob(client, tokenResponse, job)
	actions.SendLogs(client, tokenResponse, job, out)
}
