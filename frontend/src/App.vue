<template>
  <a-config-provider :locale="locale">
    <router-view/>
  </a-config-provider>
</template>
<script lang="ts" setup>
import { onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore, useMqttStore } from '@/store';
//Events.On事件需要在这里监听一次才能使他在router路由界面生效
//获取建议所有的监听都在这里写，然后用pinia同步数据
import { Events, WML } from "@wailsio/runtime";
import { AppSettingService, MqttClientService } from '/#/gofly/internal/service';
import { Message, Notification } from '@arco-design/web-vue';
import enUS from '@arco-design/web-vue/es/locale/lang/en-us';
import zhCN from '@arco-design/web-vue/es/locale/lang/zh-cn';
import zhTW from '@arco-design/web-vue/es/locale/lang/zh-tw';
import useLocale from '@/hooks/locale';

  const { currentLocale } = useLocale();
  const router = useRouter();
  const locale = computed(() => {
    switch (currentLocale.value) {
      case 'en-US':
        return enUS;
      case 'zh-CN':
        return zhCN;
      case 'zh-TW':
        return zhTW;
      default:
        return enUS;
    }
  });
  //数据监听
  const appStore = useAppStore();
  onMounted(async () => {
    Events.On('time', (timeValue: { data: string }) => {
      // On a narrow screen the full RFC1123 stamp is too wide for the footer, so
      // show just the clock time there (matching the CSS breakpoint).
      const full = timeValue.data;
      const compact = (full.match(/\d{1,2}:\d{2}:\d{2}/) || [full])[0];
      const timestr = window.matchMedia('(max-width: 640px)').matches ? compact : full;
      appStore.sysTime=timestr;
    });
    // Wire up data-wml-openURL links (logos + footer "Docs" link).
    WML.Reload();
    appStore.toggleTheme(appStore.theme=='dark')

    // 启动时自动恢复连接（设置页"启动行为"开关）
    if (localStorage.getItem('auto-restore-conn') === '1') {
      try {
        const res = await AppSettingService.GetAutoRestoreConnections();
        if (res.code === 0 && res.data?.length) {
          let ok = 0;
          for (const c of res.data) {
            const r = await MqttClientService.Connect(c.id);
            if (r.code === 0) ok++;
          }
          Message.success(`已自动恢复 ${ok}/${res.data.length} 个连接`);
        }
      } catch (e) {
        // 启动恢复失败不影响应用启动
      }
    }
  })

  // 全局快捷键事件（后端 SetHotkeys 注册后触发）
  Events.On('hotkey:publish', () => {
    Message.info('快速发布：请进入连接页选择连接后进行发布');
    router.push('/connect');
  });

  //消息监听
  const mqttStore = useMqttStore()
  Events.On('mqtt:all_message', (ev: any) => {
    if (!ev.data) return
    const data = JSON.parse(ev.data)
    const connId = data.connection_id   // 后端已带数据库连接 ID，直接用
    if (!connId) return
    if (connId === mqttStore.currentConnId) {
      // 当前连接：展示消息，不计数
    } else {
      mqttStore.incrMsgCount(connId)    // 后台连接：未读 +1
      // 消息通知（设置页"消息通知"开关）
      if (localStorage.getItem('msg-notify') === '1') {
        const payload = typeof data.payload === 'string' ? data.payload : JSON.stringify(data.payload || '');
        const preview = payload.length > 120 ? payload.slice(0, 120) + '…' : payload;
        Notification.warning({
          title: `新消息 · ${data.topic || ''}`,
          content: preview,
          position: 'bottomRight',
          duration: 3000,
        });
        if (localStorage.getItem('msg-notify-sound') === '1') {
          playBeep();
        }
      }
    }
  })

// 通知提示音（Web Audio 短提示）
function playBeep() {
  try {
    const Ctx = window.AudioContext || (window as any).webkitAudioContext;
    const ctx = new Ctx();
    const osc = ctx.createOscillator();
    const gain = ctx.createGain();
    osc.connect(gain);
    gain.connect(ctx.destination);
    osc.type = 'sine';
    osc.frequency.value = 880;
    gain.gain.value = 0.08;
    osc.start();
    setTimeout(() => { osc.stop(); ctx.close(); }, 180);
  } catch (e) { /* 音频不可用时静默 */ }
}
</script>
<style lang="less">
#app{
    color: var(--color-neutral-10);
}
</style>  