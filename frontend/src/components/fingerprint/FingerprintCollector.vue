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

// 收集并提交指纹
async function submitFingerprint() {
  try {
    // 检查Cookie中是否已有指纹ID
    const savedFingerprintId = Cookies.get('fingerprint_id')
    if (savedFingerprintId) {
      fingerprintId = parseInt(savedFingerprintId)
    }
    
    // 收集指纹信息
    const fingerprintData = await collectFingerprint()
    
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
  if (stayDuration > 0) {
    // 使用sendBeacon确保数据能发送
    recordVisit(stayDuration)
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
  
  // 收集指纹
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

