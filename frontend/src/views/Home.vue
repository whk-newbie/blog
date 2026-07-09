<template>
  <div class="home-page">
    <h2 class="section-title">最新文章</h2>
    <div v-loading="loading" class="articles-grid">
      <el-empty v-if="!loading && articles.length === 0" :description="t('home.noArticles')" />
      <ArticleCard
        v-for="article in articles"
        :key="article.id"
        :article="article"
        @click="goToArticle(article.slug)"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import api from '@/api'
import ArticleCard from '@/components/article/ArticleCard.vue'

const router = useRouter()
const { t } = useI18n()

const loading = ref(false)
const articles = ref([])

const fetchArticles = async () => {
  loading.value = true
  try {
    const res = await api.article.list({ page: 1, page_size: 20 })
    articles.value = res?.items || []
  } catch (e) {
    ElMessage.error(t('home.loadArticlesError'))
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
.home-page {
  padding: var(--spacing-md) 0;
}

.section-title {
  font-size: var(--font-size-lg);
  font-weight: 600;
  color: var(--text-heading);
  margin: 0 0 var(--spacing-lg);
}

.articles-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: var(--spacing-lg);
  min-height: 300px;
}

@media (max-width: 640px) {
  .articles-grid {
    grid-template-columns: 1fr;
  }
}
</style>
