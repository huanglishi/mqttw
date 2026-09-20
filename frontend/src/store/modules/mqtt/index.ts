import { defineStore } from 'pinia';
import { MqttState, MqttClientInfo } from './types';

const useMqttStore = defineStore('mqtt', {
  state: (): MqttState => {
    return {
      msgCount: {},
      clientIdMap: {},
      currentConnId: 0,
      currentClient: null,
    }
  },
  getters: {
    // 某连接的未读数
    getMsgCount: (state) => (connId: number) => state.msgCount[connId] || 0,
    // 全部连接未读总数（可用于菜单顶部汇总角标）
    getUnreadTotal: (state) => Object.values(state.msgCount).reduce((sum, n) => sum + n, 0),
  },
  actions: {
    // 设置当前打开/选中的连接，并清零该连接未读
    setCurrentConn(connId: number, client: MqttClientInfo | null = null) {
      this.currentConnId = connId;
      if (client) {
        this.currentClient = client;
      }
      this.clearMsgCount(connId);
    },
    // 清零某连接未读数
    clearMsgCount(connId: number) {
      if (this.msgCount[connId]) {
        this.msgCount[connId] = 0;
      }
    },
    // 未读数 +1（消息事件回调调用）
    incrMsgCount(connId: number) {
      this.msgCount[connId] = (this.msgCount[connId] || 0) + 1;
    },
    // 由连接列表构建 clientId → 连接ID 映射（跳过分组）
    buildClientIdMap(list: any[]) {
      const map: Record<string, number> = {};
      for (const item of list || []) {
        if (item.is_group == 1) continue;
        if (item.client_id) {
          map[item.client_id] = Number(item.id);
        }
      }
      this.clientIdMap = map;
    },
    // 由 clientId 反查连接ID，查不到返回 0
    connIdByClientId(clientId: string): number {
      return this.clientIdMap[clientId] || 0;
    },
  },
});

export default useMqttStore;
