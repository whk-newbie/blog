<template>
  <div class="admin-sidebar">
    <div class="sidebar-logo">
      <el-icon class="logo-icon"><Collection /></el-icon>
      <h2>博客管理系统</h2>
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
        <span>仪表盘</span>
      </el-menu-item>
      
      <el-sub-menu index="articles">
        <template #title>
          <el-icon><Document /></el-icon>
          <span>文章管理</span>
        </template>
        <el-menu-item index="/admin/articles">文章列表</el-menu-item>
        <el-menu-item index="/admin/articles/new">新建文章</el-menu-item>
        <el-menu-item index="/admin/categories">分类管理</el-menu-item>
        <el-menu-item index="/admin/tags">标签管理</el-menu-item>
      </el-sub-menu>

      <el-sub-menu index="stats">
        <template #title>
          <el-icon><DataAnalysis /></el-icon>
          <span>统计分析</span>
        </template>
        <el-menu-item index="/admin/stats/visits">访问统计</el-menu-item>
        <el-menu-item index="/admin/stats/fingerprints">指纹记录</el-menu-item>
      </el-sub-menu>

      <el-menu-item index="/admin/crawler">
        <el-icon><Monitor /></el-icon>
        <span>爬虫监控</span>
      </el-menu-item>

      <el-sub-menu index="system">
        <template #title>
          <el-icon><Setting /></el-icon>
          <span>系统管理</span>
        </template>
        <el-menu-item index="/admin/system/config">系统配置</el-menu-item>
        <el-menu-item index="/admin/system/logs">系统日志</el-menu-item>
        <el-menu-item index="/admin/system/backup">数据备份</el-menu-item>
      </el-sub-menu>

      <el-menu-item index="/admin/tools">
        <el-icon><Tools /></el-icon>
        <span>开发工具</span>
      </el-menu-item>
    </el-menu>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import {
  Monitor,
  Document,
  DataAnalysis,
  Setting,
  Tools,
  Collection
} from '@element-plus/icons-vue'

const route = useRoute()
const activeMenu = computed(() => route.path)
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

  .logo-icon {
    font-size: 36px;
    color: #ffffff;
    margin-bottom: 8px;
  }

  h2 {
    margin: 0;
    color: #ffffff;
    font-size: 18px;
    font-weight: 600;
    letter-spacing: 0.5px;
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

