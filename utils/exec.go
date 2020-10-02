package utils

import (
	"errors"
	"os/exec"
)

func ExecuteCommand(cmd string, env []string) (string, error) {
	c := exec.Command("sh", "-c", cmd)
	c.Env = env
	out, err := c.Output()
	if err != nil {
		return "", errors.New("Error executing command: " + err.Error())
	}
	return string(out), nil
}
