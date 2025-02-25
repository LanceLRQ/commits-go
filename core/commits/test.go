package commits

import (
	"fmt"
)

func TemporaryEntry() error {
	_, err := getGitStagedChanges(true, true, []string{})
	if err != nil {
		return err
	}

	status, hasUnstagedChanges, err := getGitStatus()
	if err != nil {
		return err
	}

	fmt.Println(status)
	fmt.Println(hasUnstagedChanges)
	return nil
}
