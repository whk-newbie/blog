<template>
  <div class="ai-providers-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>AI 提供方管理</span>
          <el-button type="primary" @click="showCreate">添加提供方</el-button>
        </div>
      </template>

      <el-table :data="providers" v-loading="loading">
        <el-table-column prop="name" label="名称" />
        <el-table-column prop="provider_type" label="类型" width="100" />
        <el-table-column prop="model" label="模型" width="150" />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.is_enabled ? 'success' : 'info'">
              {{ row.is_enabled ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="优先级" width="80" />
        <el-table-column label="操作" width="320">
          <template #default="{ row }">
            <el-button size="small" @click="editProvider(row)">编辑</el-button>
            <el-button size="small" @click="checkProvider(row)" :loading="checkingId === row.id">检测</el-button>
            <el-popconfirm title="确定删除此提供方？" @confirm="deleteProvider(row)">
              <template #reference>
                <el-button size="small" type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Create/Edit Dialog -->
    <el-dialog v-model="showDialog" :title="editing ? '编辑提供方' : '添加提供方'" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称">
          <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.provider_type">
            <el-option label="Claude" value="claude" />
            <el-option label="OpenAI (ChatGPT)" value="openai" />
            <el-option label="DeepSeek" value="deepseek" />
            <el-option label="自定义 (OpenAI兼容)" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="form.api_key" type="password" show-password />
        </el-form-item>
        <el-form-item label="Base URL" v-if="form.provider_type === 'custom'">
          <el-input v-model="form.base_url" placeholder="https://api.example.com/v1/chat/completions" />
        </el-form-item>
        <el-form-item label="模型">
          <el-input v-model="form.model" :placeholder="defaultModel" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="form.is_enabled" />
        </el-form-item>
        <el-form-item label="优先级">
          <el-input-number v-model="form.sort_order" :min="0" :max="100" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="saveProvider" :loading="saving">保存</el-button>
      </template>
    </el-dialog>

    <!-- Check Result Dialog -->
    <el-dialog v-model="showCheckResult" title="检测结果" width="400px">
      <el-descriptions :column="1">
        <el-descriptions-item label="可用性">
          <el-tag :type="checkResult.available ? 'success' : 'danger'">
            {{ checkResult.available ? '可用' : '不可用' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="余额">{{ checkResult.balance }}</el-descriptions-item>
        <el-descriptions-item v-if="checkResult.error" label="错误">{{ checkResult.error }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import api from '@/api'
import aiApi from '@/api/ai'

const loading = ref(false)
const saving = ref(false)
const providers = ref([])
const showDialog = ref(false)
const editing = ref(false)
const editingId = ref(null)
const checkingId = ref(null)
const showCheckResult = ref(false)
const checkResult = ref({ available: false, balance: 0, error: '' })

const form = ref({
  name: '',
  provider_type: 'deepseek',
  api_key: '',
  base_url: '',
  model: '',
  is_enabled: true,
  sort_order: 0,
})

const defaultModel = computed(() => {
  switch (form.value.provider_type) {
    case 'claude': return 'claude-sonnet-4-20250514'
    case 'openai': return 'gpt-4o'
    case 'deepseek': return 'deepseek-chat'
    default: return ''
  }
})

const fetchProviders = async () => {
  loading.value = true
  try {
    providers.value = await aiApi.listProviders()
  } catch (e) {
    ElMessage.error('获取提供方列表失败')
  } finally {
    loading.value = false
  }
}

const showCreate = () => {
  editing.value = false
  editingId.value = null
  form.value = { name: '', provider_type: 'deepseek', api_key: '', base_url: '', model: '', is_enabled: true, sort_order: 0 }
  showDialog.value = true
}

const editProvider = (row) => {
  editing.value = true
  editingId.value = row.id
  form.value = { ...row, api_key: '' }
  showDialog.value = true
}

const saveProvider = async () => {
  saving.value = true
  try {
    if (editing.value) {
      await aiApi.updateProvider(editingId.value, form.value)
      ElMessage.success('更新成功')
    } else {
      await aiApi.createProvider(form.value)
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    fetchProviders()
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    saving.value = false
  }
}

const deleteProvider = async (row) => {
  try {
    await aiApi.deleteProvider(row.id)
    ElMessage.success('删除成功')
    fetchProviders()
  } catch (e) {
    ElMessage.error('删除失败')
  }
}

const checkProvider = async (row) => {
  checkingId.value = row.id
  try {
    checkResult.value = await aiApi.checkProvider(row.id)
    showCheckResult.value = true
  } catch (e) {
    ElMessage.error('检测失败')
  } finally {
    checkingId.value = null
  }
}

onMounted(() => {
  fetchProviders()
})
</script>
