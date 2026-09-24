<template>
  <div class="gateway-page">
    <!-- 顶部：网关运行状态与控制 -->
    <a-card :bordered="false" class="gw-topbar">
      <div class="gw-topbar-row">
        <a-space size="large" align="center">
          <a-badge :status="running ? 'processing' : 'default'" :text="running ? '网关运行中' : '网关已停止'" />
          <span class="gw-stat">设备 <b>{{ status.device_count || 0 }}</b></span>
          <span class="gw-stat">总点位 <b>{{ status.total_points || 0 }}</b></span>
          <span class="gw-stat">轮询 <b>{{ status.total_polls || 0 }}</b></span>
          <span class="gw-stat" :class="{ 'gw-err': (status.total_errors || 0) > 0 }">错误 <b>{{ status.total_errors || 0 }}</b></span>
        </a-space>
        <a-space>
          <a-button type="primary" :loading="starting" :disabled="running" @click="handleStart">
            <template #icon><icon-play-circle /></template>
            启动网关
          </a-button>
          <a-button status="danger" :disabled="!running" @click="handleStop">
            <template #icon><icon-poweroff /></template>
            停止网关
          </a-button>
          <a-button @click="refreshAll">
            <template #icon><icon-refresh /></template>
            刷新
          </a-button>
        </a-space>
      </div>
    </a-card>

    <!-- 三个功能区 -->
    <a-tabs v-model:active-key="activeTab" class="gw-tabs">
      <a-tab-pane key="devices" title="设备管理">
        <DeviceTable
          :devices="devices"
          :loading="devLoading"
          :status="status"
          @edit="openDeviceForm"
          @remove="removeDevice"
          @test="testDevice"
          @add="openDeviceForm(null)"
        />
      </a-tab-pane>
      <a-tab-pane key="points" title="点位配置">
        <PointTable
          :devices="devices"
          @change-device="onChangeDevice"
        />
      </a-tab-pane>
      <a-tab-pane key="monitor" title="运行监控">
        <RealtimePanel :status="status" />
        <LogPanel />
      </a-tab-pane>
    </a-tabs>

    <!-- 设备表单抽屉 -->
    <DeviceForm
      v-model:visible="deviceFormVisible"
      :device="editingDevice"
      :broker-options="brokerOptions"
      @submit="saveDevice"
    />

    <!-- 测试读取结果 -->
    <a-modal
      v-model:visible="testVisible"
      title="在线读取测试"
      :footer="false"
      width="720px"
      unmount-on-close
    >
      <a-spin :loading="testing" style="width: 100%">
        <a-table
          :data="testResult"
          :pagination="false"
          :scroll="{ x: 680 }"
          size="small"
          :columns="testColumns"
        />
        <a-empty v-if="!testing && testResult.length === 0" description="无读取结果" />
      </a-spin>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';
import { Events } from '@wailsio/runtime';
import { Message } from '@arco-design/web-vue';
import { MqttGatewayService } from '/#/gofly/internal/service';
import DeviceTable from './components/DeviceTable.vue';
import DeviceForm from './components/DeviceForm.vue';
import PointTable from './components/PointTable.vue';
import RealtimePanel from './components/RealtimePanel.vue';
import LogPanel from './components/LogPanel.vue';

const activeTab = ref('devices');
const devices = ref<any[]>([]);
const devLoading = ref(false);
const brokerOptions = ref<any[]>([]);
const status = ref<any>({ running: false, device_count: 0 });
const running = computed(() => !!status.value.running);
const starting = ref(false);

// 设备表单
const deviceFormVisible = ref(false);
const editingDevice = ref<any>(null);

// 测试读取
const testVisible = ref(false);
const testing = ref(false);
const testResult = ref<any[]>([]);
const testColumns = [
  { title: '点位', dataIndex: 'name' },
  { title: '类型', dataIndex: 'data_type' },
  { title: '地址', dataIndex: 'address' },
  { title: '值', dataIndex: 'value' },
  { title: '单位', dataIndex: 'unit' },
  { title: '原始字节', dataIndex: 'raw' },
  { title: '时间', dataIndex: 'ts', width: 130 },
];

