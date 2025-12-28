<template>
  <header class="site-header">
    <div class="header-container">
      <div class="logo">
        <router-link to="/">
          <h1>{{ blogTitle }}</h1>
        </router-link>
      </div>
      <nav class="nav-menu">
        <router-link to="/" class="nav-item">首页</router-link>
        <router-link to="/articles" class="nav-item">文章</router-link>
        <router-link to="/about" class="nav-item">关于</router-link>
      </nav>
      <div class="header-actions">
        <ThemeSwitch />
        <LanguageSwitch />
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import ThemeSwitch from '../common/ThemeSwitch.vue'
import LanguageSwitch from '../common/LanguageSwitch.vue'

const blogTitle = ref('我的博客')
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

  &:hover,
  &.router-link-active {
    color: var(--primary-color);
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

