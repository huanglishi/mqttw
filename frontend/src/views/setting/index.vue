<template>
  <div class="container">
    <!-- 左：设置项 -->
    <div class="left-side">
      <!-- 外观 -->
      <a-card :bordered="false" class="card">
        <template #title>
          <span class="card-title"><i class="title-bar title-bar-blue"></i>外观</span>
        </template>
        <a-form-item field="theme" label="主题模式">
          <a-radio-group v-model="formData.theme" type="button" @change="handleTheme">
            <a-radio value="light">亮色</a-radio>
            <a-radio value="dark">黑色</a-radio>
            <a-radio value="auto">跟随系统</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-card>

      <!-- 窗口行为 -->
      <a-card :bordered="false" class="card">
        <template #title>
          <span class="card-title"><i class="title-bar title-bar-green"></i>窗口行为</span>
        </template>
        <div class="switch-row">
          <a-form-item label="置顶窗口">
            <a-switch type="round" @change="changeSetAlwaysOnTop">
              <template #checked>开启</template>
              <template #unchecked>关闭</template>
            </a-switch>
          </a-form-item>
          <a-form-item label="无边框窗口">
            <a-switch type="round" @change="changeSetFrameless">
              <template #checked>开启</template>
              <template #unchecked>关闭</template>
            </a-switch>
          </a-form-item>
          <a-form-item label="窗口可调整大小">
            <a-switch type="round" :default-checked="true" @change="changeSetResizable">
              <template #checked>可以</template>
              <template #unchecked>不可</template>
            </a-switch>
          </a-form-item>
        </div>
        <div class="action-row">
          <a-button type="primary" long @click="() => Window.Reload()">
            <template #icon><icon-refresh /></template>
            重置窗口
          </a-button>
        </div>
      </a-card>

      <!-- 语言 -->
      <a-card :bordered="false" class="card">
        <template #title>
          <span class="card-title"><i class="title-bar title-bar-orange"></i>语言与区域</span>
        </template>
        <a-form-item label="界面语言">
          <a-radio-group v-model="currentLocale" type="button" @change="handleLocale">
            <a-radio value="zh-CN">简体中文</a-radio>
            <a-radio value="zh-TW">繁體中文</a-radio>
            <a-radio value="en-US">English</a-radio>
          </a-radio-group>
        </a-form-item>
      </a-card>

      <!-- 数据与缓存 -->
      <a-card :bordered="false" class="card">
        <template #title>
          <span class="card-title"><i class="title-bar title-bar-red"></i>数据与缓存</span>
        </template>
        <div class="cache-row">
          <div class="cache-desc">
            清除浏览器本地存储（连接配置、压测配置、主题、语言等），并刷新窗口恢复到初始状态。
          </div>
          <a-popconfirm
            content="将清除全部本地缓存数据并刷新窗口，确定继续？"
            @ok="clearCache"
          >
            <a-button status="danger" type="primary">
              <template #icon><icon-delete /></template>
              清除缓存
            </a-button>
          </a-popconfirm>
        </div>
      </a-card>

      <!-- 启动行为 -->
      <a-card :bordered="false" class="card">
        <template #title>
          <span class="card-title"><i class="title-bar title-bar-cyan"></i>启动行为</span>
        </template>
        <div class="switch-row">
          <a-form-item label="开机自启">
            <a-switch type="round" v-model="autostartSwitch" @change="setAutostart" />
          </a-form-item>
          <a-form-item label="启动时恢复连接">
            <a-switch type="round" v-model="autoRestoreSwitch" @change="setAutoRestore" />
          </a-form-item>
          <a-form-item label="关闭时最小化到托盘">
            <a-switch type="round" v-model="trayMinimizeSwitch" @change="setTrayMinimize" />
          </a-form-item>
        </div>
      </a-card>

      <!-- 消息通知 -->
      <a-card :bordered="false" class="card">
        <template #title>
          <span class="card-title"><i class="title-bar title-bar-gold"></i>消息通知</span>
        </template>
        <div class="switch-row">
          <a-form-item label="新消息弹窗提醒">
            <a-switch type="round" v-model="notifySwitch" @change="setNotify" />
          </a-form-item>
          <a-form-item label="通知提示音">
            <a-switch type="round" v-model="notifySoundSwitch" @change="setNotifySound" />
          </a-form-item>
        </div>
      </a-card>

      <!-- 全局快捷键 -->
      <a-card :bordered="false" class="card">
        <template #title>
          <span class="card-title"><i class="title-bar title-bar-gray"></i>全局快捷键</span>
        </template>
        <div class="hotkey-row">
          <a-switch type="round" v-model="hotkeysSwitch" @change="setHotkeys" />
          <ul class="hotkey-list">
            <li><kbd>Ctrl + Shift + H</kbd> 显示 / 隐藏主窗口</li>
            <li><kbd>Ctrl + Shift + D</kbd> 断开全部连接</li>
            <li><kbd>Ctrl + Shift + P</kbd> 快速发布（跳转连接页）</li>
          </ul>
        </div>
      </a-card>

      <!-- 运行目录 -->
      <a-card :bordered="false" class="card">
        <template #title>
          <span class="card-title"><i class="title-bar title-bar-pink"></i>运行目录</span>
        </template>
        <div class="cache-row">
          <div class="cache-desc">
            运行目录存放导出文件（压测结果 CSV 等）。日志级别与保留策略将在后续版本完善。
          </div>
          <a-button @click="openRuntimeDir">
            <template #icon><icon-folder-open /></template>
            打开运行目录
          </a-button>
        </div>
      </a-card>

      <!-- 关于 -->
      <a-card :bordered="false" class="card">
        <template #title>
          <span class="card-title"><i class="title-bar title-bar-lime"></i>关于 MQTTW</span>
        </template>
        <div class="about-row">
          <div class="about-info">
            <div class="about-name">{{ appInfo?.name || 'MQTTW' }}</div>
            <div class="about-version">版本 {{ appInfo?.version || '-' }}</div>
            <div class="about-desc">{{ appInfo?.description || '' }}</div>
          </div>
          <div class="about-actions">
            <a-button size="small" @click="openURL(appInfo?.repo_releases)">
              <template #icon><icon-update /></template>
              检查更新
            </a-button>
            <a-button size="small" @click="openURL(appInfo?.community)">
              <template #icon><icon-link /></template>
              GoFly 社区
            </a-button>
            <a-button size="small" @click="openURL(appInfo?.gmqt)">
              <template #icon><icon-link /></template>
              GMQT Broker
            </a-button>
          </div>
        </div>
      </a-card>
    </div>

    <!-- 右：系统信息 -->
    <div class="right-side">
      <a-card
        :bordered="false"
        class="card"
        :style="{ '--wails-draggable': 'drag' }"
      >
        <template #title>
          <span class="card-title"><i class="title-bar title-bar-purple"></i>系统信息</span>
        </template>
        <a-descriptions :column="1" size="medium">
          <a-descriptions-item label="窗口 name">{{ windowsname }}</a-descriptions-item>
          <a-descriptions-item label="系统(OS)">{{ Osinfo?.OS }}</a-descriptions-item>
          <a-descriptions-item label="系统品牌">{{ Osinfo?.OSInfo.Branding }}</a-descriptions-item>
          <a-descriptions-item label="系统版本">{{ Osinfo?.OSInfo.Version }}</a-descriptions-item>
          <a-descriptions-item label="Arch">{{ Osinfo?.Arch }}</a-descriptions-item>
          <a-descriptions-item label="调试状态">{{ Osinfo?.Debug }}</a-descriptions-item>
          <a-descriptions-item label="WebView2">{{ Osinfo?.PlatformInfo.WebView2 }}</a-descriptions-item>
        </a-descriptions>
        <div class="drag-tip">无边框窗口下，拖动此卡片可移动窗口</div>
      </a-card>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAppStore } from '@/store';
