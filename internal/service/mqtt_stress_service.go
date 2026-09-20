// MQTT 压测服务：Wails 绑定层，转发到 internal/stress 压测核心
package service

import (
	"os"
	"path/filepath"
	"sync"

	"gofly/internal/stress"
	"gofly/internal/utils/gf"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// MqttStressService MQTT 压测（独立压测模块，不依赖调试客户端）
type MqttStressService struct {
	mu     sync.Mutex
	runner *stress.Runner
}

// Start 启动压测任务
// param: 前端 StressConfig 完整 JSON
func (s *MqttStressService) Start(param string) any {
	s.mu.Lock()
	defer s.mu.Unlock()

	cfg, err := stress.ParseConfig(param)
	if err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	if s.runner != nil && s.runner.Status() != stress.StatusIdle {
		return gf.Failed().SetMsg("已有压测任务在运行，请先停止")
	}

	runner := stress.NewRunner(cfg)
	// 指标事件推送到前端（Wails3 全局 app 在运行期必然就绪）
	runner.SetEventEmitter(func(event, payload string) {
		application.Get().Event.Emit(event, payload)
	})
	if err := runner.Start(); err != nil {
		return gf.Failed().SetMsg(err.Error())
	}
	s.runner = runner
	return gf.Success().SetMsg("压测任务已启动").SetData(map[string]any{
		"client_count": cfg.ClientCount,
		"protocol":     cfg.Protocol,
	})
}

// Stop 停止压测任务
func (s *MqttStressService) Stop() any {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.runner == nil || s.runner.Status() == stress.StatusIdle {
		return gf.Failed().SetMsg("当前没有运行中的压测任务")
	}
	s.runner.Stop()
	return gf.Success().SetMsg("正在停止压测任务…")
}

// GetStatus 获取压测状态与实时指标快照
func (s *MqttStressService) GetStatus() any {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.runner == nil {
		return gf.Success().SetMsg("压测未启动").SetData(map[string]any{
			"status": "idle",
		})
	}
	return gf.Success().SetMsg("获取压测状态").SetData(s.runner.SnapshotMap())
}

// ExportResult 导出压测结果 CSV 到 runtime/stress/
func (s *MqttStressService) ExportResult() any {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.runner == nil {
		return gf.Failed().SetMsg("压测未启动，无结果可导出")
	}
	dir := filepath.Join("runtime", "stress")
	if abs, err := filepath.Abs(dir); err == nil {
		dir = abs
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return gf.Failed().SetMsg("创建导出目录失败: " + err.Error())
		}
	}
	path, err := s.runner.ExportCSV(dir)
	if err != nil {
		return gf.Failed().SetMsg("导出失败: " + err.Error())
	}
	return gf.Success().SetMsg("导出成功: " + path).SetData(map[string]any{
		"path": path,
	})
}
