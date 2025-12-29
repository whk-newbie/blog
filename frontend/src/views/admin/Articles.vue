<template>
  <div class="admin-articles-page">
    <page-header title="文章管理">
      <el-button type="primary" @click="goToCreate">
        <el-icon><Plus /></el-icon>
        新建文章
      </el-button>
    </page-header>

    <!-- 筛选栏 -->
    <div class="filter-bar">
      <el-input
        v-model="filters.keyword"
        placeholder="搜索文章标题..."
        clearable
        style="width: 300px"
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>

      <el-select
        v-model="filters.category_id"
        placeholder="选择分类"
        clearable
        @change="handleFilterChange"
      >
        <el-option
          v-for="category in categories"
          :key="category.id"
          :label="category.name"
          :value="category.id"
        />
      </el-select>

      <el-select
        v-model="filters.status"
        placeholder="选择状态"
        clearable
        @change="handleFilterChange"
      >
        <el-option label="草稿" value="draft" />
        <el-option label="已发布" value="published" />
      </el-select>

      <el-button type="primary" @click="handleSearch">搜索</el-button>
      <el-button @click="handleReset">重置</el-button>
    </div>

    <!-- 文章表格 -->
    <el-card class="table-card" shadow="never">
      <el-table
        v-loading="loading"
        :data="articles"
        style="width: 100%"
        :stripe="true"
        :header-cell-style="{ background: 'var(--bg-secondary)', color: 'var(--text-color)' }"
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

      <el-table-column label="标签" width="200">
        <template #default="{ row }">
          <el-tag
            v-for="tag in row.tags"
            :key="tag.id"
            size="small"
            style="margin-right: 4px"
          >
            {{ tag.name }}
          </el-tag>
          <span v-if="!row.tags || row.tags.length === 0">-</span>
        </template>
      </el-table-column>

      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">
            {{ row.status === 'published' ? '已发布' : '草稿' }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="标记" width="120">
        <template #default="{ row }">
          <el-tag v-if="row.is_top" type="danger" size="small">置顶</el-tag>
          <el-tag v-if="row.is_featured" type="warning" size="small">推荐</el-tag>
        </template>
      </el-table-column>

      <el-table-column prop="view_count" label="浏览量" width="100" />

      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>

      <el-table-column label="操作" width="250" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row.id)">编辑</el-button>
          
          <el-button
            v-if="row.status === 'draft'"
            size="small"
            type="success"
            @click="handlePublish(row.id)"
          >
            发布
          </el-button>

          <el-button
            v-else
            size="small"
            type="warning"
            @click="handleUnpublish(row.id)"
          >
            取消发布
          </el-button>

          <el-popconfirm
            title="确定要删除这篇文章吗？"
            @confirm="handleDelete(row.id)"
          >
            <template #reference>
              <el-button size="small" type="danger">删除</el-button>
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
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import api from '@/api'
import PageHeader from '@/components/common/PageHeader.vue'

const router = useRouter()

const loading = ref(false)
const articles = ref([])
const categories = ref([])

const filters = reactive({
  keyword: '',
  category_id: null,
  status: null
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

// 获取文章列表
const fetchArticles = async () => {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }

    if (filters.keyword) {
      params.keyword = filters.keyword
    }

    if (filters.category_id) {
      params.category_id = filters.category_id
    }

    if (filters.status) {
      params.status = filters.status
    }

    const response = await api.article.adminList(params)
    // 后端返回的是 items，不是 list
    articles.value = response.items || []
    pagination.total = response.total || 0
  } catch (error) {
    console.error('获取文章列表失败:', error)
    ElMessage.error('获取文章列表失败')
  } finally {
    loading.value = false
  }
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await api.category.list({ page: 1, page_size: 100 })
    // 后端返回的是 items，不是 list
    categories.value = response.items || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchArticles()
}

// 筛选改变
const handleFilterChange = () => {
  pagination.page = 1
  fetchArticles()
}

// 重置筛选
const handleReset = () => {
  filters.keyword = ''
  filters.category_id = null
  filters.status = null
  pagination.page = 1
  fetchArticles()
}

// 页码改变
const handlePageChange = () => {
  fetchArticles()
}

// 页面大小改变
const handleSizeChange = () => {
  pagination.page = 1
  fetchArticles()
}

