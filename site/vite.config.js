import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// 官网独立应用，开发期将 /pcdn 代理到 GVA 后端
export default defineConfig({
  plugins: [vue()],
  base: './',
  server: {
    port: 5174,
    proxy: {
      '/pcdn': {
        target: 'http://127.0.0.1:8888',
        changeOrigin: true
      }
    }
  }
})
