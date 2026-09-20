// axios配置  可自行根据项目进行更改，只需更改该文件即可，其他文件可以不动
// The axios configuration can be changed according to the project, just change the file, other files can be left unchanged

import type { AxiosResponse } from 'axios';
import { clone } from 'lodash-es';
import type { RequestOptions, Result } from './types';
import type { AxiosTransform, CreateAxiosOptions } from './axiosTransform';
import { VAxios } from './Axios';
import { checkStatus } from './checkStatus';
import { RequestEnum, ResultEnum, ContentTypeEnum } from './httpEnum';
import { isString } from '@/utils/is';
import { getToken } from '@/utils/auth';
import { setObjToUrlParams, deepMerge } from '@/utils';
import { joinTimestamp, formatRequestDate } from './helper';
import useUserStore from '@/store/modules/user';
import { AxiosRetry } from '@/utils/http/axiosRetry';
import md5 from 'md5'

import { Message,Modal,Notification} from '@arco-design/web-vue';
const Main_url = import.meta.env.VITE_HTTP_URL;//域名前缀
const Root_url = Main_url;//域名前缀(可加模块)
type Recordable<T = any> = Record<string, T>;
/** 
 * @description: 数据处理，方便区分多种处理方式
 */
const transform: AxiosTransform = {
  /**
   * @description: 处理响应数据。如果数据不是预期格式，可直接抛出错误
   */
  transformResponseHook: (res: AxiosResponse<Result>, options: RequestOptions) => {
    const { isTransformResponse, isReturnNativeResponse } = options;
    // 是否返回原生响应头 比如：需要获取响应头时使用该属性
    if (isReturnNativeResponse) {
      return res;
    }
    // 不进行任何处理，直接返回
    // 用于页面代码可能需要直接获取code，data，message这些信息时开启
    if (!isTransformResponse) {
      return res.data;
    }
    // 错误的时候返回

    const resdata = res.data;
    if (!resdata) {
      // return '[HTTP] Request has no return value';
      throw new Error('请求出错，请稍候重试！');
    }
    //  这里 code，data，message为 后台统一的字段，需要在 types.ts内修改为项目自己的接口返回格式
    const { code, data, message,token } = resdata;
    //如果token即将过去则刷新token
    // 这里逻辑可以根据项目进行修改
    const hasSuccess = resdata && Reflect.has(resdata, 'code') && code === ResultEnum.SUCCESS;
    if (hasSuccess) {
        if (options.successMessageMode === 'modal') {
          Modal.error({ title: "提交提示", content: resdata.message });
        } else if (options.successMessageMode === 'message') {
          Message.success(resdata.message);
        }
        //请求成功后更新token
        if(token){
          const userStore = useUserStore();
          userStore.setTokenArr(token);
        }
        return data;
    }
    // 在此处根据自己项目的实际情况对不同的code执行不同的操作
    // 如果不希望中断当前请求，请return数据，否则直接抛出异常即可
    let timeoutMsg = '';
    switch (code) {
      case ResultEnum.TIMEOUT:
        options.errorMessageMode = 'message'//第一错误提示
        timeoutMsg = '登录超时,请重新登录!';
        const userStore = useUserStore();
        userStore.clearloginfo();
        break;
      default:
        if (message) {
          timeoutMsg = message;
        }
    }
    if (options.errorMessageMode === 'modal') {
      Modal.error({ title: '错误提示', content: timeoutMsg });
    } else if (options.errorMessageMode === 'message') {
      Message.error({content:timeoutMsg,id:"errmsg"});
    }else if(options.errorMessageMode === 'notification') {
      Notification.error({
        id: 'menuNotice',
        title: '错误提示',
        content: timeoutMsg,
        position: 'bottomRight'
      });
    }
    throw new Error(timeoutMsg || '请求出错，请稍候重试！');
  },

  // 请求之前处理config
  beforeRequestHook: (config, options) => {
    const { apiUrl,isRootUrl, joinPrefix, joinParamsToUrl, formatDate, joinTime = true, urlPrefix } = options;
    if (joinPrefix) {
      config.url = `${urlPrefix}${config.url}`;
    }
    //拼接域名
    if (config.url&&config.url.startsWith("./")) {
      const urlarr=config.url.split("./")
      config.url = `${Main_url}/${urlarr[1]}`;
    }else if (isRootUrl&&apiUrl && isString(apiUrl)) {
      config.url = `${apiUrl}${config.url}`;
    }
    const params = config.params || {};
    const data = config.data || false;
    formatDate && data && !isString(data) && formatRequestDate(data);
    if (config.method?.toUpperCase() === RequestEnum.GET) {
      if (!isString(params)) {
        // 给 get 请求加上时间戳参数，避免从缓存中拿数据。
        config.params = Object.assign(params || {}, joinTimestamp(joinTime, false));
      } else {
        // 兼容restful风格
        config.url = config.url + params + `${joinTimestamp(joinTime, true)}`;
        config.params = undefined;
      }
    } else {
      if (!isString(params)) {
        formatDate && formatRequestDate(params);
        if (Reflect.has(config, 'data') && config.data && (Object.keys(config.data).length > 0 || config.data instanceof FormData)) {
          config.data = data;
          config.params = params;
        } else {
          // 非GET请求如果没有提供data，则将params视为data
          config.data = params;
          config.params = undefined;
        }
        if (joinParamsToUrl) {
          config.url = setObjToUrlParams(
            config.url as string,
            Object.assign({}, config.params, config.data),
          );
        }
      } else {
        // 兼容restful风格
        config.url = config.url + params;
        config.params = undefined;
      }
    }
    return config;
  },

  /**
   * @description: 请求拦截器处理
   */
  requestInterceptors: (config, options) => {
    // 请求之前处理config
    const token = getToken();
    if (token && (config as Recordable)?.requestOptions?.withToken !== false) {
      // jwt token
      (config as Recordable).headers.Authorization = options.authenticationScheme
        ? `${options.authenticationScheme} ${token}`
        : token;
    }
    //接口验证
    const timestamp: number = Date.parse(new Date().toString())/1000;
    const passstr=md5(import.meta.env.VITE_ENCRYPT+timestamp) as any
    (config as Recordable).headers["apiverify"]=window.btoa(passstr+"#"+timestamp);//MD5加密后#拼接时间戳
    //设置语言
    (config as Recordable).headers["locale"]=localStorage.getItem('arco-locale')||"zh-CN";
    return config;
  },

  /**
   * @description: 响应拦截器处理
   */
  responseInterceptors: (res: AxiosResponse<any>) => {
    return res;
  },

  /**
   * @description: 响应错误处理
   */
  responseInterceptorsCatch: (axiosInstance: AxiosResponse, error: any) => {
    const { response, code, message, config } = error || {};
    const errorMessageMode = config?.requestOptions?.errorMessageMode || 'none';
    const msg: string = (response?.data?.message) ||(response?.data?.error?.message ?? '');
    const err: string = error?.toString?.() ?? '';
    let errMessage = '';

    try {
      if (code === 'ECONNABORTED' && message.indexOf('timeout') !== -1) {
        errMessage = '接口请求超时,请刷新页面重试!';
      }
      if (err?.includes('Network Error')) {
        errMessage = '网络异常，请检查您的网络连接是否正常(后端启动、跨域等)！';
      }

      if (errMessage) {
        if (errorMessageMode === 'modal') {
          Modal.error({ title: '错误提示', content: errMessage });
        } else if (errorMessageMode === 'message') {
          Message.error(errMessage);
        }
        return Promise.reject(error);
      }
    } catch (error) {
      throw new Error(error as unknown as string);
    }

    checkStatus(error?.response?.status, msg, errorMessageMode);

    // 添加自动重试机制 保险起见 只针对GET请求
    const retryRequest = new AxiosRetry();
    const { isOpenRetry } = config.requestOptions.retryRequest;
    config.method?.toUpperCase() === RequestEnum.GET &&
      isOpenRetry &&
      // @ts-ignore
      retryRequest.retry(axiosInstance, error);
    return Promise.reject(error);
  },
};

