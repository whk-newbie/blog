<template>
  <div class="theme-switch" @click="toggleTheme">
    <el-icon :size="20">
      <Sunny v-if="isDark" />
      <Moon v-else />
    </el-icon>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Sunny, Moon } from '@element-plus/icons-vue'
import { useThemeStore } from '@/store/theme'

const themeStore = useThemeStore()

const isDark = computed(() => themeStore.theme === 'dark')

const toggleTheme = () => {
  themeStore.toggleTheme()
}
</script>

<style scoped lang="less">
.theme-switch {
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  transition: all 0.3s;
  color: var(--text-color);

  &:hover {
    background: var(--hover-bg);
    color: var(--primary-color);
  }
}
</style>

