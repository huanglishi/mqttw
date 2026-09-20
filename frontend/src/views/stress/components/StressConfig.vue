<template>
  <a-form :model="config" layout="vertical" class="stress-config">
    <a-collapse :default-active-key="['conn','client']" :bordered="false">
      <!-- ① 连接配置 -->
      <a-collapse-item key="conn" header="① 连接配置">
        <a-form-item label="协议版本">
          <a-radio-group v-model="config.protocol" type="button">
            <a-radio value="3.1">MQTT 3.1</a-radio>
            <a-radio value="3.1.1">MQTT 3.1.1</a-radio>
            <a-radio value="5.0">MQTT 5.0</a-radio>
            <a-radio value="auto">自动协商</a-radio>
          </a-radio-group>
        </a-form-item>
        <a-form-item label="Broker 地址" required>
          <a-input v-model="config.broker" placeholder="tcp://127.0.0.1:1883" />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="用户名">
              <a-input v-model="config.username" placeholder="可选" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="密码">
              <a-input-password v-model="config.password" placeholder="可选" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="KeepAlive（秒）">
              <a-input-number v-model="config.keepalive" :min="0" :max="65535" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="连接超时（秒）">
              <a-input-number v-model="config.connectTimeout" :min="1" :max="60" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="Clean Session（干净会话）">
          <a-switch v-model="config.cleanSession" />
        </a-form-item>
      </a-collapse-item>

      <!-- ② 客户端规模 -->
      <a-collapse-item key="client" header="② 客户端规模">
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="客户端数量">
              <a-input-number v-model="config.clientCount" :min="1" :max="100000" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="连接速率（个/秒，0=不限）">
              <a-input-number v-model="config.connectRate" :min="0" :max="10000" />
            </a-form-item>
          </a-col>
        </a-row>
        <a-form-item label="ClientID 前缀">
          <a-input v-model="config.clientPrefix" placeholder="stress_test" />
        </a-form-item>
      </a-collapse-item>

      <!-- ③ 订阅配置 -->
      <a-collapse-item key="sub" header="③ 订阅配置">
        <a-form-item label="订阅主题（每行一个，支持通配符）">
          <a-textarea
            v-model="topicsText"
            :auto-size="{ minRows: 2, maxRows: 6 }"
            placeholder="stress/#"
          />
        </a-form-item>
        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="每客户端订阅数">
              <a-input-number v-model="config.subscribePerClient" :min="0" :max="100" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="订阅 QoS">
              <a-radio-group v-model="config.qos" type="button">
                <a-radio :value="0">0</a-radio>
                <a-radio :value="1">1</a-radio>
                <a-radio :value="2">2</a-radio>
              </a-radio-group>
            </a-form-item>
          </a-col>
        </a-row>
      </a-collapse-item>

      <!-- ④ 发布配置 -->
      <a-collapse-item key="pub" header="④ 发布配置">
        <a-form-item label="启用发布压测">
          <a-switch v-model="config.publishEnabled" />
        </a-form-item>
        <template v-if="config.publishEnabled">
          <a-row :gutter="12">
            <a-col :span="12">
              <a-form-item label="发布器数量（goroutine）">
                <a-input-number v-model="config.publishWorkers" :min="1" :max="64" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="发布 QoS">
                <a-radio-group v-model="config.publishQos" type="button">
                  <a-radio :value="0">0</a-radio>
                  <a-radio :value="1">1</a-radio>
                  <a-radio :value="2">2</a-radio>
                </a-radio-group>
              </a-form-item>
            </a-col>
          </a-row>
          <a-form-item label="发布主题">
            <a-input v-model="config.publishTopic" placeholder="stress/pub" />
          </a-form-item>
          <a-row :gutter="12">
            <a-col :span="12">
              <a-form-item label="Payload 大小（字节）">
                <a-input-number v-model="config.publishPayloadSize" :min="8" :max="65535" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="目标速率（msg/s，0=不限）">
                <a-input-number v-model="config.publishRate" :min="0" :max="1000000" />
              </a-form-item>
            </a-col>
          </a-row>
          <a-row :gutter="12">
            <a-col :span="12">
              <a-form-item label="总条数（0=不限）">
                <a-input-number v-model="config.publishTotal" :min="0" :max="9999999999" />
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="运行时长（秒，0=不限）">
                <a-input-number v-model="config.duration" :min="0" :max="86400" />
              </a-form-item>
            </a-col>
          </a-row>
          <a-form-item label="回显测量（发布→收到回显 计算端到端延迟 P50/P95/P99）">
            <a-switch v-model="config.echoEnabled" />
          </a-form-item>
          <a-form-item v-if="config.echoEnabled" label="回显订阅主题">
            <a-input v-model="config.publishEchoTopic" placeholder="默认=发布主题" />
          </a-form-item>
        </template>
      </a-collapse-item>
    </a-collapse>
  </a-form>
</template>

<script lang="ts" setup>
import { ref, reactive, watch } from 'vue';

// 与后端 internal/stress/config.go StressConfig 字段一致
export interface StressConfig {
  broker: string;
  protocol: string; // 3.1 / 3.1.1 / 5.0 / auto
  username: string;
  password: string;
  clientCount: number;
  clientPrefix: string;
  connectRate: number;
  connectTimeout: number;
  keepalive: number;
  cleanSession: boolean;
  topics: string[];
  subscribePerClient: number;
  qos: number;
  publishEnabled: boolean;
  publishWorkers: number;
  publishTopic: string;
  publishPayloadSize: number;
  publishQos: number;
  publishRate: number;
  publishTotal: number;
  duration: number;
  echoEnabled: boolean;
  publishEchoTopic: string;
}

const props = defineProps<{ modelValue: StressConfig }>();
const emit = defineEmits<{ 'update:modelValue': [v: StressConfig] }>();

const config = reactive<StressConfig>({ ...props.modelValue });
const topicsText = ref<string>((config.topics || []).join('\n'));

// 主题文本 → 数组
watch(topicsText, (val) => {
  config.topics = val
    .split('\n')
    .map((s) => s.trim())
    .filter(Boolean);
});

// 同步回父组件
watch(
  config,
  () => {
    emit('update:modelValue', { ...config });
  },
  { deep: true }
);

defineExpose({ getConfig: () => ({ ...config }) });
</script>

<style scoped lang="less">
.stress-config {
  :deep(.arco-collapse-item-header) {
    font-weight: 600;
  }
  :deep(.arco-form-item) {
    margin-bottom: 14px;
  }
}
</style>
