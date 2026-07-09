<template>
  <Teleport to="body">
    <transition name="overlay">
      <div v-if="visible" class="search-overlay" @click.self="close">
        <div class="search-dialog">
          <input
            ref="inputRef"
            v-model="query"
            class="search-input"
            :placeholder="t('common.searchPlaceholder')"
            @input="onSearch"
            @keydown.escape="close"
          />
          <div class="search-results" v-if="query.length > 0">
            <div v-if="results.length === 0" class="no-results">{{ t('common.noResults') }}</div>
            <div
              v-for="article in results"
              :key="article.id"
              class="result-item"
              @click="goToArticle(article.slug)"
            >
              <span class="result-title">{{ displayTitle(article) }}</span>
              <span class="result-date">{{ formatDate(article.publish_at || article.created_at) }}</span>
            </div>
          </div>
        </div>
      </div>
    </transition>
  </Teleport>
</template>

<script setup>
import { ref, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useLanguageStore } from '@/store/language'
import api from '@/api'

const { t } = useI18n()
const router = useRouter()
const languageStore = useLanguageStore()

const visible = ref(false)
const query = ref('')
const results = ref([])
const inputRef = ref(null)
let allArticles = []

const displayTitle = (article) => {
  return languageStore.language === 'en-US' && article.title_en
    ? article.title_en
    : article.title
}

const formatDate = (dateStr) => {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

const open = async () => {
  visible.value = true
  await nextTick()
  inputRef.value?.focus()
  if (allArticles.length === 0) {
    try {
      const res = await api.article.list({ page: 1, page_size: 100 })
      allArticles = res?.items || []
    } catch (e) {
      // ignore
    }
  }
}

const close = () => {
  visible.value = false
  query.value = ''
  results.value = []
}

const onSearch = () => {
  const q = query.value.toLowerCase().trim()
  if (!q) {
    results.value = []
    return
  }
  results.value = allArticles.filter(a => {
    const title = (a.title || '').toLowerCase()
    const titleEn = (a.title_en || '').toLowerCase()
    const summary = (a.summary || '').toLowerCase()
    return title.includes(q) || titleEn.includes(q) || summary.includes(q)
  }).slice(0, 10)
}

const goToArticle = (slug) => {
  close()
  router.push(`/article/${slug}`)
}

defineExpose({ open, close })
</script>

<style scoped lang="less">
.search-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  z-index: 200;
  display: flex;
  justify-content: center;
  padding-top: 15vh;
}

.search-dialog {
  width: 100%;
  max-width: 560px;
  max-height: 60vh;
  display: flex;
  flex-direction: column;
}

.search-input {
  width: 100%;
  padding: 16px 20px;
  font-size: 18px;
  border: none;
  border-radius: var(--border-radius-md);
  background: var(--card-bg);
  color: var(--text-heading);
  outline: none;
  box-shadow: var(--shadow-lg);

  &::placeholder {
    color: var(--text-secondary);
  }
}

.search-results {
  margin-top: 8px;
  background: var(--card-bg);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-lg);
  max-height: 400px;
  overflow-y: auto;
}

.no-results {
  padding: 24px;
  text-align: center;
  color: var(--text-secondary);
  font-size: var(--font-size-sm);
}

.result-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 20px;
  cursor: pointer;
  border-bottom: 1px solid var(--border-light);
  transition: background var(--transition-fast);

  &:last-child {
    border-bottom: none;
  }

  &:hover {
    background: var(--hover-bg);
  }

  .result-title {
    font-size: var(--font-size-sm);
    color: var(--text-heading);
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    margin-right: var(--spacing-md);
  }

  .result-date {
    font-size: var(--font-size-xs);
    color: var(--text-secondary);
    flex-shrink: 0;
  }
}

.overlay-enter-active,
.overlay-leave-active {
  transition: opacity 0.2s ease;
}
.overlay-enter-from,
.overlay-leave-to {
  opacity: 0;
}
</style>
