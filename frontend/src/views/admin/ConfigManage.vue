<template>
  <div class="config-manage-page">
    <div class="page-header-bar">
      <h2>{{ t('nav.config') }}</h2>
      <el-button type="primary" @click="handleCreate" size="small">
        <el-icon><Plus /></el-icon>添加配置
      </el-button>
    </div>

    <!-- Site Info -->
    <el-card class="section-card" shadow="never">
      <template #header><span class="section-title">站点信息</span></template>
      <site-info-config />
    </el-card>

    <!-- Email -->
    <el-card class="section-card" shadow="never">
      <template #header>
        <div class="section-header">
          <span class="section-title">邮箱配置</span>
          <el-tag size="small" type="info">{{ emailConfigs.length }}</el-tag>
        </div>
      </template>
      <config-list :configs="emailConfigs" :loading="loading" @edit="handleEdit" @delete="handleDelete" @toggle-active="handleToggleActive" />
    </el-card>

    <!-- API Tokens -->
    <el-card class="section-card" shadow="never">
      <template #header>
        <div class="section-header">
          <span class="section-title">API Token</span>
          <el-tag size="small" type="info">{{ apiTokenConfigs.length }}</el-tag>
        </div>
      </template>
      <config-list :configs="apiTokenConfigs" :loading="loading" @edit="handleEdit" @delete="handleDelete" @toggle-active="handleToggleActive" />
    </el-card>

    <!-- Encryption Salt -->
    <el-card class="section-card" shadow="never">
      <template #header>
        <div class="section-header">
          <span class="section-title">加密盐值</span>
          <el-tag size="small" type="info">{{ saltConfigs.length }}</el-tag>
        </div>
      </template>
      <config-list :configs="saltConfigs" :loading="loading" @edit="handleEdit" @delete="handleDelete" @toggle-active="handleToggleActive" />
    </el-card>

    <!-- IP Blacklist -->
    <el-card class="section-card" shadow="never">
      <template #header>
        <div class="section-header">
          <span class="section-title">IP 黑名单</span>
          <el-tag size="small" type="info">{{ ipBlacklistConfigs.length }}</el-tag>
        </div>
      </template>
      <config-list :configs="ipBlacklistConfigs" :loading="loading" @edit="handleEdit" @delete="handleDelete" @toggle-active="handleToggleActive" />
    </el-card>

    <!-- Other -->
    <el-card v-if="otherConfigs.length > 0" class="section-card" shadow="never">
      <template #header>
        <div class="section-header">
          <span class="section-title">{{ t('config.other') }}</span>
          <el-tag size="small" type="info">{{ otherConfigs.length }}</el-tag>
        </div>
      </template>
      <config-list :configs="otherConfigs" :loading="loading" @edit="handleEdit" @delete="handleDelete" @toggle-active="handleToggleActive" />
    </el-card>

    <!-- Create/Edit Dialog -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? t('config.editConfig') : t('config.addConfig')"
      width="520px"
    >
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item :label="t('config.configKey')" prop="config_key">
          <el-input v-model="form.config_key" :placeholder="t('config.configKeyPlaceholder')" :disabled="isEdit" />
        </el-form-item>
        <el-form-item :label="t('config.configValue')" prop="config_value">
          <el-input v-model="form.config_value" type="textarea" :rows="4" :placeholder="t('config.configValuePlaceholder')" />
        </el-form-item>
        <el-form-item :label="t('config.configType')" prop="config_type">
          <el-select v-model="form.config_type" :placeholder="t('config.configTypePlaceholder')" :disabled="isEdit" style="width:100%">
            <el-option label="邮箱" value="email" />
            <el-option label="API Token" value="api_token" />
            <el-option label="加密盐值" value="salt" />
            <el-option label="IP 黑名单" value="ip_blacklist" />
            <el-option label="站点信息" value="site_info" />
            <el-option label="应用密钥" value="application_key" />
          </el-select>
        </el-form-item>
        <el-form-item :label="t('config.description')">
          <el-input v-model="form.description" type="textarea" :rows="2" :placeholder="t('config.descriptionPlaceholder')" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item :label="t('config.isEncrypted')">
              <el-switch v-model="form.is_encrypted" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="t('config.isActive')">
              <el-switch v-model="form.is_active" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import api from '@/api'
