import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig(({mode}) => {
  const rootDir = path.resolve(__dirname, '../');//1 level up
  const env = loadEnv(mode, rootDir, '');
  //const backendUrl = import.meta.env.VITE_BACKEND_URL || 'http://127.27.27.27:8000';
  const backendUrl = env.VITE_BACKEND_URL
  return {
    plugins: [react()],
    server: {
      port: 5173,
      proxy: {
      '/api': {
        target: backendUrl,
        changeOrigin: true,
        },
      },
    },
  };
});