<template>
  <footer class="site-footer">
    <div class="footer-container">
      <div class="footer-content">
        <p>&copy; {{ currentYear }} {{ blogTitle }}. All rights reserved.</p>
        <div class="footer-links">
          <a href="https://github.com/whk-newbie/blog" target="_blank" rel="noopener noreferrer">
            GitHub
          </a>
        </div>
        <!-- 备案信息 -->
        <div v-if="icpInfo" class="icp-info">
          <a v-if="icpInfo.url" :href="icpInfo.url" target="_blank" rel="noopener noreferrer">
            {{ icpInfo.text }}
          </a>
          <span v-else>{{ icpInfo.text }}</span>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup>
import { ref, computed } from 'vue'

const blogTitle = ref('我的博客')
const currentYear = computed(() => new Date().getFullYear())

// 备案信息配置
// 如果需要显示备案信息，设置 icpInfo 对象
// 示例: { text: '京ICP备12345678号', url: 'https://beian.miit.gov.cn/' }
// 如果不需要显示，设置为 null 或 undefined
const icpInfo = ref(null)
// const icpInfo = ref({
//   text: '京ICP备12345678号',
//   url: 'https://beian.miit.gov.cn/'
// })
</script>

<style scoped lang="less">
.site-footer {
  background: var(--footer-bg);
  border-top: 1px solid var(--border-color);
  padding: 30px 0;
  margin-top: auto;
}

.footer-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
}

.footer-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 15px;
  color: var(--text-secondary);
  font-size: 14px;

  p {
    margin: 0;
  }
}

.footer-links {
  display: flex;
  align-items: center;
  gap: 10px;

  a {
    color: var(--text-secondary);
    text-decoration: none;
    transition: color 0.3s;

    &:hover {
      color: var(--primary-color);
    }
  }

  .divider {
    color: var(--border-color);
  }
}

.icp-info {
  margin-top: 10px;
  font-size: 12px;
  color: var(--text-secondary);

  a {
    color: var(--text-secondary);
    text-decoration: none;
    transition: color 0.3s;

    &:hover {
      color: var(--primary-color);
    }
  }
}

@media (max-width: 768px) {
  .footer-links {
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>

