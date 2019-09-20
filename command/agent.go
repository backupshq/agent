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
			backups = actions.ListBackups(client, tokenResponse, actions.BACKUP_TYPE_SCHEDULED)

			for id, _ := range crons {
				crons[id].Stop()
			}

			for _, backup := range backups {
				crons[backup.ID] = utils.Schedule(client, tokenResponse, backup)
			}
		})
		cr.Start()

		time.Sleep(time.Minute * 100)

		return nil
	},
}
