package main

import (
	"embed"
	"gobot/backend"
	"gobot/backend/config"
	"gobot/backend/dao"
	"gobot/backend/log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if err := config.InitConfig("./config.yml"); err != nil {
		println("Init config error:" + err.Error())
		panic(err)
	}
	log.Init()
	dao.InitKvDB()
	// Create an instance of the app structure
	app := backend.NewApp()
	// Create application with options
	err := wails.Run(&options.App{
		Title:  "GoBot",
		Width:  1440,
		Height: 720,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless:       true,
		CSSDragProperty: "user-select",
		CSSDragValue:    "none",
		MinWidth:        400,
		MinHeight:       400,
		// BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:  app.Startup,
		OnShutdown: app.Shutdown,
		Bind: []interface{}{
			app,
		},
		Debug: options.Debug{
			OpenInspectorOnStartup: false,
		},
		EnableDefaultContextMenu: false,
		Logger:                   log.Logger,
		LogLevel:                 logger.DEBUG,
	})

	if err != nil {
		log.Logger.Error("启动失败")
	}
	log.Logger.Info("服务启动成功")
}
