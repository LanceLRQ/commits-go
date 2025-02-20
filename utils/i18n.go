package utils

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

var (
	bundle    *i18n.Bundle
	localizer *i18n.Localizer
)

// normalizeLang 将系统语言环境转换为标准的 language.Tag
func normalizeLang(lang string) language.Tag {
	// 去除 .utf-8 后缀
	lang = strings.Split(lang, ".")[0]
	// 将下划线替换为减号
	lang = strings.ReplaceAll(lang, "_", "-")
	// 解析为 language.Tag
	tag, _ := language.Parse(lang)
	return tag
}

func InitI18n() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	bundle.MustLoadMessageFile("locales/en.yaml")
	bundle.MustLoadMessageFile("locales/zh.yaml")

	// 根据系统语言设置本地化
	lang := os.Getenv("LANG")
	if lang == "" {
		lang = "zh-CN"
	}
	langTag := normalizeLang(lang)
	localizer = i18n.NewLocalizer(bundle, langTag.String())
}

func Translate(key string, args ...interface{}) string {
	// 先获取翻译后的字符串
	translated := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: key,
	})

	// 如果翻译后的字符串包含格式化占位符（如 %s），则进行格式化
	if len(args) > 0 {
		return fmt.Sprintf(translated, args...)
	}

	return translated
}

func TranslateErrorF(key string, args ...interface{}) error {
	return errors.New(Translate(key, args...))
}
