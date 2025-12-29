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
  }
}

