import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import i18n from './locales'
import App from './App.vue'
import { setupLazyLoadDirective } from './utils/lazyLoad'
import performanceMonitor from './utils/performance'

// 样式
import './assets/styles/reset.less'
import './assets/styles/common.less'
import './assets/styles/element-theme.less'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(i18n)

// 注册懒加载指令
setupLazyLoadDirective(app)

// 初始化性能监控
performanceMonitor.init()

app.mount('#app')

