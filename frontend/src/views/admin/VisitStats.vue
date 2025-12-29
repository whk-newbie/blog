<template>
  <div class="visit-stats">
    <page-header :title="t('stats.visitStats')" />

    <!-- 查询条件 -->
    <el-card class="filter-card" shadow="hover">
      <el-form :inline="true" :model="queryForm">
        <el-form-item :label="t('stats.dateRange')">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            :start-placeholder="t('stats.startDate')"
            :end-placeholder="t('stats.endDate')"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            @change="handleDateRangeChange"
          />
        </el-form-item>
        <el-form-item :label="t('stats.statType')">
          <el-select v-model="queryForm.type" style="width: 150px">
            <el-option label="按日" value="daily" />
            <el-option label="按周" value="weekly" />
            <el-option label="按月" value="monthly" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadVisitStats">
            <el-icon><Search /></el-icon>
            {{ t('stats.query') }}
          </el-button>
          <el-button @click="resetQuery">
            <el-icon><Refresh /></el-icon>
            {{ t('common.reset') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据概览 -->
    <div v-loading="loading" class="overview-cards">
      <el-card class="overview-card" shadow="hover">
        <div class="overview-content">
          <div class="overview-icon pv-icon">
            <el-icon :size="32"><View /></el-icon>
          </div>
          <div class="overview-info">
            <div class="overview-label">{{ t('stats.totalPV') }}</div>
            <div class="overview-value">{{ summary.total_pv || 0 }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="overview-card" shadow="hover">
        <div class="overview-content">
          <div class="overview-icon uv-icon">
            <el-icon :size="32"><User /></el-icon>
          </div>
          <div class="overview-info">
            <div class="overview-label">{{ t('stats.totalUV') }}</div>
            <div class="overview-value">{{ summary.total_uv || 0 }}</div>
          </div>
        </div>
      </el-card>

      <el-card class="overview-card" shadow="hover">
        <div class="overview-content">
          <div class="overview-icon duration-icon">
            <el-icon :size="32"><Clock /></el-icon>
          </div>
          <div class="overview-info">
            <div class="overview-label">{{ t('stats.avgStayDuration') }}</div>
            <div class="overview-value">
              {{ formatDuration(summary.avg_stay_duration) }}
            </div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- 访问趋势图表 -->
    <el-card class="chart-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">
            <el-icon><DataAnalysis /></el-icon>
            {{ t('stats.visitTrend') }}
          </span>
        </div>
      </template>
      <div v-loading="loading" ref="trendChartRef" class="chart-container"></div>
    </el-card>

    <!-- 来源分析和热门文章 -->
    <div class="bottom-section">
      <!-- 来源分析 -->
      <el-card class="referrer-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span class="card-title">
              <el-icon><Connection /></el-icon>
              {{ t('stats.referrers') }}
            </span>
          </div>
        </template>
        <div v-loading="loading" class="referrer-content">
          <div ref="referrerChartRef" class="referrer-chart"></div>
          <div class="referrer-stats">
            <div class="referrer-stat-item">
              <span class="stat-label">{{ t('stats.directAccess') }}</span>
              <span class="stat-value">{{ referrerStats.direct || 0 }}</span>
            </div>
            <div class="referrer-stat-item">
              <span class="stat-label">{{ t('stats.searchEngine') }}</span>
              <span class="stat-value">{{ referrerStats.search_engine || 0 }}</span>
            </div>
            <div class="referrer-stat-item">
              <span class="stat-label">{{ t('stats.externalLink') }}</span>
              <span class="stat-value">{{ referrerStats.external_link || 0 }}</span>
            </div>
          </div>
        </div>
      </el-card>

      <!-- 热门文章 -->
      <el-card class="popular-articles-card" shadow="hover">
        <template #header>
          <div class="card-header">
            <span class="card-title">
              <el-icon><Star /></el-icon>
              {{ t('stats.popularArticlesTitle') }}
            </span>
          </div>
        </template>
        <el-table
          v-loading="loading"
          :data="popularArticles"
          style="width: 100%"
          :show-header="true"
          size="small"
        >
          <el-table-column type="index" label="#" width="50" />
          <el-table-column prop="title" :label="t('stats.articleTitle')" min-width="150" show-overflow-tooltip />
          <el-table-column prop="view_count" :label="t('stats.viewCount')" width="80" align="right" />
          <el-table-column prop="visit_count" :label="t('stats.visitCount')" width="80" align="right" />
          <el-table-column :label="t('stats.avgDuration')" width="100" align="right">
            <template #default="{ row }">
              {{ formatDuration(row.avg_stay_duration) }}
            </template>
          </el-table-column>
        </el-table>
        <el-empty v-if="!loading && popularArticles.length === 0" :description="t('stats.noData')" />
      </el-card>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, nextTick, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import * as echarts from 'echarts'
import dayjs from 'dayjs'
import PageHeader from '@/components/common/PageHeader.vue'
import statsApi from '@/api/stats'
import { ElMessage } from 'element-plus'
import { View, User, Clock, DataAnalysis, Connection, Star, Search, Refresh } from '@element-plus/icons-vue'

const { t } = useI18n()

// 响应式数据
const loading = ref(false)
const dateRange = ref([])
const queryForm = reactive({
  start_date: '',
  end_date: '',
  type: 'daily'
})

const summary = reactive({
  total_pv: 0,
  total_uv: 0,
  avg_stay_duration: 0
})

const dailyStats = ref([])
const referrerStats = reactive({
  direct: 0,
  search_engine: 0,
  external_link: 0,
  top_referrers: []
})

const popularArticles = ref([])

// 图表引用
const trendChartRef = ref(null)
const referrerChartRef = ref(null)
let trendChart = null
let referrerChart = null

// 初始化日期范围（最近30天）
function initDateRange() {
  const endDate = dayjs()
  const startDate = endDate.subtract(30, 'day')
  dateRange.value = [startDate.format('YYYY-MM-DD'), endDate.format('YYYY-MM-DD')]
  queryForm.start_date = startDate.format('YYYY-MM-DD')
  queryForm.end_date = endDate.format('YYYY-MM-DD')
}

// 日期范围变化
function handleDateRangeChange(value) {
  if (value && value.length === 2) {
    queryForm.start_date = value[0]
    queryForm.end_date = value[1]
  } else {
    queryForm.start_date = ''
    queryForm.end_date = ''
  }
}

// 重置查询
function resetQuery() {
  initDateRange()
  loadVisitStats()
}

// 格式化停留时间
function formatDuration(seconds) {
  if (!seconds || seconds === 0) return '0' + t('stats.seconds')
  if (seconds < 60) {
    return Math.floor(seconds) + t('stats.seconds')
  }
  const minutes = Math.floor(seconds / 60)
  const secs = Math.floor(seconds % 60)
  return `${minutes}分${secs}秒`
}

// 加载访问统计
async function loadVisitStats() {
  loading.value = true
  try {
    const params = {
      start_date: queryForm.start_date,
      end_date: queryForm.end_date,
      type: queryForm.type
    }
    
    const data = await statsApi.getVisitStats(params)
    
    // 更新概览数据
    if (data.summary) {
      summary.total_pv = data.summary.total_pv || 0
      summary.total_uv = data.summary.total_uv || 0
      summary.avg_stay_duration = data.summary.avg_stay_duration || 0
    }
    
    // 更新趋势数据
    dailyStats.value = data.daily_stats || []
    renderTrendChart()
    
    // 加载来源统计
    await loadReferrerStats()
    
    // 加载热门文章
    await loadPopularArticles()
  } catch (error) {
    console.error('加载访问统计失败:', error)
    ElMessage.error(t('stats.loadStatsError'))
  } finally {
    loading.value = false
  }
}

// 加载来源统计
async function loadReferrerStats() {
  try {
    const params = {
      start_date: queryForm.start_date,
      end_date: queryForm.end_date
    }
    
    const data = await statsApi.getReferrerStats(params)
    
    referrerStats.direct = data.direct || 0
    referrerStats.search_engine = data.search_engine || 0
    referrerStats.external_link = data.external_link || 0
    referrerStats.top_referrers = data.top_referrers || []
    
    renderReferrerChart()
  } catch (error) {
    console.error('加载来源统计失败:', error)
  }
}

// 加载热门文章
async function loadPopularArticles() {
  try {
    const data = await statsApi.getPopularArticles({
      limit: 10,
      days: 7
    })
    
    popularArticles.value = data || []
  } catch (error) {
    console.error('加载热门文章失败:', error)
  }
}

// 渲染趋势图表
function renderTrendChart() {
  if (!trendChartRef.value) return
  
  nextTick(() => {
    if (!trendChart) {
      trendChart = echarts.init(trendChartRef.value)
    }
    
    const dates = dailyStats.value.map(item => item.date)
    const pvData = dailyStats.value.map(item => item.pv)
    const uvData = dailyStats.value.map(item => item.uv)
    
    const option = {
      tooltip: {
        trigger: 'axis',
        axisPointer: {
          type: 'cross'
        }
      },
      legend: {
        data: [t('stats.pv'), t('stats.uv')]
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: dates
      },
      yAxis: {
        type: 'value'
      },
      series: [
        {
          name: t('stats.pv'),
          type: 'line',
          data: pvData,
          smooth: true,
          itemStyle: {
            color: '#409EFF'
          }
        },
        {
          name: t('stats.uv'),
          type: 'line',
          data: uvData,
          smooth: true,
          itemStyle: {
            color: '#67C23A'
          }
        }
      ]
    }
    
    trendChart.setOption(option)
  })
}

// 渲染来源图表
function renderReferrerChart() {
  if (!referrerChartRef.value) return
  
  nextTick(() => {
    if (!referrerChart) {
      referrerChart = echarts.init(referrerChartRef.value)
    }
    
    const total = referrerStats.direct + referrerStats.search_engine + referrerStats.external_link
    if (total === 0) {
      referrerChart.setOption({
        title: {
          text: t('stats.noData'),
          left: 'center',
          top: 'center',
          textStyle: {
            color: '#999',
            fontSize: 14
          }
        }
      })
      return
    }
    
    const option = {
      tooltip: {
        trigger: 'item',
        formatter: '{a} <br/>{b}: {c} ({d}%)'
      },
      legend: {
        orient: 'vertical',
        left: 'left',
        data: [
          t('stats.directAccess'),
          t('stats.searchEngine'),
          t('stats.externalLink')
        ]
      },
      series: [
        {
          name: t('stats.referrers'),
          type: 'pie',
          radius: ['40%', '70%'],
          avoidLabelOverlap: false,
          itemStyle: {
            borderRadius: 10,
            borderColor: '#fff',
            borderWidth: 2
          },
          label: {
            show: true,
            formatter: '{b}: {c}\n({d}%)'
          },
          emphasis: {
            label: {
              show: true,
              fontSize: '14',
              fontWeight: 'bold'
            }
          },
          data: [
            {
              value: referrerStats.direct,
              name: t('stats.directAccess'),
              itemStyle: { color: '#409EFF' }
            },
            {
              value: referrerStats.search_engine,
              name: t('stats.searchEngine'),
              itemStyle: { color: '#67C23A' }
            },
            {
              value: referrerStats.external_link,
              name: t('stats.externalLink'),
              itemStyle: { color: '#E6A23C' }
            }
          ]
        }
      ]
    }
    
    referrerChart.setOption(option)
  })
}

// 窗口大小变化时调整图表
function handleResize() {
  if (trendChart) {
    trendChart.resize()
  }
  if (referrerChart) {
    referrerChart.resize()
  }
}

onMounted(() => {
  initDateRange()
  loadVisitStats()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  if (trendChart) {
    trendChart.dispose()
  }
  if (referrerChart) {
    referrerChart.dispose()
  }
})
</script>

<style lang="less" scoped>
.visit-stats {
  padding: 16px;
  
  .filter-card {
    margin-bottom: 16px;
    
    :deep(.el-card__body) {
      padding: 16px;
    }
  }
  
  .overview-cards {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 10px;
    margin-bottom: 16px;
    
    .overview-card {
      :deep(.el-card__body) {
        padding: 12px 16px;
      }
      
      .overview-content {
        display: flex;
        align-items: center;
        gap: 12px;
        
        .overview-icon {
          width: 40px;
          height: 40px;
          border-radius: 6px;
          display: flex;
          align-items: center;
          justify-content: center;
          flex-shrink: 0;
          
          :deep(.el-icon) {
            font-size: 20px;
          }
          
          &.pv-icon {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: #fff;
          }
          
          &.uv-icon {
            background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
            color: #fff;
          }
          
          &.duration-icon {
            background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
            color: #fff;
          }
        }
        
        .overview-info {
          flex: 1;
          min-width: 0;
          
          .overview-label {
            font-size: 12px;
            color: #909399;
            margin-bottom: 4px;
          }
          
          .overview-value {
            font-size: 20px;
            font-weight: bold;
            color: #303133;
            word-break: break-all;
            line-height: 1.2;
          }
        }
      }
    }
  }
  
  .chart-card {
    margin-bottom: 16px;
    
    :deep(.el-card__body) {
      padding: 16px;
    }
    
    .chart-container {
      width: 100%;
      height: 320px;
    }
  }
  
  .bottom-section {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 16px;
    
    .referrer-card {
      :deep(.el-card__body) {
        padding: 16px;
      }
      
      .referrer-content {
        display: flex;
        gap: 16px;
        
        .referrer-chart {
          width: 240px;
          height: 240px;
          flex-shrink: 0;
        }
        
        .referrer-stats {
          flex: 1;
          display: flex;
          flex-direction: column;
          justify-content: center;
          gap: 12px;
          min-width: 0;
          
          .referrer-stat-item {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 10px;
            background: #f5f7fa;
            border-radius: 4px;
            
            .stat-label {
              font-size: 13px;
              color: #606266;
              flex-shrink: 0;
            }
            
            .stat-value {
              font-size: 18px;
              font-weight: bold;
              color: #303133;
            }
          }
        }
      }
    }
    
    .popular-articles-card {
      :deep(.el-card__body) {
        padding: 16px;
      }
      
      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
      }
      
      :deep(.el-table) {
        font-size: 13px;
        
        .el-table__cell {
          padding: 8px 0;
        }
      }
    }
  }
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    
    .card-title {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 15px;
      font-weight: 500;
    }
  }
}

@media (max-width: 1200px) {
  .visit-stats {
    .bottom-section {
      grid-template-columns: 1fr;
      
      .referrer-card {
        .referrer-content {
          .referrer-chart {
            width: 200px;
            height: 200px;
          }
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .visit-stats {
    padding: 12px;
    
    .overview-cards {
      grid-template-columns: 1fr;
      gap: 10px;
    }
    
    .chart-card {
      .chart-container {
        height: 280px;
      }
    }
    
    .bottom-section {
      gap: 12px;
      
      .referrer-card {
        .referrer-content {
          flex-direction: column;
          
          .referrer-chart {
            width: 100%;
            max-width: 300px;
            height: 200px;
            margin: 0 auto;
          }
        }
      }
    }
  }
}
</style>

