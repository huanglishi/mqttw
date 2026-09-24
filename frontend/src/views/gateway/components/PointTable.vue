<template>
  <div>
    <!-- 设备选择 + 操作 -->
    <div class="pt-toolbar">
      <a-space>
        <span>设备：</span>
        <a-select
          v-model="deviceId"
          style="width: 280px"
          :loading="devLoading"
          placeholder="选择要配置点位的设备"
          @change="loadPoints"
        >
          <a-option v-for="d in devices" :key="d.id" :value="d.id">{{ d.name }}（#{{ d.id }}）</a-option>
        </a-select>
        <a-button type="primary" size="small" :disabled="!deviceId" @click="openPointForm(null)">
          <template #icon><icon-plus /></template>
          新增点位
        </a-button>
        <a-button size="small" :disabled="!deviceId" @click="importVisible = true">
          <template #icon><icon-upload /></template>
          CSV 批量导入
        </a-button>
        <a-button size="small" :disabled="!deviceId" @click="loadPoints">
          <template #icon><icon-refresh /></template>
        </a-button>
      </a-space>
    </div>

    <!-- 点位表 -->
    <a-table
      :data="points"
      :loading="loading"
      :pagination="{ pageSize: 10, showTotal: true }"
      row-key="id"
      :columns="columns"
      style="margin-top: 12px"
    >
      <template #register_type="{ record }">
        <a-tag>{{ registerLabel(record.register_type) }}</a-tag>
      </template>
      <template #data_type="{ record }">
        <a-tag color="arcoblue">{{ record.data_type }}</a-tag>
      </template>
      <template #only_on_change="{ record }">
        <a-badge :status="record.only_on_change === 1 ? 'success' : 'default'" :text="record.only_on_change === 1 ? '是' : '否'" />
      </template>
      <template #enabled="{ record }">
        <a-switch :model-value="record.enabled === 1" disabled size="small" />
      </template>
      <template #ops="{ record }">
        <a-space>
          <a-button size="mini" @click="openPointForm(record)">编辑</a-button>
          <a-popconfirm content="确认删除该点位？" @ok="removePoint(record)">
            <a-button size="mini" status="danger">删除</a-button>
          </a-popconfirm>
        </a-space>
      </template>
      <template #empty>
        <a-empty :description="deviceId ? '暂无点位，可点击「新增点位」或「CSV 批量导入」' : '请先选择设备'" />
      </template>
    </a-table>

    <!-- 点位表单 -->
    <a-modal
      v-model:visible="pointFormVisible"
      :title="editingPoint?.id ? '编辑点位' : '新增点位'"
      :footer="false"
      width="620px"
      unmount-on-close
    >
      <a-form :model="pointForm" layout="vertical">
        <a-form-item label="点位名称" field="name" :rules="[{ required: true, message: '请输入点位名称' }]">
          <a-input v-model="pointForm.name" placeholder="如 temp / pressure（作为 payload key）" allow-clear />
        </a-form-item>

        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="寄存器类型" field="registerType">
              <a-select v-model="pointForm.registerType" @change="onRegisterTypeChange">
                <a-option value="coil">线圈 0x（可读写）</a-option>
                <a-option value="discrete_input">离散输入 1x（只读）</a-option>
                <a-option value="holding">保持寄存器 3x（可读写）</a-option>
                <a-option value="input">输入寄存器 4x（只读）</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="起始地址" field="address" :rules="[{ required: true, message: '请输入地址' }]">
              <a-input-number v-model="pointForm.address" :min="0" :max="65535" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="数据类型" field="dataType">
              <a-select v-model="pointForm.dataType" :disabled="isBitOnly">
                <a-option value="uint16">uint16（1 寄存器）</a-option>
                <a-option value="int16">int16（1 寄存器）</a-option>
                <a-option value="uint32">uint32（2 寄存器）</a-option>
                <a-option value="int32">int32（2 寄存器）</a-option>
                <a-option value="float32">float32（2 寄存器）</a-option>
                <a-option value="bcd">BCD（1-2 寄存器）</a-option>
                <a-option value="bit">bit（线圈/离散输入）</a-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="寄存器数量" field="quantity">
              <a-input-number v-model="pointForm.quantity" :min="1" :max="125" :disabled="isBitOnly" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="缩放系数 (scale)" field="scale">
              <a-input-number v-model="pointForm.scale" :precision="6" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="偏移 (offset)" field="offset">
              <a-input-number v-model="pointForm.offset" :precision="6" style="width: 100%" />
            </a-form-item>
          </a-col>
        </a-row>
        <div class="formula-tip">换算公式：real = raw × scale + offset（如 scale=0.1 将 0-1000 映射为 0.0-100.0）</div>

        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="单位" field="unit">
              <a-input v-model="pointForm.unit" placeholder="如 ℃ / MPa / %" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="独立发布主题" field="topicOverride">
              <a-input v-model="pointForm.topicOverride" placeholder="留空并入设备主题" />
            </a-form-item>
          </a-col>
        </a-row>

        <a-row :gutter="12">
          <a-col :span="12">
            <a-form-item label="仅变化时上报" field="onlyOnChange">
              <a-switch v-model="pointForm.onlyOnChange" :checked-value="1" :unchecked-value="0" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="启用" field="enabled">
              <a-switch v-model="pointForm.enabled" :checked-value="1" :unchecked-value="0" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
      <template #footer>
        <a-space style="float: right">
          <a-button @click="pointFormVisible = false">取消</a-button>
          <a-button type="primary" :loading="saving" @click="savePoint">保存</a-button>
        </a-space>
      </template>
    </a-modal>

    <!-- CSV 导入 -->
    <a-modal v-model:visible="importVisible" title="CSV 批量导入点位" :footer="false" width="680px" unmount-on-close>
      <a-space direction="vertical" fill>
        <div class="imp-tip">
          每行一个点位，列顺序：name, register_type, address, quantity, data_type, scale, offset, unit, topic_override, only_on_change
          <br />
          首行可带表头（自动跳过）。示例：
          <code>temp,holding,0,1,uint16,0.1,0,℃,,0</code>
          <a-button size="mini" type="text" @click="downloadTemplate">
            <template #icon><icon-download /></template>
            下载模板
          </a-button>
        </div>
        <a-textarea v-model="csvText" :auto-size="{ minRows: 8, maxRows: 16 }" placeholder="粘贴 CSV 内容…" />
        <a-space>
          <a-button type="primary" :loading="importing" :disabled="!csvText.trim()" @click="doImport">导入</a-button>
          <a-button @click="csvText = ''">清空</a-button>
        </a-space>
      </a-space>
    </a-modal>
  </div>
