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
    <el-table
      v-loading="loading"
      :data="articles"
      style="width: 100%"
      border
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
    articles.value = response.data.items || []
    pagination.total = response.data.total || 0
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
    categories.value = response.data.items || []
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
  padding: 1.5rem;
}

.filter-bar {
  display: flex;
  gap: 1rem;
  margin-bottom: 1.5rem;
  flex-wrap: wrap;
}

.no-cover {
  color: var(--text-color-secondary);
  font-size: 0.875rem;
}

.el-pagination {
  margin-top: 1.5rem;
  display: flex;
  justify-content: flex-end;
}
</style>

