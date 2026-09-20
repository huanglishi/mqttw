import { createPinia } from 'pinia';
import useAppStore from './modules/app';
import useUserStore from './modules/user';
import useMqttStore from './modules/mqtt';
import useStressStore from './modules/stress';

const pinia = createPinia();

export { useAppStore,useUserStore,useMqttStore,useStressStore };
export default pinia;
