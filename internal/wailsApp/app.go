package wailsApp

import (
	"embed"
	"fmt"
	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
	"main/pkg/utils"
)

var (
	App    *application.App
	Window *application.WebviewWindow
)

func AppInit(fs *embed.FS) {
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the web files.
	// 'Bind' is a list of Go struct instances. The web has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.

	App = application.New(application.Options{
		Name:        "GoPixiv",
		Description: "Pivix Crawler",

		Services: []application.Service{},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(fs),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: false,
		},
		Windows: application.WindowsOptions{
			WebviewUserDataPath: "",
			WebviewBrowserPath:  "",
		},
		Linux: application.LinuxOptions{
			ProgramName: "GoPixiv",
		},
	})
	Window = App.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "GoPixiv",
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		Windows: application.WindowsWindow{},
		//Windows: application.WindowsWindow{
		//	ExStyle: w32.WS_EX_NOREDIRECTIONBITMAP | w32.WS_EX_TOOLWINDOW | w32.WS_EX_TOPMOST,
		//},
		AlwaysOnTop:         false,
		MinimiseButtonState: application.ButtonEnabled,
		MaximiseButtonState: application.ButtonDisabled,
		// Frameless:        true,
		Width:            1024,
		Height:           768,
		MinWidth:         1024,
		MinHeight:        768,
		MaxWidth:         1024,
		MaxHeight:        768,
		BackgroundColour: application.NewRGBA(233, 233, 233, 128),
		URL:              "/",

		//ShouldClose: func(window *application.WebviewWindow) bool {
		//	window.Hide()
		//	return false
		//},
	})
	Window.OnWindowEvent(events.Windows.WindowClosing, func(e *application.WindowEvent) {
		fmt.Println(123456)
		utils.DebugLog.Println(114514)
		Window.Hide()
		e.Cancel()

	})
	Tray := App.SystemTray.New()
	b, err := fs.ReadFile("web/dist/appicon.png")
	if err != nil {
		fmt.Println(err)
	}
	Tray.SetTemplateIcon(b)
	Tray.SetIcon(b)
	TrayMenu := App.NewMenu()
	TrayShow := TrayMenu.Add("显示主界面")
	TrayShow.OnClick(func(*application.Context) {
		App.Show()
		Window.Show()
	})
	TrayShutDown := TrayMenu.Add("退出")
	TrayShutDown.OnClick(func(*application.Context) {
		App.Quit()
	})
	Tray.OnClick(func() {
		Window.Show()
	})
	Tray.OnRightClick(func() {
		Tray.OpenMenu()
	})
	Tray.SetMenu(TrayMenu)

}
