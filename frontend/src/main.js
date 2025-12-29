import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import i18n from './locales'
import App from './App.vue'

// 样式
import './assets/styles/reset.less'
import './assets/styles/common.less'
import './assets/styles/element-theme.less'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(i18n)

app.mount('#app')

