<template>
  <el-card class="config-list-card" shadow="never">
    <el-table
      v-loading="loading"
      :data="configs"
      style="width: 100%"
      :stripe="true"
      :header-cell-style="{ background: 'var(--bg-secondary)', color: 'var(--text-color)' }"
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="config_key" label="配置键" min-width="200" show-overflow-tooltip />
      <el-table-column label="配置值" min-width="250">
        <template #default="{ row }">
          <span v-if="row.is_encrypted" class="masked-value">
            {{ maskValue(row.config_value) }}
          </span>
          <span v-else>{{ row.config_value || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
      <el-table-column label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-switch
            v-model="row.is_active"
            @change="() => $emit('toggle-active', row)"
          />
        </template>
      </el-table-column>
      <el-table-column label="加密" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.is_encrypted ? 'success' : 'info'" size="small">
            {{ row.is_encrypted ? '是' : '否' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="$emit('edit', row)">编辑</el-button>
          <el-popconfirm
            title="确定要删除这个配置吗？"
            @confirm="$emit('delete', row.id)"
          >
            <template #reference>
              <el-button size="small" type="danger">删除</el-button>
            </template>
          </el-popconfirm>
        </template>
      </el-table-column>
    </el-table>

    <div v-if="!loading && configs.length === 0" class="empty-state">
      <el-empty description="暂无配置" />
    </div>
  </el-card>
</template>

<script setup>
import { computed } from 'vue'

defineProps({
  configs: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

defineEmits(['edit', 'delete', 'toggle-active'])

// 脱敏显示
const maskValue = (value) => {
  if (!value) return '-'
  if (value.length <= 8) return '***'
  return value.substring(0, 4) + '***' + value.substring(value.length - 4)
}

// 格式化日期
const formatDate = (dateStr) => {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN')
}
</script>

<style scoped lang="less">
.config-list-card {
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);

  :deep(.el-card__body) {
    padding: 0;
  }
}

.masked-value {
  font-family: 'Courier New', monospace;
  color: var(--text-secondary);
}

.empty-state {
  padding: 40px;
  text-align: center;
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
</style>

