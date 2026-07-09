<template>
  <div class="login-page">
    <div class="login-card">
      <h1 class="login-title">{{ blogTitle }}</h1>
      <p class="login-subtitle">后台管理</p>

      <el-form ref="formRef" :model="form" :rules="rules" @submit.prevent="handleLogin">
        <el-form-item prop="username">
          <el-input
            v-model="form.username"
            placeholder="用户名"
            prefix-icon="User"
            size="large"
          />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="密码"
            prefix-icon="Lock"
            show-password
            size="large"
            @keydown.enter="handleLogin"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="login-btn"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="login-links">
        <router-link to="/">← 返回首页</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import configApi from '@/api/config'
import { initEncryption } from '@/api/http'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const { t } = useI18n()

const formRef = ref(null)
const loading = ref(false)
const blogTitle = ref('Blog')

const form = reactive({
  username: '',
  password: '',
})

const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
}

onMounted(async () => {
  try {
    await initEncryption()
  } catch (e) {
    // continue anyway
  }
  try {
    const config = await configApi.getSiteConfig()
    if (config?.blogTitle) blogTitle.value = config.blogTitle
  } catch (e) { /* ignore */ }
})

const handleLogin = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
  } catch {
    return
  }

  loading.value = true
  try {
    const response = await api.auth.login(form.username, form.password)
    userStore.login(response)
    ElMessage.success('登录成功')

    const redirect = route.query.redirect || '/admin'
    router.push(redirect)
  } catch (error) {
    console.error('登录失败:', error)
    // http.js interceptor already shows ElMessage.error
  } finally {
    loading.value = false
  }
}
</script>

<style scoped lang="less">
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--bg-color);
}

.login-card {
  width: 380px;
  padding: 40px;
  background: var(--card-bg);
  border-radius: var(--border-radius-md);
  box-shadow: var(--shadow-md);
}

.login-title {
  font-size: 24px;
  font-weight: 700;
  color: var(--text-heading);
  text-align: center;
  margin: 0 0 4px;
}

.login-subtitle {
  font-size: 14px;
  color: var(--text-secondary);
  text-align: center;
  margin: 0 0 32px;
}

.login-btn {
  width: 100%;
}

.login-links {
  text-align: center;
  margin-top: 16px;

  a {
    color: var(--text-secondary);
    text-decoration: none;
    font-size: 14px;

    &:hover {
      color: var(--primary-color);
    }
  }
}
</style>
