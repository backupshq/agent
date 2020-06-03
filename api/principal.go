package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type Principal struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Name string `json:"name"`
}

func (c *ApiClient) GetPrincipal(id string) Principal {
	req, err := c.get(fmt.Sprintf("/principals/%s", id))
	if err != nil {
		log.Fatal("Error creating request. ", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("Unable to get principal: HTTP " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	var principal Principal
	err = json.Unmarshal(body, &principal)

	return principal
}
