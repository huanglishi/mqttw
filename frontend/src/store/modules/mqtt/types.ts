// MQTT 客户端信息（当前打开的连接）
export interface MqttClientInfo {
  id: number
  pid: number
  is_group: number
  title: string
  client_id?: string
  protocol?: string
  host?: string
  port?: string
  mqtt_version?: string
  [key: string]: any
}

// MQTT Store 状态
export interface MqttState {
  // 连接ID → 未读消息数（侧边栏菜单 badge 用）
  msgCount: Record<number, number>
  // clientId(字符串) → 连接ID，消息事件回调里反查用
  clientIdMap: Record<string, number>
  // 当前打开的连接ID
  currentConnId: number
  // 当前客户端信息
  currentClient: MqttClientInfo | null
}
