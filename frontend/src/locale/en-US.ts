//1.自动加载components组件的翻译文件
const autoLoclecomponents = import.meta.glob('@/components/**/locale/en-US.ts', {
  eager: true,
});
//2.自动加载locale下对应语言目录
const autoLocle = import.meta.glob('@/locale/en-US/**.ts', {
  eager: true,
});
//3.自动加载 views下的语言包
const autoLocleviews = import.meta.glob('@/views/**/locale/en-US.ts', {
  eager: true,
});
function formatModules(_modules: any, result: {}) {
  Object.keys(_modules).forEach((key) => {
    const defaultModule = _modules[key].default;
    if (!defaultModule) return;
    result=Object.assign({},result,defaultModule);
  });
  return result;
}
export default {
  ...formatModules(autoLoclecomponents,{}),
  ...formatModules(autoLocle,{}),
  ...formatModules(autoLocleviews,{})
};
