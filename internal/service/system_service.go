package service

import (
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type SystemService struct{}

// 设置主题-预留后续api更新使用
func (g *SystemService) SetTheme(theme string) {
	// application.OpenFileDialog().PromptForMultipleSelection()
	fmt.Println("执行设置主题函数", application.Get().GetPID(), theme)
}
