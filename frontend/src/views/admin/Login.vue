<template>
  <div class="login-container">
    <div class="login-box">
      <div class="login-header">
        <h1>博客管理系统</h1>
        <p>Blog Management System</p>
      </div>

      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="username">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入用户名"
            size="large"
            prefix-icon="User"
            clearable
          />
        </el-form-item>

        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            prefix-icon="Lock"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item>
          <el-button
            type="primary"
            size="large"
            :loading="loading"
            class="login-button"
            @click="handleLogin"
          >
            {{ loading ? '登录中...' : '登录' }}
          </el-button>
        </el-form-item>
      </el-form>

      <div class="login-footer">
        <p>默认用户名: admin</p>
        <p>默认密码: admin@123</p>
        <p class="tip">首次登录请及时修改密码</p>
      </div>
    </div>

    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="showChangePasswordDialog"
      title="修改密码"
      width="400px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="false"
    >
      <el-alert
        title="检测到您正在使用默认密码"
        type="warning"
        description="为了账户安全，请立即修改密码"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      />

      <el-form
        ref="changePasswordFormRef"
        :model="changePasswordForm"
        :rules="changePasswordRules"
        label-width="80px"
      >
        <el-form-item label="旧密码" prop="oldPassword">
          <el-input
            v-model="changePasswordForm.oldPassword"
            type="password"
            placeholder="请输入旧密码"
            show-password
          />
        </el-form-item>

        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="changePasswordForm.newPassword"
            type="password"
            placeholder="请输入新密码（至少6位）"
            show-password
          />
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="changePasswordForm.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            show-password
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="handleSkipChangePassword">稍后修改</el-button>
        <el-button
          type="primary"
          :loading="changingPassword"
          @click="handleChangePassword"
        >
          确认修改
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import api from '@/api'

const router = useRouter()
const userStore = useUserStore()

const loginFormRef = ref(null)
const changePasswordFormRef = ref(null)
const loading = ref(false)
const changingPassword = ref(false)
const showChangePasswordDialog = ref(false)

// 登录表单
const loginForm = reactive({
  username: '',
  password: ''
})

// 登录表单验证规则
const loginRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6位', trigger: 'blur' }
  ]
}

// 修改密码表单
const changePasswordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 修改密码表单验证规则
const changePasswordRules = {
  oldPassword: [
    { required: true, message: '请输入旧密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== changePasswordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 处理登录
const handleLogin = async () => {
  if (!loginFormRef.value) return

  try {
    await loginFormRef.value.validate()
    loading.value = true

    const response = await api.auth.login(loginForm.username, loginForm.password)
    
    // 保存登录信息
    userStore.login(response)

    ElMessage.success('登录成功')

    // 如果使用默认密码，提示修改
    if (response.is_default_password) {
      changePasswordForm.oldPassword = loginForm.password
      showChangePasswordDialog.value = true
    } else {
      // 跳转到管理后台
      router.push('/admin')
    }
  } catch (error) {
    console.error('登录失败:', error)
    ElMessage.error(error.message || '登录失败')
  } finally {
    loading.value = false
  }
}

// 处理修改密码
const handleChangePassword = async () => {
  if (!changePasswordFormRef.value) return

  try {
    await changePasswordFormRef.value.validate()
    changingPassword.value = true

    await api.auth.changePassword(
      changePasswordForm.oldPassword,
      changePasswordForm.newPassword
    )

    ElMessage.success('密码修改成功')
    showChangePasswordDialog.value = false

    // 更新状态
    userStore.isDefaultPassword = false
    localStorage.setItem('isDefaultPassword', 'false')

    // 跳转到管理后台
    router.push('/admin')
  } catch (error) {
    console.error('修改密码失败:', error)
    ElMessage.error(error.message || '修改密码失败')
  } finally {
    changingPassword.value = false
  }
}

// 跳过修改密码
const handleSkipChangePassword = () => {
  showChangePasswordDialog.value = false
  router.push('/admin')
}
</script>

<style scoped lang="less">
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, #2563eb 0%, #1e40af 50%, #1e3a8a 100%);
  padding: 20px;
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    width: 200%;
    height: 200%;
    background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 1px, transparent 1px);
    background-size: 50px 50px;
    animation: bgMove 60s linear infinite;
  }
}

@keyframes bgMove {
  0% {
    transform: translate(0, 0);
  }
  100% {
    transform: translate(50px, 50px);
  }
}

.login-box {
  width: 100%;
  max-width: 420px;
  padding: 48px 40px;
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3), 0 0 0 1px rgba(255, 255, 255, 0.2) inset;
  position: relative;
  z-index: 1;
}

.login-header {
  text-align: center;
  margin-bottom: 40px;

  h1 {
    margin: 0 0 12px 0;
    font-size: 32px;
    font-weight: 700;
    background: linear-gradient(135deg, var(--primary-color) 0%, var(--secondary-color) 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    letter-spacing: -0.5px;
  }

  p {
    margin: 0;
    font-size: 14px;
    color: var(--text-secondary);
    letter-spacing: 1px;
    font-weight: 500;
  }
}

.login-form {
  :deep(.el-form-item) {
    margin-bottom: 24px;
  }

  :deep(.el-input__wrapper) {
    padding: 14px 16px;
    border-radius: 10px;
    box-shadow: 0 0 0 1px var(--border-color) inset;
    transition: all 0.3s;

    &:hover {
      box-shadow: 0 0 0 1px var(--primary-light) inset;
    }

    &.is-focus {
      box-shadow: 0 0 0 2px var(--primary-color) inset, var(--shadow-blue);
    }
  }

  :deep(.el-input__inner) {
    font-size: 15px;
  }

  :deep(.el-input__prefix) {
    font-size: 18px;
    color: var(--text-secondary);
  }
}

.login-button {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  border-radius: 10px;
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--primary-dark) 100%);
  border: none;
  box-shadow: 0 4px 14px rgba(37, 99, 235, 0.4);
  transition: all 0.3s;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(37, 99, 235, 0.5);
  }

  &:active {
    transform: translateY(0);
  }
}

.login-footer {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--border-light);
  text-align: center;
  font-size: 13px;
  color: var(--text-secondary);

  p {
    margin: 6px 0;
    font-weight: 500;
  }

  .tip {
    color: var(--danger-color);
    font-weight: 600;
    margin-top: 12px;
  }
}

:deep(.el-dialog) {
  border-radius: 16px;
  box-shadow: var(--shadow-lg);

  .el-dialog__header {
    padding: 24px;
    border-bottom: 1px solid var(--border-light);
    background: var(--bg-blue-light);

    .el-dialog__title {
      font-size: 20px;
      font-weight: 600;
      color: var(--text-color);
    }
  }

  .el-dialog__body {
    padding: 24px;
  }

  .el-dialog__footer {
    padding: 16px 24px;
    border-top: 1px solid var(--border-light);
  }

  .el-button {
    border-radius: 8px;
    font-weight: 500;
    padding: 10px 20px;
  }
}
</style>

