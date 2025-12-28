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
  background: #f0f2f5;
}

.el-aside {
  background: #001529;
  color: #fff;
  height: 100vh;
  overflow-y: auto;
}

.el-header {
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.el-main {
  padding: 20px;
}
</style>

