<template>
  <header class="site-header">
    <div class="header-container">
      <nav class="nav-menu">
        <router-link to="/" class="nav-item">{{ t('nav.home') }}</router-link>
        <router-link to="/articles" class="nav-item">{{ t('nav.articles') }}</router-link>
        <router-link to="/timeline" class="nav-item">{{ t('nav.timeline') }}</router-link>
        <router-link to="/tools" class="nav-item">{{ t('nav.tools') }}</router-link>
      </nav>
      <div class="header-actions">
        <!-- Search -->
        <el-button text size="small" class="icon-btn" @click="toggleSearch" :title="t('common.search')">
          <el-icon><Search /></el-icon>
        </el-button>
        <!-- Language switch -->
        <LanguageSwitch />
        <!-- Theme switch -->
        <ThemeSwitch />
        <!-- Admin link: always show when logged in -->
        <el-tooltip
          v-if="isLoggedIn"
          :content="t('nav.adminTooltip')"
          placement="bottom"
        >
          <el-button text size="small" @click="goToAdmin" class="admin-btn">
            <el-icon><Setting /></el-icon>
          </el-button>
        </el-tooltip>
        <!-- Login button: only shown when admin link is enabled AND not logged in -->
        <el-button
          v-if="showAdminLink && !isLoggedIn"
          type="primary"
          size="small"
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
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Setting, Search } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/store/user'
import LoginDialog from '../common/LoginDialog.vue'
import LanguageSwitch from '../common/LanguageSwitch.vue'
import ThemeSwitch from '../common/ThemeSwitch.vue'
import api from '@/api'

const { t } = useI18n()
const router = useRouter()
const userStore = useUserStore()

const showLoginDialog = ref(false)
const isLoggedIn = computed(() => userStore.isLoggedIn())
const showAdminLink = ref(true)

const emit = defineEmits(['toggleSearch'])

const toggleSearch = () => {
  emit('toggleSearch')
}

const fetchAdminConfig = async () => {
  // 只在已登录时获取后台配置（需要认证），未登录用默认值
  if (!userStore.isLoggedIn()) return
  try {
    const configs = await api.config.getConfigs({ config_type: 'site_config' })
    const showLinkConfig = configs.find(c => c.config_key === 'show_admin_link')
    if (showLinkConfig) {
      showAdminLink.value = showLinkConfig.config_value === 'true'
    }
    const adminPathConfig = configs.find(c => c.config_key === 'admin_path')
    if (adminPathConfig) {
      localStorage.setItem('admin_path', adminPathConfig.config_value)
    }
  } catch (e) {
    // use defaults
  }
}

onMounted(() => {
  fetchAdminConfig()
})

const handleLoginSuccess = () => {
  showLoginDialog.value = false
  if (router.currentRoute.value.path === '/') {
    const adminPath = localStorage.getItem('admin_path') || 'admin'
    router.push(`/${adminPath}`)
  }
}

const goToAdmin = () => {
  const adminPath = localStorage.getItem('admin_path') || 'admin'
  router.push(`/${adminPath}`)
}
</script>

<style scoped lang="less">
.site-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: var(--header-height);
  background: var(--header-bg);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid var(--border-color);
  z-index: 100;
}

.header-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 var(--spacing-md);
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.nav-menu {
  display: flex;
  gap: 24px;
}

.nav-item {
  text-decoration: none;
  color: var(--text-secondary);
  font-size: var(--font-size-sm);
  font-weight: 500;
  transition: color var(--transition-fast);

  &:hover,
  &.router-link-active {
    color: var(--text-heading);
  }
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;

  .icon-btn {
    color: var(--text-secondary);
    &:hover { color: var(--text-heading); }
  }

  .admin-btn {
    background: transparent !important;
    border: none !important;
    color: var(--text-secondary);
    padding: 4px 8px;
    transition: all var(--transition-fast);

    &:hover {
      color: var(--primary-color);
      background: var(--hover-bg) !important;
    }
  }
}

@media (max-width: 768px) {
  .nav-menu {
    gap: 16px;
  }
  .nav-item {
    font-size: var(--font-size-xs);
  }
}
</style>
