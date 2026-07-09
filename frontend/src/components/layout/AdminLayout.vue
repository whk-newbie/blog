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
import { watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import Sidebar from './Sidebar.vue'
import AdminHeader from './AdminHeader.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const adminPath = () => localStorage.getItem('admin_path') || 'admin'

const checkAuth = () => {
  if (!userStore.isLoggedIn()) {
    const currentPath = route.fullPath
    router.push(`/login?redirect=${encodeURIComponent(currentPath)}`)
  }
}

watch(() => route.path, () => {
  if (route.name !== 'Login') checkAuth()
})

onMounted(() => {
  checkAuth()
})
</script>

<style scoped lang="less">
.admin-layout {
  min-height: 100vh;
  background: var(--admin-content-bg);
}

.el-aside {
  background: var(--admin-sidebar-bg);
  color: var(--admin-sidebar-text);
  min-height: 100vh;
  border-right: 1px solid var(--admin-border);
}

.el-header {
  background: var(--admin-header-bg);
  box-shadow: var(--shadow-sm);
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
