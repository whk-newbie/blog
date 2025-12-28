import { createApp } from 'vue'
import { createPinia } from 'pinia'
import router from './router'
import App from './App.vue'

// 样式
import './assets/styles/reset.less'
import './assets/styles/common.less'

const app = createApp(App)

app.use(createPinia())
app.use(router)

app.mount('#app')

