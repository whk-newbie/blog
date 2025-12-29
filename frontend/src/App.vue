<template>
  <el-config-provider :locale="elementLocale">
    <div id="app">
      <router-view />
      <FloatingActions />
      <FingerprintCollector />
    </div>
  </el-config-provider>
</template>

<script setup>
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { elementPlusLocales } from './locales'
import FloatingActions from './components/common/FloatingActions.vue'
import FingerprintCollector from './components/fingerprint/FingerprintCollector.vue'

const { locale } = useI18n()

// Element Plus 的 locale
const elementLocale = computed(() => {
  return elementPlusLocales[locale.value] || elementPlusLocales['zh-CN']
})

onMounted(() => {
  console.log('Blog App Mounted')
})
</script>

<style lang="less">
#app {
  width: 100%;
  min-height: 100vh;
}
</style>

