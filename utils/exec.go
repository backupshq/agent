package utils

import (
	"bytes"
	"errors"
	"os/exec"
	"syscall"
)

func ExecuteCommand(cmd string, env []string, cancelChannel <-chan bool) (string, error) {
	c := exec.Command("sh", "-c", cmd)
	c.Env = env
	// TODO: The syscall code for killing all child processes won't
	// work on all platforms
	c.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	var stdout bytes.Buffer
	c.Stdout = &stdout

	c.Start()

	done := make(chan error, 1)
	go func() {
		done <- c.Wait()
	}()

	select {
	case <-cancelChannel:
		// Use syscall.Kill instead of c.Process.Kill to kill the child processes too
		err := syscall.Kill(-c.Process.Pid, syscall.SIGKILL)
		if err != nil {
			return "", errors.New("Unable to kill process: " + err.Error())
		}
		return "", errors.New("cancelled")
	case err := <-done:
		if err != nil {
			return "", errors.New("Error executing command: " + err.Error())
		}
	}

	return string(stdout.Bytes()), nil
}
