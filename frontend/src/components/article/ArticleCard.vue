<template>
  <div class="article-card" @click="handleClick">
    <!-- 封面图 -->
    <div v-if="article.cover_image" class="article-cover">
      <img v-lazy="article.cover_image" :alt="article.title" loading="lazy" />
      <div v-if="article.is_top" class="top-badge">
        <el-tag type="danger" size="small">{{ t('article.isTop') }}</el-tag>
      </div>
    </div>

    <!-- 文章信息 -->
    <div class="article-info">
      <!-- 标题 -->
      <h3 class="article-title">{{ article.title }}</h3>

      <!-- 摘要 -->
      <p v-if="article.summary" class="article-summary">
        {{ article.summary }}
      </p>

      <!-- 元信息 -->
      <div class="article-meta">
        <!-- 分类 -->
        <span v-if="article.category" class="meta-item">
          <el-icon><Folder /></el-icon>
          {{ article.category.name }}
        </span>

        <!-- 标签 -->
        <span v-if="article.tags && article.tags.length > 0" class="meta-item tags">
          <el-icon><PriceTag /></el-icon>
          <el-tag
            v-for="tag in article.tags.slice(0, 3)"
            :key="tag.id"
            size="small"
            class="tag"
          >
            {{ tag.name }}
          </el-tag>
        </span>

        <!-- 浏览量 -->
        <span class="meta-item">
          <el-icon><View /></el-icon>
          {{ article.view_count || 0 }}
        </span>

        <!-- 发布时间 -->
        <span class="meta-item">
          <el-icon><Clock /></el-icon>
          {{ formatDate(article.publish_at || article.created_at) }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Folder, PriceTag, View, Clock } from '@element-plus/icons-vue'
import { useI18n } from 'vue-i18n'

const { t, locale } = useI18n()

const props = defineProps({
  article: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['click'])

const handleClick = () => {
  emit('click')
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now - date
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return t('common.today')
  if (days === 1) return t('common.yesterday')
  if (days < 7) return `${days}${t('common.daysAgo')}`
  if (days < 30) return `${Math.floor(days / 7)}${t('common.weeksAgo')}`
  if (days < 365) return `${Math.floor(days / 30)}${t('common.monthsAgo')}`
  
  return date.toLocaleDateString(locale.value === 'zh-CN' ? 'zh-CN' : 'en-US')
}
</script>

<style scoped lang="less">
.article-card {
  background: var(--bg-color);
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
  cursor: pointer;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  }
}

.article-cover {
  position: relative;
  width: 100%;
  height: 200px;
  overflow: hidden;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    transition: transform 0.3s ease;
  }

  &:hover img {
    transform: scale(1.05);
  }

  .top-badge {
    position: absolute;
    top: 10px;
    right: 10px;
  }
}

.article-info {
  padding: 1.5rem;
}

.article-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-color);
  margin: 0 0 0.75rem 0;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.article-summary {
  font-size: 0.875rem;
  color: var(--text-color-secondary);
  line-height: 1.6;
  margin: 0 0 1rem 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  text-overflow: ellipsis;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  font-size: 0.75rem;
  color: var(--text-color-secondary);

  .meta-item {
    display: flex;
    align-items: center;
    gap: 0.25rem;

    .el-icon {
      font-size: 0.875rem;
    }

    &.tags {
      flex: 1;
      gap: 0.5rem;

      .tag {
        margin: 0;
      }
    }
  }
}
</style>

