<template>
  <div class="home-page">
    <div class="home-container">
      <!-- 推荐文章 -->
      <section v-if="featuredArticles.length > 0" class="featured-section">
        <div class="section-header">
          <h2 class="section-title">
            <el-icon><Star /></el-icon>
            {{ t('home.featuredArticles') }}
          </h2>
        </div>
        <div class="featured-grid">
          <article-card
            v-for="article in featuredArticles"
            :key="article.id"
            :article="article"
            featured
            @click="goToArticle(article.slug)"
          />
        </div>
      </section>

      <div class="content-layout">
        <!-- 主内容区 - 最新文章 -->
        <main class="main-content">
          <section class="latest-section">
            <div class="section-header">
              <h2 class="section-title">
                <el-icon><Clock /></el-icon>
                {{ t('home.latestArticles') }}
              </h2>
              <el-button text @click="goToArticles">
                {{ t('home.viewAll') }} <el-icon><ArrowRight /></el-icon>
              </el-button>
            </div>
            
            <div v-loading="loading" class="articles-list">
              <el-empty v-if="!loading && latestArticles.length === 0" :description="t('home.noArticles')" />
              
              <div v-else class="article-grid">
                <article-card
                  v-for="article in latestArticles"
                  :key="article.id"
                  :article="article"
                  @click="goToArticle(article.slug)"
                />
              </div>
            </div>
          </section>
        </main>

        <!-- 侧边栏 -->
        <aside class="sidebar">
          <!-- 分类 -->
          <el-card class="sidebar-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><Folder /></el-icon>
                <span>{{ t('home.categories') }}</span>
              </div>
            </template>
            <div v-loading="loadingCategories" class="categories-list">
              <el-empty v-if="!loadingCategories && categories.length === 0" :description="t('home.noCategories')" :image-size="60" />
              <div v-else class="category-items">
                <div
                  v-for="category in categories"
                  :key="category.id"
                  class="category-item"
                  @click="goToCategoryArticles(category.id)"
                >
                  <span class="category-name">{{ category.name }}</span>
                  <el-tag size="small" type="info">{{ category.article_count || 0 }}</el-tag>
                </div>
              </div>
            </div>
          </el-card>

          <!-- 标签云 -->
          <el-card class="sidebar-card" shadow="hover">
            <template #header>
              <div class="card-header">
                <el-icon><PriceTag /></el-icon>
                <span>{{ t('home.hotTags') }}</span>
              </div>
            </template>
            <div v-loading="loadingTags" class="tags-cloud">
              <el-empty v-if="!loadingTags && tags.length === 0" :description="t('home.noTags')" :image-size="60" />
              <div v-else class="tag-items">
                <el-tag
                  v-for="tag in tags"
                  :key="tag.id"
                  class="tag-item"
                  :type="getRandomTagType()"
                  @click="goToTagArticles(tag.id)"
                >
                  {{ tag.name }} ({{ tag.article_count || 0 }})
                </el-tag>
              </div>
            </div>
          </el-card>
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import {
  Star,
  Clock,
  Folder,
  PriceTag,
  ArrowRight
} from '@element-plus/icons-vue'
import api from '@/api'
import ArticleCard from '@/components/article/ArticleCard.vue'

const router = useRouter()
const { t } = useI18n()

// 数据
const loading = ref(false)
const loadingCategories = ref(false)
const loadingTags = ref(false)
const featuredArticles = ref([])
const latestArticles = ref([])
const categories = ref([])
const tags = ref([])

// 获取推荐文章
const fetchFeaturedArticles = async () => {
  try {
    const response = await api.article.list({
      page: 1,
      page_size: 3,
      is_featured: true
    })
    // 后端返回的是 items，不是 list
    featuredArticles.value = response.items || []
  } catch (error) {
    console.error('获取推荐文章失败:', error)
  }
}

