<template>
  <div class="article-editor-page">
    <page-header :title="isEdit ? '编辑文章' : '新建文章'">
      <el-button @click="goBack">返回</el-button>
    </page-header>

    <el-form
      ref="formRef"
      v-loading="loading"
      :model="form"
      :rules="rules"
      label-width="100px"
      class="article-form"
    >
      <!-- 标题 -->
      <el-form-item label="文章标题" prop="title">
        <el-input
          v-model="form.title"
          placeholder="请输入文章标题"
          maxlength="255"
          show-word-limit
        />
      </el-form-item>

      <!-- Slug -->
      <el-form-item label="URL Slug" prop="slug">
        <el-input
          v-model="form.slug"
          placeholder="留空自动生成"
        >
          <template #append>
            <el-button @click="generateSlug">自动生成</el-button>
          </template>
        </el-input>
      </el-form-item>

      <!-- 封面图 -->
      <el-form-item label="封面图片">
        <div class="cover-upload">
          <el-image
            v-if="form.cover_image"
            :src="form.cover_image"
            class="cover-preview"
            fit="cover"
          />
          <el-upload
            :action="uploadAction"
            :headers="uploadHeaders"
            :show-file-list="false"
            :before-upload="beforeUpload"
            :on-success="handleUploadSuccess"
            :on-error="handleUploadError"
          >
            <el-button type="primary">
              <el-icon><Upload /></el-icon>
              {{ form.cover_image ? '更换图片' : '上传图片' }}
            </el-button>
          </el-upload>
          <el-button v-if="form.cover_image" @click="form.cover_image = ''">
            删除图片
          </el-button>
        </div>
      </el-form-item>

      <!-- 摘要 -->
      <el-form-item label="文章摘要">
        <el-input
          v-model="form.summary"
          type="textarea"
          :rows="3"
          placeholder="请输入文章摘要"
          maxlength="500"
          show-word-limit
        />
      </el-form-item>

      <!-- 分类 -->
      <el-form-item label="文章分类" prop="category_id">
        <el-select
          v-model="form.category_id"
          placeholder="请选择分类"
          clearable
        >
          <el-option
            v-for="category in categories"
            :key="category.id"
            :label="category.name"
            :value="category.id"
          />
        </el-select>
      </el-form-item>

      <!-- 标签 -->
      <el-form-item label="文章标签">
        <el-select
          v-model="form.tag_ids"
          multiple
          placeholder="请选择标签"
          style="width: 100%"
        >
          <el-option
            v-for="tag in tags"
            :key="tag.id"
            :label="tag.name"
            :value="tag.id"
          />
        </el-select>
      </el-form-item>

      <!-- 正文 -->
      <el-form-item label="文章正文" prop="content">
        <rich-text-editor
          v-model="form.content"
          placeholder="请输入文章内容..."
          height="500px"
        />
      </el-form-item>

      <!-- 状态设置 -->
      <el-form-item label="发布状态">
        <el-radio-group v-model="form.status">
          <el-radio label="draft">草稿</el-radio>
          <el-radio label="published">已发布</el-radio>
        </el-radio-group>
      </el-form-item>

      <!-- 发布时间 -->
      <el-form-item label="发布时间">
        <el-date-picker
          v-model="form.publish_at"
          type="datetime"
          placeholder="选择发布时间"
          format="YYYY-MM-DD HH:mm:ss"
        />
        <div class="form-tip">留空则使用当前时间</div>
      </el-form-item>

      <!-- 其他选项 -->
      <el-form-item label="其他选项">
        <el-checkbox v-model="form.is_top">置顶</el-checkbox>
        <el-checkbox v-model="form.is_featured">推荐</el-checkbox>
      </el-form-item>

      <!-- 操作按钮 -->
      <el-form-item>
        <el-button type="primary" @click="handleSubmit('published')">
          {{ form.status === 'published' ? '发布文章' : '保存并发布' }}
        </el-button>
        <el-button @click="handleSubmit('draft')">保存为草稿</el-button>
        <el-button @click="goBack">取消</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import api from '@/api'
import PageHeader from '@/components/common/PageHeader.vue'
import RichTextEditor from '@/components/editor/RichTextEditor.vue'
import { useUserStore } from '@/store/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const formRef = ref(null)
const loading = ref(false)
const categories = ref([])
const tags = ref([])

const isEdit = computed(() => !!route.params.id)

const form = reactive({
  title: '',
  slug: '',
  summary: '',
  content: '',
  cover_image: '',
  category_id: null,
  tag_ids: [],
  status: 'draft',
  publish_at: null,
  is_top: false,
  is_featured: false
})

const rules = {
  title: [
    { required: true, message: '请输入文章标题', trigger: 'blur' }
  ],
  content: [
    { required: true, message: '请输入文章内容', trigger: 'blur' }
  ]
}

const uploadAction = computed(() => {
  const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'
  return `${baseURL}/admin/upload/article-image`
})

const uploadHeaders = computed(() => {
  return {
    Authorization: `Bearer ${userStore.token}`
  }
})

