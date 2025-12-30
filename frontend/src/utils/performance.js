/**
 * 性能监控工具
 * 用于监控页面加载时间和API请求时间
 */

class PerformanceMonitor {
  constructor() {
    this.metrics = {
      pageLoad: {},
      apiRequests: [],
    }
    this.init()
  }

  init() {
    // 监控页面加载性能
    if (typeof window !== 'undefined' && window.performance) {
      this.monitorPageLoad()
    }
  }

  /**
   * 监控页面加载性能
   */
  monitorPageLoad() {
    if (document.readyState === 'complete') {
      this.calculatePageLoadTime()
    } else {
      window.addEventListener('load', () => {
        this.calculatePageLoadTime()
      })
    }
  }

  /**
   * 计算页面加载时间
   */
  calculatePageLoadTime() {
    if (!window.performance || !window.performance.timing) {
      return
    }

    const timing = window.performance.timing
    const navigation = window.performance.navigation

    const metrics = {
      // DNS查询时间
      dnsTime: timing.domainLookupEnd - timing.domainLookupStart,
      // TCP连接时间
      tcpTime: timing.connectEnd - timing.connectStart,
      // 请求时间
      requestTime: timing.responseStart - timing.requestStart,
      // 响应时间
      responseTime: timing.responseEnd - timing.responseStart,
      // DOM解析时间
      domParseTime: timing.domInteractive - timing.responseEnd,
      // DOMContentLoaded时间
      domContentLoadedTime: timing.domContentLoadedEventEnd - timing.navigationStart,
      // 页面完全加载时间
      pageLoadTime: timing.loadEventEnd - timing.navigationStart,
      // 首屏渲染时间
      firstPaint: 0,
      // 首次内容绘制时间
      firstContentfulPaint: 0,
    }

    // 获取Paint Timing API数据
    if (window.performance.getEntriesByType) {
      const paintEntries = window.performance.getEntriesByType('paint')
      paintEntries.forEach((entry) => {
        if (entry.name === 'first-paint') {
          metrics.firstPaint = Math.round(entry.startTime)
        } else if (entry.name === 'first-contentful-paint') {
          metrics.firstContentfulPaint = Math.round(entry.startTime)
        }
      })
    }

    this.metrics.pageLoad = {
      ...metrics,
      navigationType: navigation.type,
      timestamp: Date.now(),
    }

    // 在开发环境下输出性能指标
    if (import.meta.env.DEV) {
      console.log('页面加载性能指标:', this.metrics.pageLoad)
    }
  }

  /**
   * 记录API请求时间
   */
  recordApiRequest(url, method, startTime, endTime, status) {
    const duration = endTime - startTime
    const request = {
      url,
      method,
      duration,
      status,
      timestamp: Date.now(),
    }

    this.metrics.apiRequests.push(request)

    // 只保留最近100条记录
    if (this.metrics.apiRequests.length > 100) {
      this.metrics.apiRequests.shift()
    }

    // 在开发环境下输出慢请求警告
    if (import.meta.env.DEV && duration > 1000) {
      console.warn(`慢请求警告: ${method} ${url} 耗时 ${duration}ms`)
    }
  }

  /**
   * 获取性能报告
   */
  getReport() {
    const apiStats = this.calculateApiStats()
    return {
      pageLoad: this.metrics.pageLoad,
      apiStats,
      timestamp: Date.now(),
    }
  }

  /**
   * 计算API请求统计
   */
  calculateApiStats() {
    const requests = this.metrics.apiRequests
    if (requests.length === 0) {
      return {
        total: 0,
        average: 0,
        max: 0,
        min: 0,
        slowRequests: [],
      }
    }

    const durations = requests.map((r) => r.duration)
    const total = durations.reduce((sum, d) => sum + d, 0)
    const average = total / durations.length
    const max = Math.max(...durations)
    const min = Math.min(...durations)

    // 找出慢请求（超过1秒）
    const slowRequests = requests.filter((r) => r.duration > 1000)

    return {
      total: requests.length,
      average: Math.round(average),
      max,
      min,
      slowRequests: slowRequests.slice(-10), // 最近10个慢请求
    }
  }

  /**
   * 清除性能数据
   */
  clear() {
    this.metrics = {
      pageLoad: {},
      apiRequests: [],
    }
  }
}

// 创建单例
const performanceMonitor = new PerformanceMonitor()

export default performanceMonitor

