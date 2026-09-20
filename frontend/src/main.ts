import { createApp } from 'vue'
import Router from './router/index';
import store from './store';
import App from './App.vue'
import globalComponents from '@/components';
import i18n from './locale';
import '@/assets/style/global.less';

const app = createApp(App);
app.use(Router);
app.use(store);
app.use(globalComponents);
app.use(i18n);

app.mount('#app')