// 获取文章详情
const fetchArticle = async () => {
  try {
    loading.value = true
    const response = await api.article.getById(route.params.id)
    const article = response.data
    
    // 填充表单
    form.title = article.title || ''
    form.slug = article.slug || ''
    form.summary = article.summary || ''
    form.content = article.content || ''
    form.cover_image = article.cover_image || ''
    form.category_id = article.category_id || null
    form.tag_ids = article.tags ? article.tags.map(t => t.id) : []
    form.status = article.status || 'draft'
    form.publish_at = article.publish_at ? new Date(article.publish_at) : null
    form.is_top = article.is_top || false
    form.is_featured = article.is_featured || false
  } catch (error) {
    console.error('获取文章详情失败:', error)
    ElMessage.error('获取文章详情失败')
    router.push('/admin/articles')
  } finally {
    loading.value = false
  }
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await api.category.list({ page: 1, page_size: 100 })
    categories.value = response.list || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
  }
}

// 获取标签列表
const fetchTags = async () => {
  try {
    const response = await api.tag.list({ page: 1, page_size: 100 })
    tags.value = response.list || []
  } catch (error) {
    console.error('获取标签列表失败:', error)
  }
}

// 生成Slug
const generateSlug = () => {
  if (form.title) {
    // 简单的拼音转换（实际项目中应该使用专门的库）
    form.slug = form.title
      .toLowerCase()
      .replace(/\s+/g, '-')
      .replace(/[^\w\-\u4e00-\u9fa5]+/g, '')
  }
}

// 上传前检查
const beforeUpload = (file) => {
  const isImage = file.type.startsWith('image/')
  const isLt10M = file.size / 1024 / 1024 < 10

  if (!isImage) {
    ElMessage.error('只能上传图片文件！')
    return false
  }
  if (!isLt10M) {
    ElMessage.error('图片大小不能超过 10MB！')
    return false
  }
  return true
}

// 上传成功
const handleUploadSuccess = (response) => {
  if (response.code === 0) {
    form.cover_image = response.data.url
    ElMessage.success('上传成功')
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

// 上传失败
const handleUploadError = () => {
  ElMessage.error('上传失败')
}

// 提交表单
const handleSubmit = async (status) => {
  try {
    await formRef.value.validate()
    
    loading.value = true

    const data = {
      ...form,
      status,
      publish_at: form.publish_at ? form.publish_at.toISOString() : null
    }

    if (isEdit.value) {
      await api.article.update(route.params.id, data)
      ElMessage.success('更新成功')
    } else {
      await api.article.create(data)
      ElMessage.success('创建成功')
    }

    router.push('/admin/articles')
  } catch (error) {
    if (error instanceof Error) {
      console.error('提交失败:', error)
      ElMessage.error('提交失败')
    }
  } finally {
    loading.value = false
  }
}

// 返回
const goBack = () => {
  router.push('/admin/articles')
}

// 初始化
onMounted(() => {
  fetchCategories()
  fetchTags()
  
  if (isEdit.value) {
    fetchArticle()
  }
})
</script>

<style scoped lang="less">
.article-editor-page {
  padding: 0;
}

.article-form {
  max-width: 1000px;
  margin: 0 auto;
  background: var(--card-bg);
  padding: 32px;
  border-radius: 16px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);
}

.cover-upload {
  display: flex;
  align-items: center;
  gap: 16px;

  .cover-preview {
    width: 240px;
    height: 160px;
    border-radius: 10px;
    overflow: hidden;
    border: 2px solid var(--border-light);
    box-shadow: var(--shadow-sm);
  }
}

.form-tip {
  margin-top: 6px;
  font-size: 13px;
  color: var(--text-secondary);
  font-weight: 500;
}

:deep(.el-form) {
  .el-form-item__label {
    font-weight: 600;
    color: var(--text-color);
    font-size: 15px;
  }

  .el-input__wrapper,
  .el-textarea__inner {
    border-radius: 10px;
    box-shadow: 0 0 0 1px var(--border-color) inset;
    transition: all 0.3s;

    &:hover {
      box-shadow: 0 0 0 1px var(--primary-light) inset;
    }

    &.is-focus {
      box-shadow: 0 0 0 2px var(--primary-color) inset;
    }
  }

  .el-select {
    .el-select__wrapper {
      border-radius: 10px;
    }
  }

  .el-checkbox {
    font-weight: 500;

    .el-checkbox__label {
      color: var(--text-color);
    }
  }

  .el-switch {
    .el-switch__core {
      border-radius: 12px;
    }
  }

  .el-button {
    border-radius: 10px;
    font-weight: 600;
    padding: 12px 28px;
    transition: all 0.3s;

    &.el-button--primary {
      background: var(--primary-color);
      border-color: var(--primary-color);

      &:hover {
        background: var(--primary-light);
        box-shadow: 0 4px 14px rgba(37, 99, 235, 0.3);
        transform: translateY(-2px);
      }
    }

    &.el-button--success {
      background: var(--success-color);
      border-color: var(--success-color);

      &:hover {
        box-shadow: 0 4px 14px rgba(16, 185, 129, 0.3);
      }
    }
  }
}

:deep(.el-upload) {
  .el-button {
    border-radius: 10px;
    border: 2px dashed var(--border-color);
    transition: all 0.3s;

    &:hover {
      border-color: var(--primary-color);
      background: var(--bg-blue-light);
    }
  }
}

:deep(.el-card) {
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);

  .el-card__header {
    background: var(--bg-blue-light);
    border-bottom: 2px solid var(--border-blue);
    font-weight: 600;
    color: var(--text-color);
  }

  .el-card__body {
    padding: 24px;
  }
}
</style>