import { System, Window, Browser } from '@wailsio/runtime';
import { AppSettingService, SystemService } from '/#/gofly/internal/service';
import { Message } from '@arco-design/web-vue';

const appStore = useAppStore();
const { locale } = useI18n();

const formData = ref({
  theme: appStore.theme,
});
const currentLocale = ref<string>(locale.value);
const windowsname = ref('');
const Osinfo = ref<any>();

// 启动行为
const autostartSwitch = ref(false);
const trayMinimizeSwitch = ref(false);
const autoRestoreSwitch = ref(localStorage.getItem('auto-restore-conn') === '1');
// 消息通知
const notifySwitch = ref(localStorage.getItem('msg-notify') === '1');
const notifySoundSwitch = ref(localStorage.getItem('msg-notify-sound') === '1');
// 快捷键
const hotkeysSwitch = ref(false);
// 关于
const appInfo = ref<any>(null);

onMounted(async () => {
  windowsname.value = await Window.Name();
  Osinfo.value = await System.Environment();

  // 恢复后端状态
  const [autostart, tray, hotkeys, info] = await Promise.all([
    AppSettingService.GetAutostart(),
    AppSettingService.GetTrayMinimize(),
    AppSettingService.GetHotkeys(),
    AppSettingService.GetAppInfo(),
  ]);
  autostartSwitch.value = autostart.code === 0 ? !!autostart.data?.enabled : false;
  trayMinimizeSwitch.value = tray.code === 0 ? !!tray.data?.enabled : false;
  hotkeysSwitch.value = hotkeys.code === 0 ? !!hotkeys.data?.enabled : false;
  if (info.code === 0) appInfo.value = info.data;
});

