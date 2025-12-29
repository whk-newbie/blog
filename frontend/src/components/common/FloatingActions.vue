<template>
  <div class="floating-actions">
    <!-- 回到顶部按钮 -->
    <el-tooltip
      :content="t('common.backToTop')"
      placement="left"
      :show="showBackToTop"
    >
      <div
        v-show="showBackToTop"
        class="floating-btn back-to-top"
        @click="scrollToTop"
      >
        <el-icon :size="20">
          <ArrowUp />
        </el-icon>
      </div>
    </el-tooltip>

    <!-- 语言切换 -->
    <el-tooltip
      :content="t('language.switchLanguage')"
      placement="left"
    >
      <el-dropdown
        @command="changeLanguage"
        trigger="click"
        placement="left"
      >
        <div class="floating-btn language-btn">
          <el-icon :size="20">
            <Grid />
          </el-icon>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item command="zh-CN" :disabled="currentLang === 'zh-CN'">
              {{ t('language.zhCNFull') }}
            </el-dropdown-item>
            <el-dropdown-item command="en-US" :disabled="currentLang === 'en-US'">
              {{ t('language.enUSFull') }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </el-tooltip>

    <!-- 主题切换 -->
    <el-tooltip
      :content="isDark ? t('common.switchToLight') : t('common.switchToDark')"
      placement="left"
    >
      <div class="floating-btn theme-btn" @click="toggleTheme">
        <el-icon :size="20">
          <Sunny v-if="isDark" />
          <Moon v-else />
        </el-icon>
      </div>
    </el-tooltip>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Grid, Sunny, Moon, ArrowUp } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useThemeStore } from '@/store/theme'

const { locale, t } = useI18n()
const themeStore = useThemeStore()

const currentLang = computed(() => locale.value)
const isDark = computed(() => themeStore.theme === 'dark')
const showBackToTop = ref(false)

// 监听滚动，显示/隐藏回到顶部按钮
const handleScroll = () => {
  showBackToTop.value = window.scrollY > 300
}

// 回到顶部
const scrollToTop = () => {
  window.scrollTo({
    top: 0,
    behavior: 'smooth'
  })
}

// 切换主题
const toggleTheme = () => {
  themeStore.toggleTheme()
}

// 切换语言
const changeLanguage = (lang) => {
  locale.value = lang
  localStorage.setItem('language', lang)
  const messages = {
    'zh-CN': t('language.switchedToZhCN'),
    'en-US': t('language.switchedToEnUS')
  }
  ElMessage.success(messages[lang] || messages['zh-CN'])
}

onMounted(() => {
  window.addEventListener('scroll', handleScroll)
  handleScroll() // 初始化检查
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>

<style scoped lang="less">
.floating-actions {
  position: fixed;
  right: 20px;
  bottom: 20px;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.floating-btn {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: var(--bg-color);
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s;
  color: var(--text-color);
  border: 1px solid var(--border-color);

  &:hover {
    background: var(--primary-color);
    color: #fff;
    transform: translateY(-2px);
    box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
  }

  &:active {
    transform: translateY(0);
  }
}

.back-to-top {
  animation: fadeInUp 0.3s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 768px) {
  .floating-actions {
    right: 15px;
    bottom: 15px;
    gap: 10px;
  }

  .floating-btn {
    width: 40px;
    height: 40px;

    .el-icon {
      font-size: 18px;
    }
  }
}
</style>

