// 压测指标类型（与后端 internal/stress/metrics.go rawSnapshot 对应）
export interface StressSeriesPoint {
  ts: number;         // unix 毫秒
  connected: number;  // 当前已连接数
  pub_rate: number;   // 发布速率 msg/s
  recv_rate: number;  // 接收速率 msg/s
  echo_avg_ms: number; // 回显平均延迟 ms
}

export interface StressMetrics {
  status: string;
  connected: number;
  connect_total: number;
  connect_success: number;
  connect_fail: number;
  connect_avg_ms: number;
  connect_max_ms: number;
  subscribe_total: number;
  subscribe_success: number;
  subscribe_fail: number;
  publish_total: number;
  publish_success: number;
  publish_fail: number;
  received_total: number;
  echo_count: number;
  echo_avg_ms: number;
  echo_max_ms: number;
  echo_p50_ms: number;
  echo_p95_ms: number;
  echo_p99_ms: number;
  fail_reasons: [string, number][];
  start_at?: number; // 任务开始时间 unix 毫秒（事件推送）
  start_at_ms?: number; // 仅 GetStatus 恢复场景
  duration_sec?: number; // 计划时长（秒），0=不限
  series: StressSeriesPoint[];
  config?: any;
}
