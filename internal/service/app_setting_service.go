package service

import (
	"os"
	"os/exec"
	"runtime"
	"sync/atomic"

	"gofly/internal/dao"
	"gofly/internal/utils/gf"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// AppSettingService 应用设置：开机自启 / 托盘 / 快捷键 / 运行目录 / 版本信息 / 启动恢复
type AppSettingService struct{}

// minimizeOnClose 关闭窗口时最小化到托盘（main.go 的 WindowClosing 拦截读取）
var minimizeOnClose atomic.Bool

// hotkeysRegistered 快捷键是否已注册
var hotkeysRegistered atomic.Bool

const (
	appVersion   = "1.0.0"
	repoURL      = "https://gitee.com/huang_li_shi_admin/mqttw"
	communityURL = "https://goflys.cn/"
	gmqtURL      = "https://goflys.cn/gmqt"
)

// GetAppInfo 应用信息（关于区使用）
func (s *AppSettingService) GetAppInfo() any {
	return gf.Success().SetData(map[string]any{
		"name":          "MQTTW",
		"version":       appVersion,
		"description":   "MQTT 调试与压测一体化桌面工具",
		"repo":          repoURL,
		"repo_releases": repoURL + "/releases",
		"community":     communityURL,
		"gmqt":          gmqtURL,
	})
}

// ========== 开机自启 ==========

// SetAutostart 设置开机自启
func (s *AppSettingService) SetAutostart(enabled bool) any {
	app := application.Get()
	if enabled {
		if err := app.Autostart.Enable(); err != nil {
			return gf.Failed().SetMsg("开启开机自启失败：" + err.Error())
		}
		return gf.Success().SetMsg("开机自启已开启")
	}
	if err := app.Autostart.Disable(); err != nil {
		return gf.Failed().SetMsg("关闭开机自启失败：" + err.Error())
	}
	return gf.Success().SetMsg("开机自启已关闭")
}

// GetAutostart 获取开机自启状态
func (s *AppSettingService) GetAutostart() any {
	enabled, err := application.Get().Autostart.IsEnabled()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetData(map[string]any{"enabled": enabled})
}

// ========== 最小化到托盘 ==========

// SetTrayMinimize 设置关闭窗口时最小化到托盘（托盘常驻，见 main.go）
func (s *AppSettingService) SetTrayMinimize(enabled bool) any {
	minimizeOnClose.Store(enabled)
	if enabled {
		return gf.Success().SetMsg("关闭窗口时将最小化到托盘")
	}
	return gf.Success().SetMsg("关闭窗口时将直接退出")
}

// GetTrayMinimize 获取最小化到托盘状态
func (s *AppSettingService) GetTrayMinimize() any {
	return gf.Success().SetData(map[string]any{"enabled": minimizeOnClose.Load()})
}

// ShouldMinimizeOnClose 供 main.go 关闭拦截判断
func ShouldMinimizeOnClose() bool { return minimizeOnClose.Load() }

// ========== 全局快捷键 ==========

// SetHotkeys 注册/注销全局快捷键
//   - Ctrl+Shift+H 显示/隐藏主窗口
//   - Ctrl+Shift+D 断开全部连接
//   - Ctrl+Shift+P 快速发布（事件 hotkey:publish，前端监听）
func (s *AppSettingService) SetHotkeys(enabled bool) any {
	app := application.Get()
	window := app.Window.Current()

	if !enabled {
		_ = app.GlobalShortcut.Unregister("ctrl+shift+h")
		_ = app.GlobalShortcut.Unregister("ctrl+shift+d")
		_ = app.GlobalShortcut.Unregister("ctrl+shift+p")
		hotkeysRegistered.Store(false)
		return gf.Success().SetMsg("全局快捷键已关闭")
	}

	// 显示/隐藏窗口
	if err := app.GlobalShortcut.Register("ctrl+shift+h", func() {
		if window.IsVisible() {
			window.Hide()
		} else {
			window.Show()
		}
	}); err != nil {
		return gf.Failed().SetMsg("注册快捷键 Ctrl+Shift+H 失败：" + err.Error())
	}
	// 断开全部连接
	if err := app.GlobalShortcut.Register("ctrl+shift+d", func() {
		svc := &MqttClientService{}
		svc.DisconnectAll()
	}); err != nil {
		return gf.Failed().SetMsg("注册快捷键 Ctrl+Shift+D 失败：" + err.Error())
	}
	// 快速发布（事件推送，前端实现发布面板定位）
	if err := app.GlobalShortcut.Register("ctrl+shift+p", func() {
		app.Event.Emit("hotkey:publish", "triggered")
	}); err != nil {
		return gf.Failed().SetMsg("注册快捷键 Ctrl+Shift+P 失败：" + err.Error())
	}

	hotkeysRegistered.Store(true)
	return gf.Success().SetMsg("全局快捷键已开启")
}

// GetHotkeys 获取快捷键启用状态
func (s *AppSettingService) GetHotkeys() any {
	return gf.Success().SetData(map[string]any{"enabled": hotkeysRegistered.Load()})
}

// ========== 运行目录 ==========

// GetRuntimeDir 返回运行时目录（日志/导出文件所在）
func (s *AppSettingService) GetRuntimeDir() any {
	dir, err := runtimeDir()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	return gf.Success().SetData(map[string]any{"dir": dir})
}

// OpenRuntimeDir 打开运行时目录（日志/导出文件所在）
func (s *AppSettingService) OpenRuntimeDir() any {
	dir, err := runtimeDir()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	if err := openFolder(dir); err != nil {
		return gf.Failed().SetMsg("打开目录失败：" + err.Error())
	}
	return gf.Success().SetMsg("已打开运行目录：" + dir)
}

// ========== 启动时恢复连接 ==========

// GetAutoRestoreConnections 返回标记了"自动重连"的连接列表（启动恢复用）
func (s *AppSettingService) GetAutoRestoreConnections() any {
	connectionDB := dao.Query().MqttConnection
	list, err := connectionDB.Where(connectionDB.Reconnect.Eq(1)).Find()
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	type item struct {
		ID       int64  `json:"id"`
		Title    string `json:"title"`
		ClientID string `json:"client_id"`
	}
	out := make([]item, 0, len(list))
	for _, c := range list {
		out = append(out, item{ID: int64(c.ID), Title: c.Title, ClientID: c.ClientID})
	}
	return gf.Success().SetData(out)
}

// ---------- 工具函数 ----------

// runtimeDir 运行时目录（程序目录下 runtime，不存在则创建）
func runtimeDir() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", err
	}
	dir := exeDir(exe) + string(os.PathSeparator) + "runtime"
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return dir, nil
}

// exeDir 可执行文件所在目录（去掉文件名）
func exeDir(exe string) string {
	// 简单取目录：去掉最后一个分隔符后的部分
	for i := len(exe) - 1; i >= 0; i-- {
		if exe[i] == '/' || exe[i] == '\\' {
			return exe[:i]
		}
	}
	return "."
}

// openFolder 用系统资源管理器打开目录
func openFolder(dir string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("explorer", dir).Start()
	case "darwin":
		return exec.Command("open", dir).Start()
	default:
		return exec.Command("xdg-open", dir).Start()
	}
}

