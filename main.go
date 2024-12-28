/*
//go:build windows
*/
package main

import (
	"embed"
	_ "embed"
	. "main/backend/src/init"
	"time"
	// "github.com/wailsapp/wails/v3/pkg/w32"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed frontend/dist
var assets embed.FS

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {
	ServerInit()
	go func() {
		R.Run(":7234")
	}()
	AppInit()

	go func() {
		for {
			now := time.Now().Format(time.RFC1123)
			App.EmitEvent("time", now)
			time.Sleep(time.Second)
		}
	}()

	// Run the application. This blocks until the application has been exited.
	err := App.Run()
	Close()
	// If an error occurred while running the application, log it and exit.
	DebugLog.Println(err)
}
