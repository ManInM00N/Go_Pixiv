package main

import (
	"bufio"
	"bytes"
	"context"
	"io"
	. "main/backend/src/init"
	"net/http"

	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var App *application.App

type Ctl struct {
	ctx context.Context
}

func NewCtl() *Ctl {
	return &Ctl{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *Ctl) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *Ctl) DownloadByPid(text string) bool {
	InfoLog.Println("Download illust ", text)
	Download_By_Pid(text)
	return true
}

func (a *Ctl) ReturnString() string {
	return NowTaskMsg
}

func (a *Ctl) Close(ctx context.Context) {
	IsClosed = true
	P.Wait()
	defer func() {
		P.Release()
		TaskPool.Close()
		SinglePool.Release()
	}()
}

func (a *Ctl) DownloadByAuthorId(text string) bool {
	Download_By_Author(text, func(name string, data ...interface{}) {
		App.EmitEvent(name, data)
	})
	return true
}

func (a *Ctl) DownloadByRank(text, Type string) bool {
	Download_By_Rank(text, Type, func(name string, data ...interface{}) {
		App.EmitEvent(name, data)
	})
	return true
}

func (a *Ctl) DownloadByFollowPage(page, Type string) bool {
	Download_By_FollowPage(page, Type, func(name string, data ...interface{}) {
		App.EmitEvent(name, data)
	})
	return true
}

func (a *Ctl) CheckLogin() bool {
	url, ref := CheckMode("", "", 0)
	Request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		DebugLog.Println("Error creating request", err)
		App.EmitEvent("login", "False")
		return false
	}
	client := GetClient()
	Request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0")
	Request.Header.Set("referer", ref)
	Cookie := &http.Cookie{
		Name:  "PHPSESSID",
		Value: Setting.Cookie,
	}
	Request.AddCookie(Cookie)
	Request.Header.Set("PHPSESSID", Setting.Cookie)
	var res *http.Response
	for i := 0; i < 3; i++ {
		res, err = client.Do(Request)
		if err == nil {
			break
		}
	}
	if err != nil {
		App.EmitEvent("login", "False")
		return false
	}
	var buffer bytes.Buffer
	reader := bufio.NewReader(res.Body)
	io.Copy(&buffer, reader)
	data := buffer.Bytes()
	Results := gjson.ParseBytes(data)
	canbedownload := Results.Get("error").Bool()
	println(canbedownload)
	DebugLog.Println("Loading succeed? :", !canbedownload)
	if canbedownload {
		App.EmitEvent("login", "False", "Make sure your Cookie not Expired and Existed")
	} else {
		App.EmitEvent("login", "True")
	}
	return true
}
