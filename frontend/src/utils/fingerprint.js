/**
 * 指纹收集工具函数
 */

/**
 * 获取Canvas指纹
 */
export function getCanvasFingerprint() {
  try {
    const canvas = document.createElement('canvas')
    const ctx = canvas.getContext('2d')
    if (!ctx) return null
    
    canvas.width = 200
    canvas.height = 50
    ctx.textBaseline = 'top'
    ctx.font = '14px Arial'
    ctx.fillStyle = '#f60'
    ctx.fillRect(125, 1, 62, 20)
    ctx.fillStyle = '#069'
    ctx.fillText('fingerprint', 2, 15)
    ctx.fillStyle = 'rgba(102, 204, 0, 0.7)'
    ctx.fillText('fingerprint', 4, 17)
    
    return canvas.toDataURL()
  } catch (e) {
    return null
  }
}

/**
 * 获取WebGL指纹
 */
export function getWebGLFingerprint() {
  try {
    const canvas = document.createElement('canvas')
    const gl = canvas.getContext('webgl') || canvas.getContext('experimental-webgl')
    if (!gl) return null
    
    const debugInfo = gl.getExtension('WEBGL_debug_renderer_info')
    if (!debugInfo) return null
    
    return {
      vendor: gl.getParameter(gl.VENDOR),
      renderer: gl.getParameter(gl.RENDERER),
      unmaskedVendor: gl.getParameter(debugInfo.UNMASKED_VENDOR_WEBGL),
      unmaskedRenderer: gl.getParameter(debugInfo.UNMASKED_RENDERER_WEBGL),
      version: gl.getParameter(gl.VERSION),
      shadingLanguageVersion: gl.getParameter(gl.SHADING_LANGUAGE_VERSION)
    }
  } catch (e) {
    return null
  }
}

/**
 * 获取屏幕信息
 */
export function getScreenInfo() {
  return {
    width: screen.width,
    height: screen.height,
    availWidth: screen.availWidth,
    availHeight: screen.availHeight,
    colorDepth: screen.colorDepth,
    pixelDepth: screen.pixelDepth
  }
}

/**
 * 获取时区信息
 */
export function getTimezone() {
  return Intl.DateTimeFormat().resolvedOptions().timeZone
}

/**
 * 获取语言信息
 */
export function getLanguage() {
  return {
    language: navigator.language,
    languages: navigator.languages || [navigator.language]
  }
}

/**
 * 获取平台信息
 */
export function getPlatform() {
  return navigator.platform
}

/**
 * 获取插件信息
 */
export function getPlugins() {
  const plugins = []
  if (navigator.plugins) {
    for (let i = 0; i < navigator.plugins.length; i++) {
      plugins.push({
        name: navigator.plugins[i].name,
        description: navigator.plugins[i].description
      })
    }
  }
  return plugins
}

/**
 * 获取字体列表（简化版）
 */
export async function getFonts() {
  const baseFonts = ['monospace', 'sans-serif', 'serif']
  const testString = 'mmmmmmmmmmlli'
  const testSize = '72px'
  const h = document.getElementsByTagName('body')[0]
  
  const s = document.createElement('span')
  s.style.fontSize = testSize
  s.innerHTML = testString
  const defaultWidth = {}
  const defaultHeight = {}
  
  for (let i = 0; i < baseFonts.length; i++) {
    s.style.fontFamily = baseFonts[i]
    h.appendChild(s)
    defaultWidth[baseFonts[i]] = s.offsetWidth
    defaultHeight[baseFonts[i]] = s.offsetHeight
    h.removeChild(s)
  }
  
  const detected = []
  const fonts = [
    'Arial', 'Verdana', 'Times New Roman', 'Courier New', 'Georgia',
    'Palatino', 'Garamond', 'Bookman', 'Comic Sans MS', 'Trebuchet MS',
    'Arial Black', 'Impact', 'Tahoma', 'Lucida Console', 'Courier',
    'Lucida Sans Unicode', 'Franklin Gothic Medium', 'Century Gothic'
  ]
  
  for (let i = 0; i < fonts.length; i++) {
    let detected_font = false
    for (let j = 0; j < baseFonts.length; j++) {
      s.style.fontFamily = fonts[i] + ',' + baseFonts[j]
      h.appendChild(s)
      const matched = (s.offsetWidth !== defaultWidth[baseFonts[j]] || s.offsetHeight !== defaultHeight[baseFonts[j]])
      h.removeChild(s)
      if (matched) {
        detected_font = true
      }
    }
    if (detected_font) {
      detected.push(fonts[i])
    }
  }
  
  return detected
}

