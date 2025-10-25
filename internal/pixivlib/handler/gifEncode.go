package handler

import (
	"archive/zip"
	"bytes"
	"fmt"
	gifencoder "github.com/ManInM00N/nicogif"
	"image"
	. "main/internal/pixivlib/DAO"
	"main/pkg/utils"
	"runtime"

	_ "image/jpeg"
	_ "image/png"
)

// 流式提取和处理帧 - 每处理一帧就释放
func ExtractAndProcessFrames(zipData []byte, frameInfo []Frame, width, height int) ([]byte, error) {
	reader := bytes.NewReader(zipData)
	zipReader, err := zip.NewReader(reader, int64(len(zipData)))
	if err != nil {
		return nil, fmt.Errorf("创建 ZIP reader 失败: %w", err)
	}

	// 预分配帧切片
	//frames := make([]image.Image, 0, len(zipReader.File))

	var delays []int
	for _, v := range frameInfo {
		delays = append(delays, v.Delay)
	}

	var gif *gifencoder.GIFEncoder
	fmt.Println("zip files", len(zipReader.File))
	// 逐帧处理
	for idx, file := range zipReader.File {
		rc, err := file.Open()
		if err != nil {
			utils.DebugLog.Printf("打开文件 %s 失败: %v", file.Name, err)
			continue
		}

		// 解码图片
		img, _, err := image.Decode(rc)
		rc.Close()

		if err != nil {
			utils.DebugLog.Printf("解码图片 %s 失败: %v", file.Name, err)
			continue
		}

		if idx == 0 {
			width = img.Bounds().Dx()
			height = img.Bounds().Dy()
			gif = gifencoder.NewGIFEncoder(width, height)
			gif.SetRepeat(0)
			gif.SetDither(gifencoder.DitherFloydSteinberg)
			gif.SetQuality(10)
			gif.SetColorEnhancement(1.0, 1.0)
		}

		delay := 100 // default 100ms
		if idx < len(frameInfo) && frameInfo[idx].Delay > 0 {
			delay = frameInfo[idx].Delay
		}
		gif.SetDelay(delay)

		if err := gif.AddFrame(img); err != nil {
			utils.DebugLog.Println(err)
			return nil, err
		}

		//frames = append(frames, img)

		if (idx+1)%10 == 0 {
			runtime.GC()
		}
	}
	if gif == nil {
		utils.DebugLog.Println("GIF encode failed.")
		return nil, nil
	}
	gif.Finish()
	return gif.GetData(), nil
	//return frames, nil
}

func extractFramesFromZip(body []byte, width, height int64) ([]image.Image, error) {
	reader := bytes.NewReader(body)
	zipReader, err := zip.NewReader(reader, int64(len(body)))
	if err != nil {
		return nil, err
	}

	var frames []image.Image
	for _, file := range zipReader.File {
		rc, err := file.Open()
		if err != nil {
			continue
		}

		img, _, err := image.Decode(rc)
		rc.Close()

		if err == nil {
			frames = append(frames, img)
		}
	}

	return frames, nil
}

func encodeGIF(frames []image.Image, frameInfo []Frame, width, height int) ([]byte, error) {
	var delays []int
	for _, v := range frameInfo {
		delays = append(delays, v.Delay)
	}

	opt := gifencoder.EncodeOptions{
		Width:   width,
		Height:  height,
		Repeat:  0,
		Quality: 10,
		Dither:  gifencoder.DitherFloydSteinberg,
		Delays:  delays,
	}

	if len(frames) > 0 {
		opt.Width = frames[0].Bounds().Dx()
		opt.Height = frames[0].Bounds().Dy()
	}

	return gifencoder.EncodeGIFWithOptions(frames, opt)
}
