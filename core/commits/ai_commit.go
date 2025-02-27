package commits

import (
	"fmt"
	"github.com/LanceLRQ/commits-go/config"
	"github.com/urfave/cli/v2"
	"os"
)

func AICommit(ctx *cli.Context, cfg *config.Config) error {

	// check current directory
	_, err := getGitRepo()
	if err != nil {
		return err
	}

	// get staged changes
	files, err := getGitStagedChanges(true, true, []string{})
	if err != nil {
		return err
	}

	fmt.Println(files)

	err = os.Chdir("/Users/lancelrq/playground/test")
	if err != nil {
		return err
	}
	GitStashView()

	return nil
}