</template>

<script lang="ts" setup>
import { ref, reactive, computed, onMounted } from 'vue';
import { Message } from '@arco-design/web-vue';
import { MqttGatewayService } from '/#/gofly/internal/service';

const props = defineProps<{ devices: any[] }>();
defineEmits<{ (e: 'change-device'): void }>();

const deviceId = ref<number>(0);
const points = ref<any[]>([]);
const loading = ref(false);
const devLoading = ref(false);

// 点位表单
const pointFormVisible = ref(false);
const editingPoint = ref<any>(null);
const saving = ref(false);

const defaultPoint = () => ({
  id: 0,
  deviceId: 0,
  name: '',
  registerType: 'holding',
  address: 0,
  quantity: 1,
  dataType: 'uint16',
  scale: 1,
  offset: 0,
  unit: '',
  topicOverride: '',
  onlyOnChange: 0,
  enabled: 1,
  sort: 0,
});

const pointForm = reactive<any>(defaultPoint());

// CSV 导入
const importVisible = ref(false);
const csvText = ref('');
const importing = ref(false);

const isBitOnly = computed(() => pointForm.registerType === 'coil' || pointForm.registerType === 'discrete_input');

const onRegisterTypeChange = () => {
  if (isBitOnly.value) {
    pointForm.dataType = 'bit';
    pointForm.quantity = 1;
  }
};

