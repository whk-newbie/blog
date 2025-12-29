<template>
  <div class="page-header">
    <div class="page-header-title">
      <el-icon v-if="icon" :size="24">
        <component :is="icon" />
      </el-icon>
      <h2>{{ title }}</h2>
    </div>
    <p v-if="description" class="page-header-description">{{ description }}</p>
    <div v-if="$slots.extra" class="page-header-extra">
      <slot name="extra"></slot>
    </div>
  </div>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    default: ''
  },
  icon: {
    type: Object,
    default: null
  }
})
</script>

<style scoped lang="less">
.page-header {
  margin-bottom: 24px;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  flex-wrap: wrap;
  gap: 16px;
}

.page-header-title {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
  flex: 1;

  h2 {
    margin: 0;
    font-size: 20px;
    font-weight: 600;
    color: var(--text-color);
  }

  .el-icon {
    color: var(--primary-color);
  }
}

.page-header-description {
  margin: 0 0 16px;
  color: var(--text-secondary);
  font-size: 14px;
  width: 100%;
}

.page-header-extra {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 0;
}
</style>

