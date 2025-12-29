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
  Promotion
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
    if (res.data) {
      stats.value = {
        article_count: res.data.article_count || 0,
        category_count: res.data.category_count || 0,
        tag_count: res.data.tag_count || 0
      }
      recentArticles.value = res.data.recent_articles || []
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
  padding: 20px;
}

.stats-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 20px;
  margin-bottom: 30px;
}

.stat-card {
  cursor: pointer;
  transition: transform 0.3s;

  &:hover {
    transform: translateY(-4px);
  }
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
}

.article-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.category-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.tag-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-info {
  flex: 1;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 32px;
  font-weight: bold;
  color: #303133;
}

.recent-articles-card {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
}

.no-cover {
  font-size: 12px;
  color: #909399;
  display: block;
  text-align: center;
}

.empty-state {
  padding: 40px 0;
}

:deep(.el-table) {
  font-size: 14px;
}

:deep(.el-button + .el-button) {
  margin-left: 8px;
}
</style>
