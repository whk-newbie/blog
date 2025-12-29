<template>
  <div class="admin-layout">
    <el-container>
      <el-aside width="200px">
        <Sidebar />
      </el-aside>
      <el-container>
        <el-header>
          <AdminHeader />
        </el-header>
        <el-main>
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup>
import Sidebar from './Sidebar.vue'
import AdminHeader from './AdminHeader.vue'
</script>

<style scoped lang="less">
.admin-layout {
  min-height: 100vh;
  background: var(--admin-content-bg);
}

.el-aside {
  background: var(--admin-sidebar-bg);
  color: var(--text-white);
  height: 100vh;
  overflow-y: auto;
  box-shadow: 2px 0 8px rgba(37, 99, 235, 0.1);
  
  &::-webkit-scrollbar {
    width: 6px;
  }
  
  &::-webkit-scrollbar-thumb {
    background: rgba(255, 255, 255, 0.2);
    border-radius: 3px;
    
    &:hover {
      background: rgba(255, 255, 255, 0.3);
    }
  }
  
  &::-webkit-scrollbar-track {
    background: transparent;
  }
}

.el-header {
  background: var(--admin-header-bg);
  box-shadow: 0 2px 8px rgba(37, 99, 235, 0.08);
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--border-color);
}

.el-main {
  padding: 24px;
  background: var(--admin-content-bg);
  min-height: calc(100vh - 60px);
}
</style>

