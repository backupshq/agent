package utils

import (
	"fmt"
	"log"
	"net/http"

	"../actions"
	"../auth"
	"github.com/robfig/cron"
)

func Schedule(client *http.Client, tokenResponse auth.AccessTokenResponse, backup actions.Backup) {
	log.Println("Schedule " + backup.Name + " for " + backup.Cron)
	c := cron.New()
	c.AddFunc(backup.Cron, func() { RunBackup(client, tokenResponse, backup) })
	c.Start()
}

func RunBackup(client *http.Client, tokenResponse auth.AccessTokenResponse, backup actions.Backup) {
	job := actions.StartJob(client, tokenResponse, backup.ID)
	fmt.Printf("Started a new job with id %q.\n", job.ID)

	fmt.Println("Running backup...")
	fmt.Printf("Command: %s\n\n", backup.Command)
	out, err := ExecuteCommand(backup.Command)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)

	fmt.Println("Publishing results to API.")
	actions.FinishJob(client, tokenResponse, job)
	actions.SendLogs(client, tokenResponse, job, out)
}
