package cmd

import (
	"log"
	"os"

	"github.com/LanceLRQ/commits-go/config"
	"github.com/LanceLRQ/commits-go/core/commits"
	"github.com/LanceLRQ/commits-go/utils"
	"github.com/urfave/cli/v2"
)

func CommandEntry() {
	app := &cli.App{
		Name:  "commits-go",
		Usage: utils.Translate("application_welcome"),
		Commands: []*cli.Command{
			ConfigCommand(),
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "workdir",
				Aliases: []string{"d"},
				Value:   "",
				Usage:   utils.Translate("git_repository_directory"),
			},
		},
		Action: func(c *cli.Context) error {
			cfg, err := config.ReadConfig()
			if err != nil {
				return err
			}
			return commits.AICommit(c, cfg)
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
