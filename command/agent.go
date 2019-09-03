package command

import "log"
import "net/http"
import "encoding/json"
import "io/ioutil"
import "github.com/urfave/cli"
import "time"
import "../config"
import "../auth"

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
		req, err := http.NewRequest("GET", "http://localhost:8000/backups", nil)
		if err != nil {
			log.Fatal("Error reading request. ", err)
		}
		auth.AddAuthHeader(req, tokenResponse)

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Error reading response. ", err)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Error reading body. ", err)
		}

		var backups []interface{}
		err = json.Unmarshal(body, &backups)
		if err != nil {
			log.Fatal(err)
		}

		for _, backup := range backups {
			log.Println(backup)
		}

		return nil
	},
}
