import { defineConfig, loadEnv } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig(({mode}) => {
  const env = loadEnv(mode, process.cwd(), '');
  const backendUrl = env.VITE_BACKEND_URL || 'http://127.27.27.27:8000';
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