import axios from 'axios'
import { ElMessage } from 'element-plus'
import i18n from '@/locales'
import { encrypt, decrypt, getAppKey } from '@/utils/crypto'
import performanceMonitor from '@/utils/performance'

// 创建axios实例
const http = axios.create({
  baseURL: '/api/v1', // 使用相对路径，通过 Vite 代理转发
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 不需要加密的接口列表（登录等公开接口）
const NO_ENCRYPTION_PATHS = [
  '/auth/login',
  '/auth/refresh',
  '/fingerprint',
  '/visit'
]

// 请求拦截器
http.interceptors.request.use(
  async (config) => {
    // 记录请求开始时间
    config.metadata = {
      startTime: performance.now(),
    }
    
    // 添加token
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    
    // 添加请求ID
    config.headers['X-Request-ID'] = generateRequestId()
    
    // 处理请求加密
    // 只对POST/PUT/PATCH请求且不在排除列表中的请求进行加密
    const shouldEncrypt = 
      (config.method === 'post' || config.method === 'put' || config.method === 'patch') &&
      config.data &&
      !NO_ENCRYPTION_PATHS.some(path => config.url?.includes(path))
    
    if (shouldEncrypt) {
      try {
        // 获取应用密钥
        const appKey = await getAppKey()
        
        if (appKey) {
          // 将请求数据转换为JSON字符串
          const jsonData = typeof config.data === 'string' 
            ? config.data 
            : JSON.stringify(config.data)
          
          // 加密数据
          const encryptedData = await encrypt(jsonData, appKey)
          
          // 构建加密请求格式
          config.data = {
            encrypted_data: encryptedData,
            timestamp: Date.now()
          }
        }
      } catch (error) {
        console.warn('请求加密失败，使用明文:', error)
        // 加密失败，继续使用原始数据
      }
    }
    
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
    
    // 201 Created 或其他成功状态码，检查响应体
    if (!response.data) {
      return null
    }
    
    // 处理响应解密
    let responseData = response.data
    
    // 如果响应数据是字符串，先尝试解析为 JSON（可能是加密格式的字符串）
    if (typeof responseData === 'string') {
      try {
        // 清理可能的空白字符和多余内容
        const cleanedData = responseData.trim()
        // 如果字符串以 { 或 [ 开头，尝试解析为 JSON
        if (cleanedData.startsWith('{') || cleanedData.startsWith('[')) {
          // 查找第一个完整的 JSON 对象
          let jsonEnd = -1
          if (cleanedData.startsWith('{')) {
            let braceCount = 0
            for (let i = 0; i < cleanedData.length; i++) {
              if (cleanedData[i] === '{') braceCount++
              if (cleanedData[i] === '}') braceCount--
              if (braceCount === 0 && i > 0) {
                jsonEnd = i + 1
                break
              }
            }
          } else if (cleanedData.startsWith('[')) {
            let bracketCount = 0
            for (let i = 0; i < cleanedData.length; i++) {
              if (cleanedData[i] === '[') bracketCount++
              if (cleanedData[i] === ']') bracketCount--
              if (bracketCount === 0 && i > 0) {
                jsonEnd = i + 1
                break
              }
            }
          }
          
          const jsonString = jsonEnd > 0 ? cleanedData.substring(0, jsonEnd) : cleanedData
          responseData = JSON.parse(jsonString)
        }
      } catch (e) {
        // 解析失败，保持原始字符串，后续会检查是否是加密格式
        console.warn('响应数据解析失败，保持原始字符串:', e)
      }
    }
    
    // 检查是否是加密响应格式（对象格式或字符串格式）
    if (responseData && (
      (typeof responseData === 'object' && responseData.encrypted_data && typeof responseData.encrypted_data === 'string') ||
      (typeof responseData === 'string' && responseData.includes('encrypted_data'))
    )) {
      // 如果是字符串格式，先解析为对象
      if (typeof responseData === 'string') {
        try {
          responseData = JSON.parse(responseData)
        } catch (e) {
          console.warn('无法解析加密响应字符串:', e)
          return Promise.reject(new Error('响应格式错误'))
        }
      }
      
      // 确保有 encrypted_data 字段
      if (!responseData.encrypted_data || typeof responseData.encrypted_data !== 'string') {
        console.error('加密响应格式不正确:', responseData)
        return Promise.reject(new Error('响应格式错误：缺少 encrypted_data 字段'))
      }
      
      try {
        // 获取应用密钥（如果本地没有，尝试从API获取）
        let appKey = await getAppKey()
        
        // 如果本地没有密钥，尝试强制刷新从API获取
        if (!appKey) {
          console.log('本地没有应用密钥，尝试从API获取...')
          appKey = await getAppKey(true) // 强制刷新
        }
        
        if (!appKey) {
          console.error('无法获取应用密钥，无法解密响应')
          return Promise.reject(new Error('无法获取应用密钥，请重新登录'))
        }
        
        // 解密数据
        console.log('开始解密响应数据...')
        const decryptedData = await decrypt(responseData.encrypted_data, appKey)
        console.log('解密成功，解密后数据长度:', decryptedData.length)
        
        // 尝试解析为JSON
        try {
          // 清理可能的空白字符
          const cleanedData = decryptedData.trim()
          // 查找第一个完整的 JSON 对象
          let jsonEnd = -1
          if (cleanedData.startsWith('{')) {
            let braceCount = 0
            for (let i = 0; i < cleanedData.length; i++) {
              if (cleanedData[i] === '{') braceCount++
              if (cleanedData[i] === '}') braceCount--
              if (braceCount === 0 && i > 0) {
                jsonEnd = i + 1
                break
              }
            }
          } else if (cleanedData.startsWith('[')) {
            let bracketCount = 0
            for (let i = 0; i < cleanedData.length; i++) {
              if (cleanedData[i] === '[') bracketCount++
              if (cleanedData[i] === ']') bracketCount--
              if (bracketCount === 0 && i > 0) {
                jsonEnd = i + 1
                break
              }
            }
          }
          
          const jsonString = jsonEnd > 0 ? cleanedData.substring(0, jsonEnd) : cleanedData
          responseData = JSON.parse(jsonString)
          console.log('解密后JSON解析成功')
        } catch (e) {
          console.error('解密后数据解析失败:', e)
          console.error('解密后的数据（前200字符）:', decryptedData.substring(0, 200))
          return Promise.reject(new Error(`解密后数据解析失败: ${e.message}`))
        }
      } catch (error) {
        console.error('响应解密失败:', error)
        console.error('加密数据（前100字符）:', responseData.encrypted_data.substring(0, 100))
        return Promise.reject(new Error(`响应解密失败: ${error.message}`))
      }
    }
    
    // 确保 responseData 是对象
    if (!responseData || typeof responseData !== 'object') {
      console.error('响应数据格式不正确:', responseData)
      return Promise.reject(new Error('响应数据格式不正确'))
    }
    
    const { code, message, data } = responseData
    
    // 业务成功
    if (code === 0) {
      return data
    }
    
    // 业务失败
    ElMessage.error(message || i18n.global.t('common.error'))
    return Promise.reject(new Error(message || i18n.global.t('common.error')))
  },
  (error) => {
    // HTTP错误
    if (error.response) {
      const { status, data } = error.response
      
      switch (status) {
        case 401:
          ElMessage.error(i18n.global.t('common.unauthorized'))
          localStorage.removeItem('token')
          window.location.href = '/'
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

export default http