const loadPoints = async () => {
  if (!deviceId.value) {
    points.value = [];
    return;
  }
  loading.value = true;
  try {
    const res = await MqttGatewayService.ListModbusPoints(deviceId.value);
    if (res.code === 0) {
      points.value = res.data || [];
    } else {
      Message.error(res.message);
    }
  } catch (e) {
    Message.error('加载点位失败：' + e);
  } finally {
    loading.value = false;
  }
};

const openPointForm = (point: any) => {
  editingPoint.value = point;
  Object.assign(pointForm, defaultPoint(), point || {}, { deviceId: deviceId.value });
  pointFormVisible.value = true;
};

const savePoint = async () => {
  if (!pointForm.name) {
    Message.error('请填写点位名称');
    return;
  }
  saving.value = true;
  try {
    const res = await MqttGatewayService.SaveModbusPoint(JSON.stringify(pointForm));
    if (res.code === 0) {
      Message.success(res.message);
      pointFormVisible.value = false;
      await loadPoints();
    } else {
      Message.error(res.message);
    }
  } catch (e) {
    Message.error('保存失败：' + e);
  } finally {
    saving.value = false;
  }
};

const removePoint = async (point: any) => {
  try {
    const res = await MqttGatewayService.DeleteModbusPoint(point.id);
    if (res.code === 0) {
      Message.success(res.message);
      await loadPoints();
    } else {
      Message.error(res.message);
    }
  } catch (e) {
    Message.error('删除失败：' + e);
  }
};

const doImport = async () => {
  importing.value = true;
  try {
    const res = await MqttGatewayService.BatchImportModbusPoints(deviceId.value, csvText.value);
    if (res.code === 0) {
      Message.success(res.message);
      csvText.value = '';
      importVisible.value = false;
      await loadPoints();
    } else {
      Message.error(res.message);
    }
  } catch (e) {
    Message.error('导入失败：' + e);
  } finally {
    importing.value = false;
  }
};

const downloadTemplate = () => {
  const tpl = [
    'name,register_type,address,quantity,data_type,scale,offset,unit,topic_override,only_on_change',
    'temp,holding,0,1,uint16,0.1,0,℃,,0',
    'pressure,holding,1,2,float32,1,0,MPa,,0',
    'valve,coil,0,1,bit,1,0,,,1',
  ].join('\n');
  const blob = new Blob(['\ufeff' + tpl], { type: 'text/csv;charset=utf-8' });
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = 'modbus_points_template.csv';
  a.click();
  URL.revokeObjectURL(url);
};

const registerLabel = (r: string) =>
  ({ coil: '线圈 0x', discrete_input: '离散输入 1x', holding: '保持寄存器 3x', input: '输入寄存器 4x' }[r] || r);

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '名称', dataIndex: 'name' },
  { title: '寄存器', slotName: 'register_type', width: 120 },
  { title: '地址', dataIndex: 'address', width: 70 },
  { title: '数量', dataIndex: 'quantity', width: 70 },
  { title: '类型', slotName: 'data_type', width: 90 },
  { title: '倍率', dataIndex: 'scale', width: 80 },
  { title: '偏移', dataIndex: 'offset', width: 80 },
  { title: '单位', dataIndex: 'unit', width: 70 },
  { title: '仅变化', slotName: 'only_on_change', width: 80 },
  { title: '主题覆盖', dataIndex: 'topic_override', ellipsis: true },
  { title: '启用', slotName: 'enabled', width: 70 },
  { title: '操作', slotName: 'ops', width: 140, fixed: 'right' },
];

onMounted(() => {
  if (props.devices.length > 0) {
    deviceId.value = props.devices[0].id;
    loadPoints();
  }
});
</script>

<style scoped lang="less">
.pt-toolbar {
  display: flex;
  align-items: center;
}
.imp-tip {
  font-size: 12px;
  color: var(--color-text-3);
  line-height: 1.8;
  code {
    background: var(--color-fill-2);
    padding: 2px 6px;
    border-radius: 2px;
    font-size: 12px;
  }
}
.formula-tip {
  font-size: 12px;
  color: var(--color-text-3);
  margin-top: -8px;
}
</style>
