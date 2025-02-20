package cmd

import (
	"log"
	"os"

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
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
