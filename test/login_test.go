package test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	. "main/configs"
	. "main/internal/pixivlib/DAO"
	"net/http"
	"testing"

	"github.com/tidwall/gjson"
)

func TestLogin(t *testing.T) {
	url, ref := GetUrlRefer("", "", 0)
	fmt.Println(url, ref)
	Request, err := http.NewRequest("GET", url, nil)
	client := GetClient()
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0")
	Request.Header.Set("referer", ref)
	Request.Header.Set("Cookie", "PHPSESSID="+Setting.Cookie)
	var res *http.Response
	for i := 0; i < 3; i++ {

		res, err = client.Do(Request)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Error("Do request failed", err)
		return
	}
	var buffer bytes.Buffer
	reader := bufio.NewReader(res.Body)
	io.Copy(&buffer, reader)
	data := buffer.Bytes()
	Results := gjson.ParseBytes(data)
	canbedownload := Results.Get("error").Bool()
	t.Log(Results.Str, string(data))
	if canbedownload {
		t.Error("login failed check the Cookie or User-Agent")
	}
}
