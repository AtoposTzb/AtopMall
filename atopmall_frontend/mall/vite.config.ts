import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src')
    }
  },
  server: {
    port: 3000,
    open: false,
    proxy: {
      '/u/v1': {
        target: 'http://192.168.1.106:8000',
        changeOrigin: true
      },
      '/g/v1': {
        target: 'http://192.168.1.106:8000',
        changeOrigin: true
      },
      '/o/v1': {
        target: 'http://192.168.1.106:8000',
        changeOrigin: true
      },
      '/op/v1': {
        target: 'http://192.168.1.106:8000',
        changeOrigin: true
      },
      '/oss/v1': {
        target: 'http://192.168.1.106:8000',
        changeOrigin: true
      }
    }
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@use "@/assets/styles/variables.scss" as *;`
      }
    }
  }
})