// 跳转到创建页面
const goToCreate = () => {
  router.push('/admin/articles/create')
}

// 编辑
const handleEdit = (id) => {
  router.push(`/admin/articles/edit/${id}`)
}

// 发布
const handlePublish = async (id) => {
  try {
    await api.article.publish(id)
    ElMessage.success('发布成功')
    fetchArticles()
  } catch (error) {
    console.error('发布失败:', error)
    ElMessage.error('发布失败')
  }
}

// 取消发布
const handleUnpublish = async (id) => {
  try {
    await api.article.unpublish(id)
    ElMessage.success('取消发布成功')
    fetchArticles()
  } catch (error) {
    console.error('取消发布失败:', error)
    ElMessage.error('取消发布失败')
  }
}

// 删除
const handleDelete = async (id) => {
  try {
    await api.article.delete(id)
    ElMessage.success('删除成功')
    fetchArticles()
  } catch (error) {
    console.error('删除失败:', error)
    ElMessage.error('删除失败')
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

// 初始化
onMounted(() => {
  fetchArticles()
  fetchCategories()
})
</script>

<style scoped lang="less">
.admin-articles-page {
  padding: 0;
}

.filter-bar {
  display: flex;
  gap: 12px;
  margin-bottom: 24px;
  flex-wrap: wrap;
  padding: 20px 24px;
  background: var(--card-bg);
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
  transition: all 0.3s;

  &:hover {
    box-shadow: var(--shadow-md);
  }

  :deep(.el-input),
  :deep(.el-select) {
    flex: 0 0 auto;
    min-width: 200px;

    .el-input__wrapper {
      border-radius: 8px;
      box-shadow: 0 0 0 1px var(--border-color) inset;
      transition: all 0.3s;
      background: var(--bg-color);

      &:hover {
        box-shadow: 0 0 0 1px var(--primary-light) inset;
        border-color: var(--primary-light);
      }

      &.is-focus {
        box-shadow: 0 0 0 2px var(--primary-color) inset;
        border-color: var(--primary-color);
      }
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

.no-cover {
  color: var(--text-disabled);
  font-size: 0.875rem;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 60px;
  height: 60px;
  background: var(--bg-secondary);
  border-radius: 8px;
  border: 1px dashed var(--border-color);
}

:deep(.el-table) {
  border-radius: 0;

  .el-table__header-wrapper {
    .el-table__header {
      th {
        background: linear-gradient(180deg, var(--bg-secondary) 0%, var(--bg-tertiary) 100%);
        color: var(--text-color);
        font-weight: 600;
        font-size: 14px;
        padding: 16px 12px;
        border-bottom: 2px solid var(--border-color);
        text-transform: uppercase;
        letter-spacing: 0.5px;
        font-size: 12px;

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

    &--success {
      background: var(--success-color);
      border-color: var(--success-color);
      color: white;

      &:hover {
        background: #059669;
        border-color: #059669;
      }
    }

    &--warning {
      background: var(--warning-color);
      border-color: var(--warning-color);
      color: white;

      &:hover {
        background: #d97706;
        border-color: #d97706;
      }
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

:deep(.el-image) {
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid var(--border-light);
  transition: all 0.3s;

  &:hover {
    border-color: var(--primary-color);
    box-shadow: 0 2px 8px rgba(37, 99, 235, 0.2);
  }
}

:deep(.el-tag) {
  border-radius: 6px;
  font-weight: 500;
  padding: 4px 12px;
  border: none;
  margin-right: 6px;
  transition: all 0.2s;

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }

  &.el-tag--success {
    background: linear-gradient(135deg, rgba(16, 185, 129, 0.15) 0%, rgba(16, 185, 129, 0.1) 100%);
    color: var(--success-color);
  }

  &.el-tag--info {
    background: linear-gradient(135deg, rgba(37, 99, 235, 0.15) 0%, rgba(37, 99, 235, 0.1) 100%);
    color: var(--primary-color);
  }

  &.el-tag--danger {
    background: linear-gradient(135deg, rgba(239, 68, 68, 0.15) 0%, rgba(239, 68, 68, 0.1) 100%);
    color: var(--danger-color);
  }

  &.el-tag--warning {
    background: linear-gradient(135deg, rgba(245, 158, 11, 0.15) 0%, rgba(245, 158, 11, 0.1) 100%);
    color: var(--warning-color);
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

