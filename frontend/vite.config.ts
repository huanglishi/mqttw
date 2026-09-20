import { defineConfig } from "vite";
import vue from "@vitejs/plugin-vue";
import wails from "@wailsio/runtime/plugins/vite";
import { resolve } from 'path';
//按需加载与组件库主题
import { vitePluginForArco } from '@arco-plugins/vite-vue'
// https://vitejs.dev/config/
export default defineConfig({
  server: {
    host: "127.0.0.1",
    port: Number(process.env.WAILS_VITE_PORT) || 9245,
    strictPort: true,
  },
  plugins: [vue(), wails("./bindings"),vitePluginForArco({ style: 'css' })],
    //新增
  resolve: {
    alias: [
      {
        find: '@',
        replacement: resolve(import.meta.dirname, './src'),
      },
      {
        find: '/@',
        replacement: resolve(import.meta.dirname, './src'),
      },
      {
        find: '/#',
        replacement: resolve(import.meta.dirname, './bindings'),
      },
    ],
    extensions: ['.ts', '.js'],
  },
  css: {
    preprocessorOptions: {
      less: {
        modifyVars: {
          'arcoblue-6': "#165DFF",//
        },
        javascriptEnabled: true,
      },
    },
  },
});