/**
 * 获取音频指纹
 * 使用 AnalyserNode 代替已弃用的 ScriptProcessorNode
 */
export async function getAudioFingerprint() {
  return new Promise((resolve) => {
    try {
      const AudioContextClass = window.AudioContext || window.webkitAudioContext
      if (!AudioContextClass) {
        resolve(null)
        return
      }
      
      const context = new AudioContextClass()
      
      // 检查 AudioContext 状态，如果是 suspended 状态，尝试恢复
      // 注意：resume() 需要用户手势，如果失败则静默返回 null
      if (context.state === 'suspended') {
        context.resume().catch(() => {
          context.close()
          resolve(null)
          return
        })
      }
      
      const oscillator = context.createOscillator()
      const analyser = context.createAnalyser()
      const gainNode = context.createGain()
      
      // 设置分析器参数
      analyser.fftSize = 2048
      analyser.smoothingTimeConstant = 0.8
      
      // 静音输出
      gainNode.gain.value = 0
      
      // 连接节点：oscillator -> analyser -> gainNode -> destination
      oscillator.connect(analyser)
      analyser.connect(gainNode)
      gainNode.connect(context.destination)
      
      oscillator.type = 'triangle'
      oscillator.frequency.value = 10000
      
      // 启动振荡器
      oscillator.start(0)
      
      // 使用 requestAnimationFrame 读取音频数据（代替已弃用的 ScriptProcessorNode）
      let samplesCollected = 0
      const maxSamples = 3 // 收集3个样本
      const dataArray = new Float32Array(analyser.fftSize)
      
      function collectSample() {
        try {
          analyser.getFloatTimeDomainData(dataArray)
          
          // 计算样本的绝对值之和
          let sum = 0
          for (let i = 0; i < dataArray.length; i++) {
            sum += Math.abs(dataArray[i])
          }
          
          samplesCollected++
          
          // 收集足够的样本后，计算指纹并清理
          if (samplesCollected >= maxSamples) {
            const fingerprint = Math.floor(sum * 1000000).toString()
            oscillator.stop()
            context.close()
            resolve(fingerprint)
          } else {
            // 继续收集下一个样本
            requestAnimationFrame(collectSample)
          }
        } catch (e) {
          // 如果出错，清理并返回 null
          try {
            oscillator.stop()
            context.close()
          } catch (cleanupError) {
            // 忽略清理错误
          }
          resolve(null)
        }
      }
      
      // 延迟一小段时间后开始收集，确保音频已经开始
      setTimeout(() => {
        if (context.state === 'running') {
          requestAnimationFrame(collectSample)
        } else {
          // 如果 AudioContext 无法运行，清理并返回 null
          try {
            oscillator.stop()
            context.close()
          } catch (cleanupError) {
            // 忽略清理错误
          }
          resolve(null)
        }
      }, 10)
      
      // 设置超时，防止无限等待
      setTimeout(() => {
        try {
          if (oscillator) {
            oscillator.stop()
          }
          if (context) {
            context.close()
          }
        } catch (cleanupError) {
          // 忽略清理错误
        }
        if (samplesCollected < maxSamples) {
          resolve(null)
        }
      }, 1000) // 1秒超时
      
    } catch (e) {
      // 静默处理错误，不输出警告
      resolve(null)
    }
  })
}

/**
 * 收集所有指纹信息
 */
export async function collectFingerprint() {
  const fingerprint = {
    canvas: getCanvasFingerprint(),
    webgl: getWebGLFingerprint(),
    screen: getScreenInfo(),
    timezone: getTimezone(),
    language: getLanguage(),
    platform: getPlatform(),
    plugins: getPlugins(),
    fonts: await getFonts(),
    audio: await getAudioFingerprint()
  }
  
  return fingerprint
}

