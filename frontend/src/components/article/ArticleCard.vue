<template>
  <article class="article-card" @click="$emit('click')">
    <div class="card-cover">
      <img
        v-if="article.cover_image"
        :src="article.cover_image"
        :alt="displayTitle"
        loading="lazy"
      />
      <div v-else class="cover-placeholder"></div>
    </div>
    <div class="card-body">
      <h3 class="card-title">{{ displayTitle }}</h3>
      <div class="card-meta">
        <time>{{ formatDate(article.publish_at || article.created_at) }}</time>
        <span v-if="article.category" class="category-badge">{{ article.category.name }}</span>
      </div>
    </div>
  </article>
</template>

<script setup>
import { computed } from 'vue'
import { useLanguageStore } from '@/store/language'

const props = defineProps({
  article: { type: Object, required: true },
})

defineEmits(['click'])

const languageStore = useLanguageStore()
const isEnglish = computed(() => languageStore.language === 'en-US')

const displayTitle = computed(() => {
  if (isEnglish.value && props.article.title_en) {
    return props.article.title_en
  }
  return props.article.title
})

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}
</script>

<style scoped lang="less">
.article-card {
  background: var(--card-bg);
  border-radius: var(--border-radius-sm);
  overflow: hidden;
  cursor: pointer;
  transition: transform var(--transition-normal), box-shadow var(--transition-normal);
  box-shadow: var(--shadow-sm);

  &:hover {
    transform: translateY(-2px);
    box-shadow: var(--shadow-md);
  }
}

.card-cover {
  position: relative;
  padding-top: 56.25%;
  overflow: hidden;
  background: var(--bg-secondary);

  img {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .cover-placeholder {
    position: absolute;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    background: linear-gradient(135deg, var(--bg-secondary), var(--bg-tertiary));
  }
}

.card-body {
  padding: var(--spacing-md);
}

.card-title {
  margin: 0 0 10px;
  font-size: var(--font-size-md);
  font-weight: 600;
  color: var(--text-heading);
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-meta {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: var(--font-size-xs);
  color: var(--text-secondary);

  .category-badge {
    padding: 1px 8px;
    background: var(--bg-secondary);
    border-radius: 10px;
  }

  .no-translate {
    color: var(--warning-color);
  }
}
</style>
