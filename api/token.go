package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type TokenInfo struct {
	ID          string `json:"id"`
	PrincipalId string `json:"principal_id"`
	AccountId   string `json:"account_id"`
}

func (c *ApiClient) GetCurrentToken() TokenInfo {
	req, err := c.get("/tokens/self")
	if err != nil {
		log.Fatal("Error creating request. ", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("Unable to get token: HTTP " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	var token TokenInfo
	err = json.Unmarshal(body, &token)

	return token
}
