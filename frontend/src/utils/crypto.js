/**
 * 前端加密工具
 * 使用Web Crypto API实现AES-256-GCM加密/解密
 * 与后端Go加密中间件兼容
 */

/**
 * 将字符串密钥转换为ArrayBuffer
 * @param {string} keyString - 密钥字符串（32字节）
 * @returns {ArrayBuffer}
 */
function keyStringToArrayBuffer(keyString) {
  const encoder = new TextEncoder()
  return encoder.encode(keyString)
}

/**
 * 从ArrayBuffer导入密钥
 * @param {ArrayBuffer} keyData - 密钥数据
 * @returns {Promise<CryptoKey>}
 */
async function importKey(keyData) {
  return await crypto.subtle.importKey(
    'raw',
    keyData,
    { name: 'AES-GCM' },
    false,
    ['encrypt', 'decrypt']
  )
}

/**
 * 加密数据（AES-256-GCM）
 * @param {string} plaintext - 明文
 * @param {string} keyString - 密钥字符串（32字节）
 * @returns {Promise<string>} Base64编码的密文
 */
export async function encrypt(plaintext, keyString) {
  try {
    // 将密钥字符串转换为ArrayBuffer
    const keyData = keyStringToArrayBuffer(keyString)
    
    // 导入密钥
    const key = await importKey(keyData)
    
    // 生成随机nonce（12字节，GCM标准）
    const nonce = crypto.getRandomValues(new Uint8Array(12))
    
    // 将明文转换为ArrayBuffer
    const encoder = new TextEncoder()
    const plaintextBuffer = encoder.encode(plaintext)
    
    // 加密数据
    const ciphertext = await crypto.subtle.encrypt(
      {
        name: 'AES-GCM',
        iv: nonce,
        tagLength: 128 // GCM标签长度128位
      },
      key,
      plaintextBuffer
    )
    
    // 将nonce和密文合并（nonce在前）
    const combined = new Uint8Array(nonce.length + ciphertext.byteLength)
    combined.set(nonce, 0)
    combined.set(new Uint8Array(ciphertext), nonce.length)
    
    // Base64编码
    return btoa(String.fromCharCode(...combined))
  } catch (error) {
    throw new Error(`加密失败: ${error.message}`)
  }
}

/**
 * 解密数据（AES-256-GCM）
 * @param {string} ciphertextBase64 - Base64编码的密文
 * @param {string} keyString - 密钥字符串（32字节）
 * @returns {Promise<string>} 明文
 */
export async function decrypt(ciphertextBase64, keyString) {
  try {
    // Base64解码
    const combined = Uint8Array.from(atob(ciphertextBase64), c => c.charCodeAt(0))
    
    // 提取nonce（前12字节）和密文
    const nonceSize = 12
    if (combined.length < nonceSize) {
      throw new Error('密文长度不足')
    }
    
    const nonce = combined.slice(0, nonceSize)
    const ciphertext = combined.slice(nonceSize)
    
    // 将密钥字符串转换为ArrayBuffer
    const keyData = keyStringToArrayBuffer(keyString)
    
    // 导入密钥
    const key = await importKey(keyData)
    
    // 解密数据
    const plaintextBuffer = await crypto.subtle.decrypt(
      {
        name: 'AES-GCM',
        iv: nonce,
        tagLength: 128
      },
      key,
      ciphertext
    )
    
    // 将ArrayBuffer转换为字符串
    const decoder = new TextDecoder()
    return decoder.decode(plaintextBuffer)
  } catch (error) {
    throw new Error(`解密失败: ${error.message}`)
  }
}

/**
 * 从API获取应用密钥并保存到localStorage
 * @returns {Promise<string|null>} 应用密钥，获取失败返回null
 */
export async function fetchAppKeyFromAPI() {
  try {
    const configApi = (await import('@/api/config')).default
    
    // 获取配置列表，查找application_key
    const configs = await configApi.getConfigs({ config_type: 'application_key' })
    
    if (configs && configs.length > 0) {
      // 应用密钥在配置中，后端会解密后返回
      const appKey = configs[0].config_value
      
      // 保存到localStorage
      if (appKey && appKey.length === 32) {
        localStorage.setItem('app_key', appKey)
        return appKey
      }
    }
  } catch (error) {
    console.warn('从API获取应用密钥失败:', error)
  }
  
  return null
}

/**
 * 获取应用密钥（从localStorage或API）
 * @param {boolean} forceRefresh - 是否强制从API刷新
 * @returns {Promise<string|null>} 应用密钥，获取失败返回null
 */
export async function getAppKey(forceRefresh = false) {
  // 如果强制刷新，清除localStorage中的密钥
  if (forceRefresh) {
    localStorage.removeItem('app_key')
  }
  
  // 先从localStorage获取
  let appKey = localStorage.getItem('app_key')
  
  if (appKey && appKey.length === 32) {
    return appKey
  }
  
  // 如果localStorage中没有或无效，尝试从API获取
  // 注意：这需要登录后才能获取，所以首次登录请求可能无法加密
  // 但登录后的请求可以使用加密
  appKey = await fetchAppKeyFromAPI()
  
  return appKey
}

/**
 * 清除应用密钥
 */
export function clearAppKey() {
  localStorage.removeItem('app_key')
}

