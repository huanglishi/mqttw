<template>
  <div class="rt-panel">
    <!-- 设备状态卡片 -->
    <a-row :gutter="12" class="rt-cards">
      <a-col v-for="d in deviceStatus" :key="d.device_id" :span="8" :xs="24" :sm="12" :md="8">
        <a-card :bordered="false" class="rt-card">
          <div class="rt-card-head">
            <span class="rt-device-name">{{ d.name }}</span>
            <a-badge
              :status="d.running ? (d.connected ? 'processing' : 'warning') : 'default'"
              :text="d.running ? (d.connected ? '采集中' : '重连中') : '停止'"
            />
          </div>
          <a-descriptions :column="2" size="mini">
            <a-descriptions-item label="协议">
              {{ { tcp: 'TCP', rtu: 'RTU', rtu_over_tcp: 'RTU-TCP' }[d.protocol] || d.protocol }}
            </a-descriptions-item>
            <a-descriptions-item label="点位">
              {{ d.values_per_cycle }}
            </a-descriptions-item>
            <a-descriptions-item label="轮询">
              {{ d.poll_count }}
            </a-descriptions-item>
            <a-descriptions-item label="速率">
              {{ fmtRate(d.rate) }} 点/s
            </a-descriptions-item>
            <a-descriptions-item label="错误">
              <span :class="{ 'rt-err': d.error_count > 0 }">{{ d.error_count }}</span>
            </a-descriptions-item>
            <a-descriptions-item label="上次采集">
              {{ d.last_poll_at ? fmtTime(d.last_poll_at) : '-' }}
            </a-descriptions-item>
          </a-descriptions>
          <a-alert v-if="d.last_error" type="warning" :closable="false" style="margin-top: 8px">
            <template #title>{{ d.last_error }}</template>
          </a-alert>
        </a-card>
      </a-col>
      <a-col v-if="deviceStatus.length === 0" :span="24">
        <a-empty description="网关未运行或无设备状态" />
      </a-col>
    </a-row>

    <!-- 实时值表格 -->
    <a-card :bordered="false" title="实时数据" style="margin-top: 12px">
      <template #extra>
        <a-space>
          <a-select v-model="viewDeviceId" style="width: 220px" placeholder="选择设备">
            <a-option v-for="d in deviceStatus" :key="d.device_id" :value="d.device_id">{{ d.name }}</a-option>
          </a-select>
          <span class="rt-update-time">更新于 {{ lastUpdate ? fmtTime(lastUpdate) : '-' }}</span>
        </a-space>
      </template>
      <a-table
        :data="viewValues"
        :pagination="false"
        size="small"
        :columns="valueColumns"
        :loading="viewDeviceId !== null && viewValues.length === 0 && runningAny"
      >
        <template #empty>
          <a-empty :description="runningAny ? '等待设备上报数据…' : '网关未运行'" />
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue';
import { Events } from '@wailsio/runtime';

const props = defineProps<{ status: any }>();

const lastUpdate = ref(0);
const valuesMap = ref<Record<number, Record<string, number>>>({});
const viewDeviceId = ref<number | null>(null);

const deviceStatus = computed(() => {
  const devs = props.status?.devices || [];
  if (viewDeviceId.value === null && devs.length > 0) {
    viewDeviceId.value = devs[0].device_id;
  }
  return devs;
});

const runningAny = computed(() => (deviceStatus.value || []).some((d: any) => d.running));

const viewValues = computed(() => {
  if (viewDeviceId.value === null) return [];
  const m = valuesMap.value[viewDeviceId.value];
  if (!m) return [];
  return Object.entries(m).map(([name, value]) => ({ name, value }));
});

const valueColumns = [
  { title: '点位', dataIndex: 'name' },
  { title: '当前值', slotName: 'value' },
];

// 后端推送的实时值 → 更新对应设备
const onValues = (ev: any) => {
  if (!ev?.data) return;
  try {
    const payload = typeof ev.data === 'string' ? JSON.parse(ev.data) : ev.data;
    if (payload && payload.device_id != null) {
      valuesMap.value = {
        ...valuesMap.value,
        [payload.device_id]: payload.values || {},
      };
      lastUpdate.value = Date.now();
    }
  } catch (e) {
    // 忽略
  }
};

watch(deviceStatus, (devs) => {
  if (viewDeviceId.value === null && devs.length > 0) {
    viewDeviceId.value = devs[0].device_id;
  }
});

const fmtRate = (r: number) => (r ? r.toFixed(1) : '0');
const fmtTime = (t: number) => {
  const d = new Date(t * 1000);
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
};

onMounted(() => {
  Events.On('gateway:values', onValues);
});

onUnmounted(() => {
  Events.Off('gateway:values');
});
</script>

<style scoped lang="less">
.rt-panel {
  .rt-cards {
    .rt-card {
      margin-bottom: 4px;
      .rt-card-head {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 8px;
        .rt-device-name {
          font-weight: 600;
        }
      }
      .rt-err {
        color: rgb(var(--danger-6));
        font-weight: 600;
      }
    }
  }
  .rt-update-time {
    font-size: 12px;
    color: var(--color-text-3);
  }
}
</style>
