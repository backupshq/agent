package command

import "github.com/urfave/cli"

var JobStartUnmanaged = cli.Command{
	Name:  "start-unmanaged",
	Usage: "Start an unmanaged backup",
	Description: `
Start an unmanaged backup.
This command is only suitable for *unmanaged* backups that you handle yourself, for example:

backupshq start-unmanaged <backup_id> > job-id.txt && ./backup-script.sh | backupshq finish-unmanaged $(cat job-id.txt) --log-stdin

To run any other type of backup, see backupshq job run --help.
`,
	Action: func(c *cli.Context) error {
		return nil
	},
}
