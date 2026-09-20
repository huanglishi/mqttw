import { App } from 'vue';
import Chart from './chart/index.vue';
import collapseTransition from './collapseTransition/index.vue';
import {Icon} from '@/components/Icon';

export default {
  install(Vue: App) {
    Vue.component('IconFont', Icon);
    Vue.component('Chart', Chart);
    Vue.component('CollapseTransition', collapseTransition);
  },
};
