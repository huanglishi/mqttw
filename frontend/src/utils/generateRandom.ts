// generateRandom.ts
/**
 * 简易UUID v4 实现
 */
function uuidv4(): string {
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = Math.random() * 16 | 0
    const v = c === 'x' ? r : (r & 0x3 | 0x8)
    return v.toString(16)
  })
}

/**
 * MQTT 客户端ID，前缀 mqttw_，8位随机16进制字符
 */
export const getClientId = () => `mqttw_${Math.random().toString(16).substring(2, 10)}` as string

/**
 * 连接集合ID（分组）
 */
export const getCollectionId = () => `mqttw_collection_${uuidv4()}` as string

/**
 * 订阅ID
 */
export const getSubscriptionId = () => `mqttw_sub_${uuidv4()}` as string

/**
 * 消息记录ID
 */
export const getMessageId = () => `mqttw_msg_${uuidv4()}` as string

/**
 * AI助手消息ID
 */
export const getCopilotMessageId = () => `mqttw_copilot_${uuidv4()}` as string

export default {
  getClientId,
  getCollectionId,
  getSubscriptionId,
  getMessageId,
  getCopilotMessageId,
}
