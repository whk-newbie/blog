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
}

.layout-body {
  flex: 1;
  display: flex;
  max-width: 1100px;
  margin: 0 auto;
  width: 100%;
  padding: var(--spacing-lg) var(--spacing-md);
  gap: var(--spacing-lg);
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(4px);
}

html.dark .layout-body {
  background: rgba(26, 26, 46, 0.85);
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
