import legacyPlugin from '@vitejs/plugin-legacy'
import { viteLogo } from './src/core/config'
import Banner from 'vite-plugin-banner'
import * as path from 'path'
import { loadEnv } from 'vite'
import vuePlugin from '@vitejs/plugin-vue'
import vueDevTools from 'vite-plugin-vue-devtools'
import VueFilePathPlugin from './vitePlugin/componentName/index.js'
import { svgBuilder } from 'vite-auto-import-svg'
import vueRootValidator from 'vite-check-multiple-dom'
import { AddSecret } from './vitePlugin/secret'
import UnoCSS from '@unocss/vite'

// @see https://cn.vitejs.dev/config/
export default ({ mode }) => {
  AddSecret('')
  const env = loadEnv(mode, process.cwd())
  viteLogo(env)

  const timestamp = Date.parse(new Date())
  const isProduction = mode === 'production'
  const enableLegacy = env.VITE_LEGACY === 'true'
  const enableRootValidator = !isProduction && env.VITE_CHECK_MULTIPLE_DOM !== 'false'

  const optimizeDeps = {}

  const alias = {
    '@': path.resolve(__dirname, './src'),
    vue$: 'vue/dist/vue.runtime.esm-bundler.js'
  }

  const esbuild = isProduction
    ? {
        drop: ['console', 'debugger']
      }
    : {}

  const rollupOptions = {
    output: {
      entryFileNames: 'assets/087AC4D233B64EB0[name].[hash].js',
      chunkFileNames: 'assets/087AC4D233B64EB0[name].[hash].js',
      assetFileNames: 'assets/087AC4D233B64EB0[name].[hash].[ext]'
    }
  }

  const base = '/'
  const root = './'
  const outDir = 'dist'

  const config = {
    base: base, // 编译后js导入的资源路径
    root: root, // index.html文件所在位置
    publicDir: 'public', // 静态资源文件夹
    resolve: {
      alias
    },
    css: {
      preprocessorOptions: {
        scss: {
          api: 'modern-compiler' // or "modern"
        }
      }
    },
    server: {
      // 如果使用docker-compose开发模式，设置为false
      open: true,
      port: Number(env.VITE_CLI_PORT),
      proxy: {
        // 把key的路径代理到target位置
        // detail: https://cli.vuejs.org/config/#devserver-proxy
        [env.VITE_BASE_API]: {
          // 需要代理的路径   例如 '/api'
          target: `${env.VITE_BASE_PATH}:${env.VITE_SERVER_PORT}/`, // 代理到 目标路径
          changeOrigin: true,
          rewrite: (path) =>
            path.replace(new RegExp('^' + env.VITE_BASE_API), '')
        },
        '/plugin': {
          // 需要代理的路径   例如 '/api'
          target: `https://plugin.gin-vue-admin.com/api/`, // 代理到 目标路径
          changeOrigin: true,
          rewrite: (path) =>
            path.replace(new RegExp('^/plugin'), '')
        }
      }
    },
    build: {
      minify: 'esbuild', // esbuild 比 terser 快很多，适合后台系统生产构建
      manifest: false, // 是否产出manifest.json
      sourcemap: false, // 是否产出sourcemap.json
      reportCompressedSize: false, // 跳过 gzip/brotli 体积统计，减少构建等待
      outDir: outDir, // 产出目录
      rollupOptions
    },
    esbuild,
    optimizeDeps,
    plugins: [
      env.VITE_POSITION === 'open' &&
      vueDevTools({ launchEditor: env.VITE_EDITOR }),
      enableLegacy && legacyPlugin({
        targets: [
          'Android > 39',
          'Chrome >= 60',
          'Safari >= 10.1',
          'iOS >= 10.3',
          'Firefox >= 54',
          'Edge >= 15'
        ]
      }),
      vuePlugin(),
      svgBuilder(['./src/plugin/', './src/assets/icons/'], base, outDir, 'assets', mode),
      // [Banner(`\n Build based on gin-vue-admin \n Time : ${timestamp}`)],
      VueFilePathPlugin('./src/pathInfo.json'),
      UnoCSS(),
      enableRootValidator && vueRootValidator()
    ].filter(Boolean)
  }
  return config
}
