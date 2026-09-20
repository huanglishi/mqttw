<template>
  <a-card :bordered="false" class="result-card" title="结果汇总">
    <template #extra>
      <a-space>
        <a-tag v-if="metrics" color="arcoblue">{{ metrics.status }}</a-tag>
      </a-space>
    </template>
    <a-table
      v-if="metrics"
      :data="tableData"
      :pagination="false"
      :bordered="false"
      size="small"
    >
      <template #columns>
        <a-table-column title="指标" data-index="label" :width="180" />
        <a-table-column title="数值" data-index="value" />
      </template>
    </a-table>
    <a-empty v-else description="压测尚未运行，启动后展示结果" />
    <a-divider v-if="failReasons.length > 0" style="margin: 12px 0" />
    <div v-if="failReasons.length > 0" class="fail-reasons">
      <div class="fail-title">失败原因分布</div>
      <a-space wrap>
        <a-tag v-for="(item, idx) in failReasons.slice(0, 10)" :key="idx" color="red">
          {{ item[0] }}：{{ item[1] }}
        </a-tag>
      </a-space>
    </div>
  </a-card>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import { useStressStore } from '@/store';

const stressStore = useStressStore();
const metrics = computed(() => stressStore.metrics);

const tableData = computed(() => {
  const m = metrics.value;
  if (!m) return [];
  return [
    { label: '连接总数（尝试）', value: m.connect_total },
    { label: '连接成功', value: m.connect_success },
    { label: '连接失败', value: m.connect_fail },
    { label: '连接平均耗时', value: m.connect_avg_ms + ' ms' },
    { label: '连接最大耗时', value: m.connect_max_ms + ' ms' },
    { label: '订阅总数', value: m.subscribe_total },
    { label: '订阅成功', value: m.subscribe_success },
    { label: '订阅失败', value: m.subscribe_fail },
    { label: '发布总数', value: m.publish_total },
    { label: '发布成功', value: m.publish_success },
    { label: '发布失败', value: m.publish_fail },
    { label: '接收总数', value: m.received_total },
    { label: '回显次数', value: m.echo_count },
    { label: '回显平均延迟', value: m.echo_avg_ms + ' ms' },
    { label: '回显最大延迟', value: m.echo_max_ms + ' ms' },
    { label: '回显 P50 延迟', value: m.echo_p50_ms + ' ms' },
    { label: '回显 P95 延迟', value: m.echo_p95_ms + ' ms' },
    { label: '回显 P99 延迟', value: m.echo_p99_ms + ' ms' },
  ];
});

// 失败原因 Top N（表格下方展示）
const failReasons = computed(() => metrics.value?.fail_reasons || []);
</script>

<style scoped lang="less">
.result-card {
  border-radius: 6px;
  background: var(--color-bg-2);
  :deep(.arco-card-header) {
    font-size: 14px;
    padding: 14px 16px 0;
    border: none;
  }
}
.fail-reasons {
  .fail-title {
    font-size: 13px;
    color: var(--color-text-2);
    margin-bottom: 8px;
  }
}
</style>
