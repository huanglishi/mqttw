import JSONBigInt from 'json-bigint';

// 实例：启用原生 BigInt
const JSONBigIntNative = JSONBigInt({ useNativeBigInt: true });
// 普通实例，返回 BigNumber 对象
const JSONBigNumber = JSONBigInt();

/**
 * JSON parse 增强，自动处理超大整数(BigInt/BigNumber)
 */
export const jsonParse: typeof JSON.parse = (...args: any[]) => {
  try {
    // 优先使用 native BigInt
    return JSONBigIntNative.parse(...args);
  } catch {
    try {
      // 失败降级为 BigNumber 对象
      return JSONBigNumber.parse(...args);
    } catch {
      // 兜底原生 JSON.parse
      return JSON.parse(args[0], args[1]);
    }
  }
};

/**
 * JSON stringify 增强
 */
export const jsonStringify: typeof JSON.stringify = (...args: any[]) => {
  try {
    return JSONBigIntNative.stringify(...args);
  } catch (_error) {
    return JSONBigNumber.stringify(...args);
  }
};

/**
 * ECharts JSON树节点类型
 */
export interface EChartsJsonTreeNodeLike {
  raw?: any;
  [key: string]: any;
}

/**
 * 序列化节点子树，优先取 raw
 * @param nodeLike 节点对象
 * @param space 缩进空格，默认2
 */
export function stringifySubtree(nodeLike: EChartsJsonTreeNodeLike, space: number = 2): string {
  const raw = nodeLike && Object.prototype.hasOwnProperty.call(nodeLike, 'raw') ? nodeLike.raw : nodeLike;
  try {
    return jsonStringify(raw, null, space);
  } catch (_err) {
    return JSON.stringify(raw, null, space);
  }
}

/**
 * 递归把 BigInt / BigNumber 转为字符串，适配 ECharts 渲染
 * @param obj 待转换对象
 */
export function toPlainObject(obj: any): any {
  if (obj === null || obj === undefined) {
    return obj;
  }
  if (typeof obj === 'bigint') {
    return obj.toString();
  }
  if (
    typeof obj === 'object' &&
    obj !== null &&
    obj.constructor &&
    (obj.constructor.name === 'BigNumber' || obj.constructor.name === 'BN')
  ) {
    if (typeof obj.toString === 'function') {
      return obj.toString();
    }
    return String(obj);
  }
  if (Array.isArray(obj)) {
    return obj.map((item) => toPlainObject(item));
  }
  if (typeof obj === 'object') {
    const result: any = {};
    for (const key in obj) {
      if (Object.prototype.hasOwnProperty.call(obj, key)) {
        result[key] = toPlainObject(obj[key]);
      }
    }
    return result;
  }
  return obj;
}
//判断当前 payload 是否为 JSON。
export function isValidJson(str: string): boolean {
  if (!str || typeof str !== 'string') return false
  try {
    jsonParse(str)
    return true
  } catch {
    return false
  }
}