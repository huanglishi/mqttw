import { createI18n } from 'vue-i18n';
import en from './en-US';
import cn from './zh-CN';
import tw from './zh-TW';

export const LOCALE_OPTIONS = [
  { label: '简体中文', value: 'zh-CN' },
  { label: '繁體中文', value: 'zh-TW' },
  { label: 'English', value: 'en-US' },
];
const defaultLocale = localStorage.getItem('arco-locale') || 'zh-CN';

const i18n = createI18n({
  legacy: false, 
  locale: defaultLocale,
  fallbackLocale: 'zh-CN',//默认中文
  // allowComposition: true,
  messages: {
    'en-US': en,
    'zh-CN': cn,
    'zh-TW': tw,
  },
});

export default i18n;
