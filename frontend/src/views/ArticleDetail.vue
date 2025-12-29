<template>
  <div class="article-detail-page">
    <div v-loading="loading" class="article-container">
      <el-empty v-if="!loading && !article" description="文章不存在" />

      <article v-else class="article-content">
        <!-- 文章头部 -->
        <header class="article-header">
          <h1 class="article-title">{{ article.title }}</h1>

          <div class="article-meta">
            <!-- 分类 -->
            <span v-if="article.category" class="meta-item">
              <el-icon><Folder /></el-icon>
              <el-tag size="small">{{ article.category.name }}</el-tag>
            </span>

            <!-- 标签 -->
            <span v-if="article.tags && article.tags.length > 0" class="meta-item tags">
              <el-icon><PriceTag /></el-icon>
              <el-tag
                v-for="tag in article.tags"
                :key="tag.id"
                size="small"
                type="info"
              >
                {{ tag.name }}
              </el-tag>
            </span>

            <!-- 发布时间 -->
            <span class="meta-item">
              <el-icon><Clock /></el-icon>
              {{ formatDate(article.publish_at || article.created_at) }}
            </span>

            <!-- 浏览量 -->
            <span class="meta-item">
              <el-icon><View /></el-icon>
              {{ article.view_count || 0 }} 次浏览
            </span>

            <!-- 作者 -->
            <span v-if="article.author" class="meta-item">
              <el-icon><User /></el-icon>
              {{ article.author.username }}
            </span>
          </div>

          <!-- 封面图 -->
          <div v-if="article.cover_image" class="article-cover">
            <img :src="article.cover_image" :alt="article.title" />
          </div>

          <!-- 摘要 -->
          <div v-if="article.summary" class="article-summary">
            {{ article.summary }}
          </div>
        </header>

        <!-- 文章正文 -->
        <div class="article-body" v-html="article.content"></div>

        <!-- 文章底部 -->
        <footer class="article-footer">
          <el-divider />
          <div class="footer-actions">
            <el-button @click="goBack">
              <el-icon><Back /></el-icon>
              返回列表
            </el-button>
          </div>
        </footer>
      </article>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Folder, PriceTag, Clock, View, User, Back } from '@element-plus/icons-vue'
import api from '@/api'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const article = ref(null)

// 获取文章详情
const fetchArticle = async () => {
  try {
    loading.value = true
    const slug = route.params.slug
    
    if (!slug) {
      ElMessage.error('文章不存在')
      router.push('/articles')
      return
    }

    const response = await api.article.getBySlug(slug)
    article.value = response.data
  } catch (error) {
    console.error('获取文章详情失败:', error)
    ElMessage.error('获取文章详情失败')
    setTimeout(() => {
      router.push('/articles')
    }, 1500)
  } finally {
    loading.value = false
  }
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 返回
const goBack = () => {
  router.back()
}

// 初始化
onMounted(() => {
  fetchArticle()
})
</script>

<style scoped lang="less">
.article-detail-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem;
}

.article-container {
  min-height: 500px;
}

.article-content {
  background: var(--bg-color);
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
}

.article-header {
  margin-bottom: 2rem;
}

.article-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-color);
  margin: 0 0 1rem 0;
  line-height: 1.4;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  font-size: 0.875rem;
  color: var(--text-color-secondary);
  margin-bottom: 1.5rem;

  .meta-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;

    .el-icon {
      font-size: 1rem;
    }

    &.tags {
      display: flex;
      gap: 0.5rem;
    }
  }
}

.article-cover {
  width: 100%;
  margin: 1.5rem 0;
  border-radius: 8px;
  overflow: hidden;

  img {
    width: 100%;
    height: auto;
    display: block;
  }
}

.article-summary {
  padding: 1rem;
  background: var(--bg-color-secondary);
  border-left: 4px solid var(--primary-color);
  border-radius: 4px;
  font-size: 1rem;
  line-height: 1.6;
  color: var(--text-color-secondary);
  margin-top: 1.5rem;
}

.article-body {
  font-size: 1rem;
  line-height: 1.8;
  color: var(--text-color);
  
  :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
    margin: 1.5rem 0 1rem;
    font-weight: 600;
    line-height: 1.4;
  }

  :deep(h1) { font-size: 1.75rem; }
  :deep(h2) { font-size: 1.5rem; }
  :deep(h3) { font-size: 1.25rem; }
  :deep(h4) { font-size: 1.125rem; }
  
  :deep(p) {
    margin: 0.75rem 0;
  }

  :deep(img) {
    max-width: 100%;
    height: auto;
    border-radius: 4px;
    margin: 1rem 0;
  }

  :deep(pre) {
    background: #f6f8fa;
    padding: 1rem;
    border-radius: 4px;
    overflow-x: auto;
    margin: 1rem 0;
  }

  :deep(code) {
    background: #f6f8fa;
    padding: 0.2rem 0.4rem;
    border-radius: 3px;
    font-family: 'Monaco', 'Consolas', monospace;
    font-size: 0.9em;
  }

  :deep(blockquote) {
    margin: 1rem 0;
    padding: 0.5rem 1rem;
    border-left: 4px solid #ddd;
    background: #f9f9f9;
    color: #666;
  }

  :deep(ul), :deep(ol) {
    padding-left: 2rem;
    margin: 1rem 0;
  }

  :deep(li) {
    margin: 0.5rem 0;
  }

  :deep(a) {
    color: var(--primary-color);
    text-decoration: none;
    
    &:hover {
      text-decoration: underline;
    }
  }

  :deep(table) {
    width: 100%;
    border-collapse: collapse;
    margin: 1rem 0;
  }

  :deep(th), :deep(td) {
    border: 1px solid #ddd;
    padding: 0.5rem;
    text-align: left;
  }

  :deep(th) {
    background: #f6f8fa;
    font-weight: 600;
  }
}

.article-footer {
  margin-top: 3rem;

  .footer-actions {
    display: flex;
    justify-content: center;
  }
}

@media (max-width: 768px) {
  .article-detail-page {
    padding: 1rem;
  }

  .article-content {
    padding: 1.5rem;
  }

  .article-title {
    font-size: 1.5rem;
  }

  .article-meta {
    flex-direction: column;
    gap: 0.5rem;
  }
}
</style>

