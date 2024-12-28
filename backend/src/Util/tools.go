package Util

import (
	"fmt"
	"image"
	"os"
	"regexp"
	"strings"
)

var pattern = [...]string{"R-18", "r-18", "r18", "R18"}

func HasR18(raw *[]string) bool {
	for _, v := range *raw {
		if len(v) > 4 || len(v) < 3 {
			continue
		}
		for _, p := range pattern {
			if v == p {
				return true
			}
		}
	}
	return false
}

func GetWH(imagePath string) (width, height int) {
	file, _ := os.Open(imagePath)

	c, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Println("err1 = ", err)
		return
	}
	width = c.Width
	height = c.Height

	file.Close()
	return
}

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
