<template>
  <div class="empty-state">
    <el-icon :size="size" class="empty-icon">
      <component :is="icon" />
    </el-icon>
    <p class="empty-text">{{ text }}</p>
    <slot name="action"></slot>
  </div>
</template>

<script setup>
import { Document } from '@element-plus/icons-vue'

withDefaults(defineProps(), {
  icon: Document,
  text: '暂无数据',
  size: 80
})
</script>

<style scoped lang="less">
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 60px 20px;
  color: var(--text-secondary);
}

.empty-icon {
  color: var(--border-color);
  margin-bottom: 20px;
}

.empty-text {
  margin: 0 0 20px;
  font-size: 14px;
}
</style>

