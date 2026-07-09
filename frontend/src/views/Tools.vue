<template>
  <div class="tools-page">
    <div class="tools-layout">
      <!-- Left: tool list -->
      <div class="tools-nav">
        <div v-for="tool in tools" :key="tool.name" :class="['tool-nav-item', { active: activeTool === tool.name }]" @click="activeTool = tool.name">
          <el-icon class="tool-nav-icon"><component :is="tool.icon" /></el-icon>
          <span>{{ tool.label }}</span>
        </div>
      </div>
      <!-- Right: active tool panel -->
      <div class="tools-main">
        <el-card shadow="never">
          <template #header><span class="tool-title">{{ activeToolData?.label }}</span></template>
          <div v-if="activeToolData" class="tool-panel">
            <el-input v-model="activeToolData.input" type="textarea" :rows="8" :placeholder="activeToolData.placeholder" />
            <div class="tool-actions">
              <el-button v-for="btn in activeToolData.buttons" :key="btn.text" :type="btn.type||'default'" size="small" @click="btn.action">
                <el-icon v-if="btn.icon"><component :is="btn.icon" /></el-icon>
                {{ btn.text }}
              </el-button>
            </div>
            <el-input v-if="activeToolData.showOutput" v-model="activeToolData.output" type="textarea" :rows="8" readonly placeholder="结果..." />
          </div>
        </el-card>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Tools, DocumentCopy, Minus, CircleCheck, Delete, Right } from '@element-plus/icons-vue'

// JSON tool
const jsonInput = ref('')
const jsonOutput = ref('')
const formatJSON = () => { try { jsonOutput.value = JSON.stringify(JSON.parse(jsonInput.value), null, 2) } catch { ElMessage.error('JSON 格式错误') } }
const minifyJSON = () => { try { jsonOutput.value = JSON.stringify(JSON.parse(jsonInput.value)) } catch { ElMessage.error('JSON 格式错误') } }
const validateJSON = () => { try { JSON.parse(jsonInput.value); ElMessage.success('JSON 格式正确') } catch { ElMessage.error('JSON 格式错误') } }

// Header tool
const headerInput = ref('')
const headerOutput = ref('')
const formatHeader = () => {
  const lines = headerInput.value.trim().split('\n').filter(Boolean)
  const result = {}
  for (const line of lines) {
    const idx = line.indexOf(':')
    if (idx > 0) result[line.substring(0, idx).trim()] = line.substring(idx + 1).trim()
  }
  headerOutput.value = JSON.stringify(result, null, 2)
}

// Cookie tool
const cookieInput = ref('')
const cookieOutput = ref('')
const formatCookie = () => {
  const pairs = cookieInput.value.split(';')
  const result = {}
  for (const p of pairs) {
    const eq = p.indexOf('=')
    if (eq > 0) result[p.substring(0, eq).trim()] = decodeURIComponent(p.substring(eq + 1).trim())
  }
  cookieOutput.value = JSON.stringify(result, null, 2)
}

// Dict tool
const dictInput = ref('')
const dictOutput = ref('')
const formatDict = () => {
  try {
    const lines = dictInput.value.trim().split('\n').filter(Boolean)
    const result = {}
    for (const line of lines) {
      const obj = JSON.parse(line)
      Object.assign(result, obj)
    }
    dictOutput.value = JSON.stringify(result, null, 2)
  } catch { ElMessage.error('格式错误') }
}

// cURL to Python
const curlInput = ref('')
const pythonOutput = ref('')
const convertCurlToPython = () => {
  let cmd = curlInput.value.trim().replace(/\\\n/g, ' ')
  const urlMatch = cmd.match(/curl\s+(?:-[^\s]+\s+)*['"]?([^'"]+)['"]?/i)
  const url = urlMatch?.[1] || 'https://example.com'
  const methodMatch = cmd.match(/-X\s+(\w+)/i)
  const method = methodMatch?.[1]?.toUpperCase() || 'GET'
  const headers = []
  const hRe = /-H\s+['"]([^'"]+)['"]/gi; let m
  while ((m = hRe.exec(cmd)) !== null) headers.push(m[1])
  const dataMatch = cmd.match(/--data(?:-raw|-binary)?\s+['"]([^'"]+)['"]/i)
  let code = `import requests\n\nresponse = requests.${method.toLowerCase()}("${url}"`
  if (headers.length) code += `,\n    headers={${headers.map(h => { const [k,...v] = h.split(':'); return `"${k.trim()}": "${v.join(':').trim()}"` }).join(', ')}}`
  if (dataMatch) code += `,\n    data='''${dataMatch[1]}'''`
  code += '\n)\nprint(response.text)'
  pythonOutput.value = code
}

