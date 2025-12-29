<template>
  <div class="article-editor-page">
    <page-header :title="isEdit ? t('article.editArticle') : t('nav.newArticle')">
      <el-button @click="goBack">{{ t('common.back') }}</el-button>
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
      <el-form-item :label="t('article.title')" prop="title">
        <el-input
          v-model="form.title"
          :placeholder="t('article.titlePlaceholder')"
          maxlength="255"
          show-word-limit
        />
      </el-form-item>

      <!-- Slug -->
      <el-form-item label="URL Slug" prop="slug">
        <el-input
          v-model="form.slug"
          :placeholder="t('article.slugPlaceholder')"
        >
          <template #append>
            <el-button @click="generateSlug">{{ t('article.autoGenerateSlug') }}</el-button>
          </template>
        </el-input>
      </el-form-item>

      <!-- 封面图 -->
      <el-form-item :label="t('article.coverImage')">
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
              {{ form.cover_image ? t('article.changeImage') : t('article.uploadImage') }}
            </el-button>
          </el-upload>
          <el-button v-if="form.cover_image" @click="form.cover_image = ''">
            {{ t('article.deleteImage') }}
          </el-button>
        </div>
      </el-form-item>

      <!-- 摘要 -->
      <el-form-item :label="t('article.summary')">
        <el-input
          v-model="form.summary"
          type="textarea"
          :rows="3"
          :placeholder="t('article.summaryPlaceholder')"
          maxlength="500"
          show-word-limit
        />
      </el-form-item>

      <!-- 分类 -->
      <el-form-item :label="t('article.category')" prop="category_id">
        <el-select
          v-model="form.category_id"
          :placeholder="t('article.selectCategoryPlaceholder')"
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
      <el-form-item :label="t('article.tags')">
        <el-select
          v-model="form.tag_ids"
          multiple
          :placeholder="t('article.selectTagsPlaceholder')"
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
      <el-form-item :label="t('article.content')" prop="content">
        <rich-text-editor
          v-model="form.content"
          :placeholder="t('article.contentPlaceholder')"
          height="500px"
        />
      </el-form-item>

      <!-- 状态设置 -->
      <el-form-item :label="t('article.publishStatus')">
        <el-radio-group v-model="form.status">
          <el-radio label="draft">{{ t('article.draft') }}</el-radio>
          <el-radio label="published">{{ t('article.published') }}</el-radio>
        </el-radio-group>
      </el-form-item>

      <!-- 发布时间 -->
      <el-form-item :label="t('article.publishTime')">
        <el-date-picker
          v-model="form.publish_at"
          type="datetime"
          :placeholder="t('article.publishTimePlaceholder')"
          format="YYYY-MM-DD HH:mm:ss"
        />
        <div class="form-tip">{{ t('article.publishTimeTip') }}</div>
      </el-form-item>

      <!-- 其他选项 -->
      <el-form-item :label="t('article.otherOptions')">
        <el-checkbox v-model="form.is_top">{{ t('article.isTop') }}</el-checkbox>
        <el-checkbox v-model="form.is_featured">{{ t('article.isFeatured') }}</el-checkbox>
      </el-form-item>

      <!-- 操作按钮 -->
      <el-form-item>
        <el-button type="primary" @click="handleSubmit('published')">
          {{ form.status === 'published' ? t('article.publishArticle') : t('article.saveAndPublish') }}
        </el-button>
        <el-button @click="handleSubmit('draft')">{{ t('article.saveAsDraft') }}</el-button>
        <el-button @click="goBack">{{ t('common.cancel') }}</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Upload } from '@element-plus/icons-vue'
import api from '@/api'
import PageHeader from '@/components/common/PageHeader.vue'
import RichTextEditor from '@/components/editor/RichTextEditor.vue'
import { useUserStore } from '@/store/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const { t } = useI18n()

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

const rules = computed(() => ({
  title: [
    { required: true, message: t('article.titlePlaceholder'), trigger: 'blur' }
  ],
  content: [
    { required: true, message: t('article.contentPlaceholder'), trigger: 'blur' }
  ]
}))

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
    ElMessage.error(t('article.loadError'))
    router.push('/admin/articles')
  } finally {
    loading.value = false
  }
}

// 获取分类列表
const fetchCategories = async () => {
  try {
    const response = await api.category.list({ page: 1, page_size: 100 })
    categories.value = response.items || []
  } catch (error) {
    console.error('获取分类列表失败:', error)
    ElMessage.error(t('category.loadError'))
  }
}

// 获取标签列表
const fetchTags = async () => {
  try {
    const response = await api.tag.list({ page: 1, page_size: 100 })
    tags.value = response.items || []
  } catch (error) {
    console.error('获取标签列表失败:', error)
    ElMessage.error(t('tag.loadError'))
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
    ElMessage.error(t('common.uploadImageOnly'))
    return false
  }
  if (!isLt10M) {
    ElMessage.error(t('common.uploadImageSizeLimit'))
    return false
  }
  return true
}

// 上传成功
const handleUploadSuccess = (response) => {
  if (response.code === 0) {
    form.cover_image = response.data.url
    ElMessage.success(t('common.uploadSuccess'))
  } else {
    ElMessage.error(response.message || t('common.uploadError'))
  }
}

// 上传失败
const handleUploadError = () => {
  ElMessage.error(t('common.uploadError'))
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
      ElMessage.success(t('article.updateSuccess'))
    } else {
      await api.article.create(data)
      ElMessage.success(t('article.createSuccess'))
    }

    router.push('/admin/articles')
  } catch (error) {
    if (error instanceof Error) {
      console.error('提交失败:', error)
      ElMessage.error(isEdit.value ? t('article.updateError') : t('article.createError'))
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

