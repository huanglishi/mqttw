import { defineStore } from 'pinia';
import { AppState } from './types';
//获取本地保存配置
const themeval=localStorage.getItem("theme");
const useAppStore = defineStore('app', {
  state: (): AppState => {
    return {
      sysTime:"",
      theme: themeval?themeval:'light',
    }
  },
  getters: {
    //获取主题
    getTheme():string{
      return this.theme;
    },
  },

  actions: {
    //设置主题
    toggleTheme(dark:boolean) {
      if (dark) {
        this.theme = 'dark';
        document.body.setAttribute('arco-theme', 'dark');
      } else {
        this.theme = 'light';
        document.body.removeAttribute('arco-theme');
      }
      localStorage.setItem("theme",this.theme)
    },
  },
});

export default useAppStore;
