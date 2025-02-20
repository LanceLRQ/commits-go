package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LanceLRQ/commits-go/config"
	"github.com/LanceLRQ/commits-go/utils"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"
)

// 获取配置文件路径
func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".commits-go")
}

// 读取配置
func readConfig() (*config.Config, error) {
	path := getConfigPath()
	data, err := os.ReadFile(path)
	cfg := config.GetDefaultConfig()
	if os.IsNotExist(err) {
		return &cfg, nil
	} else if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}

// 写入配置
func writeConfig(config *config.Config) error {
	// 先进行校验
	// if err := config.Validate(); err != nil {
	// 	return err
	// }

	// 获取配置文件路径
	path := getConfigPath()
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

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
					cfg, err := readConfig()
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
						return fmt.Errorf("需要提供配置项和值")
					}

					key := c.Args().Get(0)
					value := c.Args().Get(1)

					cfg, err := readConfig()
					if err != nil {
						return err
					}

					err = config.SetConfigValue(cfg, key, value)
					if err != nil {
						return err
					}

					// 写入前会进行校验
					return writeConfig(cfg)
				},
			},
		},
	}
}
