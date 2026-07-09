<template>
  <div class="main-layout">
    <Header @toggleSearch="searchOverlay?.open()" />
    <div class="layout-body">
      <Sidebar />
      <main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
    <Footer />
    <SearchOverlay ref="searchOverlay" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Header from './Header.vue'
import Sidebar from './PublicSidebar.vue'
import Footer from './Footer.vue'
import SearchOverlay from '@/components/common/SearchOverlay.vue'

const searchOverlay = ref(null)
</script>

<style scoped lang="less">
.main-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--bg-color);
}

.layout-body {
  flex: 1;
  display: flex;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  padding: var(--spacing-lg) var(--spacing-md);
  gap: var(--spacing-lg);
}

.main-content {
  flex: 1;
  min-width: 0;
  padding-top: var(--header-height);
}

// 页面过渡动画
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.15s ease;
}
.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

// 响应式
@media (max-width: 768px) {
  .layout-body {
    flex-direction: column;
    padding: var(--spacing-md);
  }
}
</style>
