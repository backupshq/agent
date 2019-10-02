package command

import (
	"log"
	"time"

	"../actions"
	"../api"
	"../config"
	"github.com/robfig/cron"
	"github.com/urfave/cli"
)

func getLatestSchema(client *api.ApiClient, backups map[string]api.Backup, crons map[string]*cron.Cron) (map[string]api.Backup, map[string]*cron.Cron) {
	log.Println("Checking for changes to backups...")
	newBackups := client.ListBackups(api.BACKUP_TYPE_SCHEDULED)
	log.Println("Scheduled backups pulled from the API:", len(newBackups))

	didUpdate := false
	for id := range newBackups {
		if backups[id] != newBackups[id] {
			log.Println("Updated backup: " + newBackups[id].Name)
			didUpdate = true
			if cron, ok := crons[id]; ok { // checks if there's an existing cron job for this backup
				cron.Stop()
			}
			crons[id] = actions.Schedule(client, newBackups[id])
		}
	}

	if !didUpdate {
		log.Println("No changes detected")
	}

	return newBackups, crons
}

var Agent = cli.Command{
	Name:  "agent",
	Usage: "Run the BackupsHQ agent",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "sync",
			Value: "* * * * *",
			Usage: "Cron expression representing how frequently the agent will sync with the API.",
		},
	},
	Action: func(c *cli.Context) error {
		log.Println("Starting BackupsHQ agent with sync frequency:", c.String("sync"))

		loader := config.NewConfigLoader()
		config := loader.LoadCli(c)

		client := api.NewClient(config)

		backups, crons := getLatestSchema(client, map[string]api.Backup{}, map[string]*cron.Cron{})
		cr := cron.New()
		cr.AddFunc(c.String("sync"), func() {
			backups, crons = getLatestSchema(client, backups, crons)
		})
		cr.Start()

		time.Sleep(time.Minute * 100)

		return nil
	},
}
