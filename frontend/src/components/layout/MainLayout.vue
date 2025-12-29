<template>
  <div class="main-layout">
    <Header />
    <div class="main-container">
      <router-view />
    </div>
    <Footer />
  </div>
</template>

<script setup>
import Header from './Header.vue'
import Footer from './Footer.vue'
</script>

<style scoped lang="less">
.main-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.main-container {
  flex: 1;
  width: 100%;
}
</style>

