<template>
  <div class="home">
    <h1>欢迎来到博客系统</h1>
    <p>系统正在初始化中...</p>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'

onMounted(() => {
  console.log('Home page mounted')
})
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
  }
}
</style>

