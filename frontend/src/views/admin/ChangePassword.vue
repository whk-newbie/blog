<template>
  <div class="change-password-page">
    <PageHeader title="修改密码" />

    <el-card class="password-card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
        style="max-width: 500px"
      >
        <el-form-item label="旧密码" prop="oldPassword">
          <el-input
            v-model="form.oldPassword"
            type="password"
            placeholder="请输入旧密码"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="form.newPassword"
            type="password"
            placeholder="请输入新密码（至少6位）"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="form.confirmPassword"
            type="password"
            placeholder="请再次输入新密码"
            show-password
            clearable
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">
            确认修改
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-alert
        title="密码安全提示"
        type="info"
        :closable="false"
        style="margin-top: 20px"
      >
        <ul style="margin: 10px 0; padding-left: 20px">
          <li>密码长度至少为6位</li>
          <li>建议使用字母、数字和特殊字符的组合</li>
          <li>不要使用过于简单的密码</li>
          <li>定期更换密码以保证账户安全</li>
        </ul>
      </el-alert>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import api from '@/api'
import PageHeader from '@/components/common/PageHeader.vue'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref(null)
const loading = ref(false)

const form = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== form.newPassword) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  oldPassword: [
    { required: true, message: '请输入旧密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入新密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return

  try {
    await formRef.value.validate()
    loading.value = true

    await api.auth.changePassword(form.oldPassword, form.newPassword)

    ElMessage.success('密码修改成功，请重新登录')

    // 清空表单
    handleReset()

    // 退出登录
    setTimeout(() => {
      userStore.logout()
      router.push('/admin/login')
    }, 1500)
  } catch (error) {
    console.error('修改密码失败:', error)
    ElMessage.error(error.message || '修改密码失败')
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
  padding: 20px;
}

.password-card {
  margin-top: 20px;
}
</style>

