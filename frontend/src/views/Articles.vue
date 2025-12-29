<template>
  <div class="articles-page">
    <!-- 搜索和筛选栏 -->
    <div class="filter-bar">
      <el-input
        v-model="searchKeyword"
        :placeholder="t('home.searchPlaceholder')"
        class="search-input"
        clearable
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>

      <el-select
        v-model="selectedCategory"
        :placeholder="t('home.selectCategory')"
        clearable
        @change="handleFilterChange"
      >
        <el-option
          v-for="category in categories"
          :key="category.id"
          :label="category.name"
          :value="category.id"
        />
      </el-select>

      <el-select
        v-model="selectedTag"
        :placeholder="t('home.selectTag')"
        clearable
        @change="handleFilterChange"
      >
        <el-option
          v-for="tag in tags"
          :key="tag.id"
          :label="tag.name"
          :value="tag.id"
        />
      </el-select>
    </div>

    <!-- 文章列表 -->
    <div v-loading="loading" class="articles-list">
      <el-empty v-if="!loading && articles.length === 0" :description="t('home.noArticles')" />
      
      <div v-else class="article-grid">
        <article-card
          v-for="article in articles"
          :key="article.id"
          :article="article"
          @click="goToArticle(article.slug)"
        />
      </div>

      <!-- 分页 -->
      <el-pagination
        v-if="total > 0"
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 30, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Search } from '@element-plus/icons-vue'
import api from '@/api'
import ArticleCard from '@/components/article/ArticleCard.vue'

const router = useRouter()
const { t } = useI18n()

// 数据
const loading = ref(false)
const articles = ref([])
const categories = ref([])
const tags = ref([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const searchKeyword = ref('')
const selectedCategory = ref(null)
const selectedTag = ref(null)

// 获取文章列表
const fetchArticles = async () => {
  try {
    loading.value = true
    const params = {
      page: currentPage.value,
      page_size: pageSize.value
    }

    if (selectedCategory.value) {
      params.category_id = selectedCategory.value
    }

    if (selectedTag.value) {
      params.tag_id = selectedTag.value
    }

    if (searchKeyword.value) {
      // 使用搜索接口
      const response = await api.article.search(searchKeyword.value, params)
      // 后端返回的是 items，不是 list
      articles.value = response.items || []
      total.value = response.total || 0
    } else {
      const response = await api.article.list(params)
      // 后端返回的是 items，不是 list
      articles.value = response.items || []
      total.value = response.total || 0
    }
  } catch (error) {
    console.error('获取文章列表失败:', error)
    ElMessage.error(t('home.loadArticlesError'))
  } finally {
    loading.value = false
  }
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await api.category.list({ page: 1, page_size: 100 })
    // 后端返回的是 items，不是 list
    categories.value = response.items || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

// 获取标签列表
const fetchTags = async () => {
  try {
    const response = await api.tag.list({ page: 1, page_size: 100 })
    // 后端返回的是 items，不是 list
    tags.value = response.items || []
  } catch (error) {
    console.error('获取标签列表失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  currentPage.value = 1
  fetchArticles()
}

// 筛选改变
const handleFilterChange = () => {
  currentPage.value = 1
  fetchArticles()
}

// 页码改变
const handlePageChange = (page) => {
  currentPage.value = page
  fetchArticles()
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

// 页面大小改变
const handleSizeChange = (size) => {
  pageSize.value = size
  currentPage.value = 1
  fetchArticles()
}

// 跳转到文章详情
const goToArticle = (slug) => {
  router.push(`/article/${slug}`)
}

// 从URL参数初始化筛选条件
const initFromQuery = () => {
  const query = router.currentRoute.value.query
  
  if (query.category_id) {
    selectedCategory.value = parseInt(query.category_id)
  }
  
  if (query.tag_id) {
    selectedTag.value = parseInt(query.tag_id)
  }
  
  if (query.keyword) {
    searchKeyword.value = query.keyword
  }
}

// 初始化
onMounted(() => {
  initFromQuery()
  fetchCategories()
  fetchTags()
  fetchArticles()
})
</script>

<style scoped lang="less">
.articles-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 2rem;
}

.filter-bar {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  flex-wrap: wrap;

  .search-input {
    flex: 1;
    min-width: 300px;
  }

  .el-select {
    width: 200px;
  }
}

.articles-list {
  min-height: 400px;
}

.article-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.el-pagination {
  display: flex;
  justify-content: center;
  margin-top: 2rem;
}

@media (max-width: 768px) {
  .articles-page {
    padding: 1rem;
  }

  .filter-bar {
    flex-direction: column;

    .search-input {
      min-width: 100%;
    }

    .el-select {
      width: 100%;
    }
  }

  .article-grid {
    grid-template-columns: 1fr;
  }
}
</style>

