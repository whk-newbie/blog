import http from './http'

/**
 * 文章API
 */
export default {
  /**
   * 获取已发布文章列表（公开）
   */
  list(params = {}) {
    return http.get('/articles', { params })
  },

  /**
   * 获取文章详情（公开）
   */
  getById(id) {
    return http.get(`/articles/${id}`)
  },

  /**
   * 根据Slug获取文章（公开）
   */
  getBySlug(slug) {
    return http.get(`/articles/slug/${slug}`)
  },

  /**
   * 搜索文章（公开）
   */
  search(keyword, params = {}) {
    return http.get('/articles/search', {
      params: {
        keyword,
        ...params
      }
    })
  },

  /**
   * 获取文章列表（管理员）
   */
  adminList(params = {}) {
    return http.get('/admin/articles', { params })
  },

  /**
   * 创建文章（管理员）
   */
  create(data) {
    return http.post('/admin/articles', data)
  },

  /**
   * 更新文章（管理员）
   */
  update(id, data) {
    return http.put(`/admin/articles/${id}`, data)
  },

  /**
   * 删除文章（管理员）
   */
  delete(id) {
    return http.delete(`/admin/articles/${id}`)
  },

  /**
   * 发布文章（管理员）
   */
  publish(id) {
    return http.post(`/admin/articles/${id}/publish`)
  },

  /**
   * 取消发布（管理员）
   */
  unpublish(id) {
    return http.post(`/admin/articles/${id}/unpublish`)
  }
}

