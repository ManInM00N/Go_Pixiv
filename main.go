package main

import (
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	. "main/backend/src/init"
	"net/http"
	"os"
	"strings"
)

//go:embed all:frontend/dist
var assets embed.FS

type FileLoader struct {
	http.Handler
}

func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

func (h *FileLoader) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var err error
	requestedFilename := strings.TrimPrefix(req.URL.Path, "/")
	println("Requesting file:", requestedFilename)
	fileData, err := os.ReadFile(requestedFilename)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte(fmt.Sprintf("Could not load file %s", requestedFilename)))
	}

	res.Write(fileData)
}

func main() {
	MainApp := NewApp()
	err := wails.Run(&options.App{
		Title:       "Go!Pixiv",
		Width:       1024,
		Height:      768,
		MinWidth:    1024,
		MinHeight:   768,
		Fullscreen:  false,
		StartHidden: false,
		AssetServer: &assetserver.Options{
			Assets:  assets,
			Handler: NewFileLoader(),
		},
		BackgroundColour: &options.RGBA{R: 233, G: 233, B: 233, A: 128},
		OnStartup:        MainApp.startup,
		//OnDomReady: MainApp.startup,
		Bind: []interface{}{
			MainApp,
			&Setting,
		},
		EnumBind: []interface{}{
			WaitingTasks,
			QueueTaskMsg,
		},
		ErrorFormatter: func(err error) any { return err.Error() },
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			BackdropType:                      windows.Mica,
			DisableWindowIcon:                 false,
			DisableFramelessWindowDecorations: false,
			WebviewUserDataPath:               "",
			WebviewBrowserPath:                "",
			Theme:                             windows.SystemDefault,
			CustomTheme: &windows.ThemeSettings{
				DarkModeTitleBar:   windows.RGB(20, 20, 20),
				DarkModeTitleText:  windows.RGB(200, 200, 200),
				DarkModeBorder:     windows.RGB(20, 0, 20),
				LightModeTitleBar:  windows.RGB(200, 200, 200),
				LightModeTitleText: windows.RGB(20, 20, 20),
				LightModeBorder:    windows.RGB(200, 200, 200),
			},
			Messages:             &windows.Messages{},
			OnSuspend:            func() {},
			OnResume:             func() {},
			WebviewGpuIsDisabled: false,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  false,
			About: &mac.AboutInfo{
				Title:   "Go!Pixiv",
				Message: "© 2021 Me",
				//Icon:    Logo,
			},
		},
		Linux: &linux.Options{
			//Icon:                Logo,
			WindowIsTranslucent: false,
			WebviewGpuPolicy:    linux.WebviewGpuPolicyAlways,
			ProgramName:         "Go!Pixiv",
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
	})
	IsClosed = true
	P.Wait()
	defer func() {
		P.Release()
		TaskPool.Close()
		SinglePool.Release()
	}()
	if err != nil {
		println("Error:", err.Error())
	}
}
