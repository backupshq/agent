package actions

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
}

func GetBackup(client *http.Client, tokenResponse auth.AccessTokenResponse, backupId string) Backup {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://localhost:8000/backups/%s", backupId), nil)
	if err != nil {
		log.Fatal("Error reading request. ", err)
	}
	auth.AddAuthHeader(req, tokenResponse)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error reading response. ", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatal("Unable to start job: HTTP " + resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}

	var backup Backup
	err = json.Unmarshal(body, &backup)
	return backup
}
