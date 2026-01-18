import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],
  server: {
    proxy: {
      '/oauth': {
        target: 'http://localhost:9900',
        changeOrigin: true
      },
      '/api': {
        target: 'http://localhost:9900',
        changeOrigin: true
      }
    }
  }
});
