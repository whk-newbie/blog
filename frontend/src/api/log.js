import http from './http'

/**
 * 日志管理API
 */
export default {
  /**
   * 获取日志列表（管理员）
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.page_size - 每页数量
   * @param {string} params.level - 日志级别（可选）
   * @param {string} params.source - 日志来源（可选）
   * @param {string} params.start_date - 开始日期 (YYYY-MM-DD)（可选）
   * @param {string} params.end_date - 结束日期 (YYYY-MM-DD)（可选）
   */
  getLogs(params = {}) {
    return http.get('/admin/logs', { params })
  },

  /**
   * 获取日志详情（管理员）
   * @param {number} id - 日志ID
   */
  getLogById(id) {
    return http.get(`/admin/logs/${id}`)
  },

  /**
   * 清理旧日志（管理员）
   * @param {Object} data - 清理请求
   * @param {number} data.retention_days - 保留天数（默认90天）
   */
  cleanupLogs(data = {}) {
    return http.post('/admin/logs/cleanup', data)
  }
}

