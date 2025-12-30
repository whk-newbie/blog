<template>
  <div class="tools-page">
    <div class="tools-container">
      <div class="page-header">
        <h1 class="page-title">
          <el-icon><Tools /></el-icon>
          {{ t('nav.tools') }}
        </h1>
        <p class="page-description">{{ t('tools.description') }}</p>
      </div>

      <el-tabs v-model="activeTab" class="tools-tabs">
        <!-- JSON格式化工具 -->
        <el-tab-pane :label="t('tools.jsonFormatter')" name="json">
          <div class="tool-panel">
            <el-input
              v-model="jsonInput"
              type="textarea"
              :rows="10"
              :placeholder="t('tools.jsonPlaceholder')"
              class="tool-input"
            />
            <div class="tool-actions">
              <el-button type="primary" @click="formatJSON">
                <el-icon><DocumentCopy /></el-icon>
                {{ t('tools.format') }}
              </el-button>
              <el-button @click="minifyJSON">
                <el-icon><Minus /></el-icon>
                {{ t('tools.minify') }}
              </el-button>
              <el-button @click="validateJSON">
                <el-icon><CircleCheck /></el-icon>
                {{ t('tools.validate') }}
              </el-button>
              <el-button @click="clearJSON">
                <el-icon><Delete /></el-icon>
                {{ t('common.reset') }}
              </el-button>
            </div>
            <el-input
              v-model="jsonOutput"
              type="textarea"
              :rows="10"
              :placeholder="t('tools.result')"
              readonly
              class="tool-output"
            />
          </div>
        </el-tab-pane>

        <!-- Header格式化工具 -->
        <el-tab-pane :label="t('tools.headerFormatter')" name="header">
          <div class="tool-panel">
            <el-input
              v-model="headerInput"
              type="textarea"
              :rows="10"
              :placeholder="t('tools.headerPlaceholder')"
              class="tool-input"
            />
            <div class="tool-actions">
              <el-button type="primary" @click="formatHeader">
                <el-icon><DocumentCopy /></el-icon>
                {{ t('tools.format') }}
              </el-button>
              <el-button @click="clearHeader">
                <el-icon><Delete /></el-icon>
                {{ t('common.reset') }}
              </el-button>
            </div>
            <el-input
              v-model="headerOutput"
              type="textarea"
              :rows="10"
              :placeholder="t('tools.result')"
              readonly
              class="tool-output"
            />
          </div>
        </el-tab-pane>

        <!-- Cookie格式化工具 -->
        <el-tab-pane :label="t('tools.cookieFormatter')" name="cookie">
          <div class="tool-panel">
            <el-input
              v-model="cookieInput"
              type="textarea"
              :rows="10"
              :placeholder="t('tools.cookiePlaceholder')"
              class="tool-input"
            />
            <div class="tool-actions">
              <el-button type="primary" @click="formatCookie">
                <el-icon><DocumentCopy /></el-icon>
                {{ t('tools.format') }}
              </el-button>
              <el-button @click="clearCookie">
                <el-icon><Delete /></el-icon>
                {{ t('common.reset') }}
              </el-button>
            </div>
            <el-input
              v-model="cookieOutput"
              type="textarea"
              :rows="10"
              :placeholder="t('tools.result')"
              readonly
              class="tool-output"
            />
          </div>
        </el-tab-pane>

        <!-- Dict格式化工具 -->
        <el-tab-pane :label="t('tools.dictFormatter')" name="dict">
          <div class="tool-panel">
            <el-input
              v-model="dictInput"
              type="textarea"
              :rows="10"
              :placeholder="t('tools.dictPlaceholder')"
              class="tool-input"
            />
            <div class="tool-actions">
              <el-button type="primary" @click="formatDict">
                <el-icon><DocumentCopy /></el-icon>
                {{ t('tools.format') }}
              </el-button>
              <el-button @click="clearDict">
                <el-icon><Delete /></el-icon>
                {{ t('common.reset') }}
              </el-button>
            </div>
            <el-input
              v-model="dictOutput"
              type="textarea"
              :rows="10"
              :placeholder="t('tools.result')"
              readonly
              class="tool-output"
            />
          </div>
        </el-tab-pane>

        <!-- cURL转Python请求工具 -->
        <el-tab-pane :label="t('tools.curlToPython')" name="curl2py">
          <div class="tool-panel">
            <el-input
              v-model="curlInput"
              type="textarea"
              :rows="10"
              :placeholder="t('tools.curlPlaceholder')"
              class="tool-input"
            />
            <div class="tool-actions">
              <el-button type="primary" @click="convertCurlToPython">
                <el-icon><Right /></el-icon>
                {{ t('tools.convert') }}
              </el-button>
              <el-button @click="clearCurl">
                <el-icon><Delete /></el-icon>
                {{ t('common.reset') }}
              </el-button>
            </div>
            <el-input
              v-model="pythonOutput"
              type="textarea"
              :rows="10"
              :placeholder="t('tools.result')"
              readonly
              class="tool-output"
            />
          </div>
        </el-tab-pane>

        <!-- Python请求转cURL工具 -->
        <el-tab-pane :label="t('tools.pythonToCurl')" name="py2curl">
          <div class="tool-panel">
            <el-input
              v-model="pythonInput"
              type="textarea"
              :rows="10"
              :placeholder="t('tools.pythonPlaceholder')"
              class="tool-input"
            />
            <div class="tool-actions">
              <el-button type="primary" @click="convertPythonToCurl">
                <el-icon><Right /></el-icon>
                {{ t('tools.convert') }}
              </el-button>
              <el-button @click="clearPython">
                <el-icon><Delete /></el-icon>
                {{ t('common.reset') }}
              </el-button>
            </div>
            <el-input
              v-model="curlOutput"
              type="textarea"
              :rows="10"
              :placeholder="t('tools.result')"
              readonly
              class="tool-output"
            />
          </div>
        </el-tab-pane>

        <!-- URL格式化工具 -->
        <el-tab-pane :label="t('tools.urlFormatter')" name="url">
          <div class="tool-panel">
            <el-input
              v-model="urlInput"
              type="textarea"
              :rows="5"
              :placeholder="t('tools.urlPlaceholder')"
              class="tool-input"
            />
            <div class="tool-actions">
              <el-button type="primary" @click="formatURL">
                <el-icon><DocumentCopy /></el-icon>
                {{ t('tools.format') }}
              </el-button>
              <el-button @click="clearURL">
                <el-icon><Delete /></el-icon>
                {{ t('common.reset') }}
              </el-button>
            </div>
            <el-input
              v-model="urlOutput"
              type="textarea"
              :rows="5"
              :placeholder="t('tools.result')"
              readonly
              class="tool-output"
            />
          </div>
        </el-tab-pane>

        <!-- 加密解密工具集 -->
        <el-tab-pane :label="t('tools.encryptDecrypt')" name="crypto">
          <div class="tool-panel">
            <el-radio-group v-model="cryptoType" class="crypto-type-group">
              <el-radio-button label="url">{{ t('tools.urlEncode') }}</el-radio-button>
              <el-radio-button label="base64">{{ t('tools.base64') }}</el-radio-button>
              <el-radio-button label="md5">{{ t('tools.md5') }}</el-radio-button>
            </el-radio-group>
            <el-input
              v-model="cryptoInput"
              type="textarea"
              :rows="8"
              :placeholder="t('tools.cryptoPlaceholder')"
              class="tool-input"
            />
            <div class="tool-actions">
              <el-button type="primary" @click="handleCrypto">
                <el-icon><Lock /></el-icon>
                {{ t('tools.encrypt') }}
              </el-button>
              <el-button v-if="cryptoType === 'url' || cryptoType === 'base64'" @click="handleDecrypt">
                <el-icon><Unlock /></el-icon>
                {{ t('tools.decrypt') }}
              </el-button>
              <el-button @click="clearCrypto">
                <el-icon><Delete /></el-icon>
                {{ t('common.reset') }}
              </el-button>
            </div>
            <el-input
              v-model="cryptoOutput"
              type="textarea"
              :rows="8"
              :placeholder="t('tools.result')"
              readonly
              class="tool-output"
            />
          </div>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import {
  Tools,
  DocumentCopy,
  Minus,
  CircleCheck,
  Delete,
  Right,
  Lock,
  Unlock
} from '@element-plus/icons-vue'
import CryptoJS from 'crypto-js'

