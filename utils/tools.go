package utils

import (
	"os"
	"path/filepath"
)

// GetValidatedDirectory 检查目录是否存在并返回绝对路径
func GetValidatedDirectory(dir string) (string, error) {
	absGitDir, err := filepath.Abs(dir)
	if err == nil {
		stat, err := os.Stat(absGitDir)
		if err == nil && stat.IsDir() {
			return absGitDir, nil
		}
		return "", TranslateErrorF("directory_not_found", dir)
	}
	return "", TranslateErrorF("directory_not_found", dir)
}
