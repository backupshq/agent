package utils

import (
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
