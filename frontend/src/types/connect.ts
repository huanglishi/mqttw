/**
 * 连接数据类型
*/
  // 连接属性（MQTT5.0）
 export interface ConnectProperties {
    sessionExpiryInterval: number
    receiveMaximum: number|undefined
    maximumPacketSize: number|undefined
    topicAliasMaximum: number|undefined
    requestResponseInformation: number
    requestProblemInformation: number
  }

  // 遗嘱消息
 export interface WillMessage {
    willTopic: string
    willQos: number
    willRetain: number
    willPayload: string
    payloadFormatIndicator: number
    willDelayInterval: number|undefined
    messageExpiryInterval: number|undefined
    contentType: string
    responseTopic: string
    correlationData: string
  }

  // 用户属性键值对（KeyValueEditor 模型）
 export interface UserProperty {
    key: string
    value: string
  }

// 连接表单
export  interface MqttConnectionForm {
    id: number
    pid: number
    title: string
    protocol: string
    host: string
    port: string
    client_id: string
    path: string
    username: string
    password: string
    ssl: number
    ssl_security: number
    alpn: string
    cert_type: string
    ca: string
    cert: string
    key: string
    mqtt_version: string
    connect_timeout: number
    keepalive: number
    reconnect: number
    reconnect_period: number
    clean: number
    properties: ConnectProperties
    will: WillMessage
    user_properties: UserProperty[]
  }
