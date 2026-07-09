<template>
  <aside class="sidebar">
    <div class="sidebar-inner">
      <div class="profile">
        <el-avatar :size="80" :src="avatarUrl || undefined" class="avatar-el">{{ avatarUrl ? '' : initials }}</el-avatar>
        <h2 class="name">{{ blogTitle }}</h2>
        <p class="bio">{{ blogDescription }}</p>
      </div>

      <div class="stats">
        <div class="stat-item"><span class="stat-number">{{ articleCount }}</span><span class="stat-label">篇文章</span></div>
        <div class="stat-item"><span class="stat-number">{{ runningDays }}</span><span class="stat-label">天</span></div>
      </div>

      <div class="tag-cloud" v-if="tags.length > 0">
        <el-tag v-for="(tag,i) in tags" :key="tag.id" size="small" class="tag-item" :type="tagTypes[i % tagTypes.length]" :style="{ fontSize: getTagSize(tag.article_count) + 'px' }" @click="$router.push(`/articles?tag_id=${tag.id}`)">
          {{ tag.name }}
        </el-tag>
      </div>
    </div>
  </aside>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import api from '@/api'
import configApi from '@/api/config'

const blogTitle = ref('My Blog')
const blogDescription = ref('')
const avatarUrl = ref('')
const articleCount = ref(0)
const tags = ref([])

const tagTypes = ['', 'success', 'warning', 'danger', 'info']
const initials = computed(() => blogTitle.value.charAt(0).toUpperCase())
const runningDays = computed(() => Math.floor((new Date() - new Date('2024-01-01')) / 86400000))

const getTagSize = (count) => {
  if (!count || count <= 0) return 12
  if (count >= 20) return 18
  return 12 + (count / 20) * 6
}

onMounted(async () => {
  try {
    const [res, cfg] = await Promise.all([
      Promise.all([api.article.list({ page: 1, page_size: 1 }), api.tag.list({ page: 1, page_size: 30 })]),
      configApi.getSiteConfig().catch(() => null),
    ])
    articleCount.value = res[0]?.total || 0
    tags.value = res[1]?.items || []
    if (cfg) {
      if (cfg.blogTitle) blogTitle.value = cfg.blogTitle
      if (cfg.blogDescription) blogDescription.value = cfg.blogDescription
      if (cfg.avatarUrl) avatarUrl.value = cfg.avatarUrl
    }
  } catch (e) { /* defaults */ }
})
</script>

<style scoped lang="less">
.sidebar { width: var(--sidebar-width); flex-shrink: 0; padding-top: var(--header-height); }
.sidebar-inner { position: sticky; top: calc(var(--header-height) + var(--spacing-lg)); }
.profile { text-align: center; margin-bottom: var(--spacing-lg); }
.avatar-el { margin: 0 auto var(--spacing-md); display: block; border: 3px solid var(--border-color); }
.name { font-size: var(--font-size-xl); font-weight: 600; color: var(--text-heading); margin: 0 0 var(--spacing-xs); }
.bio { font-size: var(--font-size-sm); color: var(--text-secondary); margin: 0; }
.stats { display: flex; justify-content: center; gap: var(--spacing-lg); padding: var(--spacing-md) 0; border-top: 1px solid var(--border-color); border-bottom: 1px solid var(--border-color); margin-bottom: var(--spacing-lg); }
.stat-item { text-align: center; .stat-number { display: block; font-size: var(--font-size-xl); font-weight: 700; color: var(--text-heading); } .stat-label { font-size: var(--font-size-xs); color: var(--text-secondary); } }
.tag-cloud { margin-bottom: var(--spacing-lg); }
.tag-item { cursor: pointer; margin: 2px 4px 2px 0; }
@media (max-width: 768px) { .sidebar { display: none; } }
</style>
