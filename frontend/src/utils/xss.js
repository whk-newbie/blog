/**
 * XSS防护工具
 * 提供输入转义和输出验证功能
 */

/**
 * HTML转义字符映射
 */
const HTML_ESCAPE_MAP = {
  '&': '&amp;',
  '<': '&lt;',
  '>': '&gt;',
  '"': '&quot;',
  "'": '&#x27;',
  '/': '&#x2F;',
}

/**
 * 转义HTML特殊字符
 * @param {string} str - 需要转义的字符串
 * @returns {string} 转义后的字符串
 */
export function escapeHtml(str) {
  if (typeof str !== 'string') {
    return ''
  }
  
  return str.replace(/[&<>"'/]/g, (char) => HTML_ESCAPE_MAP[char] || char)
}

/**
 * 转义URL参数
 * @param {string} str - 需要转义的URL字符串
 * @returns {string} 转义后的URL字符串
 */
export function escapeUrl(str) {
  if (typeof str !== 'string') {
    return ''
  }
  
  return encodeURIComponent(str)
}

/**
 * 清理HTML内容，移除潜在的XSS攻击代码
 * @param {string} html - HTML内容
 * @returns {string} 清理后的HTML内容
 */
export function sanitizeHtml(html) {
  if (typeof html !== 'string') {
    return ''
  }
  
  // 创建一个临时的DOM元素
  const div = document.createElement('div')
  div.textContent = html
  
  // 获取转义后的内容
  return div.innerHTML
}

/**
 * 验证和清理用户输入
 * @param {string} input - 用户输入
 * @param {Object} options - 选项
 * @param {number} options.maxLength - 最大长度
 * @param {boolean} options.allowHtml - 是否允许HTML
 * @returns {string} 清理后的输入
 */
export function sanitizeInput(input, options = {}) {
  if (typeof input !== 'string') {
    return ''
  }
  
  const { maxLength = 10000, allowHtml = false } = options
  
  // 限制长度
  let sanitized = input.substring(0, maxLength)
  
  // 如果不允许HTML，转义HTML字符
  if (!allowHtml) {
    sanitized = escapeHtml(sanitized)
  }
  
  // 移除潜在的脚本标签
  sanitized = sanitized.replace(/<script\b[^<]*(?:(?!<\/script>)<[^<]*)*<\/script>/gi, '')
  
  // 移除事件处理器（onclick, onerror等）
  sanitized = sanitized.replace(/on\w+\s*=\s*["'][^"']*["']/gi, '')
  
  return sanitized
}

/**
 * 验证URL是否安全
 * @param {string} url - URL字符串
 * @returns {boolean} 是否安全
 */
export function isSafeUrl(url) {
  if (typeof url !== 'string') {
    return false
  }
  
  try {
    const urlObj = new URL(url)
    
    // 只允许http和https协议
    return urlObj.protocol === 'http:' || urlObj.protocol === 'https:'
  } catch (e) {
    // 相对URL也是安全的
    return !url.startsWith('javascript:') && !url.startsWith('data:')
  }
}

/**
 * 清理对象中的所有字符串值
 * @param {Object} obj - 对象
 * @param {Object} options - 选项
 * @returns {Object} 清理后的对象
 */
export function sanitizeObject(obj, options = {}) {
  if (typeof obj !== 'object' || obj === null) {
    return obj
  }
  
  if (Array.isArray(obj)) {
    return obj.map(item => sanitizeObject(item, options))
  }
  
  const sanitized = {}
  
  for (const key in obj) {
    if (Object.prototype.hasOwnProperty.call(obj, key)) {
      const value = obj[key]
      
      if (typeof value === 'string') {
        sanitized[key] = sanitizeInput(value, options)
      } else if (typeof value === 'object' && value !== null) {
        sanitized[key] = sanitizeObject(value, options)
      } else {
        sanitized[key] = value
      }
    }
  }
  
  return sanitized
}

/**
 * Vue指令：安全的HTML渲染
 * 使用方式：v-safe-html="htmlContent"
 */
export const safeHtmlDirective = {
  mounted(el, binding) {
    if (binding.value) {
      // 使用DOMPurify或简单的清理
      el.innerHTML = sanitizeHtml(binding.value)
    }
  },
  updated(el, binding) {
    if (binding.value) {
      el.innerHTML = sanitizeHtml(binding.value)
    }
  },
}

