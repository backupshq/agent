package command

import "log"
import "net/http"
import "encoding/json"
import "net/url"
import "strings"
import "io/ioutil"
import "github.com/urfave/cli"
import "time"

var clientId string
var clientSecret string

var Agent = cli.Command{
	Name:  "agent",
	Usage: "Run the BackupsHQ agent",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "client_id",
			Destination: &clientId,
		},
		cli.StringFlag{
			Name:        "client_secret",
			Destination: &clientSecret,
		},
	},
	Action: func(c *cli.Context) error {
		log.Println("Starting BackupsHQ agent")

		tokenResponse, err := getAccessToken()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Got access token")

		client := &http.Client{
			Timeout: time.Second * 3,
		}
		req, err := http.NewRequest("GET", "http://localhost:8000/backups", nil)
		if err != nil {
			log.Fatal("Error reading request. ", err)
		}
		req.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)

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

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_id"`
	Scope       string `json:"scope"`
}

func getAccessToken() (tokenResponse, error) {
	client := &http.Client{
		Timeout: time.Second * 3,
	}

	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
	}

	req, err := http.NewRequest("POST", "http://localhost:8000/auth/token", strings.NewReader(form.Encode()))
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	tokenResponse := tokenResponse{}
	err = json.Unmarshal(body, &tokenResponse)
	if err != nil {
		log.Fatal(err)
	}

	return tokenResponse, nil
}
