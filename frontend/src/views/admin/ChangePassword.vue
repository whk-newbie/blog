<template>
  <div class="change-password-page">
    <PageHeader :title="t('user.changePassword')" />

    <el-card class="password-card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        style="max-width: 500px"
      >
        <el-form-item :label="t('login.oldPassword')" prop="oldPassword">
          <el-input
            v-model="form.oldPassword"
            type="password"
            :placeholder="t('login.oldPasswordPlaceholder')"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item :label="t('login.newPassword')" prop="newPassword">
          <el-input
            v-model="form.newPassword"
            type="password"
            :placeholder="t('login.newPasswordPlaceholder')"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item :label="t('login.confirmPassword')" prop="confirmPassword">
          <el-input
            v-model="form.confirmPassword"
            type="password"
            :placeholder="t('login.confirmPasswordPlaceholder')"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">
            {{ t('common.confirm') }}
          </el-button>
          <el-button @click="handleReset">{{ t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <el-alert
        :title="t('user.passwordSecurityTip')"
        type="info"
        :closable="false"
        style="margin-top: 20px"
      >
        <ul style="margin: 10px 0; padding-left: 20px">
          <li>{{ t('user.passwordMinLength') }}</li>
          <li>{{ t('user.passwordSuggestion') }}</li>
          <li>{{ t('user.passwordSimpleWarning') }}</li>
          <li>{{ t('user.passwordChangeRegularly') }}</li>
        </ul>
      </el-alert>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import api from '@/api'
import PageHeader from '@/components/common/PageHeader.vue'

const router = useRouter()
const userStore = useUserStore()
const { t } = useI18n()
const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== form.newPassword) {
    callback(new Error(t('login.passwordMismatch')))
  } else {
    callback()
  }
}

const rules = computed(() => ({
  oldPassword: [
    { required: true, message: t('login.oldPasswordRequired'), trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: t('login.newPasswordRequired'), trigger: 'blur' },
    { min: 6, message: t('user.passwordMinLength'), trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: t('login.confirmPasswordRequired'), trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}))

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true

    await api.auth.changePassword(form.oldPassword, form.newPassword)

    ElMessage.success(t('user.changePasswordSuccessRelogin'))

    // 清空表单
    handleReset()

    // 退出登录
    setTimeout(() => {
      userStore.logout()
      router.push('/')
    }, 1500)
  } catch (error) {
    console.error('修改密码失败:', error)
    ElMessage.error(error.message || t('user.changePasswordError'))
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  formRef.value?.resetFields()
}
</script>

<style scoped lang="less">
.change-password-page {
  padding: 0;
}

.password-card {
  margin-top: 24px;
  border-radius: 12px;
  box-shadow: var(--shadow-sm);
  border: 1px solid var(--border-light);

  :deep(.el-card__body) {
    padding: 32px;
  }
}

:deep(.el-form) {
  .el-form-item__label {
    font-weight: 600;
    color: var(--text-color);
  }

  .el-input__wrapper {
    border-radius: 10px;
    box-shadow: 0 0 0 1px var(--border-color) inset;
    transition: all 0.3s;

    &:hover {
      box-shadow: 0 0 0 1px var(--primary-light) inset;
    }

    &.is-focus {
      box-shadow: 0 0 0 2px var(--primary-color) inset;
    }
  }

  .el-button {
    border-radius: 10px;
    font-weight: 600;
    padding: 10px 24px;
    transition: all 0.3s;

    &.el-button--primary {
      background: var(--primary-color);
      border-color: var(--primary-color);

      &:hover {
        background: var(--primary-light);
        box-shadow: 0 4px 14px rgba(37, 99, 235, 0.3);
        transform: translateY(-2px);
      }
    }
  }
}

:deep(.el-alert) {
  border-radius: 10px;
  border: 1px solid var(--border-blue);
  background: var(--bg-blue-light);

  .el-alert__title {
    font-weight: 600;
    color: var(--primary-color);
  }

  ul {
    li {
      color: var(--text-secondary);
      margin: 8px 0;
      font-size: 14px;
    }
  }
}
</style>

