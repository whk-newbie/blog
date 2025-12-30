import axios from 'axios'
import { ElMessage } from 'element-plus'
import i18n from '@/locales'
import { encrypt, decrypt, getAppKey } from '@/utils/crypto'

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
    
    // 检查是否是加密响应格式
    if (responseData.encrypted_data && typeof responseData.encrypted_data === 'string') {
      try {
        // 获取应用密钥
        const appKey = await getAppKey()
        
        if (appKey) {
          // 解密数据
          const decryptedData = await decrypt(responseData.encrypted_data, appKey)
          
          // 尝试解析为JSON
          try {
            responseData = JSON.parse(decryptedData)
          } catch (e) {
            // 如果不是JSON，直接使用解密后的字符串
            responseData = decryptedData
          }
        }
      } catch (error) {
        console.warn('响应解密失败，使用原始数据:', error)
        // 解密失败，继续使用原始数据
      }
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

