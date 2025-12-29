<template>
  <div class="dashboard">
    <page-header :title="t('nav.dashboard')" />

    <!-- 数据图表 -->
    <el-card class="chart-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">
            <el-icon><DataAnalysis /></el-icon>
            {{ t('stats.chartTitle') }}
          </span>
          <div class="header-stats">
            <div class="header-stat-item">
              <el-icon class="stat-icon-small article-icon-small"><Document /></el-icon>
              <span class="stat-label-small">{{ t('stats.articleCount') }}</span>
              <span class="stat-value-small">{{ stats.article_count || 0 }}</span>
            </div>
            <div class="header-stat-item">
              <el-icon class="stat-icon-small category-icon-small"><Folder /></el-icon>
              <span class="stat-label-small">{{ t('stats.categoryCount') }}</span>
              <span class="stat-value-small">{{ stats.category_count || 0 }}</span>
            </div>
            <div class="header-stat-item">
              <el-icon class="stat-icon-small tag-icon-small"><PriceTag /></el-icon>
              <span class="stat-label-small">{{ t('stats.tagCount') }}</span>
              <span class="stat-value-small">{{ stats.tag_count || 0 }}</span>
            </div>
          </div>
        </div>
      </template>
      <div v-loading="chartLoading" ref="chartRef" class="chart-container"></div>
    </el-card>

    <!-- 最近文章列表 -->
    <el-card class="recent-articles-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">
            <el-icon><Clock /></el-icon>
            {{ t('stats.recentArticles') }}
          </span>
          <el-button text @click="goToArticles">{{ t('stats.viewAll') }}</el-button>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="recentArticles"
        style="width: 100%"
        :show-header="true"
      >
        <el-table-column prop="id" label="ID" width="80" />
        
        <el-table-column :label="t('stats.cover')" width="100">
          <template #default="{ row }">
            <el-image
              v-if="row.cover_image"
              :src="row.cover_image"
              style="width: 60px; height: 60px; border-radius: 4px"
              fit="cover"
            />
            <span v-else class="no-cover">{{ t('stats.noCover') }}</span>
          </template>
        </el-table-column>

        <el-table-column prop="title" :label="t('article.title')" min-width="200" show-overflow-tooltip />

        <el-table-column :label="t('article.category')" width="120">
          <template #default="{ row }">
            <el-tag v-if="row.category" size="small">
              {{ row.category.name }}
            </el-tag>
            <span v-else>-</span>
          </template>
        </el-table-column>

        <el-table-column :label="t('article.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">
              {{ row.status === 'published' ? t('article.published') : t('article.draft') }}
            </el-tag>
          </template>
        </el-table-column>

        <el-table-column :label="t('article.viewCount')" width="100">
          <template #default="{ row }">
            <span>{{ row.view_count || 0 }}</span>
          </template>
        </el-table-column>

        <el-table-column :label="t('stats.createTime')" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>

        <el-table-column :label="t('common.operation')" width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" text size="small" @click="handleEdit(row.id)">
              <el-icon><Edit /></el-icon>
              {{ t('common.edit') }}
            </el-button>
            <el-button
              v-if="row.status === 'draft'"
              type="success"
              text
              size="small"
              @click="handlePublish(row.id)"
            >
              <el-icon><Promotion /></el-icon>
              {{ t('article.publishArticle') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="!loading && recentArticles.length === 0" class="empty-state">
        <EmptyState :description="t('stats.noArticles')" />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted, nextTick, onBeforeUnmount, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import * as echarts from 'echarts'
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
import { useThemeStore } from '@/store/theme'

const router = useRouter()
const { t } = useI18n()
const themeStore = useThemeStore()

// 数据
const loading = ref(false)
const chartLoading = ref(false)
const stats = ref({
  article_count: 0,
  category_count: 0,
  tag_count: 0,
  chart_data: null
})
const recentArticles = ref([])
const chartRef = ref(null)
let chartInstance = null

// 加载统计数据
const loadStats = async () => {
  loading.value = true
  chartLoading.value = true
  try {
    const res = await statsApi.getDashboardStats()
    if (res) {
      stats.value = {
        article_count: res.article_count || 0,
        category_count: res.category_count || 0,
        tag_count: res.tag_count || 0,
        chart_data: res.chart_data || null
      }
      recentArticles.value = res.recent_articles || []
      
      // 渲染图表
      if (res.chart_data) {
        await nextTick()
        renderChart(res.chart_data)
      }
    }
  } catch (error) {
    console.error('加载统计数据失败:', error)
    ElMessage.error(t('stats.loadStatsError'))
  } finally {
    loading.value = false
    chartLoading.value = false
  }
}

// 获取图表颜色配置（根据主题）
const getChartColors = () => {
  const isDark = themeStore.theme === 'dark'
  
  if (isDark) {
    return {
      pv: '#6b9bd2',        // 柔和的蓝色
      uv: '#6bb89c',        // 柔和的绿色
      publish: '#8b5cf6',   // 柔和的紫色
      text: '#d4d4d4',      // 柔和的浅灰色
      grid: '#3a3a3a',      // 深灰色网格
      bg: 'transparent'     // 透明背景
    }
  } else {
    return {
      pv: '#2563eb',        // 蓝色
      uv: '#10b981',        // 绿色
      publish: '#8b5cf6',   // 紫色
      text: '#1e293b',      // 深灰色
      grid: '#e2e8f0',      // 浅灰色网格
      bg: 'transparent'     // 透明背景
    }
  }
}

// 渲染图表
const renderChart = (chartData) => {
  if (!chartRef.value) return
  
  // 初始化图表
  if (!chartInstance) {
    chartInstance = echarts.init(chartRef.value)
  }
  
  // 处理访问趋势数据
  const visitTrend = chartData.visit_trend || []
  const visitDates = visitTrend.map(item => item.date)
  const pvData = visitTrend.map(item => item.pv || 0)
  const uvData = visitTrend.map(item => item.uv || 0)
  
  // 处理文章发布趋势数据
  const publishTrend = chartData.publish_trend || []
  const publishDates = publishTrend.map(item => item.date)
  const publishCounts = publishTrend.map(item => item.count || 0)
  
  // 合并日期，确保所有日期都在图表中显示
  const allDates = [...new Set([...visitDates, ...publishDates])].sort()
  
  // 创建数据映射，填充缺失的日期
  const pvMap = new Map(visitTrend.map(item => [item.date, item.pv || 0]))
  const uvMap = new Map(visitTrend.map(item => [item.date, item.uv || 0]))
  const publishMap = new Map(publishTrend.map(item => [item.date, item.count || 0]))
  
  const finalPvData = allDates.map(date => pvMap.get(date) || 0)
  const finalUvData = allDates.map(date => uvMap.get(date) || 0)
  const finalPublishData = allDates.map(date => publishMap.get(date) || 0)
  
  const colors = getChartColors()
  const isDark = themeStore.theme === 'dark'
  
  const option = {
    backgroundColor: colors.bg,
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'cross'
      },
      backgroundColor: isDark ? 'rgba(37, 37, 37, 0.9)' : 'rgba(255, 255, 255, 0.9)',
      borderColor: isDark ? '#3a3a3a' : '#e2e8f0',
      textStyle: {
        color: colors.text
      }
    },
    legend: {
      data: [t('stats.pv'), t('stats.uv'), t('stats.articleCount')],
      top: 10,
      textStyle: {
        color: colors.text
      }
    },
    grid: {
      left: '3%',
      right: '4%',
      bottom: '3%',
      top: '15%',
      containLabel: true
    },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: allDates,
      axisLabel: {
        rotate: 45,
        interval: Math.floor(allDates.length / 10), // 显示部分日期标签，避免拥挤
        color: colors.text
      },
      axisLine: {
        lineStyle: {
          color: colors.grid
        }
      },
      splitLine: {
        show: false
      }
    },
    yAxis: [
      {
        type: 'value',
        name: t('stats.pv'),
        position: 'left',
        axisLabel: {
          formatter: '{value}',
          color: colors.text
        },
        axisLine: {
          lineStyle: {
            color: colors.grid
          }
        },
        splitLine: {
          lineStyle: {
            color: isDark ? '#2d2d2d' : '#f1f5f9',
            type: 'dashed'
          }
        }
      },
      {
        type: 'value',
        name: t('stats.articleCount'),
        position: 'right',
        axisLabel: {
          formatter: '{value}',
          color: colors.text
        },
        axisLine: {
          lineStyle: {
            color: colors.grid
          }
        },
        splitLine: {
          show: false
        }
      }
    ],
    series: [
      {
        name: t('stats.pv'),
        type: 'line',
        data: finalPvData,
        smooth: true,
        itemStyle: {
          color: colors.pv
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: isDark ? [
              { offset: 0, color: 'rgba(107, 155, 210, 0.25)' },
              { offset: 1, color: 'rgba(107, 155, 210, 0.05)' }
            ] : [
              { offset: 0, color: 'rgba(37, 99, 235, 0.3)' },
              { offset: 1, color: 'rgba(37, 99, 235, 0.05)' }
            ]
          }
        }
      },
      {
        name: t('stats.uv'),
        type: 'line',
        data: finalUvData,
        smooth: true,
        itemStyle: {
          color: colors.uv
        },
        areaStyle: {
          color: {
            type: 'linear',
            x: 0,
            y: 0,
            x2: 0,
            y2: 1,
            colorStops: isDark ? [
              { offset: 0, color: 'rgba(107, 184, 156, 0.25)' },
              { offset: 1, color: 'rgba(107, 184, 156, 0.05)' }
            ] : [
              { offset: 0, color: 'rgba(16, 185, 129, 0.3)' },
              { offset: 1, color: 'rgba(16, 185, 129, 0.05)' }
            ]
          }
        }
      },
      {
        name: t('stats.articleCount'),
        type: 'bar',
        yAxisIndex: 1,
        data: finalPublishData,
        itemStyle: {
          color: colors.publish,
          borderRadius: [4, 4, 0, 0]
        }
      }
    ]
  }
  
  chartInstance.setOption(option)
  
  // 响应式调整
  window.addEventListener('resize', () => {
    if (chartInstance) {
      chartInstance.resize()
    }
  })
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
    await ElMessageBox.confirm(t('stats.publishConfirm'), t('common.tip'), {
      confirmButtonText: t('common.confirm'),
      cancelButtonText: t('common.cancel'),
      type: 'info'
    })

    await articleApi.publish(id)
    ElMessage.success(t('stats.publishSuccess'))
    loadStats() // 重新加载数据
  } catch (error) {
    if (error !== 'cancel') {
      console.error('发布失败:', error)
      ElMessage.error(error.message || t('stats.publishError'))
    }
  }
}

