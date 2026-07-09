import axios from 'axios'
import { ElMessage } from 'element-plus'
import i18n from '@/locales'
import performanceMonitor from '@/utils/performance'
import CryptoJS from 'crypto-js'
import { generateAESKey, aesEncrypt, aesDecrypt, rsaEncrypt } from '@/utils/crypto'

// 创建axios实例
const http = axios.create({
  baseURL: '/api/v1', // 使用相对路径，通过 Vite 代理转发
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Encryption session state — persisted in sessionStorage to survive page reloads
let sessionId = null
let aesKey = null
let aesKeyB64 = null
let negotiating = false
let negotiatePromise = null

const SESSION_STORE_KEY = 'enc_session'

// Restore session from sessionStorage on module load
function restoreSession() {
  try {
    const stored = sessionStorage.getItem(SESSION_STORE_KEY)
    if (stored) {
      const { sid, keyB64 } = JSON.parse(stored)
      if (sid && keyB64) {
        sessionId = sid
        aesKeyB64 = keyB64
        aesKey = CryptoJS.enc.Base64.parse(keyB64)
        return true
      }
    }
  } catch (e) { /* ignore */ }
  return false
}

function saveSession() {
  try {
    sessionStorage.setItem(SESSION_STORE_KEY, JSON.stringify({
      sid: sessionId,
      keyB64: aesKeyB64,
    }))
  } catch (e) { /* ignore */ }
}

function clearSession() {
  sessionId = null
  aesKey = null
  aesKeyB64 = null
  negotiatePromise = null
  try { sessionStorage.removeItem(SESSION_STORE_KEY) } catch (e) { /* ignore */ }
}

// Negotiate encryption key with server
async function negotiateKey() {
  if (aesKey && sessionId) return

  if (negotiating && negotiatePromise) {
    return negotiatePromise
  }

  negotiating = true
  negotiatePromise = (async () => {
    try {
      // 1. Get RSA public key
      const pubKeyRes = await axios.get('/api/v1/public-key')
      const { public_key, session_id } = pubKeyRes.data.data
      sessionId = session_id

      // 2. Generate AES key
      const rawKey = generateAESKey()
      aesKeyB64 = CryptoJS.enc.Base64.stringify(rawKey)
      aesKey = rawKey

      // 3. Encrypt AES key with RSA public key
      const encryptedKey = await rsaEncrypt(public_key, aesKeyB64)

      // 4. Send encrypted key to server
      await axios.post('/api/v1/session/key', {
        encrypted_key: encryptedKey,
        session_id: sessionId,
      })

      // 5. Save to sessionStorage for persistence across reloads
      saveSession()
    } catch (error) {
      clearSession()
      throw error
    } finally {
      negotiating = false
    }
  })()

  return negotiatePromise
}

// Try to restore session on module load — avoids re-negotiation on page reload
restoreSession()

// 请求拦截器
http.interceptors.request.use(
  async (config) => {
    // 记录请求开始时间
    config.metadata = {
      startTime: performance.now(),
    }

    // Skip encryption for whitelisted paths
    const whitelist = ['public-key', 'session/key', 'fingerprint', 'visit']
    const isWhitelisted = whitelist.some(p => config.url?.includes(p))

    if (!isWhitelisted) {
      // Ensure key is negotiated
      if (!aesKey) {
        await negotiateKey()
      }
      config.headers['X-Session-Id'] = sessionId

      // Encrypt request body
      if (config.data && typeof config.data === 'object') {
        const jsonStr = JSON.stringify(config.data)
        const encrypted = aesEncrypt(jsonStr, aesKey)
        config.data = encrypted
        config.headers['Content-Type'] = 'text/plain'
      }
    }

    // 添加token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    // 添加请求ID
    config.headers['X-Request-ID'] = generateRequestId()

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
http.interceptors.response.use(
  async (response) => {
    // 记录API请求性能
    if (response.config.metadata) {
      const endTime = performance.now()
      const duration = endTime - response.config.metadata.startTime
      performanceMonitor.recordApiRequest(
        response.config.url || '',
        response.config.method || 'get',
        response.config.metadata.startTime,
        endTime,
        response.status
      )
    }

    // 204 No Content - 删除成功，没有响应体
    if (response.status === 204) {
      return null
    }

    // Decrypt response if encrypted
    const whitelist = ['public-key', 'session/key', 'fingerprint', 'visit']
    const isWhitelisted = whitelist.some(p => response.config.url?.includes(p))

    if (!isWhitelisted && aesKey && typeof response.data === 'string') {
      try {
        const decrypted = aesDecrypt(response.data, aesKey)
        response.data = JSON.parse(decrypted)
      } catch (e) {
        // Decryption failed - might be unencrypted error response
        console.warn('Response decryption failed:', e)
        try {
          // Try parsing as JSON directly (fallback for error responses)
          response.data = JSON.parse(response.data)
        } catch (_) {
          return Promise.reject(new Error('Failed to decrypt response'))
        }
      }
    }

    // 201 Created 或其他成功状态码，检查响应体
    if (!response.data) {
      return null
    }

    // 处理响应数据
    const responseData = response.data

    // 确保 responseData 是对象
    if (!responseData || typeof responseData !== 'object') {
      console.error('响应数据格式不正确:', responseData)
      return Promise.reject(new Error('响应数据格式不正确'))
    }

    const { code, message, data } = responseData

    // Handle session expiry (40002) - re-negotiate and retry once
    if (code === 40002) {
      clearSession()
      await negotiateKey()
      const retryConfig = { ...response.config }
      return http(retryConfig)
    }

    // 业务成功
    if (code === 0) {
      return data
    }

    // 业务失败
    const errMsg = message || i18n.global.t('common.error')
    console.error('API Error:', code, errMsg)
    ElMessage.error(errMsg)
    return Promise.reject(new Error(errMsg))
  },
  async (error) => {
    // HTTP错误
    if (error.response) {
      const { status, data } = error.response

      switch (status) {
        case 401:
          ElMessage.error('登录已失效，请重新登录')
          // 清除所有本地登录状态
          localStorage.removeItem('token')
          localStorage.removeItem('username')
          localStorage.removeItem('userId')
          localStorage.removeItem('isDefaultPassword')
          // 重置加密会话
          clearSession()
          // 强制刷新页面，AdminLayout 会自动检测并弹出登录框
          window.location.reload()
          break
        case 403:
          ElMessage.error(i18n.global.t('common.forbidden'))
          break
        case 404:
          ElMessage.error(i18n.global.t('common.notFound'))
          break
        case 500:
          ElMessage.error(i18n.global.t('common.serverError'))
          break
        default:
          ElMessage.error(data?.message || i18n.global.t('common.error'))
      }
    } else if (error.request) {
      ElMessage.error(i18n.global.t('common.networkError'))
    } else {
      ElMessage.error(error.message || i18n.global.t('common.error'))
    }

    return Promise.reject(error)
  }
)

// 生成请求ID
function generateRequestId() {
  return `${Date.now()}-${Math.random().toString(36).substring(2, 9)}`
}

export function initEncryption() {
  return negotiateKey()
}

export default http
