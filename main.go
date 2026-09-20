package main

import (
	"embed"
	"gofly/internal/service"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"
)

// Wails uses Go's `embed` package to embed the frontend files into the binary.
// Any files in the frontend/dist folder will be embedded into the binary and
// made available to the frontend.
// See https://pkg.go.dev/embed for more information.

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	// Register a custom event whose associated data type is string.
	// This is not required, but the binding generator will pick up registered events
	// and provide a strongly typed JS/TS API for them.
	application.RegisterEvent[string]("time")
}

// main function serves as the application's entry point. It initializes the application, creates a window,
// and starts a goroutine that emits a time-based event every second. It subsequently runs the application and
// logs any error that might occur.
func main() {

	// var serviceList []application.Service
	// for _, inst := range svcs.All {
	// 	serviceList = append(serviceList, application.NewService(&inst))
	// }
	// Create a new Wails application by providing the necessary options.
	// Variables 'Name' and 'Description' are for application metadata.
	// 'Assets' configures the asset server with the 'FS' variable pointing to the frontend files.
	// 'Bind' is a list of Go struct instances. The frontend has access to the methods of these instances.
	// 'Mac' options tailor the application when running an macOS.
	app := application.New(application.Options{
		Name:        "MQTTW",
		Description: "MQTTW MQTT调试与压力测试桌面客户端",
		Services:    service.NewAppServices(),
		// Services: []application.Service{
		// 	application.NewService(&GreetService{}),
		// },
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// Create a new window with the necessary options.
	// 'Title' is the title of the window.
	// 'Mac' options tailor the window when running on macOS.
	// 'BackgroundColour' is the background colour of the window.
	// 'URL' is the URL that will be loaded into the webview.
	win := app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "MQTTW - MQTT调试压测工具",
		// Window sized to the golden ratio (1000 / 618 ≈ 1.618).
		Width:     1200,
		Height:    800,
		MinWidth:  900,
		MinHeight: 500,
		Mac: application.MacWindow{
			InvisibleTitleBarHeight: 50,
			Backdrop:                application.MacBackdropTranslucent,
			TitleBar:                application.MacTitleBarHiddenInset,
		},
		BackgroundType:   application.BackgroundTypeSolid,
		BackgroundColour: application.NewRGB(255, 255, 255),
		URL:              "/",
		Windows: application.WindowsWindow{
			Theme: application.SystemDefault, //跟随系统主题
		},
	})

	// 系统托盘（最小化到托盘常驻，框架在 app.Run() 后自动创建）
	tray := app.SystemTray.New()
	trayMenu := application.NewMenu()
	trayMenu.Add("显示主窗口").OnClick(func(_ *application.Context) { win.Show() })
	trayMenu.Add("隐藏窗口").OnClick(func(_ *application.Context) { win.Hide() })
	trayMenu.AddSeparator()
	trayMenu.Add("退出").OnClick(func(_ *application.Context) { app.Quit() })
	tray.SetMenu(trayMenu)
	tray.AttachWindow(win)

	// 关闭窗口拦截：设置了"最小化到托盘"则隐藏窗口而不是退出
	win.RegisterHook(events.Common.WindowClosing, func(event *application.WindowEvent) {
		if service.ShouldMinimizeOnClose() {
			event.Cancel()
			win.Hide()
		}
	})

	// Run the application. This blocks until the application has been exited.
	err := app.Run()

	// If an error occurred while running the application, log it and exit.
	if err != nil {
		log.Fatal(err)
	}
}
