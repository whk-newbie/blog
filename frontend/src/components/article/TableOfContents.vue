<template>
  <div v-if="headings.length > 0" class="table-of-contents">
    <div class="toc-header">
      <el-icon><List /></el-icon>
      <span>{{ t('toc.title') }}</span>
    </div>
    <ul class="toc-list">
      <li
        v-for="heading in headings"
        :key="heading.id"
        :class="['toc-item', `toc-level-${heading.level}`, { active: activeId === heading.id }]"
      >
        <a
          :href="`#${heading.id}`"
          @click.prevent="scrollToHeading(heading.id)"
        >
          {{ heading.text }}
        </a>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { List } from '@element-plus/icons-vue'

const { t } = useI18n()

const props = defineProps({
  container: {
    type: String,
    default: '.article-body'
  }
})

const headings = ref([])
const activeId = ref('')

// 提取文章标题
const extractHeadings = () => {
  const container = document.querySelector(props.container)
  if (!container) return

  const headingElements = container.querySelectorAll('h1, h2, h3, h4, h5, h6')
  const result = []

  headingElements.forEach((heading, index) => {
    const level = parseInt(heading.tagName.substring(1))
    const text = heading.textContent.trim()
    const id = heading.id || `heading-${index}`

    // 如果没有ID，添加一个
    if (!heading.id) {
      heading.id = id
    }

    result.push({
      id,
      text,
      level,
      element: heading
    })
  })

  headings.value = result
}

// 滚动到指定标题
const scrollToHeading = (id) => {
  const element = document.getElementById(id)
  if (element) {
    const offset = 80 // 顶部导航栏高度
    const elementPosition = element.getBoundingClientRect().top
    const offsetPosition = elementPosition + window.pageYOffset - offset

    window.scrollTo({
      top: offsetPosition,
      behavior: 'smooth'
    })

    // 更新活动状态
    activeId.value = id
  }
}

// 处理滚动，高亮当前标题
const handleScroll = () => {
  const scrollPosition = window.scrollY + 100 // 偏移量

  // 找到当前可见的标题
  for (let i = headings.value.length - 1; i >= 0; i--) {
    const heading = headings.value[i]
    const element = document.getElementById(heading.id)
    
    if (element) {
      const rect = element.getBoundingClientRect()
      const elementTop = element.offsetTop

      if (scrollPosition >= elementTop) {
        activeId.value = heading.id
        break
      }
    }
  }

  // 如果滚动到顶部，清除高亮
  if (window.scrollY < 100 && headings.value.length > 0) {
    activeId.value = ''
  }
}

// 节流函数
const throttle = (func, delay) => {
  let timeoutId
  return (...args) => {
    if (!timeoutId) {
      timeoutId = setTimeout(() => {
        func(...args)
        timeoutId = null
      }, delay)
    }
  }
}

const throttledHandleScroll = throttle(handleScroll, 100)

onMounted(() => {
  // 延迟提取标题，确保内容已渲染
  setTimeout(() => {
    extractHeadings()
  }, 300)

  // 监听滚动事件
  window.addEventListener('scroll', throttledHandleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', throttledHandleScroll)
})

defineExpose({
  extractHeadings
})
</script>

<style scoped lang="less">
.table-of-contents {
  position: sticky;
  top: 100px;
  background: var(--bg-color);
  border-radius: 8px;
  padding: 1.5rem;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  max-height: calc(100vh - 120px);
  overflow-y: auto;
}

.toc-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-color);
  margin-bottom: 1rem;
  padding-bottom: 0.75rem;
  border-bottom: 2px solid var(--border-color);

  .el-icon {
    font-size: 1.25rem;
  }
}

.toc-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.toc-item {
  margin: 0;
  padding: 0;
  transition: all 0.2s;

  a {
    display: block;
    padding: 0.5rem 0.75rem;
    color: var(--text-color-secondary);
    text-decoration: none;
    border-left: 2px solid transparent;
    transition: all 0.2s;
    font-size: 0.875rem;
    line-height: 1.6;

    &:hover {
      color: var(--primary-color);
      background: var(--bg-color-secondary);
      border-left-color: var(--primary-color);
    }
  }

  &.active {
    a {
      color: var(--primary-color);
      font-weight: 600;
      border-left-color: var(--primary-color);
      background: var(--bg-color-secondary);
    }
  }

  // 不同层级的缩进
  &.toc-level-1 {
    a {
      padding-left: 0.75rem;
      font-size: 0.9375rem;
      font-weight: 600;
    }
  }

  &.toc-level-2 {
    a {
      padding-left: 1.25rem;
    }
  }

  &.toc-level-3 {
    a {
      padding-left: 1.75rem;
      font-size: 0.8125rem;
    }
  }

  &.toc-level-4 {
    a {
      padding-left: 2.25rem;
      font-size: 0.8125rem;
    }
  }

  &.toc-level-5 {
    a {
      padding-left: 2.75rem;
      font-size: 0.75rem;
    }
  }

  &.toc-level-6 {
    a {
      padding-left: 3.25rem;
      font-size: 0.75rem;
    }
  }
}

// 滚动条样式
.table-of-contents::-webkit-scrollbar {
  width: 4px;
}

.table-of-contents::-webkit-scrollbar-track {
  background: transparent;
}

.table-of-contents::-webkit-scrollbar-thumb {
  background: #ddd;
  border-radius: 4px;

  &:hover {
    background: #bbb;
  }
}

// 响应式设计
@media (max-width: 1200px) {
  .table-of-contents {
    display: none;
  }
}
</style>

