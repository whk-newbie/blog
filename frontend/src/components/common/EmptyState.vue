<template>
  <div class="empty-state">
    <el-icon :size="size" class="empty-icon">
      <component :is="icon" />
    </el-icon>
    <p class="empty-text">{{ displayText }}</p>
    <slot name="action"></slot>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Document } from '@element-plus/icons-vue'

const props = withDefaults(defineProps({
  icon: {
    type: Object,
    default: () => Document
  },
  text: {
    type: String,
    default: '暂无数据'
  },
  description: {
    type: String,
    default: ''
  },
  size: {
    type: Number,
    default: 80
  }
}), {
  icon: Document,
  text: '暂无数据',
  size: 80
})

// 如果提供了description，优先使用description
const displayText = computed(() => props.description || props.text)
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

