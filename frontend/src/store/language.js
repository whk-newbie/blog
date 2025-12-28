import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useLanguageStore = defineStore('language', () => {
  // 从localStorage读取语言，默认为zh-CN
  const language = ref(localStorage.getItem('language') || 'zh-CN')

  // 设置语言
  const setLanguage = (newLanguage) => {
    language.value = newLanguage
    localStorage.setItem('language', newLanguage)
  }

  // 切换语言
  const toggleLanguage = () => {
    const newLanguage = language.value === 'zh-CN' ? 'en-US' : 'zh-CN'
    setLanguage(newLanguage)
  }

  return {
    language,
    setLanguage,
    toggleLanguage
  }
})

