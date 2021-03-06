package command

import (
	"fmt"
	"log"

	"github.com/backupshq/agent/api"
	"github.com/backupshq/agent/config"
	"github.com/urfave/cli"
)

var BackupStartUnmanaged = cli.Command{
	Name:  "start-unmanaged",
	Usage: "Start an unmanaged backup",
	Description: `
Start an unmanaged backup.
This command is only suitable for *unmanaged* backups that you handle yourself, for example:

backupshq backup start-unmanaged --id <backup_id> > job-id.txt && ./backup-script.sh | backupshq backup finish-unmanaged $(cat job-id.txt) --log-stdin

To run any other type of backup, see backupshq job run --help.
`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "id",
			Usage: "Output only the new job ID",
		},
	},
	Action: func(c *cli.Context) error {
		config := config.LoadCli(c)

		client := api.NewClient(config)

		backupID := c.Args().Get(0)
		backup := client.GetBackup(backupID)
		if backup.Managed {
			log.Fatal("Cannot start a managed backup using start-unmanaged command")
		}

		job := client.StartJob(backupID)
		if c.Bool("id") {
			fmt.Printf("%s\n", job.ID)
			return nil
		}
		log.Printf("Started a new job with id %q.\n", job.ID)
		log.Printf("To inform the API when this job is finished run: backupshq backup finish-unmanaged %q\n", job.ID)

		return nil
	},
}
