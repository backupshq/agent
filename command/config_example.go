package command

import "fmt"
import "github.com/urfave/cli"

var ConfigExample = cli.Command{
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
