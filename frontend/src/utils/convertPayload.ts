import { jsonParse, jsonStringify } from './jsonBigint';


interface CodeType {
  encode: (str: string) => string
  decode: (str: string) => string
}

/**
 * 纯前端 UTF-8 ↔ Hex 转换（替代 Node Buffer）
 * 使用 TextEncoder / TextDecoder，浏览器原生支持
 */
const utf8ToHex = (str: string): string =>
  Array.from(new TextEncoder().encode(str))
    .map((b) => b.toString(16).padStart(2, '0'))
    .join('');

const hexToUtf8 = (hex: string): string => {
  const clean = hex.replace(/\s+/g, '');
  // 校验是否为合法hex
  if (!/^[0-9a-fA-F]*$/.test(clean)) {
    throw new Error('Invalid hex string');
  }
  const bytes = new Uint8Array(
    clean.match(/.{1,2}/g)?.map((b) => parseInt(b, 16)) ?? []
  );
  return new TextDecoder('utf-8').decode(bytes);
};

const convertBase64 = (value: string, codeType: 'encode' | 'decode'): string => {
  const convertMap: CodeType = {
    encode(str: string): string {
      // 支持中文的 UTF-8 安全 Base64 编码
      return btoa(unescape(encodeURIComponent(str)));
    },
    decode(str: string): string {
      return decodeURIComponent(escape(atob(str)));
    },
  };
  return convertMap[codeType](value);
};

const convertHex = (value: string, codeType: 'encode' | 'decode'): string => {
  const convertMap: CodeType = {
    encode(str: string): string {
      // 每4个hex字符加空格分组，便于阅读
      return utf8ToHex(str).replace(/(.{4})/g, '$1 ');
    },
    decode(str: string): string {
      return hexToUtf8(str);
    },
  };
  return convertMap[codeType](value);
};

const convertJSON = (value: string): Promise<string> =>
  new Promise((resolve, reject) => {
    try {
      const $json = jsonParse(value);
      resolve(jsonStringify($json, null, 2));
    } catch (error) {
      reject(error);
    }
  });

/**
 * 根据来源类型解码 + 目标类型编码，转换 payload
 * @param payload 原始payload字符串
 * @param currentType 目标类型
 * @param fromType 来源类型
 */
const convertPayload = async (
  payload: string,
  currentType: string,
  fromType: string
): Promise<string> => {
  let $payload = payload;
  // 先按来源类型解码为原始文本
  switch (fromType) {
    case 'Base64':
      $payload = convertBase64(payload, 'decode');
      break;
    case 'Hex':
      $payload = convertHex(payload, 'decode');
      break;
  }
  // 再按目标类型编码
  if (currentType === 'Base64') {
    $payload = convertBase64($payload, 'encode');
  }
  if (currentType === 'JSON' || currentType === 'CBOR') {
    $payload = await convertJSON($payload);
  }
  if (currentType === 'Hex') {
    $payload = convertHex($payload, 'encode');
  }
  return $payload;
};

export default convertPayload;
