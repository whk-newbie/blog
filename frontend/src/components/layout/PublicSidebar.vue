<template>
  <aside class="sidebar">
    <div class="sidebar-inner">
      <!-- Profile -->
      <div class="profile">
        <div class="avatar">{{ initials }}</div>
        <h2 class="name">{{ blogTitle }}</h2>
        <p class="bio">{{ blogDescription }}</p>
      </div>

      <!-- Stats -->
      <div class="stats">
        <div class="stat-item">
          <span class="stat-number">{{ articleCount }}</span>
          <span class="stat-label">篇文章</span>
        </div>
        <div class="stat-item">
          <span class="stat-number">{{ runningDays }}</span>
          <span class="stat-label">天</span>
        </div>
      </div>

      <!-- Tag Cloud -->
      <div class="tag-cloud" v-if="tags.length > 0">
        <h3 class="section-title">标签</h3>
        <div class="tag-items">
          <router-link
            v-for="tag in tags"
            :key="tag.id"
            :to="`/articles?tag_id=${tag.id}`"
            class="tag-item"
            :style="{ fontSize: getTagSize(tag.article_count) + 'px' }"
          >
            #{{ tag.name }}
          </router-link>
        </div>
      </div>

      <!-- Social Links -->
      <div class="social-links">
        <a v-if="githubUrl" :href="githubUrl" target="_blank" rel="noopener" title="GitHub">
          <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/></svg>
        </a>
        <a v-if="email" :href="`mailto:${email}`" title="Email">
          <svg viewBox="0 0 24 24" width="20" height="20" fill="currentColor"><path d="M20 4H4c-1.1 0-2 .9-2 2v12c0 1.1.9 2 2 2h16c1.1 0 2-.9 2-2V6c0-1.1-.9-2-2-2zm0 4l-8 5-8-5V6l8 5 8-5v2z"/></svg>
        </a>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import configApi from '@/api/config'

const { t } = useI18n()

const blogTitle = ref('My Blog')
const blogDescription = ref('')
const articleCount = ref(0)
const tags = ref([])
const githubUrl = ref('')
const email = ref('')

const initials = computed(() => {
  return blogTitle.value.charAt(0).toUpperCase()
})

const runningDays = computed(() => {
  const start = new Date('2024-01-01')
  const now = new Date()
  return Math.floor((now - start) / (1000 * 60 * 60 * 24))
})

const getTagSize = (count) => {
  const min = 13
  const max = 22
  if (!count || count <= 0) return min
  if (count >= 20) return max
  return min + (count / 20) * (max - min)
}

const fetchSiteConfig = async () => {
  try {
    const response = await configApi.getSiteConfig()
    if (response) {
      if (response.blogTitle) blogTitle.value = response.blogTitle
      if (response.blogDescription) blogDescription.value = response.blogDescription
    }
  } catch (e) {
    // use defaults
  }
}

const fetchData = async () => {
  try {
    const [articleRes, tagRes] = await Promise.all([
      api.article.list({ page: 1, page_size: 1 }),
      api.tag.list({ page: 1, page_size: 30 }),
    ])
    articleCount.value = articleRes?.total || 0
    tags.value = tagRes?.items || []
  } catch (e) {
    // use defaults
  }
}

onMounted(() => {
  fetchSiteConfig()
  fetchData()
})
</script>

<style scoped lang="less">
.sidebar {
  width: var(--sidebar-width);
  flex-shrink: 0;
  padding-top: var(--header-height);
}

.sidebar-inner {
  position: sticky;
  top: calc(var(--header-height) + var(--spacing-lg));
}

.profile {
  text-align: center;
  margin-bottom: var(--spacing-lg);

  .avatar {
    width: 80px;
    height: 80px;
    border-radius: 50%;
    background: var(--primary-color);
    color: white;
    font-size: 32px;
    font-weight: 600;
    display: flex;
    align-items: center;
    justify-content: center;
    margin: 0 auto var(--spacing-md);
  }

  .name {
    font-size: var(--font-size-xl);
    font-weight: 600;
    color: var(--text-heading);
    margin: 0 0 var(--spacing-xs);
  }

  .bio {
    font-size: var(--font-size-sm);
    color: var(--text-secondary);
    margin: 0;
    line-height: 1.5;
  }
}

.stats {
  display: flex;
  justify-content: center;
  gap: var(--spacing-lg);
  padding: var(--spacing-md) 0;
  border-top: 1px solid var(--border-color);
  border-bottom: 1px solid var(--border-color);
  margin-bottom: var(--spacing-lg);

  .stat-item {
    text-align: center;

    .stat-number {
      display: block;
      font-size: var(--font-size-xl);
      font-weight: 700;
      color: var(--text-heading);
    }

    .stat-label {
      font-size: var(--font-size-xs);
      color: var(--text-secondary);
    }
  }
}

.section-title {
  font-size: var(--font-size-sm);
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
  margin: 0 0 var(--spacing-sm);
}

.tag-cloud {
  margin-bottom: var(--spacing-lg);
}

.tag-items {
  display: flex;
  flex-wrap: wrap;
  gap: var(--spacing-sm);

  .tag-item {
    text-decoration: none;
    color: var(--text-secondary);
    padding: 2px 6px;
    border-radius: var(--border-radius-sm);
    transition: color var(--transition-fast);

    &:hover {
      color: var(--primary-color);
    }
  }
}

.social-links {
  display: flex;
  justify-content: center;
  gap: var(--spacing-md);

  a {
    color: var(--text-secondary);
    transition: color var(--transition-fast);

    &:hover {
      color: var(--primary-color);
    }
  }
}

@media (max-width: 768px) {
  .sidebar {
    display: none;
  }
}
</style>
