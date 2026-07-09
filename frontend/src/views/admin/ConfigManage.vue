<template>
  <div class="config-manage-page">
    <div class="page-header-bar">
      <h2>系统配置</h2>
    </div>

    <el-collapse v-model="activeNames">
      <!-- Admin Entry -->
      <el-collapse-item name="admin">
        <template #title>
          <div class="collapse-title">后台入口</div>
        </template>
        <div class="section-body">
          <el-form label-width="140px">
            <el-form-item label="后台访问路径">
              <el-input v-model="adminPath" placeholder="admin" style="width:300px" />
              <div class="form-tip">修改后通过新路径访问后台</div>
            </el-form-item>
            <el-form-item label="前台显示入口">
              <el-switch v-model="showAdminLink" />
              <span style="margin-left:10px;color:var(--text-secondary);font-size:13px">关闭后已登录用户仍可访问</span>
            </el-form-item>
            <el-form-item>
              <el-button type="primary" size="small" @click="saveAdminConfig">保存</el-button>
            </el-form-item>
          </el-form>
        </div>
      </el-collapse-item>

      <!-- Site Info -->
      <el-collapse-item title="站点信息" name="site">
        <site-info-config />
      </el-collapse-item>

      <!-- AI Providers -->
      <el-collapse-item name="ai">
        <template #title>
          <div class="collapse-title">AI 提供方</div>
        </template>
        <div class="section-body">
          <div style="margin-bottom:12px;display:flex;justify-content:flex-end">
            <el-button type="primary" size="small" @click="showAiDialog = true">添加 AI 提供方</el-button>
          </div>
          <el-table :data="aiProviders" v-loading="loadingAI" size="small">
            <el-table-column prop="name" label="名称" width="140" />
            <el-table-column prop="provider_type" label="类型" width="90" />
            <el-table-column prop="model" label="模型" width="150" />
            <el-table-column label="启用" width="70">
              <template #default="{ row }">
                <el-tag :type="row.is_enabled ? 'success' : 'info'" size="small">{{ row.is_enabled ? '是' : '否' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="sort_order" label="优先级" width="70" />
            <el-table-column label="操作">
              <template #default="{ row }">
                <el-button size="small" @click="editAi(row)">编辑</el-button>
                <el-button size="small" :loading="checkingId === row.id" @click="checkAi(row)">检测</el-button>
                <el-popconfirm title="确定删除？" @confirm="deleteAi(row)">
                  <template #reference><el-button size="small" type="danger">删除</el-button></template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-collapse-item>

      <!-- IP Blacklist -->
      <el-collapse-item name="ip">
        <template #title>
          <div class="collapse-title">IP 黑名单</div>
        </template>
        <config-list :configs="ipBlacklistConfigs" :loading="loading" @edit="handleEdit" @delete="handleDelete" @toggle-active="handleToggleActive" />
        <div class="section-body" style="border-top:1px solid var(--border-color);padding-top:12px;display:flex;justify-content:flex-end">
          <el-button size="small" type="primary" @click="handleCreate('ip_blacklist')">添加黑名单</el-button>
        </div>
      </el-collapse-item>

      <!-- Other -->
      <el-collapse-item v-if="otherConfigs.length > 0" name="other">
        <template #title>
          <div class="collapse-title">其他配置</div>
        </template>
        <config-list :configs="otherConfigs" :loading="loading" @edit="handleEdit" @delete="handleDelete" @toggle-active="handleToggleActive" />
      </el-collapse-item>
    </el-collapse>

    <!-- AI Provider Dialog -->
    <el-dialog v-model="showAiDialog" :title="editingAi ? '编辑 AI 提供方' : '添加 AI 提供方'" width="480px">
      <el-form :model="aiForm" label-width="100px">
        <el-form-item label="名称"><el-input v-model="aiForm.name" /></el-form-item>
        <el-form-item label="类型">
          <el-select v-model="aiForm.provider_type" style="width:100%">
            <el-option label="Claude" value="claude" />
            <el-option label="OpenAI" value="openai" />
            <el-option label="DeepSeek" value="deepseek" />
            <el-option label="自定义" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item label="API Key"><el-input v-model="aiForm.api_key" type="password" show-password /></el-form-item>
        <el-form-item v-if="aiForm.provider_type==='custom'" label="Base URL"><el-input v-model="aiForm.base_url" /></el-form-item>
        <el-form-item label="模型"><el-input v-model="aiForm.model" /></el-form-item>
        <el-form-item label="启用"><el-switch v-model="aiForm.is_enabled" /></el-form-item>
        <el-form-item label="优先级"><el-input-number v-model="aiForm.sort_order" :min="0" :max="100" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAiDialog=false">取消</el-button>
        <el-button type="primary" @click="saveAi">保存</el-button>
      </template>
    </el-dialog>

    <!-- Generic Config Dialog -->
    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑配置' : '添加配置'" width="480px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="配置键" prop="config_key"><el-input v-model="form.config_key" :disabled="isEdit" /></el-form-item>
        <el-form-item label="配置值" prop="config_value"><el-input v-model="form.config_value" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="form.description" /></el-form-item>
        <el-row :gutter="20">
          <el-col :span="12"><el-form-item label="加密"><el-switch v-model="form.is_encrypted" /></el-form-item></el-col>
          <el-col :span="12"><el-form-item label="启用"><el-switch v-model="form.is_active" /></el-form-item></el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible=false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确认</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import api from '@/api'
import aiApi from '@/api/ai'
import ConfigList from './components/ConfigList.vue'
import SiteInfoConfig from './components/SiteInfoConfig.vue'

const loading = ref(false)
const loadingAI = ref(false)
const configs = ref([])
const activeNames = ref(['admin', 'site'])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const adminPath = ref('admin')
const showAdminLink = ref(true)
const aiProviders = ref([])
const showAiDialog = ref(false)
const editingAi = ref(false)
const editingAiId = ref(null)
const checkingId = ref(null)

const form = reactive({ id: undefined, config_key: '', config_value: '', description: '', is_encrypted: true, is_active: true })
const aiForm = reactive({ name: '', provider_type: 'deepseek', api_key: '', base_url: '', model: '', is_enabled: true, sort_order: 0 })

const rules = computed(() => ({
  config_key: [{ required: true, message: '请输入' }],
  config_value: [{ required: true, message: '请输入' }],
}))

const ipBlacklistConfigs = computed(() => configs.value.filter(c => c.config_type === 'ip_blacklist'))
const otherConfigs = computed(() => configs.value.filter(c =>
  !['email', 'api_token', 'crawler_token', 'salt', 'ip_blacklist', 'site_info', 'site_config'].includes(c.config_type) &&
  !['admin_path', 'show_admin_link'].includes(c.config_key)
))

const fetchConfigs = async () => {
  loading.value = true
  try {
    const r = await api.config.getConfigs()
    configs.value = r || []
    const ap = r?.find(c => c.config_key === 'admin_path')
    if (ap) adminPath.value = ap.config_value
    const sal = r?.find(c => c.config_key === 'show_admin_link')
    if (sal) showAdminLink.value = sal.config_value === 'true'
  } catch (e) { ElMessage.error('加载失败') } finally { loading.value = false }
}

const fetchAIProviders = async () => {
  loadingAI.value = true
  try { aiProviders.value = await aiApi.listProviders() || [] } catch (e) {} finally { loadingAI.value = false }
}

const saveAdminConfig = async () => {
  try {
    const ap = configs.value.find(c => c.config_key === 'admin_path')
    const sal = configs.value.find(c => c.config_key === 'show_admin_link')
    if (ap) await api.config.updateConfig(ap.id, { config_value: adminPath.value, config_type: 'site_config', description: '后台访问路径', is_encrypted: false, is_active: true })
    else await api.config.createConfig({ config_key: 'admin_path', config_value: adminPath.value, config_type: 'site_config', description: '后台访问路径', is_encrypted: false, is_active: true })
    if (sal) await api.config.updateConfig(sal.id, { config_value: String(showAdminLink.value), config_type: 'site_config', description: '前台入口开关', is_encrypted: false, is_active: true })
    else await api.config.createConfig({ config_key: 'show_admin_link', config_value: String(showAdminLink.value), config_type: 'site_config', description: '前台入口开关', is_encrypted: false, is_active: true })
    localStorage.setItem('admin_path', adminPath.value)
    ElMessage.success('已保存')
    fetchConfigs()
  } catch (e) { ElMessage.error('保存失败') }
}

const handleCreate = (type) => {
  isEdit.value = false
  form.id = undefined; form.config_key = ''; form.config_value = ''; form.config_type = type; form.description = ''; form.is_encrypted = true; form.is_active = true
  dialogVisible.value = true
}
const handleEdit = (c) => { isEdit.value = true; Object.assign(form, { id: c.id, config_key: c.config_key, config_value: c.config_value, config_type: c.config_type, description: c.description||'', is_encrypted: c.is_encrypted, is_active: c.is_active }); dialogVisible.value = true }
const handleSubmit = async () => {
  await formRef.value.validate()
  const d = { config_key: form.config_key, config_value: form.config_value, config_type: form.config_type || 'ip_blacklist', description: form.description, is_encrypted: form.is_encrypted, is_active: form.is_active }
  if (isEdit.value) await api.config.updateConfig(form.id, d); else await api.config.createConfig(d)
  dialogVisible.value = false; fetchConfigs(); ElMessage.success('已保存')
}
const handleDelete = async (id) => { await api.config.deleteConfig(id); fetchConfigs(); ElMessage.success('已删除') }
const handleToggleActive = async (c) => { await api.config.updateConfig(c.id, { config_value: c.config_value, is_active: !c.is_active }); fetchConfigs() }

const editAi = (row) => { editingAi.value = true; editingAiId.value = row.id; Object.assign(aiForm, row); showAiDialog.value = true }
const saveAi = async () => {
  try {
    if (editingAi.value) { await aiApi.updateProvider(editingAiId.value, aiForm) }
    else { await aiApi.createProvider(aiForm) }
    showAiDialog.value = false; fetchAIProviders(); ElMessage.success('已保存')
  } catch (e) { ElMessage.error('保存失败') }
}
const deleteAi = async (row) => { await aiApi.deleteProvider(row.id); fetchAIProviders(); ElMessage.success('已删除') }
const checkAi = async (row) => { checkingId.value = row.id; try { const r = await aiApi.checkProvider(row.id); ElMessage.info(r.available ? '可用' : '不可用: ' + (r.error||'未知')) } finally { checkingId.value = null } }

onMounted(() => { fetchConfigs(); fetchAIProviders() })
</script>

<style scoped lang="less">
.config-manage-page { width: 100%; }
.page-header-bar { display: flex; justify-content: space-between; align-items: center; margin-bottom: 24px;
  h2 { margin: 0; font-size: 20px; font-weight: 600; color: var(--text-heading); }
}
.section-body { padding: 16px 20px; }
.form-tip { font-size: 12px; color: var(--text-secondary); margin-top: 4px; }

:deep(.el-collapse) { border: none;
  .el-collapse-item { margin-bottom: 8px; border: 1px solid var(--border-color); border-radius: var(--border-radius-sm); overflow: hidden;
    .el-collapse-item__header { padding: 0 20px; height: 48px; font-size: 15px; font-weight: 600; color: var(--text-heading); background: var(--card-bg); border-bottom: none;
      &:hover { background: var(--hover-bg); }
    }
    .el-collapse-item__wrap { border-top: 1px solid var(--border-color); background: var(--card-bg); }
    .el-collapse-item__content { padding: 0; }
  }
}
.collapse-title { display: flex; align-items: center; gap: 10px; }
</style>
