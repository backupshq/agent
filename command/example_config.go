package command

import (
	"fmt"

	"github.com/urfave/cli"
)

var ExampleConfig = cli.Command{
	Name:  "example",
	Usage: "Print an example configuration file",
	Action: func(c *cli.Context) error {
		fmt.Print(`[auth]
# Client credentials for the agent to access the API.
client_id = ""
client_secret = ""
# Use an environment variable instead:
# client_secret = "{{env SOME_SECRET}}"
`)

		return nil
	},
}
