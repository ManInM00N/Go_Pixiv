package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	. "main/backend/src/init"
	"net/http"

	"github.com/tidwall/gjson"
	"github.com/wailsapp/wails/v3/pkg/application"
)

var (
	App    *application.App
	Window *application.WebviewWindow
)

func AppInit() {
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	App = application.New(application.Options{
		Name:        "GoPixiv",
		Description: "Pivix Crawler",

		Services: []application.Service{
			application.NewService(&GreetService{}),
			application.NewService(NewCtl()),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
		Windows: application.WindowsOptions{
			WebviewUserDataPath: "",
			WebviewBrowserPath:  "",
		},
		Linux: application.LinuxOptions{
			ProgramName: "GoPixiv",
		},
	})
	Window = App.NewWebviewWindowWithOptions(application.WebviewWindowOptions{
		Title: "GoPixiv",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		// Windows: application.WindowsWindow{
		// 	ExStyle: w32.WS_EX_NOREDIRECTIONBITMAP | w32.WS_EX_TOOLWINDOW | w32.WS_EX_TOPMOST,
		// },
		// Frameless:        true,
		Width:            1024,
		Height:           768,
		MinWidth:         1024,
		MinHeight:        768,
		MaxWidth:         768,
		MaxHeight:        768,
		BackgroundColour: application.NewRGBA(233, 233, 233, 128),
		URL:              "/",
		ShouldClose: func(window *application.WebviewWindow) bool {
			window.Hide()
			return false
		},
	})

	Tray := App.NewSystemTray()
	b, err := assets.ReadFile("frontend/dist/appicon.png")
	if err != nil {
		fmt.Println(err)
	}
	Tray.SetTemplateIcon(b)
	TrayMenu := App.NewMenu()
	TrayShow := TrayMenu.Add("显示主界面")
	TrayShow.OnClick(func(*application.Context) {
		App.Show()
	})
	TrayShutDown := TrayMenu.Add("退出")
	TrayShutDown.OnClick(func(*application.Context) {
		App.Quit()
	})
	Tray.OnClick(func() {
		Tray.OpenMenu()
	})
	Tray.SetMenu(TrayMenu)
}

type Ctl struct {
	ctx context.Context
}

func NewCtl() *Ctl {
	return &Ctl{}
}

func (a *Ctl) Close() {
	fmt.Println(1)
	Window.Hide()
	// App.Hide()
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

func (a *Ctl) DownloadByNovelId(text string) bool {
	InfoLog.Println("Download Novel ", text)
	Download_By_NovelId(text)
	return true
}

func (a *Ctl) ReturnString() string {
	return NowTaskMsg
}

func Close() {
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
	url, ref := GetUrlRefer("", "", 0)
	Request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		DebugLog.Println("Error creating request", err)
		App.EmitEvent("login", "False")
		return false
	}
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
		App.EmitEvent("login", "False")
		return false
	}
	var buffer bytes.Buffer
	reader := bufio.NewReader(res.Body)
	io.Copy(&buffer, reader)
	data := buffer.Bytes()
	Results := gjson.ParseBytes(data)
	canbedownload := Results.Get("error").Bool()
	DebugLog.Println("Loading succeed? :", !canbedownload)
	if canbedownload {
		App.EmitEvent("login", "False", "Make sure your Cookie not Expired and Existed")
	} else {
		App.EmitEvent("login", "True")
	}
	return true
}
