package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Secret struct {
	ID    string
	Name  string
	Value string `json:"secret"`
}

func (c *ApiClient) GetSecretByName(name string) Secret {
	req, err := c.get(fmt.Sprintf("/credentials/name/%s?with_secret=1", name))
	if err != nil {
		log.Fatal("Error creating request. ", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("Unable to get secret: HTTP " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	var secret Secret
	err = json.Unmarshal(body, &secret)

	return secret
}
