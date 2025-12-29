<template>
  <div class="admin-tags-page">
    <page-header :title="t('nav.tags')">
      <template #extra>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          {{ t('tag.createTag') }}
        </el-button>
      </template>
    </page-header>

    <!-- 标签表格 -->
    <el-card class="table-card" shadow="never">
      <el-table
        v-loading="loading"
        :data="tags"
        style="width: 100%"
        :stripe="true"
        :header-cell-style="{ background: 'var(--bg-secondary)', color: 'var(--text-color)' }"
      >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" :label="t('tag.name')" width="200" />
      <el-table-column prop="slug" label="Slug" width="200" />
      <el-table-column prop="article_count" :label="t('tag.articleCount')" width="100" />
      
      <el-table-column :label="t('stats.createTime')" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>

      <el-table-column :label="t('common.operation')" width="180" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">{{ t('common.edit') }}</el-button>
          <el-popconfirm
            :title="t('tag.deleteConfirm')"
            @confirm="handleDelete(row.id)"
          >
            <template #reference>
              <el-button size="small" type="danger">{{ t('common.delete') }}</el-button>
            </template>
          </el-popconfirm>
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

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? t('tag.editTag') : t('tag.createTag')"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item :label="t('tag.name')" prop="name">
          <el-input v-model="form.name" :placeholder="t('tag.namePlaceholder')" />
        </el-form-item>

        <el-form-item label="URL Slug" prop="slug">
          <el-input v-model="form.slug" :placeholder="t('tag.slugPlaceholder')" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import api from '@/api'
import PageHeader from '@/components/common/PageHeader.vue'

const { t } = useI18n()
const loading = ref(false)
const tags = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)

const form = reactive({
  name: '',
  slug: ''
})

const rules = computed(() => ({
  name: [
    { required: true, message: t('tag.namePlaceholder'), trigger: 'blur' }
  ]
}))

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 获取标签列表
const fetchTags = async () => {
  try {
    loading.value = true
    const response = await api.tag.list({
      page: pagination.page,
      page_size: pagination.pageSize
    })
    // 后端返回的是 items，不是 list
    tags.value = response.items || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取标签列表失败:', error)
    ElMessage.error(t('tag.loadError'))
  } finally {
    loading.value = false
  }
}

// 页码改变
const handlePageChange = () => {
  fetchTags()
}

// 页面大小改变
const handleSizeChange = () => {
  pagination.page = 1
  fetchTags()
}

// 新建
const handleCreate = () => {
  isEdit.value = false
  form.id = undefined
  form.name = ''
  form.slug = ''
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row) => {
  isEdit.value = true
  form.id = row.id
  form.name = row.name
  form.slug = row.slug
  dialogVisible.value = true
}

// 提交
const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    
    const data = {
      name: form.name,
      slug: form.slug
    }

    if (isEdit.value) {
      await api.tag.update(form.id, data)
      ElMessage.success(t('tag.updateSuccess'))
    } else {
      await api.tag.create(data)
      ElMessage.success(t('tag.createSuccess'))
    }

    dialogVisible.value = false
    fetchTags()
  } catch (error) {
    if (error instanceof Error) {
      console.error('提交失败:', error)
      ElMessage.error(isEdit.value ? t('tag.updateError') : t('tag.createError'))
    }
  }
}

// 删除
const handleDelete = async (id) => {
  try {
    await api.tag.delete(id)
    ElMessage.success(t('tag.deleteSuccess'))
    fetchTags()
  } catch (error) {
    console.error('删除失败:', error)
    ElMessage.error(t('tag.deleteError'))
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
  fetchTags()
})
</script>

<style scoped lang="less">
.admin-tags-page {
  padding: 0;
}

.table-card {
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
  overflow: hidden;
  transition: all 0.3s;

  &:hover {
    box-shadow: var(--shadow-md);
  }

  :deep(.el-card__body) {
    padding: 0;
  }
}

:deep(.el-button) {
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s;
  padding: 10px 20px;
  height: auto;

  &.el-button--primary {
    background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-light) 100%);
    border: none;
    box-shadow: 0 2px 8px rgba(37, 99, 235, 0.2);

    &:hover {
      background: linear-gradient(135deg, var(--primary-light) 0%, var(--primary-color) 100%);
      box-shadow: 0 4px 12px rgba(37, 99, 235, 0.3);
      transform: translateY(-2px);
    }

    &:active {
      transform: translateY(0);
    }
  }
}

:deep(.el-table) {
  border-radius: 0;

  .el-table__header-wrapper {
    .el-table__header {
      th {
        background: linear-gradient(180deg, var(--bg-secondary) 0%, var(--bg-tertiary) 100%);
        color: var(--text-color);
        font-weight: 600;
        font-size: 12px;
        padding: 16px 12px;
        border-bottom: 2px solid var(--border-color);
        text-transform: uppercase;
        letter-spacing: 0.5px;

        &:first-child {
          padding-left: 24px;
        }

        &:last-child {
          padding-right: 24px;
        }
      }
    }
  }

  .el-table__body-wrapper {
    .el-table__body {
      tr {
        transition: all 0.2s;

        &:hover {
          background: var(--bg-blue-light) !important;
          transform: scale(1.001);
          box-shadow: 0 2px 8px rgba(37, 99, 235, 0.05);
        }

        td {
          padding: 16px 12px;
          border-bottom: 1px solid var(--border-light);
          vertical-align: middle;

          &:first-child {
            padding-left: 24px;
          }

          &:last-child {
            padding-right: 24px;
          }
        }
      }

      .el-table__row {
        &:last-child td {
          border-bottom: none;
        }
      }
    }
  }

  .el-button {
    border-radius: 6px;
    font-weight: 500;
    padding: 6px 14px;
    margin-right: 8px;
    transition: all 0.2s;

    &:last-child {
      margin-right: 0;
    }

    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
    }

    &--small {
      padding: 5px 12px;
      font-size: 12px;
    }

    &--danger {
      background: var(--danger-color);
      border-color: var(--danger-color);
      color: white;

      &:hover {
        background: #dc2626;
        border-color: #dc2626;
      }
    }
  }
}

:deep(.el-dialog) {
  border-radius: 12px;
  box-shadow: var(--shadow-lg);

  .el-dialog__header {
    padding: 20px 24px;
    border-bottom: 1px solid var(--border-light);
    background: linear-gradient(135deg, var(--bg-secondary) 0%, var(--bg-color) 100%);
  }

  .el-dialog__body {
    padding: 24px;
  }

  .el-form-item__label {
    font-weight: 500;
    color: var(--text-color);
  }

  .el-input__wrapper {
    border-radius: 8px;
    transition: all 0.3s;

    &:hover {
      box-shadow: 0 0 0 1px var(--primary-light) inset;
    }

    &.is-focus {
      box-shadow: 0 0 0 2px var(--primary-color) inset;
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

  :deep(.el-pager) {
    li {
      border-radius: 6px;
      font-weight: 500;
      margin: 0 4px;
      transition: all 0.2s;

      &:hover {
        background: var(--primary-light);
        color: white;
      }

      &.is-active {
        background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-light) 100%);
        color: white;
        box-shadow: 0 2px 6px rgba(37, 99, 235, 0.3);
      }
    }
  }

  :deep(.btn-prev),
  :deep(.btn-next) {
    border-radius: 6px;
    margin: 0 4px;
    transition: all 0.2s;

    &:hover {
      background: var(--primary-light);
      color: white;
    }
  }

  :deep(.el-pagination__total),
  :deep(.el-pagination__jump) {
    color: var(--text-secondary);
    font-weight: 500;
  }
}
</style>

