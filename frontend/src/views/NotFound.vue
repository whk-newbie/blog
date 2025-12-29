<template>
  <div class="not-found">
    <h1>404</h1>
    <p>{{ t('app.notFound') }}</p>
    <div class="button-group">
      <el-button @click="goBack">{{ t('common.back') }}</el-button>
      <el-button type="primary" @click="goHome">{{ t('app.backToHome') }}</el-button>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const { t } = useI18n()

const goBack = () => {
  router.back()
}

const goHome = () => {
  router.push('/')
}
</script>

<style lang="less" scoped>
.not-found {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  
  h1 {
    font-size: 6rem;
    margin-bottom: 1rem;
    color: #409eff;
  }
  
  p {
    font-size: 1.5rem;
    margin-bottom: 2rem;
    color: #666;
  }
  
  .button-group {
    display: flex;
    gap: 1rem;
  }
}
</style>

