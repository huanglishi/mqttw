
import { join } from 'lodash-es';
import { isObject,isArray } from '@/utils/is';

/** @desc 问候 */
export function goodTimeText() {
    const time = new Date()
    const hour = time.getHours()
    return hour < 9 ? '早上好' : hour <= 11 ? '上午好' : hour <= 13 ? '中午好' : hour <= 18 ? '下午好' : '晚上好'
  }
  /**
 * Add the object as a parameter to the URL
 * @param baseUrl url
 * @param obj
 * @returns {string}
 * eg:
 *  let obj = {a: '3', b: '4'}
 *  setObjToUrlParams('www.baidu.com', obj)
 *  ==>www.baidu.com?a=3&b=4
 */
export function setObjToUrlParams(baseUrl: string, obj: any): string {
  let parameters = '';
  for (const key in obj) {
    parameters += key + '=' + encodeURIComponent(obj[key]) + '&';
  }
  parameters = parameters.replace(/&$/, '');
  return /\?$/.test(baseUrl) ? baseUrl + parameters : baseUrl.replace(/\/?$/, '?') + parameters;
}

// 深度合并
export function deepMerge<T = any>(src: any = {}, target: any = {}): T {
  let key: string;
  for (key in target) {
    src[key] = isObject(src[key]) ? deepMerge(src[key], target[key]) : (src[key] = target[key]);
  }
  return src;
}

/**
 *  过滤请求参数
 */
export function FilterParam(params: any): any {
  for(let key in params){
      if(isArray(params[key])){//将数组转,拼接字符串
          params[key]=join(params[key])
      }
  }
  return params;
}
// 颜色格式 hex 转 rgba
export const hexToRgba = (bgColor:string) =>  {
    let color = bgColor.slice(1);   // 去掉'#'号
    let rgba = [
        parseInt('0x'+color.slice(0, 2)),
        parseInt('0x'+color.slice(2, 4)),
        parseInt('0x'+color.slice(4, 6)),
        0.1
    ];
    return 'rgba(' + rgba.toString() + ')';
}