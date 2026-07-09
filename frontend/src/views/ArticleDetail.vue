<template>
  <div class="article-detail-page">
    <div v-loading="loading" class="article-container">
      <el-empty v-if="!loading && !article" :description="t('app.articleNotFound')" />

      <div v-else class="article-layout">
        <!-- 主要内容区 -->
        <article class="article-content">
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
                {{ article.view_count || 0 }} {{ t('common.views') }}
              </span>

              <!-- 作者 -->
              <span v-if="article.author" class="meta-item">
                <el-icon><User /></el-icon>
                {{ article.author.username }}
              </span>
            </div>

            <!-- 封面图 -->
            <div v-if="article.cover_image" class="article-cover">
              <img v-lazy="article.cover_image" :alt="article.title" loading="lazy" />
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
                {{ t('common.backToList') }}
              </el-button>
            </div>
          </footer>
        </article>

        <!-- 目录侧边栏 -->
        <aside class="article-aside">
          <table-of-contents ref="tocRef" container=".article-body" />
        </aside>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Folder, PriceTag, Clock, View, User, Back } from '@element-plus/icons-vue'
import api from '@/api'
import TableOfContents from '@/components/article/TableOfContents.vue'

const route = useRoute()
const router = useRouter()
const { t, locale } = useI18n()

const loading = ref(false)
const article = ref(null)
const tocRef = ref(null)

// 获取文章详情
const fetchArticle = async () => {
  try {
    loading.value = true
    const slug = route.params.slug
    
    if (!slug) {
      ElMessage.error(t('app.articleNotFound'))
      router.push('/articles')
      return
    }

    const response = await api.article.getBySlug(slug)
    article.value = response

    // 等待DOM更新后提取目录
    await nextTick()
    setTimeout(() => {
      if (tocRef.value) {
        tocRef.value.extractHeadings()
      }
    }, 300)
  } catch (error) {
    console.error('获取文章详情失败:', error)
    ElMessage.error(t('article.loadError'))
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
  return date.toLocaleString(locale.value === 'zh-CN' ? 'zh-CN' : 'en-US', {
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
  max-width: 1100px;
  margin: 0 auto;
  padding: var(--spacing-lg);
}

.article-container {
  min-height: 500px;
}

.article-layout {
  display: flex;
  gap: var(--spacing-lg);
  align-items: flex-start;
}

.article-content {
  flex: 1;
  min-width: 0;
  background: var(--card-bg);
  border-radius: var(--border-radius-sm);
  padding: var(--spacing-xl);
  box-shadow: var(--shadow-sm);
}

.article-aside {
  width: 260px;
  flex-shrink: 0;
}

.article-header {
  margin-bottom: var(--spacing-xl);
}

.article-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-heading);
  margin: 0 0 var(--spacing-md) 0;
  line-height: 1.4;
  font-family: var(--font-family);
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-md);
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  margin-bottom: var(--spacing-lg);

  .meta-item {
    display: flex;
    align-items: center;
    gap: var(--spacing-xs);

    .el-icon {
      font-size: 1rem;
    }

    &.tags {
      display: flex;
      gap: var(--spacing-xs);
    }
  }
}

.article-cover {
  width: 100%;
  margin: var(--spacing-lg) 0;
  border-radius: var(--border-radius-sm);
  overflow: hidden;

  img {
    width: 100%;
    height: auto;
    display: block;
  }
}

.article-summary {
  padding: var(--spacing-md);
  background: var(--bg-secondary);
  border-left: 4px solid var(--primary-color);
  border-radius: var(--border-radius-sm);
  font-size: var(--font-size-md);
  line-height: 1.6;
  color: var(--text-secondary);
  margin-top: var(--spacing-lg);
  font-family: var(--font-family);
}

.article-body {
  font-size: var(--font-size-md);
  line-height: 1.8;
  color: var(--text-color);
  font-family: var(--font-family);

  :deep(h1), :deep(h2), :deep(h3), :deep(h4), :deep(h5), :deep(h6) {
    margin: var(--spacing-lg) 0 var(--spacing-md);
    font-weight: 600;
    line-height: 1.4;
    color: var(--text-heading);
  }

  :deep(h1) { font-size: 1.75rem; }
  :deep(h2) { font-size: 1.5rem; }
  :deep(h3) { font-size: 1.25rem; }
  :deep(h4) { font-size: 1.125rem; }

  :deep(p) {
    margin: var(--spacing-sm) 0;
  }

  :deep(img) {
    max-width: 100%;
    height: auto;
    border-radius: var(--border-radius-sm);
    margin: var(--spacing-md) 0;
  }

  :deep(pre) {
    background: var(--bg-secondary);
    padding: var(--spacing-md);
    border-radius: var(--border-radius-sm);
    overflow-x: auto;
    margin: var(--spacing-md) 0;
    font-family: var(--font-mono);
  }

  :deep(code) {
    background: var(--bg-secondary);
    padding: 0.2rem 0.4rem;
    border-radius: var(--border-radius-sm);
    font-family: var(--font-mono);
    font-size: 0.9em;
    color: var(--text-color);
  }

  :deep(blockquote) {
    margin: var(--spacing-md) 0;
    padding: 0.5rem var(--spacing-md);
    border-left: 4px solid var(--border-color);
    background: var(--bg-secondary);
    color: var(--text-secondary);
    font-family: var(--font-family);
  }

  :deep(ul), :deep(ol) {
    padding-left: 2rem;
    margin: var(--spacing-md) 0;
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
    margin: var(--spacing-md) 0;
  }

  :deep(th), :deep(td) {
    border: 1px solid var(--border-color);
    padding: 0.5rem;
    text-align: left;
  }

  :deep(th) {
    background: var(--bg-secondary);
    font-weight: 600;
  }
}

.article-footer {
  margin-top: var(--spacing-xl);

  .footer-actions {
    display: flex;
    justify-content: center;
  }
}

@media (max-width: 1200px) {
  .article-layout {
    flex-direction: column;
  }

  .article-aside {
    width: 100%;
    order: -1;
  }
}

@media (max-width: 768px) {
  .article-detail-page {
    padding: var(--spacing-md);
  }

  .article-content {
    padding: var(--spacing-lg);
  }

  .article-title {
    font-size: 1.5rem;
  }

  .article-meta {
    flex-direction: column;
    gap: var(--spacing-xs);
  }
}
</style>

