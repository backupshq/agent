package utils

import "os/exec"

func ExecuteCommand(cmd string) string {
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return "Error executing command"
	}
	return string(out)
}
