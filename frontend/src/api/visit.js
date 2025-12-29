import http from './http'

/**
 * 访问记录API
 */
export default {
  /**
   * 记录访问
   */
  recordVisit(data) {
    return http.post('/visit', data)
  }
}

