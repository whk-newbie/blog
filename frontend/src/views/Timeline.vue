<template>
  <div class="timeline-page">
    <h1 class="page-title">时间轴</h1>
    <div class="timeline" v-loading="loading">
      <el-empty v-if="!loading && timelineData.length === 0" description="暂无文章" />
      <div v-for="year in timelineData" :key="year.year" class="year-group">
        <div class="year-dot">{{ year.year }}</div>
        <div v-for="month in year.months" :key="`${year.year}-${month.month}`" class="month-group">
          <h3 class="month-label">{{ month.month }}月</h3>
          <div v-for="article in month.articles" :key="article.id" class="timeline-item" @click="goToArticle(article.slug)">
            <span class="item-dot"></span>
            <span class="item-title">{{ displayTitle(article) }}</span>
            <span class="item-date">{{ formatDay(article.publish_at) }}</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useLanguageStore } from '@/store/language'
import api from '@/api'

const router = useRouter()
const languageStore = useLanguageStore()
const loading = ref(false)
const timelineData = ref([])

const displayTitle = (article) => {
  return languageStore.language === 'en-US' && article.title_en ? article.title_en : article.title
}

const formatDay = (dateStr) => {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return `${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const fetchArticles = async () => {
  loading.value = true
  try {
    const res = await api.article.list({ page: 1, page_size: 500 })
    const items = res?.items || []

    // Group by year then month
    const map = {}
    for (const article of items) {
      const date = new Date(article.publish_at || article.created_at)
      const year = date.getFullYear()
      const month = date.getMonth() + 1

      if (!map[year]) map[year] = {}
      if (!map[year][month]) map[year][month] = []
      map[year][month].push(article)
    }

    const result = Object.entries(map)
      .sort(([a], [b]) => Number(b) - Number(a))
      .map(([year, months]) => ({
        year: Number(year),
        months: Object.entries(months)
          .sort(([a], [b]) => Number(b) - Number(a))
          .map(([month, articles]) => ({
            month: Number(month),
            articles,
          })),
      }))

    timelineData.value = result
  } catch (e) {
    // silently fail
  } finally {
    loading.value = false
  }
}

const goToArticle = (slug) => {
  router.push(`/article/${slug}`)
}

onMounted(() => {
  fetchArticles()
})
</script>

<style scoped lang="less">
.timeline-page {
  padding: var(--spacing-md) 0;
}

.page-title {
  font-size: var(--font-size-xxl);
  font-weight: 700;
  color: var(--text-heading);
  margin: 0 0 var(--spacing-xl);
}

.timeline {
  position: relative;
  padding-left: 40px;

  &::before {
    content: '';
    position: absolute;
    left: 8px;
    top: 0;
    bottom: 0;
    width: 2px;
    background: var(--border-color);
  }
}

.year-group {
  position: relative;
  margin-bottom: var(--spacing-xl);

  .year-dot {
    position: absolute;
    left: -36px;
    top: 0;
    width: 10px;
    height: 10px;
    border-radius: 50%;
    background: var(--primary-color);
    border: 2px solid var(--bg-color);
  }

  .year-dot {
    font-size: var(--font-size-lg);
    font-weight: 700;
    color: var(--text-heading);
    background: none;
    border: none;
    width: auto;
    height: auto;
    left: -48px;
    top: -2px;
  }
}

.month-group {
  margin: var(--spacing-md) 0;

  .month-label {
    font-size: var(--font-size-sm);
    font-weight: 600;
    color: var(--text-secondary);
    margin: 0 0 var(--spacing-sm);
  }
}

.timeline-item {
  display: flex;
  align-items: center;
  gap: var(--spacing-md);
  padding: 8px 0;
  cursor: pointer;
  transition: color var(--transition-fast);
  position: relative;

  &::before {
    content: '';
    position: absolute;
    left: -32px;
    top: 50%;
    transform: translateY(-50%);
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: var(--border-color);
  }

  &:hover {
    color: var(--primary-color);
  }

  .item-title {
    flex: 1;
    font-size: var(--font-size-sm);
    color: inherit;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .item-date {
    font-size: var(--font-size-xs);
    color: var(--text-secondary);
    flex-shrink: 0;
  }
}
</style>
