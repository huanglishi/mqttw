<template>
  <a-card :bordered="false" title="采集日志" style="margin-top: 12px">
    <template #extra>
      <a-space>
        <a-select v-model="levelFilter" size="small" style="width: 110px" allow-clear placeholder="全部级别">
          <a-option value="info">info</a-option>
          <a-option value="warn">warn</a-option>
          <a-option value="error">error</a-option>
          <a-option value="alarm">alarm</a-option>
        </a-select>
        <a-button size="small" @click="logs = []">清空</a-button>
      </a-space>
    </template>
    <div ref="logBox" class="log-box">
      <div v-for="(item, i) in filteredLogs" :key="i" class="log-line" :class="'log-' + item.level">
        <span class="log-time">{{ fmtTime(item.ts) }}</span>
        <a-tag size="mini" :color="tagColor(item.level)">{{ item.level }}</a-tag>
        <span v-if="item.device" class="log-device">{{ item.device }}</span>
        <span class="log-msg">{{ item.msg }}</span>
      </div>
      <a-empty v-if="filteredLogs.length === 0" description="暂无日志" :style="{ padding: '20px 0' }" />
    </div>
  </a-card>
</template>

<script lang="ts" setup>
import { ref, computed, nextTick, onMounted, onUnmounted } from 'vue';
import { Events } from '@wailsio/runtime';

interface LogItem {
  level: string;
  msg: string;
  device?: string;
  ts: number;
}

const logs = ref<LogItem[]>([]);
const levelFilter = ref<string>();
const logBox = ref<HTMLDivElement>();

const filteredLogs = computed(() => {
  if (!levelFilter.value) return logs.value;
  return logs.value.filter((l) => l.level === levelFilter.value);
});

const MAX_LOGS = 500;

const onLog = (ev: any) => {
  if (!ev?.data) return;
  try {
    const payload = typeof ev.data === 'string' ? JSON.parse(ev.data) : ev.data;
    if (!payload || !payload.msg) return;
    logs.value.push({
      level: payload.level || 'info',
      msg: payload.msg,
      device: payload.device || '',
      ts: payload.ts || Math.floor(Date.now() / 1000),
    });
    if (logs.value.length > MAX_LOGS) {
      logs.value.splice(0, logs.value.length - MAX_LOGS);
    }
    scrollToBottom();
  } catch (e) {
    // 忽略
  }
};

const scrollToBottom = async () => {
  await nextTick();
  if (logBox.value) {
    logBox.value.scrollTop = logBox.value.scrollHeight;
  }
};

const fmtTime = (t: number) => {
  const d = new Date(t * 1000);
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${pad(d.getHours())}:${pad(d.getMinutes())}:${pad(d.getSeconds())}`;
};

const tagColor = (level: string) =>
  ({ info: 'arcoblue', warn: 'orange', error: 'red', alarm: 'gold' }[level] || 'gray');

onMounted(() => {
  Events.On('gateway:log', onLog);
});

onUnmounted(() => {
  Events.Off('gateway:log');
});
</script>

<style scoped lang="less">
.log-box {
  height: 320px;
  overflow-y: auto;
  background: var(--color-fill-1);
  border-radius: 4px;
  padding: 8px 12px;
  font-size: 12px;

  .log-line {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 2px 0;
    line-height: 1.6;
    border-bottom: 1px dashed var(--color-fill-3);

    .log-time {
      color: var(--color-text-3);
      flex-shrink: 0;
      font-variant-numeric: tabular-nums;
    }
    .log-device {
      color: rgb(var(--primary-6));
      flex-shrink: 0;
    }
    .log-msg {
      word-break: break-all;
    }
    &.log-error .log-msg,
    &.log-alarm .log-msg {
      color: rgb(var(--danger-6));
    }
    &.log-warn .log-msg {
      color: rgb(var(--warning-6));
    }
  }
}
</style>