// 切换主题
const handleTheme = async () => {
  let isdark = false;
  if (formData.value.theme == 'dark') {
    isdark = true;
  } else if (formData.value.theme == 'auto') {
    // 跟随系统
    isdark = await System.IsDarkMode();
  } else {
    isdark = false;
  }
  appStore.toggleTheme(isdark);
  if (isdark) {
    await Window.SetBackgroundColour(6, 7, 15, 255);
    document.documentElement.setAttribute('data-theme', 'dark');
    localStorage.setItem('app-theme', 'dark');
  } else {
    document.body.classList.add('light-theme');
    document.body.classList.remove('dark-theme');
  }
  SystemService.SetTheme(formData.value.theme);
};

// 切换界面语言（简体中文 / 繁體中文 / English）
const handleLocale = (value: any) => {
  if (locale.value === value) return;
  locale.value = value;
  localStorage.setItem('arco-locale', value);
  Message.success(
    value === 'en-US' ? 'Language switched' : value === 'zh-TW' ? '語言已切換' : '语言已切换'
  );
};

// 清除本地缓存并刷新窗口
const clearCache = () => {
  localStorage.clear();
  sessionStorage.clear();
  Message.success('缓存已清除，正在刷新窗口...');
  setTimeout(() => {
    Window.Reload();
  }, 800);
};

// 窗口置顶
const changeSetAlwaysOnTop = (value: any) => {
  Window.SetAlwaysOnTop(value);
  Message.success({ content: value ? '窗口已位于顶部' : '窗口已取消置顶', id: 'setting' });
};
// 无边框
const changeSetFrameless = (value: any) => {
  Window.SetFrameless(value);
  Message.success({ content: value ? '窗口已无边框' : '窗口已取消无边框', id: 'setting' });
};
// 可调整大小
const changeSetResizable = (value: any) => {
  Window.SetResizable(value);
  Message.success({ content: value ? '窗口可调整大小' : '窗口不可调整大小', id: 'setting' });
};

// ========== 启动行为 ==========
const setAutostart = async (val: any) => {
  const res = await AppSettingService.SetAutostart(val);
  res.code === 0 ? Message.success(res.message) : Message.error(res.message);
};
const setAutoRestore = (val: any) => {
  localStorage.setItem('auto-restore-conn', val ? '1' : '0');
  Message.success(val ? '下次启动将自动恢复连接' : '已关闭启动时恢复连接');
};
const setTrayMinimize = async (val: any) => {
  const res = await AppSettingService.SetTrayMinimize(val);
  res.code === 0 ? Message.success(res.message) : Message.error(res.message);
};

// ========== 消息通知 ==========
const setNotify = (val: any) => {
  localStorage.setItem('msg-notify', val ? '1' : '0');
  Message.success(val ? '新消息提醒已开启' : '新消息提醒已关闭');
};
const setNotifySound = (val: any) => {
  localStorage.setItem('msg-notify-sound', val ? '1' : '0');
  Message.success(val ? '通知提示音已开启' : '通知提示音已关闭');
};

// ========== 全局快捷键 ==========
const setHotkeys = async (val: any) => {
  const res = await AppSettingService.SetHotkeys(val);
  res.code === 0 ? Message.success(res.message) : Message.error(res.message);
};

// ========== 运行目录 ==========
const openRuntimeDir = async () => {
  const res = await AppSettingService.OpenRuntimeDir();
  res.code === 0 ? Message.success(res.message) : Message.error(res.message);
};

// ========== 关于 ==========
const openURL = (url?: string) => {
  if (url) Browser.OpenURL(url);
};
</script>

