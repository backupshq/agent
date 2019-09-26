package utils

import (
	"fmt"
	"log"

	"../actions"
	"github.com/robfig/cron"
)

func Schedule(client *actions.ApiClient, backup actions.Backup) *cron.Cron {
	log.Println("Schedule " + backup.Name + " for " + backup.Schedule)
	c := cron.New()
	c.AddFunc(backup.Schedule, func() { RunBackup(client, backup) })
	c.Start()
	return c
}

func RunBackup(client *actions.ApiClient, backup actions.Backup) {
	job := client.StartJob(backup.ID)
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
	client.FinishJob(job)
	client.SendLogs(job, out)
}
