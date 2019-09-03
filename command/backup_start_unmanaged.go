package command

import "fmt"
import "github.com/urfave/cli"
import "time"
import "log"
import "net/http"
import "../config"
import "../auth"
import "../actions"

var BackupStartUnmanaged = cli.Command{
	Name:  "start-unmanaged",
	Usage: "Start an unmanaged backup",
	Description: `
Start an unmanaged backup.
This command is only suitable for *unmanaged* backups that you handle yourself, for example:

backupshq start-unmanaged <backup_id> > job-id.txt && ./backup-script.sh | backupshq finish-unmanaged $(cat job-id.txt) --log-stdin

To run any other type of backup, see backupshq job run --help.
`,
	Action: func(c *cli.Context) error {
		config := config.LoadCli(c)

		tokenResponse, err := auth.GetAccessToken(config)
		if err != nil {
			log.Fatal(err)
		}

		client := &http.Client{
			Timeout: time.Second * 3,
		}
		job := actions.StartJob(client, tokenResponse, c.Args().Get(0))
		fmt.Printf("Started a new job with id %q.\n", job.ID)
		fmt.Printf("To inform the API when this job is finished run: backupshq finish-unmanaged %q\n", job.ID)

		return nil
	},
}
