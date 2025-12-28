<template>
  <div class="admin-sidebar">
    <div class="sidebar-logo">
      <h2>博客管理</h2>
    </div>
    <el-menu
      :default-active="activeMenu"
      router
      background-color="#001529"
      text-color="#fff"
      active-text-color="#1890ff"
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
  Tools
} from '@element-plus/icons-vue'

const route = useRoute()
const activeMenu = computed(() => route.path)
</script>

<style scoped lang="less">
.admin-sidebar {
  height: 100%;
}

.sidebar-logo {
  padding: 20px;
  text-align: center;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);

  h2 {
    margin: 0;
    color: #fff;
    font-size: 18px;
  }
}

.el-menu {
  border-right: none;
}
</style>