import ConfigList from './components/ConfigList.vue'
import SiteInfoConfig from './components/SiteInfoConfig.vue'

const { t } = useI18n()
const loading = ref(false)
const configs = ref([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)

const form = reactive({
  id: undefined,
  config_key: '',
  config_value: '',
  config_type: 'email',
  description: '',
  is_encrypted: true,
  is_active: true
})

const rules = computed(() => ({
  config_key: [{ required: true, message: '请输入配置键', trigger: 'blur' }],
  config_value: [{ required: true, message: '请输入配置值', trigger: 'blur' }],
  config_type: [{ required: true, message: '请选择配置类型', trigger: 'change' }]
}))

const emailConfigs = computed(() => configs.value.filter(c => c.config_type === 'email'))
const apiTokenConfigs = computed(() => configs.value.filter(c => c.config_type === 'api_token'))
const saltConfigs = computed(() => configs.value.filter(c => c.config_type === 'salt'))
const ipBlacklistConfigs = computed(() => configs.value.filter(c => c.config_type === 'ip_blacklist'))
const otherConfigs = computed(() => configs.value.filter(c =>
  !['email', 'api_token', 'salt', 'ip_blacklist', 'site_info'].includes(c.config_type)
))

const fetchConfigs = async () => {
  loading.value = true
  try {
    const response = await api.config.getConfigs()
    configs.value = response || []
  } catch (error) {
    ElMessage.error(t('config.loadError'))
  } finally {
    loading.value = false
  }
}

const handleCreate = () => {
  isEdit.value = false
  form.id = undefined
  form.config_key = ''
  form.config_value = ''
  form.config_type = 'email'
  form.description = ''
  form.is_encrypted = true
  form.is_active = true
  dialogVisible.value = true
}

const handleEdit = (config) => {
  isEdit.value = true
  Object.assign(form, {
    id: config.id, config_key: config.config_key, config_value: config.config_value,
    config_type: config.config_type, description: config.description || '',
    is_encrypted: config.is_encrypted, is_active: config.is_active,
  })
  dialogVisible.value = true
}

const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    const data = { config_key: form.config_key, config_value: form.config_value, config_type: form.config_type, description: form.description, is_encrypted: form.is_encrypted, is_active: form.is_active }
    if (isEdit.value) {
      await api.config.updateConfig(form.id, data)
      ElMessage.success(t('config.updateSuccess'))
    } else {
      await api.config.createConfig(data)
      ElMessage.success(t('config.createSuccess'))
    }
    dialogVisible.value = false
    fetchConfigs()
  } catch (error) {
    if (error instanceof Error) {
      ElMessage.error(isEdit.value ? t('config.updateError') : t('config.createError'))
    }
  }
}

const handleDelete = async (id) => {
  try {
    await api.config.deleteConfig(id)
    ElMessage.success(t('config.deleteSuccess'))
    fetchConfigs()
  } catch (error) {
    ElMessage.error(t('config.deleteError'))
  }
}

const handleToggleActive = async (config) => {
  try {
    await api.config.updateConfig(config.id, { config_value: config.config_value, is_active: !config.is_active })
    ElMessage.success(t('config.updateSuccess'))
    fetchConfigs()
  } catch (error) {
    ElMessage.error(t('config.updateError'))
  }
}

onMounted(() => { fetchConfigs() })
</script>

<style scoped lang="less">
.config-manage-page {
  max-width: 960px;
}

.page-header-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;

  h2 {
    margin: 0;
    font-size: 20px;
    font-weight: 600;
    color: var(--text-heading);
  }
}

.section-card {
  margin-bottom: 16px;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-sm);

  :deep(.el-card__header) {
    padding: 14px 20px;
    background: var(--bg-secondary);
    border-bottom: 1px solid var(--border-color);
  }

  :deep(.el-card__body) {
    padding: 0;
  }
}

.section-header {
  display: flex;
  align-items: center;
  gap: 12px;
}

.section-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text-heading);
}
</style>
