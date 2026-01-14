<template>
  <div class="config-manage-page">
    <page-header :title="t('nav.config')">
      <template #extra>
        <el-button type="primary" @click="handleCreate">
          <el-icon><Plus /></el-icon>
          {{ t('config.addConfig') }}
        </el-button>
      </template>
    </page-header>

    <!-- 配置类型标签页 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane :label="t('config.emailConfig')" name="email">
        <config-list
          :configs="emailConfigs"
          :loading="loading"
          @edit="handleEdit"
          @delete="handleDelete"
          @toggle-active="handleToggleActive"
        />
      </el-tab-pane>
      <el-tab-pane :label="t('config.apiToken')" name="api_token">
        <config-list
          :configs="apiTokenConfigs"
          :loading="loading"
          @edit="handleEdit"
          @delete="handleDelete"
          @toggle-active="handleToggleActive"
        />
      </el-tab-pane>
      <el-tab-pane :label="t('config.crawlerToken')" name="crawler_token">
        <div class="crawler-token-section">
          <div class="section-header">
            <el-button type="primary" @click="handleGenerateCrawlerToken">
              <el-icon><Plus /></el-icon>
              {{ t('config.generateCrawlerToken') }}
            </el-button>
          </div>
          <config-list
            :configs="crawlerTokenConfigs"
            :loading="loading"
            @edit="handleEdit"
            @delete="handleDelete"
            @toggle-active="handleToggleActive"
          />
        </div>
      </el-tab-pane>
      <el-tab-pane :label="t('config.encryptionSalt')" name="salt">
        <config-list
          :configs="saltConfigs"
          :loading="loading"
          @edit="handleEdit"
          @delete="handleDelete"
          @toggle-active="handleToggleActive"
        />
      </el-tab-pane>
      <el-tab-pane :label="t('config.ipBlacklist')" name="ip_blacklist">
        <config-list
          :configs="ipBlacklistConfigs"
          :loading="loading"
          @edit="handleEdit"
          @delete="handleDelete"
          @toggle-active="handleToggleActive"
        />
      </el-tab-pane>
      <el-tab-pane :label="t('config.siteInfo')" name="site_info">
        <site-info-config />
      </el-tab-pane>
      <el-tab-pane :label="t('config.other')" name="other">
        <config-list
          :configs="otherConfigs"
          :loading="loading"
          @edit="handleEdit"
          @delete="handleDelete"
          @toggle-active="handleToggleActive"
        />
      </el-tab-pane>
    </el-tabs>

    <!-- 创建/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? t('config.editConfig') : t('config.addConfig')"
      width="600px"
    >
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
      >
        <el-form-item :label="t('config.configKey')" prop="config_key">
          <el-input
            v-model="form.config_key"
            :placeholder="t('config.configKeyPlaceholder')"
            :disabled="isEdit"
          />
        </el-form-item>

        <el-form-item :label="t('config.configValue')" prop="config_value">
          <el-input
            v-model="form.config_value"
            type="textarea"
            :rows="4"
            :placeholder="t('config.configValuePlaceholder')"
          />
        </el-form-item>

        <el-form-item :label="t('config.configType')" prop="config_type">
          <el-select
            v-model="form.config_type"
            :placeholder="t('config.configTypePlaceholder')"
            :disabled="isEdit"
            style="width: 100%"
          >
            <el-option :label="t('config.emailConfig')" value="email" />
            <el-option :label="t('config.apiToken')" value="api_token" />
            <el-option :label="t('config.crawlerToken')" value="crawler_token" />
            <el-option :label="t('config.encryptionSalt')" value="salt" />
            <el-option :label="t('config.ipBlacklist')" value="ip_blacklist" />
            <el-option :label="t('config.siteInfo')" value="site_info" />
            <el-option :label="t('config.applicationKey')" value="application_key" />
          </el-select>
        </el-form-item>

        <el-form-item :label="t('config.description')">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="2"
            :placeholder="t('config.descriptionPlaceholder')"
          />
        </el-form-item>

        <el-form-item :label="t('config.isEncrypted')">
          <el-switch v-model="form.is_encrypted" />
        </el-form-item>

        <el-form-item :label="t('config.isActive')">
          <el-switch v-model="form.is_active" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 生成爬虫Token对话框 -->
    <el-dialog
      v-model="generateTokenDialogVisible"
      :title="t('config.generateCrawlerToken')"
      width="500px"
    >
      <el-form
        ref="tokenFormRef"
        :model="tokenForm"
        :rules="tokenRules"
        label-width="100px"
      >
        <el-form-item :label="t('config.tokenName')" prop="name">
          <el-input
            v-model="tokenForm.name"
            :placeholder="t('config.tokenNamePlaceholder')"
          />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="generateTokenDialogVisible = false">{{ t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleGenerateTokenSubmit">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- Token显示对话框 -->
    <el-dialog
      v-model="tokenDisplayDialogVisible"
      :title="t('config.tokenGenerated')"
      width="500px"
    >
      <el-alert
        :title="t('config.tokenWarning')"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      />
      <el-input
        v-model="generatedToken"
        readonly
        type="textarea"
        :rows="3"
      />
      <template #footer>
        <el-button @click="handleCopyToken">
          <el-icon><DocumentCopy /></el-icon>
          {{ t('config.copyToken') }}
        </el-button>
        <el-button type="primary" @click="tokenDisplayDialogVisible = false">{{ t('common.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus, DocumentCopy } from '@element-plus/icons-vue'
import api from '@/api'
import PageHeader from '@/components/common/PageHeader.vue'
import ConfigList from './components/ConfigList.vue'
import SiteInfoConfig from './components/SiteInfoConfig.vue'

const { t } = useI18n()
const loading = ref(false)
const configs = ref([])
const activeTab = ref('email')
const dialogVisible = ref(false)
const generateTokenDialogVisible = ref(false)
const tokenDisplayDialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const tokenFormRef = ref(null)
const generatedToken = ref('')

const form = reactive({
  id: undefined,
  config_key: '',
  config_value: '',
  config_type: 'email',
  description: '',
  is_encrypted: true,
  is_active: true
})

const tokenForm = reactive({
  name: ''
})

const rules = computed(() => ({
  config_key: [
    { required: true, message: t('config.configKeyPlaceholder'), trigger: 'blur' }
  ],
  config_value: [
    { required: true, message: t('config.configValuePlaceholder'), trigger: 'blur' }
  ],
  config_type: [
    { required: true, message: t('config.configTypePlaceholder'), trigger: 'change' }
  ]
}))

const tokenRules = computed(() => ({
  name: [
    { required: true, message: t('config.tokenNamePlaceholder'), trigger: 'blur' }
  ]
}))

// 按类型分组配置
const emailConfigs = computed(() => configs.value.filter(c => c.config_type === 'email'))
const apiTokenConfigs = computed(() => configs.value.filter(c => c.config_type === 'api_token'))
const crawlerTokenConfigs = computed(() => configs.value.filter(c => c.config_type === 'crawler_token'))
const saltConfigs = computed(() => configs.value.filter(c => c.config_type === 'salt'))
const ipBlacklistConfigs = computed(() => configs.value.filter(c => c.config_type === 'ip_blacklist'))
const otherConfigs = computed(() => configs.value.filter(c => 
  !['email', 'api_token', 'crawler_token', 'salt', 'ip_blacklist', 'site_info'].includes(c.config_type)
))

// 获取配置列表
const fetchConfigs = async () => {
  try {
    loading.value = true
    // 获取所有配置,前端按类型分组
    const response = await api.config.getConfigs()
    configs.value = response || []
  } catch (error) {
    console.error('获取配置列表失败:', error)
    ElMessage.error(t('config.loadError'))
  } finally {
    loading.value = false
  }
}

// 标签页切换
const handleTabChange = () => {
  // 不需要重新获取,因为已经获取了所有配置
}

// 新建
const handleCreate = () => {
  isEdit.value = false
  form.id = undefined
  form.config_key = ''
  form.config_value = ''
  form.config_type = activeTab.value === 'other' ? 'application_key' : activeTab.value
  form.description = ''
  form.is_encrypted = true
  form.is_active = true
  dialogVisible.value = true
}

// 编辑
const handleEdit = (config) => {
  isEdit.value = true
  form.id = config.id
  form.config_key = config.config_key
  form.config_value = config.config_value
  form.config_type = config.config_type
  form.description = config.description || ''
  form.is_encrypted = config.is_encrypted
  form.is_active = config.is_active
  dialogVisible.value = true
}

// 提交
const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    
    const data = {
      config_key: form.config_key,
      config_value: form.config_value,
      config_type: form.config_type,
      description: form.description,
      is_encrypted: form.is_encrypted,
      is_active: form.is_active
    }

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
      console.error('提交失败:', error)
      ElMessage.error(isEdit.value ? t('config.updateError') : t('config.createError'))
    }
  }
}

// 删除
const handleDelete = async (id) => {
  try {
    await api.config.deleteConfig(id)
    ElMessage.success(t('config.deleteSuccess'))
    fetchConfigs()
  } catch (error) {
    console.error('删除失败:', error)
    ElMessage.error(t('config.deleteError'))
  }
}

// 切换启用状态
const handleToggleActive = async (config) => {
  try {
    await api.config.updateConfig(config.id, {
      config_value: config.config_value,
      is_active: !config.is_active
    })
    ElMessage.success(t('config.updateSuccess'))
    fetchConfigs()
  } catch (error) {
    console.error('更新失败:', error)
    ElMessage.error(t('config.updateError'))
  }
}

// 生成爬虫Token
const handleGenerateCrawlerToken = () => {
  tokenForm.name = ''
  generateTokenDialogVisible.value = true
}

// 提交生成Token
const handleGenerateTokenSubmit = async () => {
  try {
    await tokenFormRef.value.validate()
    const response = await api.config.generateCrawlerToken({ name: tokenForm.name })
    generatedToken.value = response.token
    generateTokenDialogVisible.value = false
    tokenDisplayDialogVisible.value = true
    fetchConfigs()
  } catch (error) {
    console.error('生成Token失败:', error)
    ElMessage.error(t('config.generateTokenError'))
  }
}

// 复制Token
const handleCopyToken = async () => {
  try {
    await navigator.clipboard.writeText(generatedToken.value)
    ElMessage.success(t('config.copyTokenSuccess'))
  } catch (error) {
    console.error('复制失败:', error)
    ElMessage.error(t('config.copyTokenError'))
  }
}

// 初始化
onMounted(() => {
  fetchConfigs()
})
</script>

<style scoped lang="less">
.config-manage-page {
  padding: 0;
}

.crawler-token-section {
  .section-header {
    margin-bottom: 20px;
  }
}

:deep(.el-tabs) {
  .el-tabs__header {
    margin-bottom: 20px;
  }

  .el-tabs__item {
    font-weight: 500;
    padding: 0 24px;
  }
}
</style>

