<template>
  <!-- 指纹收集组件，无UI，自动收集 -->
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { collectFingerprint } from '@/utils/fingerprint'
import fingerprintApi from '@/api/fingerprint'
import visitApi from '@/api/visit'
import Cookies from 'js-cookie'

const route = useRoute()
let enterTime = Date.now()
let fingerprintId = null
let audioFingerprintCollected = false

// 用户交互后尝试收集音频指纹（如果之前失败）
async function tryCollectAudioFingerprint() {
  if (audioFingerprintCollected || !fingerprintId) {
    return
  }
  
  try {
    const { getAudioFingerprint } = await import('@/utils/fingerprint')
    const audioFingerprint = await getAudioFingerprint()
    if (audioFingerprint) {
      audioFingerprintCollected = true
      // 音频指纹已收集，但不更新已提交的指纹（因为只是辅助信息）
    }
  } catch (error) {
    // 静默处理错误
  }
}

// 用户交互处理
function handleUserInteraction() {
  // 用户交互后，尝试收集音频指纹（如果之前失败）
  tryCollectAudioFingerprint()
}

// 收集并提交指纹
async function submitFingerprint() {
  try {
    // 检查Cookie中是否已有指纹ID
    const savedFingerprintId = Cookies.get('fingerprint_id')
    if (savedFingerprintId) {
      fingerprintId = parseInt(savedFingerprintId)
    }
    
    // 收集指纹信息（包括音频，如果失败会返回 null）
    const fingerprintData = await collectFingerprint()
    
    // 标记音频指纹是否已收集
    audioFingerprintCollected = fingerprintData.audio !== null
    
    // 提交指纹
    const result = await fingerprintApi.collectFingerprint({
      ...fingerprintData,
      user_agent: navigator.userAgent
    })
    
    // 保存指纹ID到Cookie（365天）
    if (result && result.fingerprint_id) {
      fingerprintId = result.fingerprint_id
      Cookies.set('fingerprint_id', fingerprintId.toString(), { expires: 365 })
    }
  } catch (error) {
    console.error('指纹收集失败:', error)
  }
}

// 记录访问
async function recordVisit(stayDuration = null) {
  try {
    if (!fingerprintId) {
      const savedFingerprintId = Cookies.get('fingerprint_id')
      if (savedFingerprintId) {
        fingerprintId = parseInt(savedFingerprintId)
      }
    }
    
    if (!fingerprintId) {
      return
    }
    
    const referrer = document.referrer || ''
    const pageTitle = document.title
    const url = route.fullPath
    
    // 提取文章ID（如果有）
    let articleId = null
    if (route.name === 'ArticleDetail' && route.params.slug) {
      // 这里可以根据实际情况获取文章ID
      // 暂时不处理，因为需要先获取文章信息
    }
    
    await visitApi.recordVisit({
      fingerprint_id: fingerprintId,
      url: url,
      referrer: referrer,
      page_title: pageTitle,
      article_id: articleId,
      stay_duration: stayDuration,
      user_agent: navigator.userAgent
    })
  } catch (error) {
    console.error('访问记录失败:', error)
  }
}

// 页面卸载时记录停留时间
function handleBeforeUnload() {
  const stayDuration = Math.floor((Date.now() - enterTime) / 1000)
  if (stayDuration > 0 && fingerprintId) {
    // 使用sendBeacon确保数据能发送（页面卸载时）
    const referrer = document.referrer || ''
    const pageTitle = document.title
    const url = route.fullPath
    
    const data = JSON.stringify({
      fingerprint_id: fingerprintId,
      url: url,
      referrer: referrer,
      page_title: pageTitle,
      article_id: null,
      stay_duration: stayDuration,
      user_agent: navigator.userAgent
    })
    
    // 使用sendBeacon发送数据（不需要等待响应）
    try {
      navigator.sendBeacon('/api/v1/visit', new Blob([data], { type: 'application/json' }))
    } catch (e) {
      // 如果sendBeacon失败，尝试使用fetch（keepalive）
      fetch('/api/v1/visit', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: data,
        keepalive: true
      }).catch(() => {}) // 忽略错误
    }
  }
}

// 页面可见性变化时记录
function handleVisibilityChange() {
  if (document.hidden) {
    const stayDuration = Math.floor((Date.now() - enterTime) / 1000)
    if (stayDuration > 0) {
      recordVisit(stayDuration)
    }
  } else {
    enterTime = Date.now()
  }
}

onMounted(async () => {
  enterTime = Date.now()
  
  // 监听用户交互事件（用于在用户交互后重试音频指纹收集）
  const interactionEvents = ['click', 'touchstart', 'keydown', 'mousedown']
  interactionEvents.forEach(eventType => {
    document.addEventListener(eventType, handleUserInteraction, { once: true, passive: true })
  })
  
  // 收集指纹（音频指纹如果失败会返回 null，不会显示警告）
  await submitFingerprint()
  
  // 记录页面进入
  await recordVisit()
  
  // 监听页面卸载
  window.addEventListener('beforeunload', handleBeforeUnload)
  
  // 监听页面可见性变化
  document.addEventListener('visibilitychange', handleVisibilityChange)
})

onUnmounted(() => {
  const stayDuration = Math.floor((Date.now() - enterTime) / 1000)
  if (stayDuration > 0) {
    recordVisit(stayDuration)
  }
  
  window.removeEventListener('beforeunload', handleBeforeUnload)
  document.removeEventListener('visibilitychange', handleVisibilityChange)
})
</script>

