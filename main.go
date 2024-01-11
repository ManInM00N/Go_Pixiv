package main

import (
	"embed"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:       "Go!Pixiv",
		Width:       1024,
		Height:      768,
		Fullscreen:  false,
		StartHidden: false,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 233, G: 233, B: 233, A: 128},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
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
			// User messages that can be customised
			Messages: &windows.Messages{},
			// OnSuspend is called when Windows enters low power mode
			OnSuspend: func() {},
			// OnResume is called when Windows resumes from low power mode
			OnResume:             func() {},
			WebviewGpuIsDisabled: false,
		},
		//Mac: &mac.Options{
		//	TitleBar: &mac.TitleBar{
		//		TitlebarAppearsTransparent: true,
		//		HideTitle:                  false,
		//		HideTitleBar:               false,
		//		FullSizeContent:            false,
		//		UseToolbar:                 false,
		//		HideToolbarSeparator:       true,
		//	},
		//	Appearance:           mac.NSAppearanceNameDarkAqua,
		//	WebviewIsTransparent: true,
		//	WindowIsTranslucent:  false,
		//	About: &mac.AboutInfo{
		//		Title:   "My Application",
		//		Message: "© 2021 Me",
		//		Icon:    icon,
		//	},
		//},
		//Linux: &linux.Options{
		//	Icon: icon,
		//	WindowIsTranslucent: false,
		//	WebviewGpuPolicy: linux.WebviewGpuPolicyAlways,
		//	ProgramName: "wails",
		//},
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
