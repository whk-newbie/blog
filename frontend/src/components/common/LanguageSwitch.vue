<template>
  <el-dropdown @command="changeLanguage">
    <div class="language-switch">
      <el-icon :size="20">
        <Grid />
      </el-icon>
      <span class="current-lang">{{ currentLangLabel }}</span>
    </div>
    <template #dropdown>
      <el-dropdown-menu>
        <el-dropdown-item command="zh-CN" :disabled="currentLang === 'zh-CN'">
          简体中文 (Simplified Chinese)
        </el-dropdown-item>
        <el-dropdown-item command="en-US" :disabled="currentLang === 'en-US'">
          English (英语)
        </el-dropdown-item>
      </el-dropdown-menu>
    </template>
  </el-dropdown>
</template>

<script setup>
import { computed } from 'vue'
import { Grid } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'

const { locale, t } = useI18n()

const currentLang = computed(() => locale.value)

const currentLangLabel = computed(() => {
  const labels = {
    'zh-CN': '中文',
    'en-US': 'EN'
  }
  return labels[currentLang.value] || '中文'
})

const changeLanguage = (lang) => {
  locale.value = lang
  localStorage.setItem('language', lang)
  // 提示语言切换成功 - 使用当前语言显示
  const messages = {
    'zh-CN': '已切换到简体中文',
    'en-US': 'Switched to English'
  }
  ElMessage.success(messages[lang] || messages['zh-CN'])
}
</script>

<style scoped lang="less">
.language-switch {
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 5px;
  padding: 5px 10px;
  border-radius: 4px;
  transition: all 0.3s;
  color: var(--text-color);

  &:hover {
    background: var(--hover-bg);
    color: var(--primary-color);
  }

  .current-lang {
    font-size: 14px;
  }
}
</style>