<style scoped lang="less">
.container {
  padding: 12px;
  display: flex;
  gap: 16px;
  align-items: flex-start;

  .left-side {
    flex: 1;
    min-width: 0;
    display: flex;
    flex-direction: column;
    gap: 14px;
  }

  .right-side {
    width: 300px;
    flex: none;
  }

  .card {
    border-radius: 8px;
    background: var(--color-bg-2);
    box-shadow: 0 10px 11px rgb(var(--arcoblue-3), 0.08),
      0 6px 4px rgb(var(--arcoblue-3), 0.06),
      0 0 0 1px rgb(var(--arcoblue-3), 0.05);
    transition: box-shadow 0.2s ease, transform 0.2s ease;

    &:hover {
      box-shadow: 0 14px 28px rgb(var(--arcoblue-3), 0.12),
        0 8px 10px rgb(var(--arcoblue-3), 0.08),
        0 0 0 1px rgb(var(--arcoblue-3), 0.06);
      transform: translateY(-1px);
    }

    :deep(.arco-card-header) {
      font-size: 15px;
      font-weight: 600;
      padding: 12px 16px;
      border-bottom: 1px solid var(--color-border-2);
    }

    :deep(.arco-card-body) {
      padding: 16px;
    }

    :deep(.arco-form-item) {
      margin-bottom: 0;
    }

    :deep(.arco-form-item-label-col) {
      margin-bottom: 4px;
    }
  }

  // 卡片标题装饰条
  .card-title {
    display: inline-flex;
    align-items: center;
    gap: 8px;
  }

  .title-bar {
    width: 3px;
    height: 15px;
    border-radius: 2px;
    background: #165dff;
  }

  .title-bar-blue {
    background: linear-gradient(180deg, #165dff, #0fc6c2);
  }
  .title-bar-green {
    background: linear-gradient(180deg, #00b42a, #23c343);
  }
  .title-bar-orange {
    background: linear-gradient(180deg, #ff7d00, #ffb400);
  }
  .title-bar-red {
    background: linear-gradient(180deg, #f53f3f, #ff7d00);
  }
  .title-bar-purple {
    background: linear-gradient(180deg, #722ed1, #b37feb);
  }
  .title-bar-cyan {
    background: linear-gradient(180deg, #0fc6c2, #14c9c9);
  }
  .title-bar-gold {
    background: linear-gradient(180deg, #f7ba1e, #ffb400);
  }
  .title-bar-gray {
    background: linear-gradient(180deg, #86909c, #c9cdd4);
  }
  .title-bar-pink {
    background: linear-gradient(180deg, #f5319d, #eb0aa4);
  }
  .title-bar-lime {
    background: linear-gradient(180deg, #9fdb1d, #c3e32e);
  }

  .switch-row {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;

    @media (max-width: 560px) {
      grid-template-columns: 1fr;
    }

    :deep(.arco-form-item) {
      display: flex;
      flex-direction: column;
      padding: 12px 14px;
      border-radius: 8px;
      background: var(--color-fill-1);

      :deep(.arco-form-item-label-col) {
        margin-bottom: 10px;
      }
    }
  }

  .action-row {
    margin-top: 14px;

    .arco-btn {
      display: inline-flex;
      align-items: center;
    }
  }

  .cache-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 16px;
    flex-wrap: wrap;

    .cache-desc {
      flex: 1;
      min-width: 240px;
      font-size: 13px;
      line-height: 1.7;
      color: var(--color-text-3);
    }
  }

  .hotkey-row {
    display: flex;
    align-items: flex-start;
    gap: 24px;

    @media (max-width: 480px) {
      flex-direction: column;
    }

    .hotkey-list {
      margin: 0;
      padding-left: 2px;
      list-style: none;

      li {
        font-size: 13px;
        line-height: 2.1;
        color: var(--color-text-2);
      }

      kbd {
        display: inline-block;
        padding: 1px 8px;
        margin-right: 6px;
        border-radius: 4px;
        font-size: 12px;
        background: var(--color-fill-2);
        border: 1px solid var(--color-border-2);
        color: var(--color-text-1);
      }
    }
  }

  .about-row {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 16px;
    flex-wrap: wrap;

    .about-name {
      font-size: 17px;
      font-weight: 600;
    }

    .about-version {
      margin: 3px 0 4px;
      font-size: 12.5px;
      color: var(--color-text-3);
    }

    .about-desc {
      font-size: 12.5px;
      line-height: 1.6;
      color: var(--color-text-3);
    }

    .about-actions {
      display: flex;
      gap: 8px;
      flex-wrap: wrap;

      .arco-btn {
        display: inline-flex;
        align-items: center;
      }
    }
  }

  .drag-tip {
    margin-top: 12px;
    padding-top: 10px;
    border-top: 1px dashed var(--color-border-2);
    font-size: 12px;
    color: var(--color-text-4);
  }

  @media (max-width: 900px) {
    flex-direction: column;

    .right-side {
      width: 100%;
    }
  }
}
</style>
