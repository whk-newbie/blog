import http from './http'

/**
 * 统计API
 */
export default {
  /**
   * 获取仪表盘统计数据（管理员）
   */
  getDashboardStats() {
    return http.get('/admin/stats/dashboard')
  },
  
  /**
   * 获取访问统计
   * @param {Object} params - 查询参数
   * @param {string} params.start_date - 开始日期 (YYYY-MM-DD)
   * @param {string} params.end_date - 结束日期 (YYYY-MM-DD)
   * @param {string} params.type - 统计类型 (daily/weekly/monthly)
   */
  getVisitStats(params = {}) {
    return http.get('/admin/stats/visits', { params })
  },
  
  /**
   * 获取热门文章
   * @param {Object} params - 查询参数
   * @param {number} params.limit - 返回数量
   * @param {number} params.days - 统计天数
   */
  getPopularArticles(params = {}) {
    return http.get('/admin/stats/popular-articles', { params })
  },
  
  /**
   * 获取访问来源统计
   * @param {Object} params - 查询参数
   * @param {string} params.start_date - 开始日期 (YYYY-MM-DD)
   * @param {string} params.end_date - 结束日期 (YYYY-MM-DD)
   */
  getReferrerStats(params = {}) {
    return http.get('/admin/stats/referrers', { params })
  }
}

