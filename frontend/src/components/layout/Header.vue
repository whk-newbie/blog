<template>
  <header class="site-header">
    <div class="header-container">
      <div class="logo">
        <router-link to="/">
          <h1>{{ blogTitle }}</h1>
        </router-link>
      </div>
      <nav class="nav-menu">
        <router-link to="/" class="nav-item">{{ t('nav.home') }}</router-link>
        <router-link to="/articles" class="nav-item">{{ t('nav.articles') }}</router-link>
        <router-link v-if="isLoggedIn" to="/admin" class="nav-item admin-link">
          <el-icon><Setting /></el-icon>
          {{ t('nav.dashboard') }}
        </router-link>
      </nav>
      <div class="header-actions">
        <ThemeSwitch />
        <LanguageSwitch />
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { Setting } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/store/user'
import ThemeSwitch from '../common/ThemeSwitch.vue'
import LanguageSwitch from '../common/LanguageSwitch.vue'

const { t } = useI18n()
const userStore = useUserStore()

// 博客标题 - 可以根据语言切换
const blogTitle = computed(() => t('app.title'))

// 判断是否已登录
const isLoggedIn = computed(() => userStore.isLoggedIn())
</script>

<style scoped lang="less">
.site-header {
  background: var(--header-bg);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.logo {
  h1 {
    margin: 0;
    font-size: 24px;
    font-weight: bold;
    color: var(--primary-color);
  }

  a {
    text-decoration: none;
  }
}

.nav-menu {
  display: flex;
  gap: 30px;
}

.nav-item {
  text-decoration: none;
  color: var(--text-color);
  font-size: 16px;
  transition: color 0.3s;
  display: flex;
  align-items: center;
  gap: 4px;

  &:hover,
  &.router-link-active {
    color: var(--primary-color);
  }

  &.admin-link {
    .el-icon {
      font-size: 14px;
    }
  }
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 15px;
}

@media (max-width: 768px) {
  .nav-menu {
    display: none;
  }

  .logo h1 {
    font-size: 20px;
  }
}
</style>

