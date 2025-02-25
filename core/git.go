package core

import (
	"github.com/LanceLRQ/commits-go/utils"
)

func getGitRepo() (string, error) {
	repoDir, err := execCommand("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", utils.TranslateErrorF("not_a_git_repository")
	}

	return repoDir, nil
}
