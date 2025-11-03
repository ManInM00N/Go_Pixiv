package test

import (
	"bytes"
	_ "image/jpeg"
	_ "image/png"
	"main/configs"
	"main/internal/imageService"
	"os"
	"testing"
)

func TestImgUpload(t *testing.T) {
	file, err := os.Open("./124613508_p1.png")
	if err != nil {
		t.Error(err)
	}
	//img, _, _ := image.Decode(file)
	bytes := bytes.NewBuffer(nil)
	bytes.ReadFrom(file)
	upload, err := imageService.ImgManager.Upload(bytes.Bytes(), "124613508_p1.png")
	if err != nil {
		return
	}
	t.Error(upload.ServiceName, upload.URL)
}

// has upload  http://tmpfiles.org/dl/5975186/124613508_p1.png

func TestImgSearch(t *testing.T) {
	url := "http://tmpfiles.org/dl/5975186/124613508_p1.png"
	setting := configs.NowSetting()
	service := imageService.NewSauceNAOService(setting.SearchEngine.SauceNaoConf.ApiKey)
	byURL, err := service.SearchByURL(url)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(byURL.Header)
	t.Error(byURL.Results)

}
