<template>
  <div v-if="visible" class="loading-overlay" :class="{ fullscreen }">
    <div class="loading-content">
      <el-icon class="loading-icon" :size="size">
        <Loading />
      </el-icon>
      <p v-if="text" class="loading-text">{{ text }}</p>
    </div>
  </div>
</template>

<script setup>
import { Loading } from '@element-plus/icons-vue'

defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  text: {
    type: String,
    default: ''
  },
  fullscreen: {
    type: Boolean,
    default: true
  },
  size: {
    type: Number,
    default: 40
  }
})
</script>

<style scoped lang="less">
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.9);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;

  &:not(.fullscreen) {
    position: absolute;
    background: rgba(255, 255, 255, 0.7);
  }
}

.loading-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
}

.loading-icon {
  color: var(--primary-color);
  animation: rotate 1s linear infinite;
}

.loading-text {
  margin: 0;
  color: var(--text-color);
  font-size: 14px;
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

:root[data-theme='dark'] {
  .loading-overlay {
    background: rgba(0, 0, 0, 0.9);

    &:not(.fullscreen) {
      background: rgba(0, 0, 0, 0.7);
    }
  }
}
</style>

