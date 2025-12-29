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
  }
}

