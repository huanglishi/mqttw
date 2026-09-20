const TOKEN_KEY = 'biztoken';
const CODESTORE_TOKEN_KEY = 'codestoretoken';

const isLogin = () => {
  return !!localStorage.getItem(TOKEN_KEY);
};

const getToken = () => {
  return localStorage.getItem(TOKEN_KEY);
};

const setToken = (token: string| undefined) => {
  if(token){
    localStorage.setItem(TOKEN_KEY, token);
  }else{
    localStorage.removeItem(TOKEN_KEY);
  }
};

const clearToken = () => {
  localStorage.removeItem(TOKEN_KEY);
};
//代码仓库登录
const cstoreIsLogin = () => {
  return !!localStorage.getItem(CODESTORE_TOKEN_KEY);
};
const getCstoreToken = () => {
  return localStorage.getItem(CODESTORE_TOKEN_KEY);
};
const setCstoreToken = (token: string| undefined) => {
  if(token){
    localStorage.setItem(CODESTORE_TOKEN_KEY, token);
  }else{
    localStorage.removeItem(CODESTORE_TOKEN_KEY);
  }
};
const clearCstoreToken = () => {
  localStorage.removeItem(CODESTORE_TOKEN_KEY);
};

export { isLogin, getToken, setToken, clearToken,
  cstoreIsLogin,getCstoreToken,setCstoreToken,clearCstoreToken};
