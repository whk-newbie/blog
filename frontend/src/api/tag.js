import http from './http'

/**
 * 标签API
 */
export default {
  /**
   * 获取标签列表
   */
  list(params = {}) {
    return http.get('/tags', { params })
  },

  /**
   * 获取标签详情
   */
  getById(id) {
    return http.get(`/tags/${id}`)
  },

  /**
   * 根据Slug获取标签
   */
  getBySlug(slug) {
    return http.get(`/tags/slug/${slug}`)
  },

  /**
   * 创建标签（管理员）
   */
  create(data) {
    return http.post('/admin/tags', data)
  },

  /**
   * 更新标签（管理员）
   */
  update(id, data) {
    return http.put(`/admin/tags/${id}`, data)
  },

  /**
   * 删除标签（管理员）
   */
  delete(id) {
    return http.delete(`/admin/tags/${id}`)
  }
}

