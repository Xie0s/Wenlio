import path from 'node:path'
import { defineConfig } from 'vite'
import tailwindcss from '@tailwindcss/vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig(() => {
  const frontendOutDir = process.env.DOCMATHER_FRONTEND_OUT_DIR?.trim() || 'dist'

  return {
    plugins: [vue(), tailwindcss()],
    server: {
      port: 5173,
      proxy: {
        '/api': {
          target: 'http://localhost:1501',
          changeOrigin: true,
        },
        '/uploads': {
          target: 'http://localhost:1501',
          changeOrigin: true,
        },
        '/raw': {
          target: 'http://localhost:1501',
          changeOrigin: true,
        },
      },
    },
    build: {
      outDir: frontendOutDir,
      emptyOutDir: true,
      rollupOptions: {
        output: {
          entryFileNames: 'assets/[name]-[hash].js',
          chunkFileNames: 'assets/[name]-[hash].js',
          assetFileNames: 'assets/[name]-[hash][extname]',
        },
      },
    },
    resolve: {
      alias: {
        '@': path.resolve(__dirname, './src'),
      },
    },
  }
})
