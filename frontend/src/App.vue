<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'

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

