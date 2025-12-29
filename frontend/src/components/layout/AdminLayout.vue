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
    <!-- 登录对话框 - 未登录时自动显示 -->
    <LoginDialog v-model="showLoginDialog" @success="handleLoginSuccess" />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import Sidebar from './Sidebar.vue'
import AdminHeader from './AdminHeader.vue'
import LoginDialog from '../common/LoginDialog.vue'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const showLoginDialog = ref(false)

// 检查是否需要显示登录对话框
const checkLoginStatus = () => {
  if (route.meta.requiresAuth && !userStore.isLoggedIn()) {
    showLoginDialog.value = true
  }
}

// 监听路由变化
watch(() => route.path, () => {
  checkLoginStatus()
})

// 监听登录状态变化
watch(() => userStore.isLoggedIn(), (isLoggedIn) => {
  if (isLoggedIn) {
    showLoginDialog.value = false
  } else if (route.meta.requiresAuth) {
    showLoginDialog.value = true
  }
})

// 登录成功回调
const handleLoginSuccess = () => {
  showLoginDialog.value = false
}

// 组件挂载时检查
onMounted(() => {
  checkLoginStatus()
})
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

