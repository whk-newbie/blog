<template>
  <div class="config-manage-page">
    <div class="page-header-bar">
      <h2>{{ t('nav.config') }}</h2>
      <el-button type="primary" @click="handleCreate" size="small">
        <el-icon><Plus /></el-icon>添加配置
      </el-button>
    </div>

    <el-collapse v-model="activeNames">
      <!-- Site Info -->
      <el-collapse-item title="站点信息" name="site">
        <site-info-config />
      </el-collapse-item>

      <!-- Email -->
      <el-collapse-item name="email">
        <template #title>
          <div class="collapse-title">邮箱配置 <el-tag size="small">{{ emailConfigs.length }}</el-tag></div>
        </template>
        <config-list :configs="emailConfigs" :loading="loading" @edit="handleEdit" @delete="handleDelete" @toggle-active="handleToggleActive" />
      </el-collapse-item>

      <!-- API Tokens -->
      <el-collapse-item name="api_token">
        <template #title>
          <div class="collapse-title">API Token <el-tag size="small">{{ apiTokenConfigs.length }}</el-tag></div>
        </template>
        <config-list :configs="apiTokenConfigs" :loading="loading" @edit="handleEdit" @delete="handleDelete" @toggle-active="handleToggleActive" />
      </el-collapse-item>

      <!-- Salt -->
      <el-collapse-item name="salt">
        <template #title>
          <div class="collapse-title">加密盐值 <el-tag size="small">{{ saltConfigs.length }}</el-tag></div>
        </template>
        <config-list :configs="saltConfigs" :loading="loading" @edit="handleEdit" @delete="handleDelete" @toggle-active="handleToggleActive" />
      </el-collapse-item>

      <!-- IP Blacklist -->
      <el-collapse-item name="ip">
        <template #title>
          <div class="collapse-title">IP 黑名单 <el-tag size="small">{{ ipBlacklistConfigs.length }}</el-tag></div>
        </template>
        <config-list :configs="ipBlacklistConfigs" :loading="loading" @edit="handleEdit" @delete="handleDelete" @toggle-active="handleToggleActive" />
      </el-collapse-item>

      <!-- Other -->
      <el-collapse-item v-if="otherConfigs.length > 0" name="other">
        <template #title>
          <div class="collapse-title">{{ t('config.other') }} <el-tag size="small">{{ otherConfigs.length }}</el-tag></div>
        </template>
        <config-list :configs="otherConfigs" :loading="loading" @edit="handleEdit" @delete="handleDelete" @toggle-active="handleToggleActive" />
      </el-collapse-item>
    </el-collapse>

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
const activeNames = ref(['site'])
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

:deep(.el-collapse) {
  border: none;

  .el-collapse-item {
    margin-bottom: 8px;
    border: 1px solid var(--border-color);
    border-radius: var(--border-radius-sm);
    overflow: hidden;

    .el-collapse-item__header {
      padding: 0 20px;
      height: 48px;
      font-size: 15px;
      font-weight: 600;
      color: var(--text-heading);
      background: var(--card-bg);
      border-bottom: none;

      &:hover {
        background: var(--hover-bg);
      }
    }

    .el-collapse-item__wrap {
      border-top: 1px solid var(--border-color);
      background: var(--card-bg);
    }

    .el-collapse-item__content {
      padding: 0;
    }
  }
}

.collapse-title {
  display: flex;
  align-items: center;
  gap: 10px;
}
</style>
