package commits

import (
	"github.com/LanceLRQ/commits-go/config"
	"github.com/LanceLRQ/commits-go/utils"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli/v2"
	"os"
)

func AICommit(c *cli.Context, cfg *config.Config) error {

	// check workdir
	gitDir, err := utils.GetValidatedDirectory(c.String("workdir"))
	if err != nil {
		return err
	}
	// change working directory
	err = os.Chdir(gitDir)
	if err != nil {
		return err
	}

	// check current directory
	_, err = getGitRepo()
	if err != nil {
		return utils.TranslateErrorF("not_a_git_repository", gitDir)
	}

	//// get staged changes
	//files, err := getGitStagedChanges(true, true, []string{})
	//if err != nil {
	//	return err
	//}

	model, err := NewGitStatusView()
	if err != nil {
		return err
	}

	if _, err := tea.NewProgram(model).Run(); err != nil {
		return err
	}

	return nil
}
