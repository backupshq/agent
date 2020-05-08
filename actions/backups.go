package actions

import (
	"fmt"
	"log"

	"github.com/backupshq/agent/api"
	"github.com/backupshq/agent/utils"
)

func RunBackup(client *api.ApiClient, backup api.Backup) {
	job := client.StartJob(backup.ID)
	log.Printf("Started a new job with id %q.\n", job.ID)

	log.Println("Running backup...")
	log.Printf("Command: %s\n\n", backup.Command)
	out, err := utils.ExecuteCommand(backup.Command)
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
