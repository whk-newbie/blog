import axios from 'axios'
import { ElMessage } from 'element-plus'
import i18n from '@/locales'

// 创建axios实例
const http = axios.create({
  baseURL: '/api/v1', // 使用相对路径，通过 Vite 代理转发
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
http.interceptors.request.use(
  (config) => {
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
  (response) => {
    // 204 No Content - 删除成功，没有响应体
    if (response.status === 204) {
      return null
    }
    
    // 201 Created 或其他成功状态码，检查响应体
    if (!response.data) {
      return null
    }
    
    const { code, message, data } = response.data
    
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

