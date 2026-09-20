<template>
  <div class="metrics-chart">
    <a-card :bordered="false" class="chart-card" title="实时曲线（连接数 / 发布速率 / 接收速率）">
      <Chart height="260px" :option="lineOption" />
    </a-card>
    <a-card :bordered="false" class="chart-card" title="连接结果分布">
      <Chart height="220px" :option="pieOption" />
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import useChartOption from '@/hooks/chart-option';
import { useStressStore } from '@/store';

const stressStore = useStressStore();

// 时间序列 → echarts 数据集
const seriesData = computed(() => stressStore.metrics?.series || []);
const times = computed(() => seriesData.value.map((p) => {
  const d = new Date(p.ts);
  return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}:${String(d.getSeconds()).padStart(2, '0')}`;
}));
const connectedData = computed(() => seriesData.value.map((p) => p.connected));
const pubRateData = computed(() => seriesData.value.map((p) => p.pub_rate));
const recvRateData = computed(() => seriesData.value.map((p) => p.recv_rate));

// 折线：左轴=连接数，右轴=速率
const { chartOption: lineOption } = useChartOption((isDark) => {
  const axisColor = isDark ? 'rgba(255,255,255,.65)' : '#4E5969';
  const splitColor = isDark ? 'rgba(255,255,255,.08)' : '#E5E8EF';
  return {
    tooltip: { trigger: 'axis' },
    legend: { data: ['连接数', '发布速率', '接收速率'], top: 0, icon: 'circle', itemWidth: 8 },
    grid: { left: 48, right: 56, top: 34, bottom: 28 },
    xAxis: {
      type: 'category',
      data: times.value,
      boundaryGap: false,
      axisLabel: { color: axisColor, fontSize: 11 },
      axisLine: { show: false },
      axisTick: { show: false },
      splitLine: { show: true, lineStyle: { color: splitColor } },
    },
    yAxis: [
      {
        type: 'value',
        name: '连接数',
        axisLabel: { color: axisColor, fontSize: 11 },
        splitLine: { show: true, lineStyle: { type: 'dashed', color: splitColor } },
      },
      {
        type: 'value',
        name: 'msg/s',
        axisLabel: { color: axisColor, fontSize: 11 },
        splitLine: { show: false },
      },
    ],
    series: [
      {
        name: '连接数',
        type: 'line',
        data: connectedData.value,
        smooth: true,
        showSymbol: false,
        yAxisIndex: 0,
        lineStyle: { width: 2, color: '#165DFF' },
        areaStyle: { opacity: 0.12, color: '#165DFF' },
      },
      {
        name: '发布速率',
        type: 'line',
        data: pubRateData.value,
        smooth: true,
        showSymbol: false,
        yAxisIndex: 1,
        lineStyle: { width: 2, color: '#00B42A' },
      },
      {
        name: '接收速率',
        type: 'line',
        data: recvRateData.value,
        smooth: true,
        showSymbol: false,
        yAxisIndex: 1,
        lineStyle: { width: 2, color: '#FF7D00' },
      },
    ],
  };
});

// 环形：连接成功/失败
const pieData = computed(() => {
  const m = stressStore.metrics;
  return [
    { value: m?.connect_success || 0, name: '成功' },
    { value: m?.connect_fail || 0, name: '失败' },
  ];
});
const { chartOption: pieOption } = useChartOption((isDark) => {
  const total = pieData.value.reduce((s, i) => s + i.value, 0);
  const labelColor = isDark ? 'rgba(255,255,255,.7)' : '#4E5969';
  return {
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    legend: { bottom: 0, icon: 'circle', itemWidth: 8, textStyle: { color: labelColor } },
    graphic: {
      elements: [
        {
          type: 'text',
          left: 'center',
          top: '38%',
          style: { text: '连接总数', textAlign: 'center', fill: labelColor, fontSize: 13 },
        },
        {
          type: 'text',
          left: 'center',
          top: '50%',
          style: {
            text: String(total),
            textAlign: 'center',
            fill: isDark ? '#fff' : '#1D2129',
            fontSize: 18,
            fontWeight: 600,
          },
        },
      ],
    },
    series: [
      {
        type: 'pie',
        radius: ['48%', '68%'],
        center: ['50%', '50%'],
        label: { formatter: '{d}%', fontSize: 12, color: labelColor },
        itemStyle: { borderColor: isDark ? '#232324' : '#fff', borderWidth: 1 },
        color: ['#00B42A', '#F53F3F'],
        data: pieData.value,
      },
    ],
  };
});
</script>

<style scoped lang="less">
.metrics-chart {
  display: grid;
  grid-template-columns: 1fr 320px;
  gap: 10px;

  .chart-card {
    border-radius: 6px;
    background: var(--color-bg-2);

    :deep(.arco-card-header) {
      font-size: 14px;
      padding: 14px 16px 0;
      border: none;
    }
  }

  @media (max-width: 1100px) {
    grid-template-columns: 1fr;
  }
}
</style>
