package utils

import (
	"bufio"
	"io"
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

// WriteToFile 写入文件
func WriteToFile(filepath string, reader io.Reader) error {
	f, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	bufWriter := bufio.NewWriter(f)
	if _, err := io.Copy(bufWriter, reader); err != nil {
		return err
	}

	return bufWriter.Flush()
}

// 写入 GIF 文件
func WriteBytesToFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	bufWriter := bufio.NewWriter(f)
	defer bufWriter.Flush()

	_, err = bufWriter.Write(data)
	return err
}