const { t } = useI18n()

const activeTab = ref('json')

// JSON格式化工具
const jsonInput = ref('')
const jsonOutput = ref('')

const formatJSON = () => {
  try {
    const obj = JSON.parse(jsonInput.value)
    jsonOutput.value = JSON.stringify(obj, null, 2)
    ElMessage.success(t('tools.formatSuccess'))
  } catch (e) {
    ElMessage.error(t('tools.invalidJSON'))
  }
}

const minifyJSON = () => {
  try {
    const obj = JSON.parse(jsonInput.value)
    jsonOutput.value = JSON.stringify(obj)
    ElMessage.success(t('tools.minifySuccess'))
  } catch (e) {
    ElMessage.error(t('tools.invalidJSON'))
  }
}

const validateJSON = () => {
  try {
    JSON.parse(jsonInput.value)
    ElMessage.success(t('tools.validJSON'))
  } catch (e) {
    ElMessage.error(t('tools.invalidJSON'))
  }
}

const clearJSON = () => {
  jsonInput.value = ''
  jsonOutput.value = ''
}

// Header格式化工具
const headerInput = ref('')
const headerOutput = ref('')

const formatHeader = () => {
  const lines = headerInput.value.split('\n').filter(line => line.trim())
  const formatted = lines.map(line => {
    const colonIndex = line.indexOf(':')
    if (colonIndex === -1) return line
    const key = line.substring(0, colonIndex).trim()
    const value = line.substring(colonIndex + 1).trim()
    return `${key}: ${value}`
  }).join('\n')
  headerOutput.value = formatted
  ElMessage.success(t('tools.formatSuccess'))
}

