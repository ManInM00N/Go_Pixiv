package test

import (
	"bufio"
	"bytes"
	"io"
	. "main/backend/src/init"
	"net/http"
	"testing"

	"github.com/tidwall/gjson"
)

func TestLogin(t *testing.T) {
	url, ref := CheckMode("", "", 0)
	Request, err := http.NewRequest("GET", url, nil)
	client := GetClient()
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0")
	Request.Header.Set("referer", ref)
	Cookie := &http.Cookie{
		Name:  "PHPSESSID",
		Value: Setting.Cookie,
	}
	Request.AddCookie(Cookie)
	Request.Header.Set("PHPSESSID", Setting.Cookie)
	res, err := client.Do(Request)
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
