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
        <el-button
          v-if="!isLoggedIn"
          type="primary"
          size="default"
          @click="showLoginDialog = true"
        >
          {{ t('login.title') }}
        </el-button>
        <LoginDialog v-model="showLoginDialog" @success="handleLoginSuccess" />
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Setting } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/store/user'
import ThemeSwitch from '../common/ThemeSwitch.vue'
import LanguageSwitch from '../common/LanguageSwitch.vue'
import LoginDialog from '../common/LoginDialog.vue'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()

const showLoginDialog = ref(false)

// 博客标题 - 可以根据语言切换
const blogTitle = computed(() => t('app.title'))

// 判断是否已登录
const isLoggedIn = computed(() => userStore.isLoggedIn())

// 登录成功回调
const handleLoginSuccess = () => {
  showLoginDialog.value = false
  // 如果当前在首页，可以选择跳转到管理后台
  if (router.currentRoute.value.path === '/') {
    router.push('/admin')
  }
}
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

  .el-button {
    border-radius: 8px;
    font-weight: 500;
    padding: 10px 20px;
  }
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

