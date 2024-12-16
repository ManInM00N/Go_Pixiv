package Util

import (
	"fmt"
	"image"
	"os"
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
