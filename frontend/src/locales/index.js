import { createI18n } from 'vue-i18n'
import zhCN from './zh-CN.json'
import enUS from './en-US.json'

// Element Plus 国际化
import zhCnLocale from 'element-plus/dist/locale/zh-cn.mjs'
import enLocale from 'element-plus/dist/locale/en.mjs'

const messages = {
  'zh-CN': {
    ...zhCN,
    el: zhCnLocale.el
  },
  'en-US': {
    ...enUS,
    el: enLocale.el
  }
}

// 获取默认语言
const getDefaultLocale = () => {
  const savedLanguage = localStorage.getItem('language')
  if (savedLanguage && messages[savedLanguage]) {
    return savedLanguage
  }
  
  // 从浏览器获取语言
  const browserLanguage = navigator.language || navigator.userLanguage
  if (browserLanguage.startsWith('zh')) {
    return 'zh-CN'
  }
  return 'en-US'
}

const i18n = createI18n({
  legacy: false, // 使用 Composition API 模式
  locale: getDefaultLocale(),
  fallbackLocale: 'zh-CN',
  messages
})

export default i18n

