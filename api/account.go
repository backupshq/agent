package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Account struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

func (c *ApiClient) GetAccount(id string) Account {
	req, err := c.get(fmt.Sprintf("/accounts/%s", id))
	if err != nil {
		log.Fatal("Error creating request. ", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("Unable to get account: HTTP " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	var account Account
	err = json.Unmarshal(body, &account)

	return account
}