// 监听主题变化，重新渲染图表
watch(() => themeStore.theme, () => {
  if (chartInstance && stats.value.chart_data) {
    renderChart(stats.value.chart_data)
  }
})

// 初始化
onMounted(() => {
  loadStats()
})

// 组件卸载时销毁图表
onBeforeUnmount(() => {
  if (chartInstance) {
    chartInstance.dispose()
    chartInstance = null
  }
  // 移除窗口大小调整监听器
  window.removeEventListener('resize', () => {
    if (chartInstance) {
      chartInstance.resize()
    }
  })
})
</script>

<style lang="less" scoped>
.dashboard {
  padding: 0;
}

.header-stats {
  display: flex;
  align-items: center;
  gap: 24px;
  
  @media (max-width: 768px) {
    gap: 16px;
  }
}

.header-stat-item {
  display: flex;
  align-items: center;
  gap: 6px;
}

.stat-icon-small {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 14px;
  flex-shrink: 0;
}

.article-icon-small {
  background: linear-gradient(135deg, #2563eb 0%, #3b82f6 100%);
}

.category-icon-small {
  background: linear-gradient(135deg, #06b6d4 0%, #0891b2 100%);
}

.tag-icon-small {
  background: linear-gradient(135deg, #8b5cf6 0%, #6366f1 100%);
}

.stat-label-small {
  font-size: 12px;
  color: var(--text-secondary);
  font-weight: 500;
  white-space: nowrap;
}

.stat-value-small {
  font-size: 16px;
  font-weight: 700;
  color: var(--primary-color);
  white-space: nowrap;
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

.chart-container {
  width: 100%;
  height: 400px;
  min-height: 400px;
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
  padding: 16px 24px;
  border-bottom: 1px solid var(--border-light);
  background: var(--bg-blue-light);
  flex-wrap: wrap;
  gap: 16px;
  
  @media (max-width: 768px) {
    padding: 12px 16px;
  }
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
