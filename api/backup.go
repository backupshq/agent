package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"../auth"
)

const BACKUP_TYPE_UNMANAGED = 0
const BACKUP_TYPE_UNSCHEDULED = 1
const BACKUP_TYPE_SCHEDULED = 2

type Backup struct {
	ID          string
	Name        string
	Description string
	Type        int
	Command     string
	Schedule    string
}

func (c *ApiClient) GetBackup(backupId string) Backup {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8000/backups/%s", backupId), nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	auth.AddAuthHeader(req, c.tokenResponse)

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

func (c *ApiClient) ListBackups(backupType int) map[string]Backup {
	req, err := http.NewRequest("GET", "http://localhost:8000/backups/", nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	auth.AddAuthHeader(req, c.tokenResponse)

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
	json.Unmarshal(body, &allBackups)

	for i, _ := range allBackups {
		if allBackups[i].Type != backupType {
			continue
		}
		filteredBackupMap[allBackups[i].ID] = allBackups[i]
	}

	return filteredBackupMap
}
