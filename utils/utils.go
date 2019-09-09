package utils

import (
	"encoding/base64"
	"errors"
	"os/exec"
)

func ExecuteCommand(cmd string) (string, error) {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "", errors.New("Error executing command: " + err.Error())
	}
	return string(out), nil
}

func Base64Encode(json string) string {
	return base64.StdEncoding.EncodeToString([]byte(json))
}
