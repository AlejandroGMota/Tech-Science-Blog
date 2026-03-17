import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import path from 'path'

export default defineConfig({
  plugins: [react()],
  base: '/admin/',
  root: '.',
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
    },
  },
  build: {
    outDir: '../backend/admin-dist',
    emptyOutDir: true,
    rollupOptions: {
      input: path.resolve(__dirname, 'admin.html'),
    },
  },
  server: {
    port: 5174,
    proxy: {
      '/api': 'http://localhost:8080',
    },
  },
})
