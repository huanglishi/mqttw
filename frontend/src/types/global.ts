import type {
  PropType as VuePropType,
} from 'vue';
//全局类型
declare global {
  type Nullable<T> = T | null;
  type TimeoutHandle = ReturnType<typeof setTimeout>;
  type IntervalHandle = ReturnType<typeof setInterval>;
   type Recordable<T = any> = Record<string, T>;
     // vue
   type PropType<T> = VuePropType<T>;
}

export interface AnyObject {
  [key: string]: unknown;
}

export interface Options {
  value: unknown;
  label: string;
}

export interface NodeOptions extends Options {
  children?: NodeOptions[];
}

export interface GetParams {
  body: null;
  type: string;
  url: string;
}

export interface PostData {
  body: string;
  type: string;
  url: string;
}

export interface Pagination {
  current: number;
  pageSize: number;
  total?: number;
}

export type TimeRanger = [string, string];

export interface GeneralChart {
  xAxis: string[];
  data: Array<{ name: string; value: number[] }>;
}

// MQTT 5 feature
// 用户属性键值对，与后端解析的 [{key,value}] 对象数组格式对齐
export interface UserProperty {
  key: string
  value: string
}

export interface PushPropertiesModel {
  user_properties?: UserProperty[] 
  response_topic?: string 
  payload_format?: boolean 
  message_expiry?: number | undefined
  topic_alias?: number | undefined
  content_type?: string 
  correlation_data?: string 
  subscription_identifier?: number | undefined
}
// 对话消息
export interface MessageModel {
  id?: number
  connection_id: number
  msg_type: string
  qos: 0 | 1 | 2
  retain: boolean 
  topic: string
  color?: string
  payload: string
  //mqtt v5
  user_properties?: UserProperty[] | null
  response_topic?: string | null
  content_type?: string | null
  correlation_data?: string | null
  message_expiry?: number | null
  create_at: string
}
// 消息记录
export interface MessageHistoryModel {
  id?: number
  payload: string
  payloadType: string
  create_at: string
}
// Topic记录
export interface TopicHistoryModel {
  id?: number
  topic: string
  qos: 0 | 1 | 2
  retain: boolean
  create_at: string
}