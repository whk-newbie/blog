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
        // 禁用可能导致循环依赖的优化
        passes: 1, // 减少压缩轮数
        unsafe: false, // 禁用不安全的优化
        unsafe_comps: false,
        unsafe_math: false,
        unsafe_methods: false,
        unsafe_proto: false,
        unsafe_regexp: false,
        unsafe_undefined: false,
      },
      format: {
        comments: false, // 移除注释
      },
      mangle: {
        // 禁用顶级作用域混淆，避免循环依赖问题
        toplevel: false,
        // 保留类名和函数名，避免初始化顺序问题
        keep_classnames: true,
        keep_fnames: true,
        properties: {
          // 禁用属性混淆
          regex: /^$/,
        },
      },
    },
    rollupOptions: {
      output: {
        manualChunks: (id) => {
          // 更精细的代码分割策略，确保加载顺序
          if (id.includes('node_modules')) {
            if (id.includes('element-plus')) {
              return 'element-plus'
            }
            if (id.includes('vue') || id.includes('vue-router') || id.includes('pinia')) {
              return 'vue-vendor'
            }
            // 其他第三方库
            return 'vendor'
          }
        },
        // 混淆文件名
        chunkFileNames: 'js/[name]-[hash].js',
        entryFileNames: 'js/[name]-[hash].js',
        assetFileNames: 'assets/[name]-[hash].[ext]',
      },
    },
  },
})