// Python to cURL
const pythonInput = ref('')
const curlOutput = ref('')
const convertPythonToCurl = () => {
  const code = pythonInput.value
  const urlMatch = code.match(/\.(get|post|put|delete)\(['"]([^'"]+)['"]/i)
  const method = urlMatch?.[1]?.toUpperCase() || 'GET'
  const url = urlMatch?.[2] || 'https://example.com'
  let curl = `curl -X ${method} "${url}"`
  const hMatch = code.match(/headers\s*=\s*\{([^}]+)\}/s)
  if (hMatch) {
    const pairs = hMatch[1].match(/"([^"]+)":\s*"([^"]+)"/g)
    if (pairs) for (const p of pairs) { const pm = p.match(/"([^"]+)":\s*"([^"]+)"/); curl += ` \\\n  -H "${pm[1]}: ${pm[2]}"` }
  }
  const dMatch = code.match(/data\s*=\s*'''([^']+)'''/s) || code.match(/data\s*=\s*"([^"]+)"/s) || code.match(/json\s*=\s*(\{[^}]+\})/s)
  if (dMatch) curl += ` \\\n  --data '${dMatch[1].trim()}'`
  curlOutput.value = curl
}

const clear = (input, output) => { input.value = ''; output.value = '' }

const activeTool = ref('json')
const tools = [
  { name: 'json', label: 'JSON 格式化', icon: DocumentCopy, input: jsonInput, placeholder: '粘贴 JSON...', showOutput: true, get output() { return jsonOutput },
    buttons: [{ text: '格式化', type: 'primary', icon: DocumentCopy, action: formatJSON }, { text: '压缩', icon: Minus, action: minifyJSON }, { text: '验证', icon: CircleCheck, action: validateJSON }, { text: '清空', icon: Delete, action: () => clear(jsonInput, jsonOutput) }] },
  { name: 'header', label: 'Header 格式化', icon: Right, input: headerInput, placeholder: '粘贴 HTTP Headers...\nContent-Type: application/json\nAuthorization: Bearer xxx', showOutput: true, get output() { return headerOutput },
    buttons: [{ text: '格式化', type: 'primary', icon: DocumentCopy, action: formatHeader }, { text: '清空', icon: Delete, action: () => clear(headerInput, headerOutput) }] },
  { name: 'cookie', label: 'Cookie 格式化', icon: Right, input: cookieInput, placeholder: '粘贴 Cookie 字符串...', showOutput: true, get output() { return cookieOutput },
    buttons: [{ text: '格式化', type: 'primary', icon: DocumentCopy, action: formatCookie }, { text: '清空', icon: Delete, action: () => clear(cookieInput, cookieOutput) }] },
  { name: 'dict', label: 'Dict 格式化', icon: Right, input: dictInput, placeholder: '每行一个 JSON 对象...\n{"key1":"value1"}\n{"key2":"value2"}', showOutput: true, get output() { return dictOutput },
    buttons: [{ text: '格式化', type: 'primary', icon: DocumentCopy, action: formatDict }, { text: '清空', icon: Delete, action: () => clear(dictInput, dictOutput) }] },
  { name: 'curl2py', label: 'cURL → Python', icon: Right, input: curlInput, placeholder: '粘贴 cURL 命令...', showOutput: true, get output() { return pythonOutput },
    buttons: [{ text: '转换', type: 'primary', icon: Right, action: convertCurlToPython }, { text: '清空', icon: Delete, action: () => clear(curlInput, pythonOutput) }] },
  { name: 'py2curl', label: 'Python → cURL', icon: Right, input: pythonInput, placeholder: '粘贴 Python requests 代码...', showOutput: true, get output() { return curlOutput },
    buttons: [{ text: '转换', type: 'primary', icon: Right, action: convertPythonToCurl }, { text: '清空', icon: Delete, action: () => clear(pythonInput, curlOutput) }] },
]

const activeToolData = computed(() => tools.find(t => t.name === activeTool.value))
</script>

<style scoped lang="less">
.tools-page { padding: var(--spacing-md) 0; }
.tools-layout { display: flex; gap: var(--spacing-lg); align-items: flex-start; }
.tools-nav { width: 180px; flex-shrink: 0; background: var(--card-bg); border-radius: var(--border-radius-sm); border: 1px solid var(--border-color); overflow: hidden; }
.tool-nav-item { display: flex; align-items: center; gap: 8px; padding: 12px 16px; font-size: 14px; color: var(--text-secondary); cursor: pointer; border-left: 2px solid transparent; transition: all var(--transition-fast);
  &:hover { background: var(--hover-bg); color: var(--text-heading); }
  &.active { color: var(--primary-color); background: var(--active-bg); border-left-color: var(--primary-color); font-weight: 500; }
}
.tool-nav-icon { font-size: 16px; }
.tools-main { flex: 1; min-width: 0; }
.tool-title { font-size: 16px; font-weight: 600; color: var(--text-heading); }
.tool-panel { display: flex; flex-direction: column; gap: 12px; }
.tool-actions { display: flex; gap: 8px; flex-wrap: wrap; }

@media (max-width: 768px) {
  .tools-layout { flex-direction: column; }
  .tools-nav { width: 100%; display: flex; flex-wrap: wrap; gap: 4px; padding: 8px; }
  .tool-nav-item { padding: 8px 12px; font-size: 13px; border-left: none; border-radius: var(--border-radius-sm); }
}
</style>
