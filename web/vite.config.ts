import path from 'node:path';
import { defineConfig } from 'vite';
import vue from '@vitejs/plugin-vue';

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src')
    }
  },
  build: {
    outDir: '../assets/web',
    emptyOutDir: true,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (!id.includes('node_modules')) {
            return undefined;
          }

          if (id.includes('vuetify/components')) {
            return 'vuetify-components';
          }

          if (id.includes('vuetify/directives')) {
            return 'vuetify-directives';
          }

          if (id.includes('vuetify/iconsets')) {
            return 'vuetify-icons';
          }

          if (id.includes('vuetify')) {
            return 'vuetify-core';
          }

          if (id.includes('@mdi') || id.includes('materialdesignicons')) {
            return 'mdi';
          }

          if (
            id.includes('/vue/') ||
            id.includes('/vue-router/') ||
            id.includes('/pinia/')
          ) {
            return 'vue-core';
          }

          return 'vendor';
        }
      }
    }
  },
  server: {
    proxy: {
      '/oauth': {
        target: 'http://localhost:9900',
        changeOrigin: true
      },
      '/api': {
        target: 'http://localhost:9900',
        changeOrigin: true,
        ws: true
      }
    }
  }
});
