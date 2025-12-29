import http from './http'

/**
 * 爬虫任务API
 */
export default {
  /**
   * 获取爬虫任务列表（管理员）
   * @param {Object} params - 查询参数
   * @param {number} params.page - 页码
   * @param {number} params.page_size - 每页数量
   * @param {string} params.status - 状态筛选 (running/completed/failed)
   * @param {string} params.task_id - 任务ID
   */
  getTasks(params = {}) {
    return http.get('/admin/crawler/tasks', { params })
  },

  /**
   * 获取任务详情（管理员）
   * @param {string|number} taskId - 任务ID或数字ID
   */
  getTaskById(taskId) {
    return http.get(`/admin/crawler/tasks/${taskId}`)
  }
}

