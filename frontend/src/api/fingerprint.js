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
  },
  
  /**
   * 获取指纹详情（管理员）
   * @param {number} id - 指纹ID
   */
  getById(id) {
    return http.get(`/admin/fingerprints/${id}`)
  },
  
  /**
   * 更新指纹（管理员）
   * @param {number} id - 指纹ID
   * @param {Object} data - 更新数据
   * @param {string} data.user_agent - User-Agent
   */
  update(id, data) {
    return http.put(`/admin/fingerprints/${id}`, data)
  },
  
  /**
   * 删除指纹（管理员）
   * @param {number} id - 指纹ID
   */
  delete(id) {
    return http.delete(`/admin/fingerprints/${id}`)
  }
}

