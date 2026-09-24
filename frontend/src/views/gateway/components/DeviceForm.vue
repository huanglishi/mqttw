<template>
  <a-drawer
    :visible="visible"
    :width="560"
    :title="form.id ? '编辑设备' : '新增设备'"
    :footer="false"
    unmount-on-close
    @cancel="$emit('update:visible', false)"
    @before-ok="handleSubmit"
  >
    <a-form :model="form" layout="vertical">
      <a-form-item label="设备名称" field="name" :rules="[{ required: true, message: '请输入设备名称' }]">
        <a-input v-model="form.name" placeholder="如 plc01 / 温控仪表" allow-clear />
      </a-form-item>

      <a-form-item label="连接协议" field="protocol" :rules="[{ required: true }]">
        <a-select v-model="form.protocol">
          <a-option value="tcp">Modbus TCP（以太网）</a-option>
          <a-option value="rtu">Modbus RTU（串口 RS485/RS232）</a-option>
          <a-option value="rtu_over_tcp">RTU over TCP（网关透传）</a-option>
        </a-select>
      </a-form-item>

      <template v-if="form.protocol === 'tcp' || form.protocol === 'rtu_over_tcp'">
        <a-row :gutter="12">
          <a-col :span="16">
            <a-form-item label="主机地址" field="host" :rules="[{ required: true, message: '请输入 IP/域名' }]">
              <a-input v-model="form.host" placeholder="192.168.1.10" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="端口" field="port">
              <a-input-number v-model="form.port" :min="1" :max="65535" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
      </template>

      <template v-else>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="串口" field="serialPort" :rules="[{ required: true, message: '请输入串口名' }]">
              <a-input v-model="form.serialPort" placeholder="COM3" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="波特率" field="baudRate">
              <a-select v-model="form.baudRate">
                <a-option v-for="b in [4800, 9600, 19200, 38400, 57600, 115200]" :key="b" :value="b">{{ b }}</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="12">
          <a-col :span="8">
            <a-form-item label="数据位" field="dataBits">
              <a-select v-model="form.dataBits">
                <a-option :value="7">7</a-option>
                <a-option :value="8">8</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="停止位" field="stopBits">
              <a-select v-model="form.stopBits">
                <a-option :value="1">1</a-option>
                <a-option :value="2">2</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="校验位" field="parity">
              <a-select v-model="form.parity">
                <a-option value="none">无(N)</a-option>
                <a-option value="even">偶(E)</a-option>
                <a-option value="odd">奇(O)</a-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
      </template>

      <a-row :gutter="12">
        <a-col :span="12">
          <a-form-item label="从站 ID" field="slaveId">
            <a-input-number v-model="form.slaveId" :min="1" :max="247" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="采集周期 (ms)" field="pollInterval">
            <a-input-number v-model="form.pollInterval" :min="100" :step="100" style="width: 100%" />
          </a-form-item>
        </a-col>
      </a-row>
      <a-row :gutter="12">
        <a-col :span="12">
          <a-form-item label="请求超时 (ms)" field="timeout">
            <a-input-number v-model="form.timeout" :min="100" :step="100" style="width: 100%" />
          </a-form-item>
        </a-col>
        <a-col :span="12">
          <a-form-item label="失败重试次数" field="retryCount">
            <a-input-number v-model="form.retryCount" :min="0" :max="10" style="width: 100%" />
          </a-form-item>
        </a-col>
      </a-row>

      <a-form-item label="寄存器字节序 / 字序（32 位浮点、双寄存器关键）">
        <a-space>
          <a-select v-model="form.byteOrder" style="width: 180px">
            <a-option value="big">高字节在前 (Big)</a-option>
            <a-option value="little">低字节在前 (Little)</a-option>
          </a-select>
          <a-select v-model="form.wordOrder" style="width: 200px">
            <a-option value="high_first">首寄存器为高字</a-option>
            <a-option value="low_first">首寄存器为低字</a-option>
          </a-select>
        </a-space>
      </a-form-item>

      <a-divider orientation="left">MQTT 输出</a-divider>

      <a-form-item label="MQTT Broker 连接" field="brokerConnectionId">
        <a-select v-model="form.brokerConnectionId" allow-clear placeholder="选择「连接」页已配置的 Broker（不选则不发布）">
          <a-option
            v-for="b in brokerOptions"
            :key="b.id"
            :value="b.id"
          >
            {{ b.title }}（{{ b.host }}:{{ b.port }}）
          </a-option>
        </a-select>
      </a-form-item>

      <a-row :gutter="12">
        <a-col :span="16">
          <a-form-item label="主题前缀" field="topicPrefix">
            <a-input v-model="form.topicPrefix" placeholder="留空默认 modbus/设备名" />
          </a-form-item>
        </a-col>
        <a-col :span="8">
          <a-form-item label="QoS" field="qos">
            <a-select v-model="form.qos">
              <a-option :value="0">0</a-option>
              <a-option :value="1">1</a-option>
              <a-option :value="2">2</a-option>
            </a-select>
          </a-form-item>
        </a-col>
      </a-row>
      <a-row :gutter="12">
        <a-col :span="16">
          <a-form-item label="告警主题（离线/采集失败）" field="alarmTopic">
            <a-input v-model="form.alarmTopic" placeholder="如 modbus/alarm，留空不推送告警" />
          </a-form-item>
        </a-col>
        <a-col :span="8">
          <a-form-item label="保留消息" field="retain">
            <a-switch v-model="form.retain" :checked-value="1" :unchecked-value="0" />
          </a-form-item>
        </a-col>
      </a-row>

      <a-form-item label="启用采集" field="enabled">
        <a-switch v-model="form.enabled" :checked-value="1" :unchecked-value="0" />
      </a-form-item>
    </a-form>

    <template #footer>
      <a-space style="float: right">
        <a-button @click="$emit('update:visible', false)">取消</a-button>
        <a-button type="primary" :loading="saving" @click="handleSubmit">保存</a-button>
      </a-space>
    </template>
  </a-drawer>
</template>

<script lang="ts" setup>
import { reactive, ref, watch } from 'vue';
import { Message } from '@arco-design/web-vue';

const props = defineProps<{
  visible: boolean;
  device: any;
  brokerOptions: any[];
}>();

const emit = defineEmits<{
  (e: 'update:visible', v: boolean): void;
  (e: 'submit', form: any): void;
}>();

const defaultForm = () => ({
  id: 0,
  name: '',
  protocol: 'tcp',
  host: '',
  port: 502,
  slaveId: 1,
  serialPort: '',
  baudRate: 9600,
  dataBits: 8,
  stopBits: 1,
  parity: 'none',
  pollInterval: 1000,
  timeout: 1000,
  retryCount: 3,
  byteOrder: 'big',
  wordOrder: 'high_first',
  brokerConnectionId: 0,
  topicPrefix: '',
  qos: 0,
  retain: 0,
  alarmTopic: '',
  enabled: 1,
});

const form = reactive<any>(defaultForm());
const saving = ref(false);

watch(
  () => props.visible,
  (v) => {
    if (!v) return;
    Object.assign(form, defaultForm(), props.device || {});
  },
  { immediate: true }
);

const handleSubmit = async () => {
  saving.value = true;
  try {
    emit('submit', { ...form });
  } finally {
    saving.value = false;
  }
};
</script>
