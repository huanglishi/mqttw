<template>
  <div class="stress-page">
    <!-- 左：配置面板 + 控制 -->
    <div class="stress-left">
      <a-card :bordered="false" class="config-card" title="压测配置">
        <StressConfig v-model="config" />
        <a-space direction="vertical" fill class="stress-actions">
          <a-button
            type="primary"
            long
            :loading="starting"
            :disabled="stressStore.running"
            @click="handleStart"
          >
            启动压测
          </a-button>
          <a-button
            status="danger"
            long
            :disabled="!stressStore.running"
            @click="handleStop"
          >
            停止压测
          </a-button>
        </a-space>
      </a-card>
    </div>

    <!-- 右：仪表盘 -->
    <div class="stress-right">
      <MetricsCards :config="config" />
      <MetricsChart />
      <ResultTable />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, onMounted, onUnmounted } from 'vue';
import { Events } from '@wailsio/runtime';
import { Message } from '@arco-design/web-vue';
import StressConfig, { StressConfig as StressConfigType } from './components/StressConfig.vue';
import MetricsCards from './components/MetricsCards.vue';
import MetricsChart from './components/MetricsChart.vue';
import ResultTable from './components/ResultTable.vue';
import { useStressStore } from '@/store';
import { MqttStressService } from '/#/gofly/internal/service';

const stressStore = useStressStore();
const starting = ref(false);

// ========== 配置持久化（localStorage） ==========
const STORAGE_KEY = 'mqttw_stress_config';

// 默认配置（与后端 normalize 默认值一致）
const defaultConfig = (): StressConfigType => ({
  broker: 'tcp://127.0.0.1:1883',
  protocol: '3.1',
  username: '',
  password: '',
  clientCount: 1000,
  clientPrefix: 'stress_test',
  connectRate: 500,
  connectTimeout: 5,
  keepalive: 30,
  cleanSession: true,
  topics: ['stress/#'],
  subscribePerClient: 1,
  qos: 0,
  publishEnabled: false,
  publishWorkers: 8,
  publishTopic: 'stress/pub',
  publishPayloadSize: 128,
  publishQos: 0,
  publishRate: 5000,
  publishTotal: 0,
  duration: 60,
  echoEnabled: false,
  publishEchoTopic: '',
});

// 从 localStorage 恢复，无则用默认值
const loadConfig = (): StressConfigType => {
  try {
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved) {
      return { ...defaultConfig(), ...JSON.parse(saved) };
    }
  } catch (e) {
    // 解析失败忽略，用默认值
  }
  return defaultConfig();
};

const config = ref<StressConfigType>(loadConfig());

// 配置变化 → 自动持久化
watch(
  config,
  () => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(config.value));
  },
  { deep: true }
);

// 启动压测
const handleStart = async () => {
  starting.value = true;
  try {
    stressStore.reset();
    const payload = JSON.stringify({ ...config.value, topics: config.value.topics.filter(Boolean) });
    const res = await MqttStressService.Start(payload);
    if (res.code == 0) {
      Message.success(res.message);
    } else {
      Message.error(res.message);
    }
  } catch (e) {
    Message.error('启动失败：' + e);
  } finally {
    starting.value = false;
  }
};

// 停止压测
const handleStop = async () => {
  const res = await MqttStressService.Stop();
  if (res.code == 0) {
    Message.success(res.message);
  } else {
    Message.warning(res.message);
  }
};

// 后端每 500ms 推送 stress:metrics 事件
const handleMetricsEvent = (ev: any) => {
  if (!ev?.data) return;
  try {
    stressStore.handleMetrics(JSON.parse(ev.data));
  } catch (e) {
    // 忽略解析失败
  }
};

onMounted(async () => {
  Events.On('stress:metrics', handleMetricsEvent);
  // 页面打开时恢复运行中的任务快照
  const res = await MqttStressService.GetStatus();
  if (res.code == 0 && res.data && res.data.status !== 'idle') {
    stressStore.handleMetrics(res.data);
  }
});

onUnmounted(() => {
  Events.Off('stress:metrics');
});
</script>

<style scoped lang="less">
.stress-page {
  display: flex;
  gap: 12px;
  padding: 0px;
  height: 100vh;
  width: 100%;
  box-sizing: border-box;
  .stress-left {
    flex: 0 0 360px;
    overflow-y: auto;

    .config-card {
      border-radius: 6px;
      background: var(--color-bg-2);

      :deep(.arco-card-header) {
        font-size: 15px;
        font-weight: 600;
        padding: 14px 16px;
        border: none;
      }

      .stress-actions {
        margin-top: 4px;
      }
    }
  }

  .stress-right {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 10px;
    overflow-y: auto;
  }

  @media (max-width: 900px) {
    flex-direction: column;
    .stress-left {
      flex: none;
    }
  }
}
</style>
