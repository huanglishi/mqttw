<template>
  <div>
    <a-table
      :data="devices"
      :loading="loading"
      :pagination="{ pageSize: 10, showTotal: true }"
      row-key="id"
      :columns="columns"
    >
      <template #protocol="{ record }">
        <a-tag :color="protocolColor(record.protocol)">
          {{ protocolLabel(record.protocol) }}
        </a-tag>
      </template>
      <template #endpoint="{ record }">
        <span v-if="record.protocol === 'rtu'">{{ record.serial_port || '-' }}</span>
        <span v-else>{{ record.host }}:{{ record.port }}</span>
      </template>
      <template #enabled="{ record }">
        <a-switch :model-value="record.enabled === 1" disabled size="small" />
      </template>
      <template #running="{ record }">
        <a-badge
          :status="isRunning(record) ? 'processing' : 'default'"
          :text="isRunning(record) ? '采集中' : '停止'"
        />
      </template>
      <template #ops="{ record }">
        <a-space>
          <a-button size="mini" @click="$emit('edit', record)">编辑</a-button>
          <a-button size="mini" type="outline" @click="$emit('test', record)">测试读取</a-button>
          <a-popconfirm content="删除设备会同时删除其全部点位，确认？" @ok="$emit('remove', record)">
            <a-button size="mini" status="danger">删除</a-button>
          </a-popconfirm>
        </a-space>
      </template>
      <template #empty>
        <a-empty description="暂无设备，点击右下角新增设备按钮创建" />
      </template>
    </a-table>

    <div style="margin-top: 12px; text-align: right">
      <a-button type="primary" @click="$emit('add')">
        <template #icon><icon-plus /></template>
        新增设备
      </a-button>
    </div>

    <!-- 网关功能说明（实验性） -->
    <a-collapse class="gw-faq" :bordered="false" style="margin-top: 12px" :default-active-key="['faq']">
      <a-collapse-item key="faq" header="网关功能说明（实验性）">
        <ul class="gw-faq-body">
          <li><b>协议</b>：支持 Modbus TCP、Modbus RTU（RS485/RS232）、RTU over TCP 三种接入方式；</li>
          <li><b>寄存器</b>：线圈(0x)、离散输入(1x)、保持寄存器(3x)、输入寄存器(4x)，按功能码 01/02/03/04 批量读取；</li>
          <li><b>数据解析</b>：uint16 / int16 / uint32 / int32 / float32 / BCD / bit，字节序与字序可配置，支持 scale / offset 换算（real = raw × scale + offset）；</li>
          <li><b>采集策略</b>：每台设备独立协程串行轮询，断线自动重连，支持失败重试、变化才上报（减少 MQTT 报文）；</li>
          <li><b>MQTT 输出</b>：复用 MQTTW 的 Broker 连接发布消息，可自定义主题前缀、QoS、保留消息与告警主题；</li>
          <li><b>注意</b>：本模块为实验性功能；采集任务随 MQTTW 进程退出而停止，本期暂不支持写寄存器指令。</li>
        </ul>
      </a-collapse-item>
    </a-collapse>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';

const props = defineProps<{
  devices: any[];
  loading: boolean;
  status?: any;
}>();

defineEmits<{
  (e: 'edit', device: any): void;
  (e: 'remove', device: any): void;
  (e: 'test', device: any): void;
  (e: 'add'): void;
}>();

const runningDevices = computed(() => {
  const map: Record<number, boolean> = {};
  const devs = props.status?.devices || [];
  devs.forEach((d: any) => {
    map[d.device_id] = d.running;
  });
  return map;
});

const isRunning = (record: any) => !!runningDevices.value[record.id];

const protocolLabel = (p: string) =>
  ({ tcp: 'Modbus TCP', rtu: 'Modbus RTU', rtu_over_tcp: 'RTU over TCP' }[p] || p);

const protocolColor = (p: string) =>
  ({ tcp: 'arcoblue', rtu: 'green', rtu_over_tcp: 'orange' }[p] || 'gray');

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '名称', dataIndex: 'name' },
  { title: '协议', slotName: 'protocol', width: 130 },
  { title: '地址', slotName: 'endpoint' },
  { title: '从站', dataIndex: 'slave_id', width: 70 },
  { title: '采集周期(ms)', dataIndex: 'poll_interval', width: 120 },
  { title: '主题前缀', dataIndex: 'topic_prefix', ellipsis: true },
  { title: '启用', slotName: 'enabled', width: 70 },
  { title: '运行', slotName: 'running', width: 90 },
  { title: '操作', slotName: 'ops', width: 220, fixed: 'right' },
];
</script>
<style lang="less" scoped>
:deep(.arco-collapse-item-content){
  padding-left: 0px;
}
</style>