const clearHeader = () => {
  headerInput.value = ''
  headerOutput.value = ''
}

// Cookie格式化工具
const cookieInput = ref('')
const cookieOutput = ref('')

const formatCookie = () => {
  const cookies = cookieInput.value.split(';').map(c => c.trim()).filter(c => c)
  const formatted = cookies.map(cookie => {
    const parts = cookie.split('=')
    if (parts.length >= 2) {
      const key = parts[0].trim()
      const value = parts.slice(1).join('=').trim()
      return `${key}: ${value}`
    }
    return cookie
  }).join('\n')
  cookieOutput.value = formatted
  ElMessage.success(t('tools.formatSuccess'))
}

const clearCookie = () => {
  cookieInput.value = ''
  cookieOutput.value = ''
}

// Dict格式化工具
const dictInput = ref('')
const dictOutput = ref('')

const formatDict = () => {
  try {
    const lines = dictInput.value.split('\n').filter(line => line.trim())
    const dict = {}
    lines.forEach(line => {
      const colonIndex = line.indexOf(':')
      if (colonIndex !== -1) {
        const key = line.substring(0, colonIndex).trim()
        const value = line.substring(colonIndex + 1).trim()
        dict[key] = value
      }
    })
    dictOutput.value = JSON.stringify(dict, null, 2)
    ElMessage.success(t('tools.formatSuccess'))
  } catch (e) {
    ElMessage.error(t('tools.formatError'))
  }
}

const clearDict = () => {
  dictInput.value = ''
  dictOutput.value = ''
}

// cURL转Python
const curlInput = ref('')
const pythonOutput = ref('')

const convertCurlToPython = () => {
  try {
    const curl = curlInput.value.trim()
    let pythonCode = 'import requests\n\n'
    
    // 提取URL
    const urlMatch = curl.match(/curl\s+['"]?([^'"]+)['"]?/)
    if (!urlMatch) {
      throw new Error('Invalid cURL command')
    }
    const url = urlMatch[1]
    
    // 提取方法
    const methodMatch = curl.match(/-X\s+(\w+)/i)
    const method = methodMatch ? methodMatch[1].toUpperCase() : 'GET'
    
    // 提取headers
    const headerMatches = curl.matchAll(/-H\s+['"]([^'"]+)['"]/g)
    const headers = {}
    for (const match of headerMatches) {
      const [key, value] = match[1].split(':').map(s => s.trim())
      headers[key] = value
    }
    
    // 提取data
    const dataMatch = curl.match(/-d\s+['"]([^'"]+)['"]/)
    const data = dataMatch ? dataMatch[1] : null
    
    pythonCode += `url = "${url}"\n`
    pythonCode += `headers = ${JSON.stringify(headers, null, 2)}\n`
    if (data) {
      pythonCode += `data = "${data}"\n`
    }
    pythonCode += `\nresponse = requests.${method.toLowerCase()}(${data ? 'url, headers=headers, data=data' : 'url, headers=headers'})\n`
    pythonCode += `print(response.text)`
    
    pythonOutput.value = pythonCode
    ElMessage.success(t('tools.convertSuccess'))
  } catch (e) {
    ElMessage.error(t('tools.convertError'))
  }
}

const clearCurl = () => {
  curlInput.value = ''
  pythonOutput.value = ''
}

// Python转cURL
const pythonInput = ref('')
const curlOutput = ref('')

