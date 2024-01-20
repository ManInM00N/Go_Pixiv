package main

import (
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
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
func (a *App) PreloadRank(text, Type string) bool {
	LoadingNow = true
	DownloadRankMsg(text, Type, func(name string, data ...interface{}) {
		runtime.EventsEmit(a.ctx, name, data)
	})
	return true
}
func (a *App) PopLoadPool() {
	LoadingNow = false
	RankPool.Wait()
	LoadingNow = true
}
