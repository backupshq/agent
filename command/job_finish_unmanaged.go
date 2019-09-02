package command

import "github.com/urfave/cli"

var JobFinishUnmanaged = cli.Command{
	Name:  "finish-unmanaged",
	Usage: "Finish an unmanaged backup",
	Action: func(c *cli.Context) error {
		return nil
	},
}
