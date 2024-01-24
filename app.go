package main

import (
	"context"
	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"main/backend/src/DAO"
	. "main/backend/src/init"
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
	Setting.UpdateSettings(data)
	UpdateSettings()

}

//	 func (a *App) UpdateSetting(data string) {
//		newdata := Setting
//		jsonmsg, err := json.Marshal(data)
//		if err != nil {
//			return
//		}
//		json.Unmarshal(jsonmsg, newdata)
//		Setting.UpdateSettings(newdata)
//	}
func (a *App) CheckLogin() bool {
	data, err := GetWebpageData("", "", 0)
	if err != nil {
		return false
	}
	Results := gjson.ParseBytes(data)
	canbedownload := Results.Get("error").Bool()
	println(canbedownload)
	if canbedownload {
		runtime.EventsEmit(a.ctx, "login", "False")
	} else {
		runtime.EventsEmit(a.ctx, "login", "True")

	}
	return true
}
