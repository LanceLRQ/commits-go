package commits

import (
	"fmt"
	"strings"

	"github.com/LanceLRQ/commits-go/utils"
)

type GitFileStatus struct {
	StagedStatus   rune
	UnstagedStatus rune
	Path           string
	TargetPath     string
}

type GitChangeStatus struct {
	hasStagedChanges   bool
	hasUnstagedChanges bool
	hasUntrackedFiles  bool
}

func getGitRepo() (string, error) {
	repoDir, err := utils.ExecCommand("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return "", err
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

func getGitStatus() ([]GitFileStatus, GitChangeStatus, error) {
	gitChangeStatus := GitChangeStatus{}

	statusTextRaw, err := utils.ExecCommand("git", "status", "-u", "--porcelain")
	if err != nil {
		return nil, gitChangeStatus, err
	}

	if strings.TrimSpace(statusTextRaw) == "" {
		return nil, gitChangeStatus, nil
	}

	gitFileResult := make([]GitFileStatus, 0)

	statusText := strings.Split(statusTextRaw, "\n")

	for _, status := range statusText {
		if len(status) < 3 {
			continue
		}
		flags := []rune(status[:2])
		split := strings.Fields(status[2:])
		gitFileStatus := GitFileStatus{
			Path: split[0],
		}

		if flags[0] != ' ' {
			gitFileStatus.StagedStatus = flags[0]
			gitChangeStatus.hasStagedChanges = true
		}
		if flags[1] != ' ' {
			gitFileStatus.UnstagedStatus = flags[1]
			gitChangeStatus.hasUnstagedChanges = true
		}

		if gitFileStatus.StagedStatus == '?' || gitFileStatus.UnstagedStatus == '?' {
			gitChangeStatus.hasUntrackedFiles = true
		}

		if len(split) == 3 && split[1] == "->" {
			gitFileStatus.TargetPath = split[2]
		}

		gitFileResult = append(gitFileResult, gitFileStatus)
	}

	return gitFileResult, gitChangeStatus, nil
}
