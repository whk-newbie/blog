<template>
  <div class="crawler-monitor">
    <page-header :title="t('crawler.title')" />

    <!-- 连接状态提示 -->
    <el-alert
      v-if="!wsStore.connected"
      :title="t('crawler.wsDisconnected')"
      type="warning"
      :closable="false"
      show-icon
      style="margin-bottom: 20px"
    />

    <!-- 查询条件 -->
    <el-card class="filter-card" shadow="hover">
      <el-form :inline="true" :model="queryForm">
        <el-form-item :label="t('crawler.status')">
          <el-select v-model="queryForm.status" style="width: 150px" clearable>
            <el-option :label="t('crawler.statusRunning')" value="running" />
            <el-option :label="t('crawler.statusCompleted')" value="completed" />
            <el-option :label="t('crawler.statusFailed')" value="failed" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('crawler.taskId')">
          <el-input
            v-model="queryForm.task_id"
            :placeholder="t('crawler.taskIdPlaceholder')"
            style="width: 200px"
            clearable
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadTasks">
            <el-icon><Search /></el-icon>
            {{ t('common.search') }}
          </el-button>
          <el-button @click="resetQuery">
            <el-icon><Refresh /></el-icon>
            {{ t('common.reset') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 任务列表 -->
    <el-card class="task-list-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">
            <el-icon><List /></el-icon>
            {{ t('crawler.taskList') }}
          </span>
          <el-button type="primary" size="small" @click="loadTasks">
            <el-icon><Refresh /></el-icon>
            {{ t('common.refresh') }}
          </el-button>
        </div>
      </template>

      <el-table
        v-loading="loading"
        :data="displayTasks"
        style="width: 100%"
        :show-header="true"
        stripe
      >
        <el-table-column type="index" label="#" width="60" />
        <el-table-column prop="task_id" :label="t('crawler.taskId')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="task_name" :label="t('crawler.taskName')" min-width="200" show-overflow-tooltip />
        <el-table-column :label="t('crawler.status')" width="120" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="t('crawler.progress')" width="200" align="center">
          <template #default="{ row }">
            <el-progress
              :percentage="row.progress"
              :status="getProgressStatus(row.status)"
              :stroke-width="8"
            />
            <span style="font-size: 12px; color: #909399; margin-left: 8px">
              {{ row.progress }}%
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="message" :label="t('crawler.message')" min-width="200" show-overflow-tooltip />
        <el-table-column :label="t('crawler.startTime')" width="180" align="center">
          <template #default="{ row }">
            {{ formatTime(row.start_time) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('crawler.duration')" width="120" align="center">
          <template #default="{ row }">
            {{ formatDuration(row.duration) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('common.operation')" width="120" align="center" fixed="right">
          <template #default="{ row }">
            <el-button
              type="primary"
              link
              size="small"
              @click="viewTaskDetail(row)"
            >
              {{ t('common.view') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>

      <el-empty v-if="!loading && displayTasks.length === 0" :description="t('common.noData')" />
    </el-card>

    <!-- 任务详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="t('crawler.taskDetail')"
      width="800px"
    >
      <div v-if="currentTask" class="task-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item :label="t('crawler.taskId')">
            {{ currentTask.task_id }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('crawler.taskName')">
            {{ currentTask.task_name }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('crawler.status')">
            <el-tag :type="getStatusType(currentTask.status)">
              {{ getStatusText(currentTask.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item :label="t('crawler.progress')">
            <el-progress
              :percentage="currentTask.progress"
              :status="getProgressStatus(currentTask.status)"
            />
          </el-descriptions-item>
          <el-descriptions-item :label="t('crawler.message')" :span="2">
            {{ currentTask.message || '-' }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('crawler.startTime')">
            {{ formatTime(currentTask.start_time) }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('crawler.endTime')">
            {{ currentTask.end_time ? formatTime(currentTask.end_time) : '-' }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('crawler.duration')">
            {{ formatDuration(currentTask.duration) }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('crawler.createdAt')">
            {{ formatTime(currentTask.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('crawler.updatedAt')" :span="2">
            {{ formatTime(currentTask.updated_at) }}
          </el-descriptions-item>
          <el-descriptions-item v-if="currentTask.metadata" :label="t('crawler.metadata')" :span="2">
            <pre class="metadata-display">{{ JSON.stringify(currentTask.metadata, null, 2) }}</pre>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Search, Refresh, List } from '@element-plus/icons-vue'
import PageHeader from '@/components/common/PageHeader.vue'
import { useWebSocketStore } from '@/store/websocket'
import crawlerApi from '@/api/crawler'
import dayjs from 'dayjs'

const { t } = useI18n()
const wsStore = useWebSocketStore()

// 查询表单
const queryForm = ref({
  status: '',
  task_id: ''
})

// 分页
const pagination = ref({
  page: 1,
  pageSize: 20,
  total: 0
})

// 数据
const loading = ref(false)
const tasks = ref([])
const currentTask = ref(null)
const detailDialogVisible = ref(false)

// 显示的任务列表（分页后）
const displayTasks = computed(() => {
  const start = (pagination.value.page - 1) * pagination.value.pageSize
  const end = start + pagination.value.pageSize
  return tasks.value.slice(start, end)
})

// 加载任务列表
const loadTasks = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.value.page,
      page_size: pagination.value.pageSize
    }
    
    if (queryForm.value.status) {
      params.status = queryForm.value.status
    }
    
    if (queryForm.value.task_id) {
      params.task_id = queryForm.value.task_id
    }

    const data = await crawlerApi.getTasks(params)
    tasks.value = data.items || []
    pagination.value.total = data.total || 0
  } catch (error) {
    console.error('加载任务列表失败:', error)
    ElMessage.error(t('crawler.loadTasksError'))
  } finally {
    loading.value = false
  }
}

// 重置查询
const resetQuery = () => {
  queryForm.value = {
    status: '',
    task_id: ''
  }
  pagination.value.page = 1
  loadTasks()
}

// 查看任务详情
const viewTaskDetail = async (task) => {
  try {
    const data = await crawlerApi.getTaskById(task.task_id)
    currentTask.value = data
    detailDialogVisible.value = true
  } catch (error) {
    console.error('加载任务详情失败:', error)
    ElMessage.error(t('crawler.loadTaskDetailError'))
  }
}

// 获取状态类型
const getStatusType = (status) => {
  const statusMap = {
    running: 'primary',
    completed: 'success',
    failed: 'danger'
  }
  return statusMap[status] || 'info'
}

// 获取状态文本
const getStatusText = (status) => {
  const statusMap = {
    running: t('crawler.statusRunning'),
    completed: t('crawler.statusCompleted'),
    failed: t('crawler.statusFailed')
  }
  return statusMap[status] || status
}

// 获取进度条状态
const getProgressStatus = (status) => {
  if (status === 'failed') return 'exception'
  if (status === 'completed') return 'success'
  return null
}

// 格式化时间
const formatTime = (time) => {
  if (!time) return '-'
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

// 格式化时长
const formatDuration = (duration) => {
  if (!duration) return '-'
  const hours = Math.floor(duration / 3600)
  const minutes = Math.floor((duration % 3600) / 60)
  const seconds = duration % 60
  
  if (hours > 0) {
    return `${hours}${t('crawler.hour')}${minutes}${t('crawler.minute')}${seconds}${t('crawler.second')}`
  } else if (minutes > 0) {
    return `${minutes}${t('crawler.minute')}${seconds}${t('crawler.second')}`
  } else {
    return `${seconds}${t('crawler.second')}`
  }
}

// 分页处理
const handlePageChange = (page) => {
  pagination.value.page = page
  loadTasks()
}

const handleSizeChange = (size) => {
  pagination.value.pageSize = size
  pagination.value.page = 1
  loadTasks()
}

// 监听WebSocket任务更新
watch(
  () => wsStore.tasks,
  (newTasks) => {
    // 合并WebSocket推送的任务更新到当前任务列表
    newTasks.forEach((wsTask) => {
      const index = tasks.value.findIndex(t => t.task_id === wsTask.task_id)
      if (index !== -1) {
        // 更新现有任务
        tasks.value[index] = { ...tasks.value[index], ...wsTask }
      } else {
        // 如果是新任务且符合筛选条件，添加到列表
        if (!queryForm.value.status || queryForm.value.status === wsTask.status) {
          tasks.value.unshift(wsTask)
        }
      }
    })
  },
  { deep: true }
)

// 组件挂载
onMounted(() => {
  // 连接WebSocket
  wsStore.connect()
  
  // 加载任务列表
  loadTasks()
  
  // 定期刷新任务列表（每30秒）
  const refreshInterval = setInterval(() => {
    if (!loading.value) {
      loadTasks()
    }
  }, 30000)
  
  // 保存定时器以便清理
  onUnmounted(() => {
    clearInterval(refreshInterval)
  })
})

// 组件卸载
onUnmounted(() => {
  // 断开WebSocket连接
  wsStore.disconnect()
})
</script>

<style scoped lang="less">
.crawler-monitor {
  .filter-card {
    margin-bottom: 20px;
  }

  .task-list-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .card-title {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 16px;
        font-weight: 500;
      }
    }

    .pagination-wrapper {
      margin-top: 20px;
      display: flex;
      justify-content: flex-end;
    }
  }

  .task-detail {
    .metadata-display {
      background: #f5f7fa;
      padding: 12px;
      border-radius: 4px;
      font-size: 12px;
      max-height: 300px;
      overflow-y: auto;
      margin: 0;
    }
  }
}
</style>

