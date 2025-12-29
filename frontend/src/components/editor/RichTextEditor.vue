<template>
  <div class="rich-text-editor">
    <QuillEditor
      v-model:content="content"
      content-type="html"
      :options="editorOptions"
      :toolbar="toolbar"
      theme="snow"
      @ready="onEditorReady"
    />
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { QuillEditor } from '@vueup/vue-quill'
import '@vueup/vue-quill/dist/vue-quill.snow.css'
import { ElMessage } from 'element-plus'
import api from '@/api'

const { t } = useI18n()

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: '请输入文章内容...'
  },
  height: {
    type: String,
    default: '500px'
  }
})

const emit = defineEmits(['update:modelValue'])

const content = ref(props.modelValue)
const quillInstance = ref(null)

// 工具栏配置
const toolbar = [
  [{ 'header': [1, 2, 3, 4, 5, 6, false] }],
  [{ 'size': ['small', false, 'large', 'huge'] }],
  ['bold', 'italic', 'underline', 'strike'],
  [{ 'color': [] }, { 'background': [] }],
  [{ 'list': 'ordered' }, { 'list': 'bullet' }],
  [{ 'indent': '-1' }, { 'indent': '+1' }],
  [{ 'align': [] }],
  ['blockquote', 'code-block'],
  ['link', 'image'],
  ['clean']
]

// 编辑器选项
const editorOptions = {
  placeholder: props.placeholder,
  modules: {
    toolbar: {
      container: toolbar,
      handlers: {
        image: imageHandler
      }
    }
  }
}

// 自定义图片上传处理
function imageHandler() {
  const input = document.createElement('input')
  input.setAttribute('type', 'file')
  input.setAttribute('accept', 'image/*')
  input.click()

  input.onchange = async () => {
    const file = input.files[0]
    if (!file) return

    // 检查文件类型
    if (!file.type.startsWith('image/')) {
      ElMessage.error(t('common.uploadImageOnly'))
      return
    }

    // 检查文件大小（10MB）
    if (file.size > 10 * 1024 * 1024) {
      ElMessage.error(t('common.uploadImageSizeLimit'))
      return
    }

    try {
      ElMessage.info(t('common.uploading'))
      
      // 上传图片
      const response = await api.upload.uploadArticleImage(file)
      
      if (response.code === 0) {
        const imageUrl = response.data.url
        
        // 获取编辑器实例
        const quill = quillInstance.value.getQuill()
        
        // 获取当前光标位置
        const range = quill.getSelection(true)
        
        // 插入图片
        quill.insertEmbed(range.index, 'image', imageUrl)
        
        // 光标移动到图片后面
        quill.setSelection(range.index + 1)
        
        ElMessage.success(t('common.uploadSuccess'))
      } else {
        ElMessage.error(response.message || t('common.uploadError'))
      }
    } catch (error) {
      console.error('图片上传失败:', error)
      ElMessage.error(t('common.uploadError'))
    }
  }
}

// 编辑器准备就绪
const onEditorReady = (quill) => {
  quillInstance.value = { getQuill: () => quill }
  
  // 设置编辑器高度
  const editor = quill.root
  editor.style.minHeight = props.height
}

// 监听内容变化
watch(content, (newValue) => {
  emit('update:modelValue', newValue)
})

// 监听props变化
watch(() => props.modelValue, (newValue) => {
  if (newValue !== content.value) {
    content.value = newValue
  }
})
</script>

