<template>
  <div class="home">
    <h1>欢迎来到博客系统</h1>
    <p>系统正在初始化中...</p>
    <div class="actions">
      <el-button type="primary" size="large" @click="goToLogin">
        管理员登录
      </el-button>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

onMounted(() => {
  console.log('Home page mounted')
})

const goToLogin = () => {
  router.push('/admin/login')
}
</script>

<style lang="less" scoped>
.home {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  
  h1 {
    font-size: 2.5rem;
    margin-bottom: 1rem;
  }
  
  p {
    font-size: 1.2rem;
    color: #666;
    margin-bottom: 2rem;
  }
  
  .actions {
    margin-top: 2rem;
  }
}
</style>

