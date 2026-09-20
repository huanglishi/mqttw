import { defineStore } from 'pinia';
import { UserState} from './types';
import { setToken,getToken, clearToken } from '@/utils/auth';
const useUserStore = defineStore('user', {
  state:(): UserState => ({
    name: undefined,
    nickname: undefined,
    mobile: undefined,
    email: undefined,
    avatar: undefined,
    id: 0,
    city: '',
    company: '',
    sessionTimeout: false,//登录是否已过期
    createtime:"",
  }),

  getters: {
    userInfo(state: UserState): UserState {
      return { ...state };
    },
  },

  actions: {
    setTokenData(info: string | undefined) {
      setToken(info);
    },
    setSessionTimeout(flag: boolean) {
      this.sessionTimeout = flag;
    },
    // Set user's information
    setInfo(partial: Partial<UserState>) {
      this.$patch(partial);
    },
    // Reset user's information
    resetInfo() {
      this.$reset();
    },
  
   //设置token
   async setTokenArr(token:any) {
      setToken(token);
    },
    // 这里请求用户信息数据
    async info() {
      // const res = await getUserInfo();
      // this.setInfo(res);
    },

    // Login
    async login() {
      setToken(undefined);
    },
    // Logout
    async logout(goLogin = false) {
    },
     //清除登录信息
     clearloginfo() {
      clearToken();
     }
  },
});

export default useUserStore;
