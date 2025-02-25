package utils

import (
	"os/exec"
)

func ExecCommand(Program string, args ...string) (string, error) {
	cmd := exec.Command(Program, args...)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}
