const opt = Object.prototype.toString;

export function isArray(obj: any): obj is any[] {
  return opt.call(obj) === '[object Array]';
}

export function inArray(arr: any,val:any):boolean {
    for (var i = 0; i < arr.length; i++) {
      if (arr[i] === val) {
          return true;
      }
    }
    return false;
}

export function isObject(obj: any): obj is { [key: string]: any } {
  return opt.call(obj) === '[object Object]';
}

export function isString(obj: any): obj is string {
  return opt.call(obj) === '[object String]';
}

export function isNumber(obj: any): obj is number {
  return opt.call(obj) === '[object Number]' && obj === obj; // eslint-disable-line
}

export function isRegExp(obj: any) {
  return opt.call(obj) === '[object RegExp]';
}

export function isFile(obj: any): obj is File {
  return opt.call(obj) === '[object File]';
}

export function isBlob(obj: any): obj is Blob {
  return opt.call(obj) === '[object Blob]';
}

export function isUndefined(obj: any): obj is undefined {
  return obj === undefined;
}

export function isNull(obj: any): obj is null {
  return obj === null;
}

export function isFunction(obj: any): obj is (...args: any[]) => any {
  return typeof obj === 'function';
}

export function isEmptyObject(obj: any): boolean {
  return isObject(obj) && Object.keys(obj).length === 0;
}

export function isExist(obj: any): boolean {
  return obj || obj === 0;
}
/*
* @三目运算符 valbloom是判断布尔值，yesstr对正确值，otherstr不等于的值
*/
export function threeOperation(valbloom: any,yesstr:string,otherstr:string): string {
  return valbloom?yesstr:otherstr;
}

export function isWindow(el: any): el is Window {
  return el === window;
}
//判断是否为https地址
export function isUrl(path: string): boolean {
  const reg =
    /^(((^https?:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+(?::\d+)?|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?(\/#\/)?(?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[\w]*))?)$/;
  return reg.test(path);
}
//判断是否为域名包含端口
export function isAllUrl(path: string): boolean {
  const reg =/^(http(s)?:\/\/)\w+[^\s]+(\.[^\s]+){1,}$/;
  return reg.test(path);
}
