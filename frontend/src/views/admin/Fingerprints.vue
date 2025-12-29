<template>
  <div class="fingerprints-page">
    <page-header :title="t('nav.fingerprints')" />

    <!-- 指纹列表表格 -->
    <el-card class="table-card" shadow="never">
      <el-table
        v-loading="loading"
        :data="fingerprints"
        style="width: 100%"
        :stripe="true"
        :header-cell-style="{ background: 'var(--bg-secondary)', color: 'var(--text-color)' }"
      >
        <el-table-column prop="id" label="ID" width="80" />
        
        <el-table-column :label="t('stats.fingerprintHash')" min-width="200">
          <template #default="{ row }">
            <el-tooltip :content="row.fingerprint_hash" placement="top">
              <span class="hash-text">{{ truncateHash(row.fingerprint_hash) }}</span>
            </el-tooltip>
          </template>
        </el-table-column>

        <el-table-column :label="t('stats.userAgent')" min-width="300" show-overflow-tooltip>
          <template #default="{ row }">
            <span>{{ row.user_agent || '-' }}</span>
          </template>
        </el-table-column>

        <el-table-column :label="t('stats.visitCount')" width="120" align="right">
          <template #default="{ row }">
            <el-tag type="info">{{ row.visit_count || 0 }}</el-tag>
          </template>
        </el-table-column>

        <el-table-column :label="t('stats.firstSeen')" width="180">
          <template #default="{ row }">
            {{ formatDate(row.first_seen_at) }}
          </template>
        </el-table-column>

        <el-table-column :label="t('stats.lastSeen')" width="180">
          <template #default="{ row }">
            {{ formatDate(row.last_seen_at) }}
          </template>
        </el-table-column>

        <el-table-column :label="t('common.operation')" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleViewDetail(row)">
              {{ t('common.view') }}
            </el-button>
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

    <!-- 详情对话框 -->
    <el-dialog
      v-model="detailDialogVisible"
      :title="t('stats.fingerprintInfo')"
      width="800px"
    >
      <div v-if="selectedFingerprint" class="fingerprint-detail">
        <el-descriptions :column="1" border>
          <el-descriptions-item :label="t('stats.fingerprintHash')">
            <code>{{ selectedFingerprint.fingerprint_hash }}</code>
          </el-descriptions-item>
          <el-descriptions-item :label="t('stats.userAgent')">
            {{ selectedFingerprint.user_agent || '-' }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('stats.visitCount')">
            {{ selectedFingerprint.visit_count || 0 }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('stats.firstSeen')">
            {{ formatDate(selectedFingerprint.first_seen_at) }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('stats.lastSeen')">
            {{ formatDate(selectedFingerprint.last_seen_at) }}
          </el-descriptions-item>
          <el-descriptions-item :label="t('stats.fingerprintInfo')">
            <pre class="fingerprint-data">{{ formatFingerprintData(selectedFingerprint.fingerprint_data) }}</pre>
          </el-descriptions-item>
        </el-descriptions>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import PageHeader from '@/components/common/PageHeader.vue'
import fingerprintApi from '@/api/fingerprint'

const { t } = useI18n()

const loading = ref(false)
const fingerprints = ref([])
const detailDialogVisible = ref(false)
const selectedFingerprint = ref(null)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 获取指纹列表
const fetchFingerprints = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }

    const response = await fingerprintApi.list(params)
    fingerprints.value = response.items || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取指纹列表失败:', error)
    ElMessage.error(t('common.error'))
  } finally {
    loading.value = false
  }
}

// 页码改变
const handlePageChange = () => {
  fetchFingerprints()
}

// 页面大小改变
const handleSizeChange = () => {
  pagination.page = 1
  fetchFingerprints()
}

// 查看详情
const handleViewDetail = (row) => {
  selectedFingerprint.value = row
  detailDialogVisible.value = true
}

// 格式化日期
const formatDate = (date) => {
  if (!date) return '-'
  return dayjs(date).format('YYYY-MM-DD HH:mm:ss')
}

// 截断哈希值
const truncateHash = (hash) => {
  if (!hash) return '-'
  if (hash.length <= 20) return hash
  return hash.substring(0, 10) + '...' + hash.substring(hash.length - 10)
}

// 格式化指纹数据
const formatFingerprintData = (data) => {
  if (!data) return '-'
  try {
    if (typeof data === 'string') {
      return JSON.stringify(JSON.parse(data), null, 2)
    }
    return JSON.stringify(data, null, 2)
  } catch (e) {
    return String(data)
  }
}

onMounted(() => {
  fetchFingerprints()
})
</script>

<style lang="less" scoped>
.fingerprints-page {
  padding: 20px;

  .table-card {
    margin-bottom: 20px;

    .hash-text {
      font-family: 'Courier New', monospace;
      font-size: 12px;
      color: #409eff;
      cursor: pointer;
    }
  }

  .fingerprint-detail {
    .fingerprint-data {
      background: #f5f7fa;
      padding: 12px;
      border-radius: 4px;
      font-family: 'Courier New', monospace;
      font-size: 12px;
      line-height: 1.6;
      max-height: 400px;
      overflow-y: auto;
      white-space: pre-wrap;
      word-break: break-all;
    }

    code {
      background: #f5f7fa;
      padding: 2px 6px;
      border-radius: 3px;
      font-family: 'Courier New', monospace;
      font-size: 12px;
    }
  }
}
</style>

