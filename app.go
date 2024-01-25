package main

import (
	"bufio"
	"bytes"
	"context"
	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io"
	"main/backend/src/DAO"
	. "main/backend/src/init"
	"net/http"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}
func (a *App) DownloadByPid(text string) bool {
	println(text)
	Download_By_Pid(text)
	return true
}
func (a *App) ReturnString() string {
	return NowTaskMsg
}
func (a *App) Close(ctx context.Context) {
	IsClosed = true
	P.Wait()
	defer func() {
		P.Release()
		TaskPool.Close()
		SinglePool.Release()
	}()

}
func (a *App) DownloadByAuthorId(text string) bool {
	Download_By_Author(text, func(name string, data ...interface{}) {
		runtime.EventsEmit(a.ctx, name, data)
	})
	return true
}
func (a *App) DownloadByRank(text, Type string) bool {
	Download_By_Rank(text, Type, func(name string, data ...interface{}) {
		runtime.EventsEmit(a.ctx, name, data)
	})
	return true
}
func (a *App) DownloadByFollowPage(page, Type string) bool {
	Download_By_FollowPage(page, Type, func(name string, data ...interface{}) {
		runtime.EventsEmit(a.ctx, name, data)
	})
	return true
}
func (a *App) PreloadRank(text, Type, page string) bool {
	RankLoadingNow = true
	DownloadRankMsg(text, Type, page, func(name string, data ...interface{}) {
		runtime.EventsEmit(a.ctx, name, data)
	})
	return true
}
func (a *App) PreloadFollow(page, Type string) bool {
	FollowLoadingNow = true
	DownloadFollowMsg(page, Type, func(name string, data ...interface{}) {
		runtime.EventsEmit(a.ctx, name, data)
	})
	return true
}
func (a *App) PopFollowPool() {
	FollowLoadingNow = false
	FollowLoadPool.Wait()
	FollowLoadingNow = true
	runtime.EventsEmit(a.ctx, "PopUp")
}
func (a *App) PopLoadPool() {
	RankLoadingNow = false
	RankloadPool.Wait()
	RankLoadingNow = true
	runtime.EventsEmit(a.ctx, "RankmsgPopUp")
}
func (a *App) GetSetting() DAO.Settings {
	return Setting
}

func (a *App) UpdateSetting(data DAO.Settings) {
	//UpdateSettings()
	DebugLog.Println(data.MsgDetail())
	Setting.UpdateSettings(data)
	UpdateSettings()
	a.CheckLogin()
}

//	func (a *App) UpdateSetting(data string) {
//		newdata := Setting
//		jsonmsg, err := json.Marshal(data)
//		if err != nil {
//			return
//		}
//		json.Unmarshal(jsonmsg, newdata)
//		Setting.UpdateSettings(newdata)
//		UpdateSettings()
//		a.CheckLogin()
//	}
func (a *App) CheckLogin() bool {
	url, ref := CheckMode("", "", 0)
	Request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		DebugLog.Println("Error creating request", err)
		runtime.EventsEmit(a.ctx, "login", "False")

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
	res, err := client.Do(Request)
	if err != nil {
		runtime.EventsEmit(a.ctx, "login", "False")
		return false
	}
	var buffer bytes.Buffer
	reader := bufio.NewReader(res.Body)
	io.Copy(&buffer, reader)
	data := buffer.Bytes()
	Results := gjson.ParseBytes(data)
	canbedownload := Results.Get("error").Bool()
	println(canbedownload)
	DebugLog.Println("Loading results:", canbedownload)
	if canbedownload {
		runtime.EventsEmit(a.ctx, "login", "False")
	} else {
		runtime.EventsEmit(a.ctx, "login", "True")

	}
	return true
}
