<template>
  <footer class="site-footer">
    <div class="footer-container">
      <div class="footer-content">
        <p>&copy; {{ currentYear }} {{ blogTitle }}. All rights reserved.</p>
        <div class="footer-links">
          <a href="https://github.com" target="_blank" rel="noopener noreferrer">
            GitHub
          </a>
          <span class="divider">|</span>
          <router-link to="/privacy">隐私政策</router-link>
          <span class="divider">|</span>
          <router-link to="/terms">服务条款</router-link>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup>
import { ref, computed } from 'vue'

const blogTitle = ref('我的博客')
const currentYear = computed(() => new Date().getFullYear())
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

@media (max-width: 768px) {
  .footer-links {
    flex-wrap: wrap;
    justify-content: center;
  }
}
</style>

