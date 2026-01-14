<template>
  <el-card class="site-info-card" shadow="never">
    <template #header>
      <div class="card-header">
        <span>{{ t('config.siteInfo') }}</span>
        <el-button type="primary" @click="handleSave" :loading="saving">
          {{ t('common.save') }}
        </el-button>
      </div>
    </template>

    <el-form
      ref="formRef"
      :model="form"
      :rules="rules"
      label-width="120px"
    >
      <el-form-item :label="t('config.blogTitle')" prop="blogTitle">
        <el-input
          v-model="form.blogTitle"
          :placeholder="t('config.blogTitlePlaceholder')"
          maxlength="100"
          show-word-limit
        />
      </el-form-item>

      <el-divider />

      <el-form-item :label="t('config.icpInfo')">
        <el-switch
          v-model="form.enableIcp"
          :active-text="t('config.enableIcp')"
          @change="handleIcpToggle"
        />
      </el-form-item>

      <template v-if="form.enableIcp">
        <el-form-item :label="t('config.icpText')" prop="icpText">
          <el-input
            v-model="form.icpText"
            :placeholder="t('config.icpTextPlaceholder')"
            maxlength="50"
            show-word-limit
          />
          <div class="form-tip">{{ t('config.icpTextTip') }}</div>
        </el-form-item>

        <el-form-item :label="t('config.icpUrl')" prop="icpUrl">
          <el-input
            v-model="form.icpUrl"
            :placeholder="t('config.icpUrlPlaceholder')"
            type="url"
          />
          <div class="form-tip">{{ t('config.icpUrlTip') }}</div>
        </el-form-item>
      </template>
    </el-form>
  </el-card>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import api from '@/api'

const { t } = useI18n()
const formRef = ref(null)
const saving = ref(false)
const siteInfoConfig = ref(null)

const form = reactive({
  blogTitle: '我的博客',
  enableIcp: false,
  icpText: '',
  icpUrl: ''
})

const rules = computed(() => ({
  blogTitle: [
    { required: true, message: t('config.blogTitlePlaceholder'), trigger: 'blur' }
  ],
  icpText: [
    {
      validator: (rule, value, callback) => {
        if (form.enableIcp && !value) {
          callback(new Error(t('config.icpTextPlaceholder')))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  icpUrl: [
    {
      validator: (rule, value, callback) => {
        if (form.enableIcp && value && !isValidUrl(value)) {
          callback(new Error(t('validation.url')))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}))

// URL验证
const isValidUrl = (url) => {
  try {
    new URL(url)
    return true
  } catch {
    return false
  }
}

// 加载站点信息配置
const loadSiteInfo = async () => {
  try {
    const response = await api.config.getConfigs({ config_type: 'site_info' })
    const configs = Array.isArray(response) ? response : []
    
    if (configs.length > 0) {
      siteInfoConfig.value = configs[0]
      try {
        // 站点信息配置设置为不加密，所以列表返回的值就是完整的JSON字符串
        // 如果配置是加密的，需要通过GetConfigByID获取完整值
        let configValue = configs[0].config_value
        
        // 如果配置是加密的，需要通过详情接口获取完整值
        if (configs[0].is_encrypted) {
          const detailResponse = await api.config.getConfigById(configs[0].id)
          if (detailResponse && detailResponse.config_value) {
            configValue = detailResponse.config_value
          }
        }
        
        const configData = JSON.parse(configValue)
        form.blogTitle = configData.blogTitle || '我的博客'
        
        if (configData.icpInfo) {
          form.enableIcp = true
          form.icpText = configData.icpInfo.text || ''
          form.icpUrl = configData.icpInfo.url || ''
        } else {
          form.enableIcp = false
          form.icpText = ''
          form.icpUrl = ''
        }
      } catch (error) {
        console.error('解析站点配置失败:', error)
      }
    }
  } catch (error) {
    console.error('加载站点信息失败:', error)
  }
}

// ICP开关切换
const handleIcpToggle = (value) => {
  if (!value) {
    form.icpText = ''
    form.icpUrl = ''
  }
}

// 保存配置
const handleSave = async () => {
  try {
    await formRef.value.validate()
    
    saving.value = true
    
    // 构建配置数据
    const configData = {
      blogTitle: form.blogTitle
    }
    
    if (form.enableIcp && form.icpText) {
      configData.icpInfo = {
        text: form.icpText
      }
      if (form.icpUrl) {
        configData.icpInfo.url = form.icpUrl
      }
    } else {
      configData.icpInfo = null
    }
    
    const configValue = JSON.stringify(configData)
    
    if (siteInfoConfig.value) {
      // 更新现有配置
      await api.config.updateConfig(siteInfoConfig.value.id, {
        config_key: 'site_info',
        config_value: configValue,
        config_type: 'site_info',
        description: t('config.siteInfoDescription'),
        is_encrypted: false,
        is_active: true
      })
    } else {
      // 创建新配置
      await api.config.createConfig({
        config_key: 'site_info',
        config_value: configValue,
        config_type: 'site_info',
        description: t('config.siteInfoDescription'),
        is_encrypted: false,
        is_active: true
      })
    }
    
    ElMessage.success(t('config.updateSuccess'))
    await loadSiteInfo()
  } catch (error) {
    if (error instanceof Error) {
      console.error('保存失败:', error)
      ElMessage.error(t('config.updateError'))
    }
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadSiteInfo()
})
</script>

<style scoped lang="less">
.site-info-card {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}

.form-tip {
  font-size: 12px;
  color: var(--text-color-secondary);
  margin-top: 4px;
  line-height: 1.5;
}

:deep(.el-divider) {
  margin: 20px 0;
}
</style>
