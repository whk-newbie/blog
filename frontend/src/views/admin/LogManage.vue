<template>
  <div class="log-manage-page">
    <page-header :title="t('nav.logs')">
      <template #extra>
        <el-button type="warning" @click="handleCleanup">
          <el-icon><Delete /></el-icon>
          {{ t('log.cleanupLogs') }}
        </el-button>
      </template>
    </page-header>

    <!-- 查询条件 -->
    <el-card class="filter-card" shadow="hover">
      <el-form :inline="true" :model="queryForm">
        <el-form-item :label="t('log.level')">
          <el-select v-model="queryForm.level" style="width: 150px" clearable>
            <el-option :label="t('log.levelDebug')" value="DEBUG" />
            <el-option :label="t('log.levelInfo')" value="INFO" />
            <el-option :label="t('log.levelWarn')" value="WARN" />
            <el-option :label="t('log.levelError')" value="ERROR" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('log.source')">
          <el-input
            v-model="queryForm.source"
            :placeholder="t('log.sourcePlaceholder')"
            style="width: 200px"
            clearable
          />
        </el-form-item>
        <el-form-item :label="t('log.dateRange')">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            :start-placeholder="t('log.startDate')"
            :end-placeholder="t('log.endDate')"
            format="YYYY-MM-DD"
            value-format="YYYY-MM-DD"
            @change="handleDateRangeChange"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadLogs">
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

    <!-- 日志列表 -->
    <el-card class="table-card" shadow="never">
      <el-table
        v-loading="loading"
        :data="logs"
        style="width: 100%"
        :stripe="true"
        :header-cell-style="{ background: 'var(--bg-secondary)', color: 'var(--text-color)' }"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column :label="t('log.level')" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getLevelType(row.level)" size="small">
              {{ row.level }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="message" :label="t('log.message')" min-width="300" show-overflow-tooltip />
        <el-table-column prop="source" :label="t('log.source')" width="150" show-overflow-tooltip />
        <el-table-column prop="ip_address" :label="t('log.ipAddress')" width="150" show-overflow-tooltip />
        <el-table-column :label="t('log.createTime')" width="180">
          <template #default="{ row }">
            {{ formatDate(row.created_at) }}
          </template>
        </el-table-column>
        <el-table-column :label="t('common.operation')" width="120" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleViewDetail(row)">{{ t('common.view') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 分页 -->
    <el-pagination
      v-model:current-page="pagination.page"
      v-model:page-size="pagination.pageSize"
      :total="pagination.total"
      :page-sizes="[10, 20, 50, 100]"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handlePageChange"
    />

    <!-- 日志详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="t('log.logDetail')"
      width="800px"
    >
      <el-descriptions :column="2" border>
        <el-descriptions-item :label="t('log.id')">{{ currentLog.id }}</el-descriptions-item>
        <el-descriptions-item :label="t('log.level')">
          <el-tag :type="getLevelType(currentLog.level)" size="small">
            {{ currentLog.level }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="t('log.source')" :span="2">
          {{ currentLog.source || '-' }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('log.ipAddress')">
          {{ currentLog.ip_address || '-' }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('log.createTime')">
          {{ formatDate(currentLog.created_at) }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('log.message')" :span="2">
          {{ currentLog.message }}
        </el-descriptions-item>
        <el-descriptions-item :label="t('log.context')" :span="2">
          <pre class="context-content">{{ formatContext(currentLog.context) }}</pre>
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button type="primary" @click="detailDialogVisible = false">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 清理日志对话框 -->
    <el-dialog
      v-model="cleanupDialogVisible"
      :title="t('log.cleanupLogs')"
      width="500px"
    >
      <el-form
        ref="cleanupFormRef"
        :model="cleanupForm"
        :rules="cleanupRules"
        label-width="120px"
      >
        <el-form-item :label="t('log.retentionDays')" prop="retention_days">
          <el-input-number
            v-model="cleanupForm.retention_days"
            :min="1"
            :max="365"
            style="width: 100%"
          />
          <div class="form-tip">{{ t('log.retentionDaysTip') }}</div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="cleanupDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="warning" @click="handleCleanupSubmit">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Delete } from '@element-plus/icons-vue'
import api from '@/api'
import PageHeader from '@/components/common/PageHeader.vue'

const { t } = useI18n()
const loading = ref(false)
const logs = ref([])
const dateRange = ref([])
const detailDialogVisible = ref(false)
const cleanupDialogVisible = ref(false)
const currentLog = ref({})
const cleanupFormRef = ref(null)

const queryForm = reactive({
  level: '',
  source: '',
  start_date: '',
  end_date: ''
})

const cleanupForm = reactive({
  retention_days: 90
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const cleanupRules = computed(() => ({
  retention_days: [
    { required: true, message: t('log.retentionDaysPlaceholder'), trigger: 'blur' }
  ]
}))

// 获取日志级别类型
const getLevelType = (level) => {
  const typeMap = {
    DEBUG: 'info',
    INFO: '',
    WARN: 'warning',
    ERROR: 'danger'
  }
  return typeMap[level] || ''
}

// 获取日志列表
const loadLogs = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize,
      level: queryForm.level || undefined,
      source: queryForm.source || undefined,
      start_date: queryForm.start_date || undefined,
      end_date: queryForm.end_date || undefined
    }
    const response = await api.log.getLogs(params)
    logs.value = response.list || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取日志列表失败:', error)
    ElMessage.error(t('log.loadError'))
  } finally {
    loading.value = false
  }
}

// 日期范围改变
const handleDateRangeChange = (dates) => {
  if (dates && dates.length === 2) {
    queryForm.start_date = dates[0]
    queryForm.end_date = dates[1]
  } else {
    queryForm.start_date = ''
    queryForm.end_date = ''
  }
}

// 重置查询
const resetQuery = () => {
  queryForm.level = ''
  queryForm.source = ''
  queryForm.start_date = ''
  queryForm.end_date = ''
  dateRange.value = []
  pagination.page = 1
  loadLogs()
}

// 页码改变
const handlePageChange = () => {
  loadLogs()
}

// 页面大小改变
const handleSizeChange = () => {
  pagination.page = 1
  loadLogs()
}

// 查看详情
const handleViewDetail = async (log) => {
  try {
    const response = await api.log.getLogById(log.id)
    currentLog.value = response || log
    detailDialogVisible.value = true
  } catch (error) {
    console.error('获取日志详情失败:', error)
    ElMessage.error(t('log.loadDetailError'))
    // 如果获取详情失败,使用列表中的数据
    currentLog.value = log
    detailDialogVisible.value = true
  }
}

// 格式化上下文
const formatContext = (context) => {
  if (!context) return '-'
  if (typeof context === 'string') {
    try {
      return JSON.stringify(JSON.parse(context), null, 2)
    } catch {
      return context
    }
  }
  return JSON.stringify(context, null, 2)
}

// 清理日志
const handleCleanup = () => {
  cleanupForm.retention_days = 90
  cleanupDialogVisible.value = true
}

// 提交清理
const handleCleanupSubmit = async () => {
  try {
    await cleanupFormRef.value.validate()
    // 先关闭对话框，避免确认弹窗位置错误
    cleanupDialogVisible.value = false
    // 等待对话框关闭动画完成后再显示确认弹窗
    await new Promise(resolve => setTimeout(resolve, 300))
    await ElMessageBox.confirm(
      t('log.cleanupConfirm', { days: cleanupForm.retention_days }),
      t('log.cleanupLogs'),
      {
        confirmButtonText: t('common.confirm'),
        cancelButtonText: t('common.cancel'),
        type: 'warning'
      }
    )
    const response = await api.log.cleanupLogs({
      retention_days: cleanupForm.retention_days
    })
    ElMessage.success(t('log.cleanupSuccess', { count: response.deleted_count || 0 }))
    loadLogs()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('清理日志失败:', error)
      ElMessage.error(t('log.cleanupError'))
      // 如果取消或出错，重新打开对话框
      cleanupDialogVisible.value = true
    }
  }
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}

// 初始化
onMounted(() => {
  loadLogs()
})
</script>

<style scoped lang="less">
.log-manage-page {
  padding: 0;
}

.filter-card {
  margin-bottom: 20px;
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
}

.table-card {
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
  overflow: hidden;

  :deep(.el-card__body) {
    padding: 0;
  }
}

.context-content {
  background: var(--bg-secondary);
  padding: 12px;
  border-radius: 6px;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  line-height: 1.6;
  max-height: 400px;
  overflow: auto;
  white-space: pre-wrap;
  word-break: break-all;
}

.form-tip {
  margin-top: 4px;
  font-size: 12px;
  color: var(--text-secondary);
}

:deep(.el-table) {
  .el-table__header-wrapper {
    .el-table__header {
      th {
        background: linear-gradient(180deg, var(--bg-secondary) 0%, var(--bg-tertiary) 100%);
        color: var(--text-color);
        font-weight: 600;
        font-size: 12px;
        padding: 16px 12px;
        border-bottom: 2px solid var(--border-color);
      }
    }
  }

  .el-table__body-wrapper {
    .el-table__body {
      tr {
        transition: all 0.2s;

        &:hover {
          background: var(--bg-blue-light) !important;
        }

        td {
          padding: 16px 12px;
          border-bottom: 1px solid var(--border-light);
        }
      }
    }
  }
}

.el-pagination {
  margin-top: 24px;
  display: flex;
  justify-content: flex-end;
  padding: 16px 24px;
  background: var(--bg-secondary);
  border-radius: 8px;
}
</style>

