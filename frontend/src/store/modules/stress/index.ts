import { defineStore } from 'pinia';
import { StressMetrics } from './types';

const useStressStore = defineStore('stress', {
  state: () => ({
    running: false,          // 是否有压测任务在运行
    metrics: null as StressMetrics | null, // 最新指标快照
  }),
  getters: {
    // 连接成功率（0-100）
    connectRate: (state) => {
      if (!state.metrics || state.metrics.connect_total === 0) return 0;
      const rate = (state.metrics.connect_success / state.metrics.connect_total) * 100;
      return Number(Math.min(100, Math.max(0, rate)).toFixed(1));
    },
    // 当前状态文案
    statusText: (state) => state.metrics?.status || 'idle',
  },
  actions: {
    // 后端 stress:metrics 事件回调 / GetStatus 拉取
    handleMetrics(data: StressMetrics) {
      // GetStatus 恢复场景：start_at_ms 合并为 start_at
      if (typeof data.start_at_ms === 'number' && data.start_at === undefined) {
        data.start_at = data.start_at_ms;
      }
      this.metrics = data;
      this.running = data.status === 'running' || data.status === 'stopping';
    },
    // 压测启动时清空旧数据
    reset() {
      this.metrics = null;
      this.running = false;
    },
  },
});

export default useStressStore;