function createAxios(opt?: Partial<CreateAxiosOptions>) {
  return new VAxios(
    // 深度合并
    deepMerge(
      {
        authenticationScheme: '',
        // timeout: 10 * 1000,//超时
        // 基础接口地址
        // baseURL: globSetting.apiUrl,
        headers: { 'Content-Type': ContentTypeEnum.JSON },
        // 如果是form-data格式
        // headers: { 'Content-Type': ContentTypeEnum.FORM_URLENCODED },
        // 数据处理方式
        transform: clone(transform),
        // 配置项，下面的选项都可以在独立的接口请求中覆盖
        requestOptions: {
          // 默认将prefix 添加到url
          joinPrefix: true,
          // 是否返回原生响应头 比如：需要获取响应头时使用该属性
          isReturnNativeResponse: false,
          // 需要对返回数据进行处理
          isTransformResponse: true,
           // 是否添加接口地址
           isRootUrl: true,
          // post请求的时候添加参数到url
          joinParamsToUrl: false,
          // 格式化提交参数时间
          formatDate: true,
          // 消息提示类型
          errorMessageMode: 'message',
          // 接口地址
          apiUrl: Root_url,
          // 接口拼接地址
          urlPrefix: "",
          //  是否加入时间戳
          joinTime: true,
          // 忽略重复请求
          ignoreCancelToken: true,
          // 是否携带token
          withToken: true,
          retryRequest: {
            isOpenRetry: true,
            count: 5,
            waitTime: 100,
          },
        },
      },
      opt || {},
    ),
  );
}
export const defHttp = createAxios();

