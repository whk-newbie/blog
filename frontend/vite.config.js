import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'

export default defineConfig({
  plugins: [
    vue(),
    AutoImport({
      resolvers: [ElementPlusResolver()],
      imports: ['vue', 'vue-router', 'pinia'],
      dts: process.env.NODE_ENV === 'development' ? 'src/auto-imports.d.ts' : false,
    }),
    Components({
      resolvers: [ElementPlusResolver()],
      dts: process.env.NODE_ENV === 'development' ? 'src/components.d.ts' : false,
    }),
  ],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
    },
  },
  server: {
    port: 5173,
    host: '0.0.0.0',
    proxy: {
      '/api': {
        // 在Docker环境中使用服务名backend-dev，本地开发可通过环境变量覆盖为127.0.0.1
        target: process.env.VITE_API_PROXY_TARGET || 'http://backend-dev:8080',
        changeOrigin: true,
        rewrite: (path) => path,
      },
    },
  },
  // 优化依赖预构建
  optimizeDeps: {
    include: ['vue', 'vue-router', 'pinia', 'element-plus'],
    esbuildOptions: {
      // 禁用 source map 以减少警告
      sourcemap: false,
    },
  },
  // 开发环境 source map 配置
  css: {
    devSourcemap: false,
  },
  build: {
    outDir: 'dist',
    sourcemap: false,
    minify: 'terser', // 使用terser进行代码混淆
    terserOptions: {
      compress: {
        drop_console: true, // 生产环境移除console
        drop_debugger: true, // 移除debugger
        pure_funcs: ['console.log', 'console.info'], // 移除指定的函数调用
      },
      format: {
        comments: false, // 移除注释
      },
      mangle: {
        toplevel: true, // 混淆顶级作用域
        properties: {
          regex: /^_/, // 混淆以下划线开头的属性
        },
      },
    },
    rollupOptions: {
      output: {
        manualChunks: {
          'element-plus': ['element-plus'],
          'vue-vendor': ['vue', 'vue-router', 'pinia'],
        },
        // 混淆文件名
        chunkFileNames: 'js/[name]-[hash].js',
        entryFileNames: 'js/[name]-[hash].js',
        assetFileNames: 'assets/[name]-[hash].[ext]',
      },
    },
  },
})