const convertPythonToCurl = () => {
  try {
    const code = pythonInput.value.trim()
    let curlCmd = 'curl'
    
    // 提取URL
    const urlMatch = code.match(/url\s*=\s*['"]([^'"]+)['"]/)
    if (!urlMatch) {
      throw new Error('URL not found')
    }
    const url = urlMatch[1]
    
    // 提取方法
    const methodMatch = code.match(/requests\.(\w+)\(/)
    const method = methodMatch ? methodMatch[1].toUpperCase() : 'GET'
    
    // 提取headers
    const headersMatch = code.match(/headers\s*=\s*({[^}]+})/)
    let headers = {}
    if (headersMatch) {
      try {
        // 使用 JSON.parse 替代 eval，更安全
        const headersStr = headersMatch[1].replace(/'/g, '"')
        headers = JSON.parse(headersStr)
      } catch (e) {
        // 如果解析失败，尝试手动解析简单的字典格式
        const headerLines = headersMatch[1].match(/'([^']+)':\s*'([^']+)'/g)
        if (headerLines) {
          headerLines.forEach(line => {
            const match = line.match(/'([^']+)':\s*'([^']+)'/)
            if (match) {
              headers[match[1]] = match[2]
            }
          })
        }
      }
    }
    
    // 提取data
    const dataMatch = code.match(/data\s*=\s*['"]([^'"]+)['"]/)
    const data = dataMatch ? dataMatch[1] : null
    
    if (method !== 'GET') {
      curlCmd += ` -X ${method}`
    }
    
    curlCmd += ` "${url}"`
    
    Object.entries(headers).forEach(([key, value]) => {
      curlCmd += ` -H "${key}: ${value}"`
    })
    
    if (data) {
      curlCmd += ` -d "${data}"`
    }
    
    curlOutput.value = curlCmd
    ElMessage.success(t('tools.convertSuccess'))
  } catch (e) {
    ElMessage.error(t('tools.convertError'))
  }
}

const clearPython = () => {
  pythonInput.value = ''
  curlOutput.value = ''
}

// URL格式化工具
const urlInput = ref('')
const urlOutput = ref('')

const formatURL = () => {
  try {
    const url = new URL(urlInput.value)
    const formatted = {
      protocol: url.protocol,
      host: url.host,
      hostname: url.hostname,
      port: url.port || '(default)',
      pathname: url.pathname,
      search: url.search || '(none)',
      hash: url.hash || '(none)',
      params: Object.fromEntries(url.searchParams)
    }
    urlOutput.value = JSON.stringify(formatted, null, 2)
    ElMessage.success(t('tools.formatSuccess'))
  } catch (e) {
    ElMessage.error(t('tools.invalidURL'))
  }
}

const clearURL = () => {
  urlInput.value = ''
  urlOutput.value = ''
}

// 加密解密工具
const cryptoType = ref('url')
const cryptoInput = ref('')
const cryptoOutput = ref('')

const handleCrypto = () => {
  try {
    switch (cryptoType.value) {
      case 'url':
        cryptoOutput.value = encodeURIComponent(cryptoInput.value)
        break
      case 'base64':
        cryptoOutput.value = btoa(unescape(encodeURIComponent(cryptoInput.value)))
        break
      case 'md5':
        cryptoOutput.value = CryptoJS.MD5(cryptoInput.value).toString()
        break
    }
    ElMessage.success(t('tools.encryptSuccess'))
  } catch (e) {
    ElMessage.error(t('tools.encryptError'))
  }
}

const handleDecrypt = () => {
  try {
    switch (cryptoType.value) {
      case 'url':
        cryptoOutput.value = decodeURIComponent(cryptoInput.value)
        break
      case 'base64':
        cryptoOutput.value = decodeURIComponent(escape(atob(cryptoInput.value)))
        break
    }
    ElMessage.success(t('tools.decryptSuccess'))
  } catch (e) {
    ElMessage.error(t('tools.decryptError'))
  }
}

const clearCrypto = () => {
  cryptoInput.value = ''
  cryptoOutput.value = ''
}
</script>

<style scoped lang="less">
.tools-page {
  min-height: calc(100vh - 200px);
  padding: 40px 20px;
  background: var(--bg-color);
}

.tools-container {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  text-align: center;
  margin-bottom: 40px;

  .page-title {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 12px;
    font-size: 32px;
    font-weight: 600;
    color: var(--text-color);
    margin: 0 0 12px;

    .el-icon {
      font-size: 36px;
      color: var(--primary-color);
    }
  }

  .page-description {
    color: var(--text-color-secondary);
    font-size: 16px;
    margin: 0;
  }
}

.tools-tabs {
  background: var(--card-bg);
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);

  :deep(.el-tabs__header) {
    margin-bottom: 24px;
  }

  :deep(.el-tabs__item) {
    font-size: 16px;
    font-weight: 500;
  }
}

.tool-panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.tool-input,
.tool-output {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', 'Consolas', monospace;
  font-size: 14px;
}

.tool-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.crypto-type-group {
  margin-bottom: 16px;
  width: 100%;
  display: flex;
  justify-content: center;
}

@media (max-width: 768px) {
  .tools-page {
    padding: 20px 12px;
  }

  .page-header .page-title {
    font-size: 24px;
  }

  .tools-tabs {
    padding: 16px;
  }

  .tool-actions {
    flex-direction: column;
    
    .el-button {
      width: 100%;
    }
  }
}
</style>

