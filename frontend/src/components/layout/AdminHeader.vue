<template>
  <div class="admin-header">
    <div class="header-left">
      <el-breadcrumb separator="/">
        <el-breadcrumb-item
          v-for="item in breadcrumbs"
          :key="item.path"
          :to="item.path"
        >
          {{ item.title }}
        </el-breadcrumb-item>
      </el-breadcrumb>
    </div>
    <div class="header-right">
      <ThemeSwitch />
      <el-dropdown @command="handleCommand">
        <div class="user-info">
          <el-avatar :size="32">
            {{ username.charAt(0).toUpperCase() }}
          </el-avatar>
          <span class="username">{{ username }}</span>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="profile">个人资料</el-dropdown-item>
            <el-dropdown-item command="password">修改密码</el-dropdown-item>
            <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import ThemeSwitch from '../common/ThemeSwitch.vue'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const username = computed(() => userStore.username || 'Admin')

const breadcrumbs = computed(() => {
  const matched = route.matched.filter(item => item.meta?.title)
  return matched.map(item => ({
    path: item.path,
    title: item.meta.title
  }))
})

const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      router.push('/admin/profile')
      break
    case 'password':
      router.push('/admin/password')
      break
    case 'logout':
      userStore.logout()
      router.push('/')
      ElMessage.success('退出登录成功')
      break
  }
}
</script>

<style scoped lang="less">
.admin-header {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px;
}

.header-left {
  flex: 1;

  :deep(.el-breadcrumb) {
    font-size: 14px;

    .el-breadcrumb__item {
      .el-breadcrumb__inner {
        color: var(--text-secondary);
        font-weight: 400;
        transition: color 0.3s;

        &:hover {
          color: var(--primary-color);
        }

        &.is-link {
          color: var(--text-secondary);
        }
      }

      &:last-child .el-breadcrumb__inner {
        color: var(--text-color);
        font-weight: 500;
      }
    }

    .el-breadcrumb__separator {
      color: var(--text-disabled);
    }
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  cursor: pointer;
  padding: 8px 16px;
  border-radius: 8px;
  transition: all 0.3s;
  border: 1px solid transparent;

  &:hover {
    background: var(--bg-blue-light);
    border-color: var(--border-blue);
  }

  :deep(.el-avatar) {
    background: linear-gradient(135deg, var(--primary-color) 0%, var(--secondary-color) 100%);
    font-weight: 600;
    box-shadow: var(--shadow-sm);
  }

  .username {
    font-size: 14px;
    color: var(--text-color);
    font-weight: 500;
  }
}
</style>

