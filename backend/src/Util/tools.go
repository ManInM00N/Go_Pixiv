package Util

import (
	"fmt"
	"image"
	"os"
)

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
