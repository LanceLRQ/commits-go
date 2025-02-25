package core

import (
	"os/exec"
)

func execCommand(Program string, args ...string) (string, error) {
	cmd := exec.Command(Program, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
