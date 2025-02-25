package commits

import (
	"fmt"
	"strings"

	"github.com/LanceLRQ/commits-go/utils"
)

type GitFileStatus struct {
	StagedStatus   string
	UnstagedStatus string
	Path           string
	TargetPath     string
}

func getGitRepo() (string, error) {
	repoDir, err := utils.ExecCommand("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", utils.TranslateErrorF("not_a_git_repository")
	}

	return repoDir, nil
}

func getGitStagedChanges(cached bool, nameOnly bool, excludeFiles []string) ([]string, error) {
	baseArgs := []string{"diff", "--diff-algorithm=minimal"}
	if cached {
		baseArgs = append(baseArgs, "--cached")
	}

	if nameOnly {
		baseArgs = append(baseArgs, "--name-only")
	}

	fileExclusions := []string{
		":(exclude)package-lock.json",
		":(exclude)pnpm-lock.yaml",
		":(exclude)*.lock",
	}

	for _, file := range excludeFiles {
		fileExclusions = append(fileExclusions, fmt.Sprintf(":(exclude)%s", file))
	}

	diff, err := utils.ExecCommand("git", append(baseArgs, fileExclusions...)...)
	if err != nil {
		return nil, err
	}

	return strings.Split(strings.TrimSpace(diff), "\n"), nil
}

func getGitStatus() ([]GitFileStatus, bool, error) {
	statusTextRaw, err := utils.ExecCommand("git", "status", "-u", "--porcelain")
	if err != nil {
		return nil, false, err
	}

	statusText := strings.Split(strings.TrimSpace(statusTextRaw), "\n")

	gitFileResult := []GitFileStatus{}
	hasUnstagedChanges := false

	for _, status := range statusText {
		split := strings.Fields(status)
		gitFileStatus := GitFileStatus{
			Path: split[1],
		}
		stat := strings.Split(split[0], "")

		if len(stat) == 2 {
			gitFileStatus.StagedStatus = stat[0]
			gitFileStatus.UnstagedStatus = stat[1]
			if stat[0] == "?" || stat[1] == "?" {
				hasUnstagedChanges = true
			}
		} else {
			gitFileStatus.StagedStatus = stat[0]
			if stat[0] == "?" {
				hasUnstagedChanges = true
			}
		}
		if len(split) == 4 && split[2] == "->" {
			gitFileStatus.TargetPath = split[3]
		}

		gitFileResult = append(gitFileResult, gitFileStatus)
	}

	return gitFileResult, hasUnstagedChanges, nil
}
