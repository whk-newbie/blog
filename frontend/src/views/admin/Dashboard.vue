<template>
  <div class="dashboard">
    <page-header title="仪表盘" />

    <!-- 统计卡片 -->
    <div v-loading="loading" class="stats-cards">
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon article-icon">
            <el-icon :size="32"><Document /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">文章总数</div>
            <div class="stat-value">{{ stats.article_count || 0 }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon category-icon">
            <el-icon :size="32"><Folder /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">分类总数</div>
            <div class="stat-value">{{ stats.category_count || 0 }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon tag-icon">
            <el-icon :size="32"><PriceTag /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">标签总数</div>
            <div class="stat-value">{{ stats.tag_count || 0 }}</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 数据图表占位 -->
    <el-card class="chart-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">
            <el-icon><DataAnalysis /></el-icon>
            数据统计图表
          </span>
        </div>
      </template>
      <div class="chart-placeholder">
        <el-icon :size="64" class="placeholder-icon"><DataAnalysis /></el-icon>
        <p class="placeholder-text">图表占位区域</p>
        <p class="placeholder-desc">未来可集成 ECharts 等图表库展示访问趋势、文章发布统计等数据</p>
      </div>
    </el-card>

    <!-- 最近文章列表 -->
    <el-card class="recent-articles-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">
            <el-icon><Clock /></el-icon>
            最近文章
          </span>
          <el-button text @click="goToArticles">查看全部</el-button>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="recentArticles"
        style="width: 100%"
        :show-header="true"
      >
        <el-table-column prop="id" label="ID" width="80" />
        
        <el-table-column label="封面" width="100">
          <template #default="{ row }">
            <el-image
              v-if="row.cover_image"
              :src="row.cover_image"
              style="width: 60px; height: 60px; border-radius: 4px"
              fit="cover"
            />
            <span v-else class="no-cover">无封面</span>
          </template>
        </el-table-column>

        <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />

        <el-table-column label="分类" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.category" size="small">
              {{ row.category.name }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">
              {{ row.status === 'published' ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column label="浏览量" width="100">
          <template #default="{ row }">
            <span>{{ row.view_count || 0 }}</span>
          </template>
        </el-table-column>

        <el-table-column label="创建时间" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text size="small" @click="handleEdit(row.id)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
            <el-button
              v-if="row.status === 'draft'"
              type="success"
              text
              size="small"
              @click="handlePublish(row.id)"
            >
              <el-icon><Promotion /></el-icon>
              发布
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="!loading && recentArticles.length === 0" class="empty-state">
        <empty-state description="暂无文章" />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Document,
  Folder,
  PriceTag,
  Clock,
  Edit,
  Promotion,
  DataAnalysis
} from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import statsApi from '@/api/stats'
import articleApi from '@/api/article'

const router = useRouter()

// 数据
const loading = ref(false)
const stats = ref({
  article_count: 0,
  category_count: 0,
  tag_count: 0
})
const recentArticles = ref([])

// 加载统计数据
const loadStats = async () => {
  loading.value = true
  try {
    const res = await statsApi.getDashboardStats()
    if (res) {
      stats.value = {
        article_count: res.article_count || 0,
        category_count: res.category_count || 0,
        tag_count: res.tag_count || 0
      }
      recentArticles.value = res.recent_articles || []
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
    ElMessage.error('加载统计数据失败')
  } finally {
    loading.value = false
  }
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 跳转到文章列表
const goToArticles = () => {
  router.push('/admin/articles')
}

// 编辑文章
const handleEdit = (id) => {
  router.push(`/admin/articles/edit/${id}`)
}

// 发布文章
const handlePublish = async (id) => {
  try {
    await ElMessageBox.confirm('确定要发布这篇文章吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'info'
    })

    await articleApi.publish(id)
    ElMessage.success('发布成功')
    loadStats() // 重新加载数据
  } catch (error) {
    if (error !== 'cancel') {
      console.error('发布失败:', error)
      ElMessage.error(error.message || '发布失败')
    }
  }
}

// 初始化
onMounted(() => {
  loadStats()
})
</script>

<style lang="less" scoped>
.dashboard {
  padding: 0;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 24px;
  margin-bottom: 32px;
}

.stat-card {
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border-radius: 12px;
  overflow: hidden;
  position: relative;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
    opacity: 0;
    transition: opacity 0.3s;
  }

  &:hover {
    transform: translateY(-6px);
    box-shadow: var(--shadow-blue);

    &::before {
      opacity: 1;
    }

    .stat-icon {
      transform: scale(1.1) rotate(5deg);
    }
  }
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 24px;
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 28px;
  transition: all 0.3s;
}

.article-icon {
  background: linear-gradient(135deg, #2563eb 0%, #3b82f6 100%);
  box-shadow: 0 4px 14px rgba(37, 99, 235, 0.3);

  &::before {
    background: linear-gradient(90deg, #2563eb, #3b82f6);
  }
}

.category-icon {
  background: linear-gradient(135deg, #06b6d4 0%, #0891b2 100%);
  box-shadow: 0 4px 14px rgba(6, 182, 212, 0.3);

  &::before {
    background: linear-gradient(90deg, #06b6d4, #0891b2);
  }
}

.tag-icon {
  background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
  box-shadow: 0 4px 14px rgba(139, 92, 246, 0.3);

  &::before {
    background: linear-gradient(90deg, #8b5cf6, #6366f1);
  }
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: var(--text-secondary);
  margin-bottom: 8px;
  font-weight: 500;
}

.stat-value {
  font-size: 36px;
  font-weight: 700;
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--secondary-color) 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  letter-spacing: -1px;
}

.chart-card {
  margin-top: 32px;
  margin-bottom: 32px;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid var(--border-light);

  &:hover {
    box-shadow: var(--shadow-md);
    border-color: var(--border-blue);
  }
}

.chart-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
  min-height: 300px;
  background: var(--bg-secondary);
  border-radius: 8px;
}

.placeholder-icon {
  color: var(--text-disabled);
  margin-bottom: 16px;
  opacity: 0.5;
}

.placeholder-text {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-secondary);
  margin: 0 0 8px 0;
}

.placeholder-desc {
  font-size: 14px;
  color: var(--text-disabled);
  margin: 0;
  text-align: center;
  max-width: 400px;
}

.recent-articles-card {
  margin-top: 0;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid var(--border-light);

  &:hover {
    box-shadow: var(--shadow-md);
    border-color: var(--border-blue);
  }
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px;
  border-bottom: 1px solid var(--border-light);
  background: var(--bg-blue-light);
}

.card-title {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 600;
  color: var(--text-color);

  .el-icon {
    color: var(--primary-color);
    font-size: 20px;
  }
}

.no-cover {
  font-size: 12px;
  color: var(--text-disabled);
  display: block;
  text-align: center;
}

.empty-state {
  padding: 60px 0;
}

:deep(.el-card__body) {
  padding: 24px;
}

:deep(.el-table) {
  font-size: 14px;

  .el-table__header th {
    background: var(--bg-secondary);
    color: var(--text-color);
    font-weight: 600;
  }

  .el-table__row:hover {
    background: var(--bg-blue-light);
  }

  .el-table__cell {
    padding: 14px 0;
  }
}

:deep(.el-button) {
  border-radius: 6px;
  font-weight: 500;
  transition: all 0.3s;

  &.el-button--primary {
    background: var(--primary-color);
    border-color: var(--primary-color);

    &:hover {
      background: var(--primary-light);
      box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
    }
  }

  &.el-button--success {
    background: var(--success-color);
    border-color: var(--success-color);
  }

  & + .el-button {
    margin-left: 8px;
  }
}

:deep(.el-image) {
  border-radius: 6px;
  overflow: hidden;
}

:deep(.el-tag) {
  border-radius: 6px;
  font-weight: 500;
}
</style>
