<template>
  <div class="metrics-cards">
    <a-card v-if="metrics" :bordered="false" class="metric-card status-card">
      <div class="status-title">运行状态</div>
      <div class="status-value">
        <a-tag :color="statusColor" size="large">{{ statusText }}</a-tag>
        <span v-if="stressStore.running" class="elapsed-text">{{ formatElapsed(elapsed) }}</span>
      </div>
    </a-card>
    <a-card :bordered="false" class="metric-card">
      <a-statistic
        title="已连接"
        :value="metrics ? metrics.connected : 0"
        :suffix="'/ ' + (config.clientCount || 0)"
        :precision="0"
      />
    </a-card>
    <a-card :bordered="false" class="metric-card">
      <a-statistic title="连接成功率" :value="connectRate" suffix="%" :precision="1" />
    </a-card>
    <a-card :bordered="false" class="metric-card">
      <a-statistic title="发布速率" :value="pubRate" suffix="msg/s" :precision="0" />
    </a-card>
    <a-card :bordered="false" class="metric-card">
      <a-statistic title="接收速率" :value="recvRate" suffix="msg/s" :precision="0" />
    </a-card>
    <a-card :bordered="false" class="metric-card">
      <a-statistic title="发布总数" :value="metrics?.publish_total || 0" :precision="0" />
    </a-card>
    <a-card :bordered="false" class="metric-card">
      <a-statistic title="接收总数" :value="metrics?.received_total || 0" :precision="0" />
    </a-card>
    <a-card v-if="metrics?.echo_count > 0" :bordered="false" class="metric-card">
      <a-statistic title="平均延迟" :value="metrics.echo_avg_ms" suffix="ms" :precision="1" />
    </a-card>
    <a-card v-if="metrics?.echo_count > 0" :bordered="false" class="metric-card">
      <a-statistic title="P95 延迟" :value="metrics.echo_p95_ms" suffix="ms" :precision="0" />
    </a-card>
  </div>

  <!-- 进行中进度 -->
  <div v-if="stressStore.running || hasProgress" class="progress-cards">
    <a-card v-if="config.clientCount > 0" :bordered="false" class="progress-card">
      <div class="progress-label">连接进度：{{ metrics?.connected || 0 }} / {{ config.clientCount }}</div>
      <a-progress :percent="connectPercent" :stroke-width="8" />
    </a-card>
    <a-card v-if="config.duration > 0" :bordered="false" class="progress-card duration-card">
      <div class="progress-label">运行时长</div>
      <div class="duration-value">
        <span class="duration-now">{{ formatElapsed(elapsed) }}</span>
        <span class="duration-total">/ {{ config.duration }} 秒</span>
      </div>
    </a-card>
    <a-card
      v-if="config.publishEnabled && config.publishTotal > 0"
      :bordered="false"
      class="progress-card"
    >
      <div class="progress-label">
        发布进度：{{ metrics?.publish_total || 0 }} / {{ config.publishTotal }}
      </div>
      <a-progress :percent="publishPercent" :stroke-width="8" color="#00B42A" />
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, watch, onUnmounted } from 'vue';
import { useStressStore } from '@/store';
import type { StressConfig } from './StressConfig.vue';

const props = defineProps<{ config: StressConfig }>();
const stressStore = useStressStore();

const metrics = computed(() => stressStore.metrics);
const statusText = computed(() => {
  const map: Record<string, string> = { running: '运行中', stopping: '停止中', idle: '空闲' };
  return map[stressStore.statusText] || stressStore.statusText;
});
const statusColor = computed(() => {
  const map: Record<string, string> = { running: 'arcoblue', stopping: 'orange', idle: 'gray' };
  return map[stressStore.statusText] || 'gray';
});
const connectRate = computed(() => stressStore.connectRate);
const pubRate = computed(() => {
  const s = metrics.value?.series;
  return s && s.length ? s[s.length - 1].pub_rate : 0;
});
const recvRate = computed(() => {
  const s = metrics.value?.series;
  return s && s.length ? s[s.length - 1].recv_rate : 0;
});

// ========== 运行计时（前端本地计时，秒） ==========
const startAt = ref(0);
const elapsed = ref(0);
let timer: number | null = null;

watch(
  () => stressStore.running,
  (running) => {
    if (running) {
      // 优先用后端推送的任务真实开始时间；缺失时回退当前时刻
      startAt.value = Number(metrics.value?.start_at) || Date.now();
      elapsed.value = 0;
      timer = window.setInterval(() => {
        elapsed.value = Math.floor((Date.now() - startAt.value) / 1000);
      }, 1000);
    } else if (timer) {
      window.clearInterval(timer);
      timer = null;
    }
  },
  { immediate: true }
);

onUnmounted(() => {
  if (timer) window.clearInterval(timer);
});

// ========== 进度计算 ==========
const connectPercent = computed(() => {
  const n = metrics.value?.connected || 0;
  if (!props.config.clientCount) return 0;
  return (n / props.config.clientCount).toFixed(2);
});
const publishPercent = computed(() => {
  const total = props.config.publishTotal;
  if (!total) return 0;
  const done = metrics.value?.publish_total || 0;
  return Math.min(100, Math.round((done / total) * 100));
});
const hasProgress = computed(() => !!metrics.value); // 有快照即显示（含结束后最终值）

// 秒 → mm:ss / hh:mm:ss
const formatElapsed = (sec: number) => {
  if (!Number.isFinite(sec)) sec = 0;
  const s = Math.max(0, sec);
  const h = Math.floor(s / 3600);
  const m = Math.floor((s % 3600) / 60);
  const ss = s % 60;
  const pad = (n: number) => String(n).padStart(2, '0');
  return h > 0 ? `${pad(h)}:${pad(m)}:${pad(ss)}` : `${pad(m)}:${pad(ss)}`;
};
</script>

<style scoped lang="less">
.metrics-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 10px;

  .metric-card {
    border-radius: 6px;
    background: var(--color-bg-2);

    :deep(.arco-statistic-title) {
      font-size: 13px;
      color: var(--color-text-3);
    }
  }

  .status-card {
    .status-title {
      font-size: 13px;
      color: var(--color-text-3);
      margin-bottom: 4px;
    }
    .status-value {
      display: flex;
      align-items: center;
      gap: 8px;
      min-height: 32px;
    }
    .elapsed-text {
      font-size: 14px;
      font-weight: 600;
      color: var(--color-text-1);
      font-variant-numeric: tabular-nums;
    }
  }
}

.progress-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 10px;

  .progress-card {
    border-radius: 6px;
    background: var(--color-bg-2);

    .progress-label {
      font-size: 13px;
      color: var(--color-text-2);
      margin-bottom: 8px;
      font-variant-numeric: tabular-nums;
    }
  }

    .duration-value {
      display: flex;
      align-items: baseline;
      gap: 6px;

      .duration-now {
        font-size: 24px;
        font-weight: 700;
        font-variant-numeric: tabular-nums;
        color: var(--color-text-1);
      }

      .duration-total {
        font-size: 13px;
        color: var(--color-text-3);
      }
    }
}
</style>
