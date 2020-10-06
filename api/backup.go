package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type StepDefinition struct {
	Name   string
	Script struct {
		Script   string
		Checksum string
	}
	SortOrder   int `json:"sort_order"`
	Expressions map[string]string
}

type Backup struct {
	ID              string
	Name            string
	Description     string
	Managed         bool
	Schedule        string
	UpdatedAt       string           `json:"updated_at"`
	StepDefinitions []StepDefinition `json:"step_definitions"`
}

func (c *ApiClient) GetBackup(backupId string) Backup {
	req, err := c.get(fmt.Sprintf("/backups/%s", backupId))
	if err != nil {
		log.Fatal("Error creating request. ", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("Unable to get backup: HTTP " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	var backup Backup
	err = json.Unmarshal(body, &backup)
	return backup
}

func (c *ApiClient) ListBackups(principalId string) map[string]Backup {
	req, err := c.get("/backups")
	if err != nil {
		log.Fatal("Error creating request. ", err)
	}

	q := req.URL.Query()
	q.Add("assigned_to", principalId)
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("Unable to get backups: HTTP " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	allBackups := make([]Backup, 0)
	filteredBackupMap := map[string]Backup{}
	err = json.Unmarshal(body, &allBackups)
	if err != nil {
		log.Fatal("Error decoding json. ", err)
	}

	for i, _ := range allBackups {
		if !allBackups[i].Managed {
			continue
		}
		filteredBackupMap[allBackups[i].ID] = allBackups[i]
	}

	return filteredBackupMap
}
