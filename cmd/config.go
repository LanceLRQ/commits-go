package cmd

import (
	"fmt"

	"github.com/LanceLRQ/commits-go/config"
	"github.com/LanceLRQ/commits-go/utils"
	"github.com/urfave/cli/v2"
)

func ConfigCommand() *cli.Command {
	return &cli.Command{
		Name:    "config",
		Aliases: []string{"cfg"},
		Usage:   utils.Translate("config_command_usage"),
		Subcommands: []*cli.Command{
			{
				Name:  "get",
				Usage: utils.Translate("config_get_command_usage"),
				Action: func(c *cli.Context) error {
					cfg, err := config.ReadConfig()
					if err != nil {
						return err
					}
					fmt.Println(cfg)
					return nil
				},
			},
			{
				Name:  "set",
				Usage: utils.Translate("config_set_command_usage"),
				Action: func(c *cli.Context) error {
					if c.NArg() < 2 {
						return utils.TranslateErrorF("wrong_arguments")
					}

					key := c.Args().Get(0)
					value := c.Args().Get(1)

					cfg, err := config.ReadConfig()
					if err != nil {
						return err
					}

					err = config.SetConfigValue(cfg, key, value)
					if err != nil {
						return err
					}

					// 写入前会进行校验
					return config.WriteConfig(cfg)
				},
			},
		},
	}
}
