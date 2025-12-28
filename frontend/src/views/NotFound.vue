<template>
  <div class="not-found">
    <h1>404</h1>
    <p>页面未找到</p>
    <el-button type="primary" @click="goHome">返回首页</el-button>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

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
}
</style>

