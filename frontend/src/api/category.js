import http from './http'

/**
 * 分类API
 */
export default {
  /**
   * 获取分类列表
   */
  list(params = {}) {
    return http.get('/categories', { params })
  },

  /**
   * 获取分类详情
   */
  getById(id) {
    return http.get(`/categories/${id}`)
  },

  /**
   * 根据Slug获取分类
   */
  getBySlug(slug) {
    return http.get(`/categories/slug/${slug}`)
  },

  /**
   * 创建分类（管理员）
   */
  create(data) {
    return http.post('/admin/categories', data)
  },

  /**
   * 更新分类（管理员）
   */
  update(id, data) {
    return http.put(`/admin/categories/${id}`, data)
  },

  /**
   * 删除分类（管理员）
   */
  delete(id) {
    return http.delete(`/admin/categories/${id}`)
  }
}