const loadDevices = async () => {
  devLoading.value = true;
  try {
    const res = await MqttGatewayService.ListModbusDevices();
    if (res.code === 0) {
      devices.value = res.data || [];
    } else {
      Message.error(res.message);
    }
  } catch (e) {
    Message.error('加载设备失败：' + e);
  } finally {
    devLoading.value = false;
  }
};

const loadBrokerOptions = async () => {
  try {
    const res = await MqttGatewayService.ListGatewayBrokerOptions();
    if (res.code === 0) {
      brokerOptions.value = res.data || [];
    }
  } catch (e) {
    // Broker 列表加载失败不影响设备管理
  }
};

const fetchStatus = async () => {
  try {
    const res = await MqttGatewayService.GetModbusGatewayStatus();
    if (res.code === 0) {
      status.value = res.data || {};
    }
  } catch (e) {
    // 忽略
  }
};

const handleStart = async () => {
  starting.value = true;
  try {
    const res = await MqttGatewayService.StartModbusGateway();
    if (res.code === 0) {
      Message.success(res.message);
      await fetchStatus();
    } else {
      Message.error(res.message);
    }
  } catch (e) {
    Message.error('启动失败：' + e);
  } finally {
    starting.value = false;
  }
};

const handleStop = async () => {
  try {
    const res = await MqttGatewayService.StopModbusGateway();
    if (res.code === 0) {
      Message.success(res.message);
      status.value = { running: false, device_count: 0 };
    } else {
      Message.error(res.message);
    }
  } catch (e) {
    Message.error('停止失败：' + e);
  }
};

const refreshAll = async () => {
  await Promise.all([loadDevices(), fetchStatus()]);
};

const openDeviceForm = (device: any) => {
  editingDevice.value = device;
  deviceFormVisible.value = true;
};

const saveDevice = async (form: any) => {
  try {
    const res = await MqttGatewayService.SaveModbusDevice(JSON.stringify(form));
    if (res.code === 0) {
      Message.success(res.message);
      deviceFormVisible.value = false;
      await loadDevices();
    } else {
      Message.error(res.message);
    }
  } catch (e) {
    Message.error('保存失败：' + e);
  }
};

const removeDevice = async (device: any) => {
  try {
    const res = await MqttGatewayService.DeleteModbusDevice(device.id);
    if (res.code === 0) {
      Message.success(res.message);
      await loadDevices();
    } else {
      Message.error(res.message);
    }
  } catch (e) {
    Message.error('删除失败：' + e);
  }
};

const testDevice = async (device: any) => {
  testResult.value = [];
  testVisible.value = true;
  testing.value = true;
  try {
    const res = await MqttGatewayService.TestReadModbus(device.id);
    if (res.code === 0) {
      testResult.value = res.data || [];
      Message.success(res.message);
    } else {
      Message.error(res.message);
    }
  } catch (e) {
    Message.error('测试读取失败：' + e);
  } finally {
    testing.value = false;
  }
};

const onChangeDevice = () => {
  // 点位页自行加载，无需处理
};

// 后端推送的状态事件 → 同步顶部状态条
const onStatusEvent = (ev: any) => {
  if (!ev?.data) return;
  try {
    const payload = typeof ev.data === 'string' ? JSON.parse(ev.data) : ev.data;
    if (payload && typeof payload.running === 'boolean') {
      status.value = payload;
    }
  } catch (e) {
    // 忽略解析失败
  }
};

onMounted(() => {
  loadDevices();
  loadBrokerOptions();
  fetchStatus();
  Events.On('gateway:status', onStatusEvent);
});

onUnmounted(() => {
  Events.Off('gateway:status');
});
</script>

<style scoped lang="less">
.gateway-page {
  padding: 12px;

  .gw-topbar {
    margin-bottom: 0px;
    .gw-topbar-row {
      display: flex;
      justify-content: space-between;
      align-items: center;
      flex-wrap: wrap;
      gap: 8px;
    }

    .gw-stat {
      color: var(--color-text-3);
      b {
        color: var(--color-text-1);
      }
      &.gw-err b {
        color: rgb(var(--danger-6));
      }
    }
  }

  .gw-tabs {
    background: var(--color-bg-2);
    border-radius: 4px;
    padding: 4px 12px 12px;
  }
}
</style>