<style scoped lang="less">
.rich-text-editor {
  :deep(.ql-toolbar) {
    background-color: #fafafa;
    border-radius: 4px 4px 0 0;
    border-color: #dcdfe6;
  }

  :deep(.ql-container) {
    border-radius: 0 0 4px 4px;
    border-color: #dcdfe6;
    font-size: 14px;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  }

  :deep(.ql-editor) {
    min-height: 300px;
    max-height: 800px;
    overflow-y: auto;
    line-height: 1.8;

    &.ql-blank::before {
      font-style: normal;
      color: #999;
      left: 15px;
    }

    // 图片样式
    img {
      max-width: 100%;
      height: auto;
      display: block;
      margin: 10px 0;
    }

    // 代码块样式
    pre.ql-syntax {
      background-color: #f6f8fa;
      color: #333;
      padding: 15px;
      border-radius: 4px;
      overflow-x: auto;
      font-family: 'Monaco', 'Consolas', 'Courier New', monospace;
      font-size: 13px;
      line-height: 1.6;
    }

    // 引用样式
    blockquote {
      border-left: 4px solid #409eff;
      padding-left: 16px;
      margin: 10px 0;
      color: #666;
      font-style: italic;
    }

    // 链接样式
    a {
      color: #409eff;
      text-decoration: none;

      &:hover {
        text-decoration: underline;
      }
    }

    // 标题样式
    h1, h2, h3, h4, h5, h6 {
      margin: 15px 0 10px;
      font-weight: 600;
      line-height: 1.4;
    }

    h1 { font-size: 2em; }
    h2 { font-size: 1.75em; }
    h3 { font-size: 1.5em; }
    h4 { font-size: 1.25em; }
    h5 { font-size: 1.1em; }
    h6 { font-size: 1em; }

    // 列表样式
    ul, ol {
      padding-left: 1.5em;
      margin: 10px 0;
    }

    li {
      margin: 5px 0;
    }

    // 表格样式
    table {
      border-collapse: collapse;
      width: 100%;
      margin: 10px 0;

      td, th {
        border: 1px solid #ddd;
        padding: 8px;
      }

      th {
        background-color: #f6f8fa;
        font-weight: 600;
      }
    }
  }

  // 工具栏按钮样式
  :deep(.ql-toolbar button:hover),
  :deep(.ql-toolbar button:focus),
  :deep(.ql-toolbar .ql-picker-label:hover),
  :deep(.ql-toolbar .ql-picker-label.ql-active),
  :deep(.ql-toolbar button.ql-active) {
    color: #409eff;
  }

  :deep(.ql-toolbar button:hover .ql-stroke),
  :deep(.ql-toolbar button:focus .ql-stroke),
  :deep(.ql-toolbar button.ql-active .ql-stroke),
  :deep(.ql-toolbar .ql-picker-label:hover .ql-stroke),
  :deep(.ql-toolbar .ql-picker-label.ql-active .ql-stroke) {
    stroke: #409eff;
  }

  :deep(.ql-toolbar button:hover .ql-fill),
  :deep(.ql-toolbar button:focus .ql-fill),
  :deep(.ql-toolbar button.ql-active .ql-fill),
  :deep(.ql-toolbar .ql-picker-label:hover .ql-fill),
  :deep(.ql-toolbar .ql-picker-label.ql-active .ql-fill) {
    fill: #409eff;
  }
}

// 暗色主题样式
:root[data-theme='dark'] .rich-text-editor {
  :deep(.ql-toolbar) {
    background-color: #252525;
    border-color: #3a3a3a;
  }

  :deep(.ql-container) {
    border-color: #3a3a3a;
    background-color: #1a1a1a;
  }

  :deep(.ql-editor) {
    color: #d4d4d4;

    &.ql-blank::before {
      color: #6e6e6e;
    }

    // 代码块样式
    pre.ql-syntax {
      background-color: #2d2d2d;
      color: #d4d4d4;
    }

    // 引用样式
    blockquote {
      border-left-color: #6b9bd2;
      color: #a8a8a8;
    }

    // 链接样式
    a {
      color: #6b9bd2;
    }

    // 表格样式
    table {
      td, th {
        border-color: #3a3a3a;
      }

      th {
        background-color: #252525;
      }
    }
  }

  // 工具栏按钮样式
  :deep(.ql-toolbar button:hover),
  :deep(.ql-toolbar button:focus),
  :deep(.ql-toolbar .ql-picker-label:hover),
  :deep(.ql-toolbar .ql-picker-label.ql-active),
  :deep(.ql-toolbar button.ql-active) {
    color: #6b9bd2;
  }

  :deep(.ql-toolbar button:hover .ql-stroke),
  :deep(.ql-toolbar button:focus .ql-stroke),
  :deep(.ql-toolbar button.ql-active .ql-stroke),
  :deep(.ql-toolbar .ql-picker-label:hover .ql-stroke),
  :deep(.ql-toolbar .ql-picker-label.ql-active .ql-stroke) {
    stroke: #6b9bd2;
  }

  :deep(.ql-toolbar button:hover .ql-fill),
  :deep(.ql-toolbar button:focus .ql-fill),
  :deep(.ql-toolbar button.ql-active .ql-fill),
  :deep(.ql-toolbar .ql-picker-label:hover .ql-fill),
  :deep(.ql-toolbar .ql-picker-label.ql-active .ql-fill) {
    fill: #6b9bd2;
  }

  :deep(.ql-toolbar .ql-stroke) {
    stroke: #a8a8a8;
  }

  :deep(.ql-toolbar .ql-fill) {
    fill: #a8a8a8;
  }

  :deep(.ql-toolbar .ql-picker-label) {
    color: #d4d4d4;
  }
}
</style>

