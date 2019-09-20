package command

import (
	"log"
	"net/http"
	"time"

	"../actions"
	"../auth"
	"../config"
	"../utils"
	"github.com/robfig/cron"
	"github.com/urfave/cli"
)

var Agent = cli.Command{
	Name:  "agent",
	Usage: "Run the BackupsHQ agent",
	Action: func(c *cli.Context) error {
		log.Println("Starting BackupsHQ agent")

		config := config.LoadCli(c)

		tokenResponse, err := auth.GetAccessToken(config)
		if err != nil {
			log.Fatal(err)
		}

		client := &http.Client{
			Timeout: time.Second * 3,
		}

		backups := map[string]actions.Backup{}
		crons := map[string]*cron.Cron{}

		cr := cron.New()
		cr.AddFunc("* * * * *", func() {
			log.Println("Checking for changes to backups...")
			newBackups := actions.ListBackups(client, tokenResponse, actions.BACKUP_TYPE_SCHEDULED)

			for id := range newBackups {
				if backups[id] != newBackups[id] {
					if cron, ok := crons[id]; ok { // checks if there's an existing cron job for this backup
						cron.Stop()
					}
					crons[id] = utils.Schedule(client, tokenResponse, newBackups[id])
				}
			}

			backups = newBackups
		})
		cr.Start()

		time.Sleep(time.Minute * 100)

		return nil
	},
}
