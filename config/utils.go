package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"gopkg.in/yaml.v3"

	"github.com/LanceLRQ/commits-go/utils"
)

func GetDefaultConfig() Config {
	return Config{
		Locale: "zh-CN",
		Proxy: ProxyConfig{
			Enabled:  false,
			Url:      "",
			UserName: "",
			Password: "",
		},
		LLM: LLMConfig{
			Temperature:      0.7,
			TopP:             1.0,
			MaxTokens:        150,
			FrequencyPenalty: 0.0,
			PresencePenalty:  0.0,
		},
		Server: OpenAIConfig{
			ApiUrl:  "https://api.deepseek.com/v1",
			ApiKey:  "",
			Model:   "deepseek-chat",
			Timeout: 30,
		},
		Git: GitCommitConfig{
			MaxLength:     50,
			SuggestLength: 1,
			CommitStyle:   "",
		},
	}
}

// equal 函数支持将驼峰式变量名转换为下划线分隔形式，并进行不区分大小写的比较
func equalFieldName(fieldName, key string) bool {
	// 将驼峰式变量名转换为下划线分隔形式
	normalizedFieldName := camelToSnake(fieldName)
	// 进行不区分大小写的比较
	return strings.EqualFold(normalizedFieldName, key)
}

// camelToSnake 将驼峰式变量名转换为下划线分隔形式
func camelToSnake(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// SetConfigValue 设置配置值
func SetConfigValue(cfg interface{}, key, value string) error {
	keys := strings.Split(key, ".")
	if len(keys) == 0 {
		return fmt.Errorf("invalid configuration key format")
	}

	v := reflect.ValueOf(cfg)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	field := v.FieldByNameFunc(func(s string) bool {
		return equalFieldName(s, string(keys[0]))
	})

	if !field.IsValid() {
		return utils.TranslateErrorF("unknown_configuration_section", keys[0])
	}

	if len(keys) == 1 {
		if !field.CanSet() {
			return utils.TranslateErrorF("cannot_set_configuration_key", keys[0])
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(value)
		case reflect.Int:
			intValue, err := strconv.Atoi(value)
			if err != nil {
				return utils.TranslateErrorF("invalid_value_for", keys[0], value)
			}
			field.SetInt(int64(intValue))
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return utils.TranslateErrorF("invalid_value_for", keys[0], value)
			}
			field.SetBool(boolValue)
		case reflect.Float64:
			floatValue, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return utils.TranslateErrorF("invalid_value_for", keys[0], value)
			}
			field.SetFloat(floatValue)
		default:
			return utils.TranslateErrorF("unsupported_configuration_key_type", keys[0])
		}

		return nil
	}

	if field.Kind() != reflect.Struct {
		return utils.TranslateErrorF("cannot_set_nested_configuration_key", key)
	}

	return SetConfigValue(field.Addr().Interface(), strings.Join(keys[1:], "."), value)
}

// GetConfigPath 获取配置文件路径
func getConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".commits-go")
}

// ReadConfig 读取配置
func ReadConfig() (*Config, error) {
	path := getConfigPath()
	data, err := os.ReadFile(path)
	cfg := GetDefaultConfig()
	if os.IsNotExist(err) {
		return &cfg, nil
	} else if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &cfg)
	return &cfg, err
}

// WriteConfig 写入配置
func WriteConfig(config *Config) error {
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
