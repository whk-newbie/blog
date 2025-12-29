<template>
  <el-dialog
    v-model="visible"
    :title="t('login.title')"
    width="420px"
    :close-on-click-modal="false"
    @close="handleClose"
  >
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
          :placeholder="t('login.username')"
          size="large"
          prefix-icon="User"
          clearable
        />
      </el-form-item>

      <el-form-item prop="password">
        <el-input
          v-model="loginForm.password"
          type="password"
          :placeholder="t('login.password')"
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
          {{ loading ? t('common.loading') : t('login.submit') }}
        </el-button>
      </el-form-item>
    </el-form>

    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="showChangePasswordDialog"
      :title="t('user.changePassword')"
      width="400px"
      :close-on-click-modal="false"
      :close-on-press-escape="false"
      :show-close="false"
      append-to-body
    >
      <el-alert
        :title="t('login.changePasswordWarning')"
        type="warning"
        :description="t('login.changePasswordDescription')"
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
        <el-form-item :label="t('login.oldPassword')" prop="oldPassword">
          <el-input
            v-model="changePasswordForm.oldPassword"
            type="password"
            :placeholder="t('login.oldPasswordPlaceholder')"
            show-password
          />
        </el-form-item>

        <el-form-item :label="t('login.newPassword')" prop="newPassword">
          <el-input
            v-model="changePasswordForm.newPassword"
            type="password"
            :placeholder="t('login.newPasswordPlaceholder')"
            show-password
          />
        </el-form-item>

        <el-form-item :label="t('login.confirmPassword')" prop="confirmPassword">
          <el-input
            v-model="changePasswordForm.confirmPassword"
            type="password"
            :placeholder="t('login.confirmPasswordPlaceholder')"
            show-password
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="handleSkipChangePassword">{{ t('common.cancel') }}</el-button>
        <el-button
          type="primary"
          :loading="changingPassword"
          @click="handleChangePassword"
        >
          {{ t('common.confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </el-dialog>
</template>

<script setup>
import { ref, reactive, watch, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { useUserStore } from '@/store/user'
import api from '@/api'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'success'])

const router = useRouter()
const userStore = useUserStore()
const { t } = useI18n()

const visible = ref(props.modelValue)
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

// 登录表单验证规则（使用computed使其响应语言变化）
const loginRules = computed(() => ({
  username: [
    { required: true, message: t('login.usernameRequired'), trigger: 'blur' }
  ],
  password: [
    { required: true, message: t('login.passwordRequired'), trigger: 'blur' },
    { min: 6, message: t('validation.minLength', { min: 6 }), trigger: 'blur' }
  ]
}))

// 修改密码表单
const changePasswordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 修改密码表单验证规则（使用computed使其响应语言变化）
const changePasswordRules = computed(() => ({
  oldPassword: [
    { required: true, message: t('login.oldPasswordRequired'), trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: t('login.newPasswordRequired'), trigger: 'blur' },
    { min: 6, message: t('validation.minLength', { min: 6 }), trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: t('login.confirmPasswordRequired'), trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== changePasswordForm.newPassword) {
          callback(new Error(t('login.passwordMismatch')))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}))

// 监听 modelValue 变化
watch(() => props.modelValue, (newVal) => {
  visible.value = newVal
  if (newVal) {
    // 打开对话框时重置表单
    loginForm.username = ''
    loginForm.password = ''
    if (loginFormRef.value) {
      loginFormRef.value.clearValidate()
    }
  }
})

// 监听 visible 变化
watch(visible, (newVal) => {
  emit('update:modelValue', newVal)
})

// 处理登录
const handleLogin = async () => {
  if (!loginFormRef.value) return

  try {
    await loginFormRef.value.validate()
    loading.value = true

    const response = await api.auth.login(loginForm.username, loginForm.password)
    
    // 保存登录信息
    userStore.login(response)

    ElMessage.success(t('login.success'))

    // 如果使用默认密码，提示修改
    if (response.is_default_password) {
      changePasswordForm.oldPassword = loginForm.password
      showChangePasswordDialog.value = true
    } else {
      // 关闭登录对话框
      visible.value = false
      emit('success')
      // 跳转到管理后台
      router.push('/admin')
    }
  } catch (error) {
    console.error('登录失败:', error)
    ElMessage.error(error.message || t('login.error'))
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

    ElMessage.success(t('user.changePasswordSuccess'))
    showChangePasswordDialog.value = false
    visible.value = false

    // 更新状态
    userStore.isDefaultPassword = false
    localStorage.setItem('isDefaultPassword', 'false')

    emit('success')
    // 跳转到管理后台
    router.push('/admin')
  } catch (error) {
    console.error('修改密码失败:', error)
    ElMessage.error(error.message || t('user.changePasswordError'))
  } finally {
    changingPassword.value = false
  }
}

// 跳过修改密码
const handleSkipChangePassword = () => {
  showChangePasswordDialog.value = false
  visible.value = false
  emit('success')
  router.push('/admin')
}

// 关闭对话框
const handleClose = () => {
  visible.value = false
  loginForm.username = ''
  loginForm.password = ''
  if (loginFormRef.value) {
    loginFormRef.value.clearValidate()
  }
}
</script>

<style scoped lang="less">
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

