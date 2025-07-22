/*
//go:build windows
*/
package main

import (
	"context"
	_ "embed"
	"github.com/wailsapp/wails/v3/pkg/w32"
	. "main/init"
	. "main/internal/pixivSvr"
	"main/internal/pixivlib"
	. "main/internal/wailsApp"
	. "main/pkg/utils"
	"net/http"
	"time"
)

// Wails uses Go's `embed` package to embed the web files into the binary.
// Any files in the web/dist folder will be embedded into the binary and
// made available to the web.
// See https://pkg.go.dev/embed for more information.

func init() {
	// 尝试设置 DPI 感知（Windows 10 1607+ 推荐）
	err := w32.SetProcessDpiAwareness(w32.PROCESS_PER_MONITOR_DPI_AWARE)
	if err != nil {
		// 回退到旧版 API（Windows 8.1+）
		w32.SetProcessDPIAware()
	}
}
func main() {
	ServerInit()
	CacheInit()
	server := &http.Server{
		Addr:    ":7234",
		Handler: R,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			DebugLog.Fatalln(err)
		}
		//R.Run(":7234")
	}()
	AppInit(&Assets)
	pixivlib.NewCtl().RegisterService(App)
	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			App.Event.Emit("time", now)
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := App.Run()
	if err != nil {
		DebugLog.Fatalln(err)
		return
	}
	err = server.Shutdown(context.Background())
	if err != nil {
		DebugLog.Fatalln(err)
		return
	}
	Close()
	// If an error occurred while running the application, log it and exit.
	DebugLog.Println(err)
}