// 获取最新文章
const fetchLatestArticles = async () => {
  try {
    loading.value = true
    const response = await api.article.list({
      page: 1,
      page_size: 6
    })
    // 后端返回的是 items，不是 list
    latestArticles.value = response.items || []
  } catch (error) {
    console.error('获取最新文章失败:', error)
    ElMessage.error(t('home.loadArticlesError'))
  } finally {
    loading.value = false
  }
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    loadingCategories.value = true
    const response = await api.category.list({ page: 1, page_size: 20 })
    // 后端返回的是 items，不是 list
    categories.value = response.items || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
  } finally {
    loadingCategories.value = false
  }
}

// 获取标签列表
const fetchTags = async () => {
  try {
    loadingTags.value = true
    const response = await api.tag.list({ page: 1, page_size: 30 })
    // 后端返回的是 items，不是 list
    tags.value = response.items || []
  } catch (error) {
    console.error('获取标签列表失败:', error)
  } finally {
    loadingTags.value = false
  }
}

// 跳转到文章列表
const goToArticles = () => {
  router.push('/articles')
}

// 跳转到文章详情
const goToArticle = (slug) => {
  router.push(`/article/${slug}`)
}

// 跳转到分类文章列表
const goToCategoryArticles = (categoryId) => {
  router.push(`/articles?category_id=${categoryId}`)
}

// 跳转到标签文章列表
const goToTagArticles = (tagId) => {
  router.push(`/articles?tag_id=${tagId}`)
}

// 随机标签类型（用于标签云样式）
const tagTypes = ['success', 'info', 'warning', 'danger']
let tagTypeIndex = 0
const getRandomTagType = () => {
  const type = tagTypes[tagTypeIndex % tagTypes.length]
  tagTypeIndex++
  return type
}

// 初始化
onMounted(() => {
  fetchFeaturedArticles()
  fetchLatestArticles()
  fetchCategories()
  fetchTags()
})
</script>

<style lang="less" scoped>
.home-page {
  min-height: 100vh;
  background: var(--bg-color-secondary);
}

// 主容器
.home-container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 2rem;
}

// 推荐文章区域
.featured-section {
  margin-bottom: 3rem;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-color);
  margin: 0;

  .el-icon {
    font-size: 1.75rem;
    color: var(--primary-color);
  }
}

.featured-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 2rem;
}

// 内容布局
.content-layout {
  display: flex;
  gap: 2rem;
  align-items: flex-start;
}

.main-content {
  flex: 1;
  min-width: 0;
}

.latest-section {
  background: var(--bg-color);
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.articles-list {
  min-height: 300px;
}

.article-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

// 侧边栏
.sidebar {
  width: 320px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.sidebar-card {
  .card-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 1.125rem;
    font-weight: 600;
    color: var(--text-color);

    .el-icon {
      font-size: 1.25rem;
      color: var(--primary-color);
    }
  }
}

// 分类列表
.categories-list {
  min-height: 150px;
}

.category-items {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.category-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background: var(--bg-color-secondary);
  border-radius: 4px;
  cursor: pointer;
  transition: all 0.2s;

  &:hover {
    background: var(--primary-color);
    color: white;
    transform: translateX(4px);

    .el-tag {
      background: white;
      color: var(--primary-color);
    }
  }

  .category-name {
    font-size: 0.9375rem;
  }
}

// 标签云
.tags-cloud {
  min-height: 150px;
}

.tag-items {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.tag-item {
  cursor: pointer;
  transition: all 0.2s;
  font-size: 0.875rem;

  &:hover {
    transform: scale(1.1);
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  }
}

// 响应式设计
@media (max-width: 1200px) {
  .content-layout {
    flex-direction: column;
  }

  .sidebar {
    width: 100%;
  }
}

@media (max-width: 768px) {
  .home-container {
    padding: 1rem;
  }

  .featured-grid,
  .article-grid {
    grid-template-columns: 1fr;
  }

  .latest-section {
    padding: 1.5rem;
  }

  .section-title {
    font-size: 1.25rem;
  }
}
</style>

