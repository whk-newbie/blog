import http from './http'

/**
 * 配置管理API
 */
export default {
  /**
   * 获取配置列表（管理员）
   * @param {Object} params - 查询参数
   * @param {string} params.config_type - 配置类型（可选）
   */
  getConfigs(params = {}) {
    return http.get('/admin/configs', { params })
  },

  /**
   * 获取配置详情（管理员）
   * @param {number} id - 配置ID
   */
  getConfigById(id) {
    return http.get(`/admin/configs/${id}`)
  },

  /**
   * 创建配置（管理员）
   * @param {Object} data - 配置数据
   * @param {string} data.config_key - 配置键
   * @param {string} data.config_value - 配置值
   * @param {string} data.config_type - 配置类型
   * @param {boolean} data.is_encrypted - 是否加密
   * @param {boolean} data.is_active - 是否启用
   * @param {string} data.description - 描述
   */
  createConfig(data) {
    return http.post('/admin/configs', data)
  },

  /**
   * 更新配置（管理员）
   * @param {number} id - 配置ID
   * @param {Object} data - 配置数据
   */
  updateConfig(id, data) {
    return http.put(`/admin/configs/${id}`, data)
  },

  /**
   * 删除配置（管理员）
   * @param {number} id - 配置ID
   */
  deleteConfig(id) {
    return http.delete(`/admin/configs/${id}`)
  },

  /**
   * 生成爬虫Token（管理员）
   * @param {Object} data - 生成Token请求
   * @param {string} data.name - Token名称
   */
  generateCrawlerToken(data) {
    return http.post('/admin/configs/generate-crawler-token', data)
  },

  /**
   * 获取公开的站点配置
   * @returns {Promise} 站点配置数据(博客标题、备案信息等)
   */
  getSiteConfig() {
    return http.get('/site/config')
  }
}

