<template>
  <div class="admin-sidebar">
    <div class="sidebar-logo" @click="goToHome">
      <el-icon class="logo-icon"><Collection /></el-icon>
      <h2>{{ t('app.adminSystem') }}</h2>
      <p class="logo-subtitle">Admin Dashboard</p>
    </div>
    <el-menu
      :default-active="activeMenu"
      router
      background-color="transparent"
      text-color="rgba(255, 255, 255, 0.85)"
      active-text-color="#ffffff"
      class="sidebar-menu"
    >
      <el-menu-item index="/admin">
        <el-icon><Monitor /></el-icon>
        <span>{{ t('nav.dashboard') }}</span>
      </el-menu-item>
      
      <el-sub-menu index="articles">
        <template #title>
          <el-icon><Document /></el-icon>
          <span>{{ t('nav.articleManage') }}</span>
        </template>
        <el-menu-item index="/admin/articles">{{ t('nav.articleList') }}</el-menu-item>
        <el-menu-item index="/admin/articles/new">{{ t('nav.newArticle') }}</el-menu-item>
        <el-menu-item index="/admin/categories">{{ t('nav.categories') }}</el-menu-item>
        <el-menu-item index="/admin/tags">{{ t('nav.tags') }}</el-menu-item>
      </el-sub-menu>

      <el-sub-menu index="stats">
        <template #title>
          <el-icon><DataAnalysis /></el-icon>
          <span>{{ t('nav.stats') }}</span>
        </template>
        <el-menu-item index="/admin/stats">{{ t('nav.visits') }}</el-menu-item>
        <el-menu-item index="/admin/fingerprints">{{ t('nav.fingerprints') }}</el-menu-item>
      </el-sub-menu>

      <el-menu-item index="/admin/crawler">
        <el-icon><Monitor /></el-icon>
        <span>{{ t('nav.crawler') }}</span>
      </el-menu-item>

      <el-sub-menu index="system">
        <template #title>
          <el-icon><Setting /></el-icon>
          <span>{{ t('nav.system') }}</span>
        </template>
        <el-menu-item index="/admin/config">{{ t('nav.config') }}</el-menu-item>
        <el-menu-item index="/admin/logs">{{ t('nav.logs') }}</el-menu-item>
        <el-menu-item index="/admin/system/backup">{{ t('nav.backup') }}</el-menu-item>
      </el-sub-menu>

      <el-menu-item index="/admin/tools">
        <el-icon><Tools /></el-icon>
        <span>{{ t('nav.tools') }}</span>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import {
  Monitor,
  Document,
  DataAnalysis,
  Setting,
  Tools,
  Collection
} from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const { t } = useI18n()
const activeMenu = computed(() => route.path)

// 跳转到首页
const goToHome = () => {
  router.push('/')
}
</script>

<style scoped lang="less">
.admin-sidebar {
  height: 100%;
}

.sidebar-logo {
  padding: 24px 16px;
  text-align: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.15);
  background: rgba(255, 255, 255, 0.05);
  cursor: pointer;
  transition: all 0.3s;
  border-radius: 0;

  &:hover {
    background: rgba(255, 255, 255, 0.1);
    
    h2 {
      color: #ffffff;
    }
    
    .logo-icon {
      transform: scale(1.1);
    }
  }

  .logo-icon {
    font-size: 36px;
    color: #ffffff;
    margin-bottom: 8px;
    transition: transform 0.3s;
  }

  h2 {
    margin: 0;
    color: #ffffff;
    font-size: 18px;
    font-weight: 600;
    letter-spacing: 0.5px;
    transition: color 0.3s;
  }

  .logo-subtitle {
    margin: 4px 0 0;
    color: rgba(255, 255, 255, 0.6);
    font-size: 12px;
    letter-spacing: 1px;
  }
}

:deep(.sidebar-menu) {
  border-right: none;
  padding: 8px 0;

  .el-menu-item,
  .el-sub-menu__title {
    height: 48px;
    line-height: 48px;
    margin: 4px 8px;
    border-radius: 8px;
    transition: all 0.3s;

    &:hover {
      background: var(--admin-sidebar-hover) !important;
    }

    &.is-active {
      background: var(--admin-sidebar-active) !important;
      color: #ffffff !important;
      font-weight: 500;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    }
  }

  .el-sub-menu {
    .el-menu-item {
      min-width: auto;
      padding-left: 48px !important;
      height: 40px;
      line-height: 40px;
    }
  }

  .el-icon {
    font-size: 18px;
    margin-right: 8px;
  }
}
</style>

