package cmd

import (
	"log"
	"os"

	"github.com/LanceLRQ/commits-go/config"
	"github.com/LanceLRQ/commits-go/core/commits"
	"github.com/LanceLRQ/commits-go/utils"
	"github.com/urfave/cli/v2"
)

func CmmandEntry() {
	app := &cli.App{
		Name:  "commits-go",
		Usage: utils.Translate("application_welcome"),
		Commands: []*cli.Command{
			ConfigCommand(),
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
