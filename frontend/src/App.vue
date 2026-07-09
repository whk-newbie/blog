<template>
  <el-config-provider :locale="elementLocale">
    <div id="app">
      <!-- 加密通道加载中 -->
      <div v-if="!encryptionReady" class="app-loading">
        <div class="loading-content">
          <div class="loading-spinner"></div>
          <p class="loading-text">{{ loadingText }}</p>
          <el-button v-if="loadError" type="primary" size="small" @click="retryInit">
            重试
          </el-button>
        </div>
      </div>
      <!-- 正常页面 -->
      <template v-else>
        <router-view />
        <FloatingActions />
        <FingerprintCollector />
      </template>
    </div>
  </el-config-provider>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { elementPlusLocales } from './locales'
import FloatingActions from './components/common/FloatingActions.vue'
import FingerprintCollector from './components/fingerprint/FingerprintCollector.vue'
import { initEncryption } from '@/api/http'

const { locale } = useI18n()

const elementLocale = computed(() => {
  return elementPlusLocales[locale.value] || elementPlusLocales['zh-CN']
})

const encryptionReady = ref(false)
const loadingText = ref('正在建立安全连接...')
const loadError = ref(false)

async function initApp() {
  loadError.value = false
  loadingText.value = '正在建立安全连接...'
  try {
    await initEncryption()
    encryptionReady.value = true
  } catch (e) {
    console.error('Encryption init failed:', e)
    loadingText.value = '安全连接失败：' + (e.message || '未知错误')
    loadError.value = true
  }
}

function retryInit() {
  initApp()
}

onMounted(() => {
  initApp()
})
</script>

<style lang="less">
#app {
  width: 100%;
  min-height: 100vh;
}

.app-loading {
  position: fixed;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-color, #fafafa);
  z-index: 9999;

  .loading-content {
    text-align: center;
  }

  .loading-spinner {
    width: 40px;
    height: 40px;
    margin: 0 auto 20px;
    border: 3px solid var(--border-color, #e8e8e8);
    border-top-color: var(--primary-color, #4078c0);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }

  .loading-text {
    color: var(--text-secondary, #888);
    font-size: 14px;
  }
}
</style>
