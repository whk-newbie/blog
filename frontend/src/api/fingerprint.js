import http from './http'

/**
 * 指纹API
 */
export default {
  /**
   * 收集指纹
   */
  collectFingerprint(data) {
    return http.post('/fingerprint', data)
  },
  
  /**
   * 获取指纹列表（管理员）
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.page_size - 每页数量
   */
  list(params = {}) {
    return http.get('/admin/fingerprints', { params })
  }
}

