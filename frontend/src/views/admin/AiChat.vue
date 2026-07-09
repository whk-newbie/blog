<template>
  <div class="ai-chat-page">
    <el-card class="chat-container">
      <template #header>
        <div class="chat-header">
          <span>AI 助手</span>
          <el-select v-model="selectedProviderId" placeholder="选择 AI 提供方" style="width: 200px">
            <el-option v-for="p in enabledProviders" :key="p.id" :label="`${p.name} (${p.model})`" :value="p.id" />
          </el-select>
        </div>
      </template>

      <div class="chat-messages" ref="messagesContainer">
        <el-empty v-if="messages.length === 0 && !streaming" description="开始对话" :image-size="80" />
        <div v-for="(msg, i) in messages" :key="i" :class="['message', msg.role]">
          <div class="message-role">{{ msg.role === 'user' ? '你' : 'AI' }}</div>
          <div class="message-content" v-html="renderContent(msg.content)"></div>
        </div>
        <div v-if="streaming" class="message assistant">
          <div class="message-role">AI</div>
          <div class="message-content" v-html="renderContent(streamContent)"></div>
          <span class="streaming-indicator">●</span>
        </div>
      </div>

      <div class="chat-input">
        <div class="quick-actions">
          <el-button size="small" @click="quickAction('请帮我润色以下内容：')">润色</el-button>
          <el-button size="small" @click="quickAction('请帮我扩写以下内容：')">扩写</el-button>
          <el-button size="small" @click="quickAction('请帮我缩写以下内容：')">缩写</el-button>
          <el-button size="small" @click="quickAction('请帮我总结以下内容：')">总结</el-button>
          <el-button size="small" @click="quickAction('请将以下内容翻译成英文：')">翻译</el-button>
        </div>
        <div class="input-row">
          <el-input
            v-model="inputText"
            type="textarea"
            :rows="3"
            placeholder="输入消息... (Enter 发送, Shift+Enter 换行)"
            @keydown.enter.exact.prevent="sendMessage"
          />
          <el-button type="primary" @click="sendMessage" :loading="streaming" :disabled="!selectedProviderId">
            发送
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick, watch } from 'vue'
import { ElMessage } from 'element-plus'
import aiApi from '@/api/ai'

const messages = ref([])
const inputText = ref('')
const selectedProviderId = ref(null)
const providers = ref([])
const streaming = ref(false)
const streamContent = ref('')
const messagesContainer = ref(null)
let abortController = null

const enabledProviders = computed(() =>
  providers.value.filter(p => p.is_enabled)
)

const fetchProviders = async () => {
  try {
    providers.value = await aiApi.listProviders()
    if (!selectedProviderId.value && enabledProviders.value.length > 0) {
      selectedProviderId.value = enabledProviders.value[0].id
    }
  } catch (e) {
    // silently fail
  }
}

const sendMessage = async () => {
  const text = inputText.value.trim()
  if (!text || !selectedProviderId.value || streaming.value) return

  inputText.value = ''
  messages.value.push({ role: 'user', content: text })

  const allMessages = [
    { role: 'system', content: 'You are a helpful assistant. Respond in the same language as the user.' },
    ...messages.value.map(m => ({ role: m.role, content: m.content })),
  ]

  streaming.value = true
  streamContent.value = ''

  try {
    abortController = new AbortController()
    const response = await aiApi.chat(selectedProviderId.value, allMessages, abortController.signal)

    const reader = response.body.getReader()
    const decoder = new TextDecoder()
    let buffer = ''

    while (true) {
      const { done, value } = await reader.read()
      if (done) break

      buffer += decoder.decode(value, { stream: true })
      const lines = buffer.split('\n')
      buffer = lines.pop() || ''

      for (const line of lines) {
        if (line.startsWith('data: ')) {
          const data = line.slice(6)
          if (data === '[DONE]') {
            break
          }
          try {
            const parsed = JSON.parse(data)
            const delta = parsed.choices?.[0]?.delta?.content || ''
            streamContent.value += delta
          } catch (e) {
            // ignore parse errors for partial chunks
          }
        }
      }
    }

    if (streamContent.value) {
      messages.value.push({ role: 'assistant', content: streamContent.value })
    }
  } catch (e) {
    if (e.name !== 'AbortError') {
      ElMessage.error('聊天请求失败')
    }
  } finally {
    streaming.value = false
    streamContent.value = ''
    abortController = null
  }

  await nextTick()
  scrollToBottom()
}

const quickAction = (prefix) => {
  inputText.value = prefix
}

const renderContent = (content) => {
  return content.replace(/\n/g, '<br>')
}

const scrollToBottom = () => {
  if (messagesContainer.value) {
    messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight
  }
}

watch(() => messages.value.length, () => {
  nextTick(() => scrollToBottom())
})

onMounted(() => {
  fetchProviders()
})
</script>

<style scoped lang="less">
.chat-container {
  height: calc(100vh - 180px);
  display: flex;
  flex-direction: column;

  :deep(.el-card__body) {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }
}

.chat-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 16px;
  background: var(--bg-color-secondary, #f5f5f5);
  border-radius: 8px;
  margin-bottom: 16px;
}

.message {
  margin-bottom: 16px;

  .message-role {
    font-weight: 600;
    font-size: 13px;
    margin-bottom: 4px;
    color: var(--text-secondary, #888);
  }

  .message-content {
    padding: 10px 14px;
    border-radius: 8px;
    line-height: 1.6;
    word-break: break-word;
  }

  &.user .message-content {
    background: var(--primary-color, #4078c0);
    color: white;
    margin-left: 40px;
  }

  &.assistant .message-content {
    background: white;
    border: 1px solid var(--border-color, #e8e8e8);
    margin-right: 40px;
  }
}

.streaming-indicator {
  color: var(--primary-color, #4078c0);
  animation: blink 1s infinite;
  margin-left: 8px;
}

@keyframes blink {
  50% { opacity: 0; }
}

.chat-input {
  .quick-actions {
    display: flex;
    gap: 8px;
    margin-bottom: 10px;
    flex-wrap: wrap;
  }

  .input-row {
    display: flex;
    gap: 10px;
    align-items: flex-end;

    :deep(.el-textarea) {
      flex: 1;
    }
  }
}
</style>
