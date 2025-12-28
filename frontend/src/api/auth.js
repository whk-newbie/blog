import http from './http'

/**
 * 认证相关API
 */
export const authAPI = {
  /**
   * 登录
   * @param {string} username - 用户名
   * @param {string} password - 密码
   */
  login(username, password) {
    return http.post('/auth/login', {
      username,
      password
    })
  },

  /**
   * 登出（前端清除token）
   */
  logout() {
    // 后端不需要登出接口，只需要清除本地token
    return Promise.resolve()
  },

  /**
   * 修改密码
   * @param {string} oldPassword - 旧密码
   * @param {string} newPassword - 新密码
   */
  changePassword(oldPassword, newPassword) {
    return http.put('/auth/password', {
      old_password: oldPassword,
      new_password: newPassword
    })
  },

  /**
   * 验证Token
   */
  verifyToken() {
    return http.get('/auth/verify')
  },

  /**
   * 刷新Token
   */
  refreshToken() {
    return http.post('/auth/refresh')
  }
}

