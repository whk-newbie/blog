<template>
  <div class="admin-tags-page">
    <page-header title="标签管理">
      <el-button type="primary" @click="handleCreate">
        <el-icon><Plus /></el-icon>
        新建标签
      </el-button>
    </page-header>

    <!-- 标签表格 -->
    <el-table
      v-loading="loading"
      :data="tags"
      style="width: 100%"
      border
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" label="标签名称" width="200" />
      <el-table-column prop="slug" label="Slug" width="200" />
      <el-table-column prop="article_count" label="文章数" width="100" />
      
      <el-table-column label="创建时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>

      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button size="small" @click="handleEdit(row)">编辑</el-button>
          <el-popconfirm
            title="确定要删除这个标签吗？"
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

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑标签' : '新建标签'"
      width="500px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="标签名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入标签名称" />
        </el-form-item>

        <el-form-item label="URL Slug" prop="slug">
          <el-input v-model="form.slug" placeholder="留空自动生成" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import api from '@/api'
import PageHeader from '@/components/common/PageHeader.vue'

const loading = ref(false)
const tags = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)

const form = reactive({
  name: '',
  slug: ''
})

const rules = {
  name: [
    { required: true, message: '请输入标签名称', trigger: 'blur' }
  ]
}

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
    tags.value = response.data.items || []
    pagination.total = response.data.total || 0
  } catch (error) {
    console.error('获取标签列表失败:', error)
    ElMessage.error('获取标签列表失败')
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
      ElMessage.success('更新成功')
    } else {
      await api.tag.create(data)
      ElMessage.success('创建成功')
    }

    dialogVisible.value = false
    fetchTags()
  } catch (error) {
    if (error instanceof Error) {
      console.error('提交失败:', error)
      ElMessage.error('提交失败')
    }
  }
}

// 删除
const handleDelete = async (id) => {
  try {
    await api.tag.delete(id)
    ElMessage.success('删除成功')
    fetchTags()
  } catch (error) {
    console.error('删除失败:', error)
    ElMessage.error('删除失败')
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
  padding: 1.5rem;
}

.el-pagination {
  margin-top: 1.5rem;
  display: flex;
  justify-content: flex-end;
}
</style>

