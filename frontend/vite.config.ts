import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'

// https://vite.dev/config/
export default defineConfig({
  plugins: [
    vue(),
    vueJsx(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    host: true, // 允许外部访问 (代替了之前的 --host 参数)
    proxy: {
      // 当请求路径以 /api 开头时，Vite 会自动把它转发到后端的 8080 端口
      '/api': {
        target: 'http://127.0.0.1:8080', // 注意：这里的 127 是 Vite 在 Ubuntu 上去访问 Ubuntu 本地的后端，所以是对的！
        changeOrigin: true
      }
    }
  },
})
