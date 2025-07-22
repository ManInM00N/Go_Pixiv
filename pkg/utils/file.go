package utils

import (
	"os"
	"regexp"
	"strings"
)

func GetFileType(path string) string {
	index := strings.LastIndex(path, ".")
	path = path[index+1:]
	return path
}

const maxLength = 210 // 通常文件名的最大长度限制
func Cut(filename string) string {
	if len(filename) > maxLength {
		filename = filename[:maxLength]
	}
	return filename
}

func filterFilename(filename string) string {
	re := regexp.MustCompile(`[?*:\x20\"<>|/\\]`)
	return re.ReplaceAllString(filename, "")
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
