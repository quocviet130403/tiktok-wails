package main

import (
	"embed"
	"tiktok-wails/backend"
	"tiktok-wails/backend/initialize"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure

	err := initialize.InitServer()
	if err != nil {
		println("Error initializing server:", err.Error())
		return
	}

	app := backend.NewApp()

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "tiktok-wails",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
