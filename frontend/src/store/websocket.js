import { defineStore } from 'pinia'
import { ref } from 'vue'
import { useUserStore } from './user'

// WebSocket连接配置
const WS_BASE_URL = import.meta.env.VITE_WS_BASE_URL || (() => {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = window.location.host
  return `${protocol}//${host}`
})()

export const useWebSocketStore = defineStore('websocket', () => {
  const ws = ref(null)
  const connected = ref(false)
  const tasks = ref([])
  const reconnectAttempts = ref(0)
  const reconnectDelay = ref(3000) // 初始重连延迟3秒
  const shouldReconnect = ref(true) // 是否应该重连（页面可见时）
  let reconnectTimer = null
  let heartbeatTimer = null

  /**
   * 连接WebSocket
   */
  const connect = () => {
    const userStore = useUserStore()
    const token = userStore.token

    if (!token) {
      console.error('WebSocket连接失败: 缺少Token')
      return
    }

    // 设置允许重连
    shouldReconnect.value = true

    // 如果已经连接，先断开
    if (ws.value && ws.value.readyState === WebSocket.OPEN) {
      disconnect()
    }

    try {
      const wsUrl = `${WS_BASE_URL}/ws/crawler/tasks?token=${token}`
      ws.value = new WebSocket(wsUrl)

      ws.value.onopen = () => {
        console.log('WebSocket连接成功')
        connected.value = true
        reconnectAttempts.value = 0
        reconnectDelay.value = 3000 // 连接成功后重置延迟
        startHeartbeat()
      }

      ws.value.onmessage = (event) => {
        try {
          const message = JSON.parse(event.data)
          handleMessage(message)
        } catch (error) {
          console.error('WebSocket消息解析失败:', error)
        }
      }

      ws.value.onerror = (error) => {
        console.error('WebSocket错误:', error)
        connected.value = false
      }

      ws.value.onclose = () => {
        console.log('WebSocket连接关闭')
        connected.value = false
        stopHeartbeat()
        
        // 自动重连（只要页面可见且允许重连，就无限重试）
        if (shouldReconnect.value && isPageVisible()) {
          scheduleReconnect()
        }
      }
    } catch (error) {
      console.error('WebSocket连接失败:', error)
      connected.value = false
      // 如果页面可见且允许重连，则安排重连
      if (shouldReconnect.value && isPageVisible()) {
        scheduleReconnect()
      }
    }
  }

  /**
   * 断开WebSocket连接
   */
  const disconnect = () => {
    shouldReconnect.value = false // 停止重连
    if (reconnectTimer) {
      clearTimeout(reconnectTimer)
      reconnectTimer = null
    }
    stopHeartbeat()
    
    if (ws.value) {
      ws.value.close()
      ws.value = null
    }
    connected.value = false
  }

  /**
   * 处理WebSocket消息
   */
  const handleMessage = (message) => {
    if (message.type === 'task_update') {
      updateTask(message.data)
    } else if (message.type === 'pong') {
      // 心跳响应
      console.log('收到心跳响应')
    }
  }

  /**
   * 更新任务
   */
  const updateTask = (taskData) => {
    const index = tasks.value.findIndex(t => t.task_id === taskData.task_id)
    if (index !== -1) {
      // 更新现有任务
      tasks.value[index] = { ...tasks.value[index], ...taskData }
    } else {
      // 添加新任务
      tasks.value.push(taskData)
    }
  }

  /**
   * 启动心跳
   */
  const startHeartbeat = () => {
    stopHeartbeat()
    heartbeatTimer = setInterval(() => {
      if (ws.value && ws.value.readyState === WebSocket.OPEN) {
        ws.value.send(JSON.stringify({ type: 'ping' }))
      }
    }, 30000) // 每30秒发送一次心跳
  }

  /**
   * 停止心跳
   */
  const stopHeartbeat = () => {
    if (heartbeatTimer) {
      clearInterval(heartbeatTimer)
      heartbeatTimer = null
    }
  }

  /**
   * 检查页面是否可见
   */
  const isPageVisible = () => {
    if (typeof document === 'undefined') return true
    return !document.hidden
  }

  /**
   * 安排重连
   */
  const scheduleReconnect = () => {
    // 如果已经安排了重连，或者不应该重连，或者页面不可见，则返回
    if (reconnectTimer || !shouldReconnect.value || !isPageVisible()) {
      return
    }

    reconnectAttempts.value++
    reconnectTimer = setTimeout(() => {
      reconnectTimer = null
      console.log(`尝试重连 (第${reconnectAttempts.value}次)...`)
      
      // 再次检查是否应该重连和页面是否可见
      if (shouldReconnect.value && isPageVisible()) {
        connect()
      }
    }, reconnectDelay.value)

    // 指数退避：每次重连延迟翻倍，最大30秒
    reconnectDelay.value = Math.min(reconnectDelay.value * 2, 30000)
  }

  /**
   * 重置任务列表
   */
  const resetTasks = () => {
    tasks.value = []
  }

  /**
   * 添加任务到列表
   */
  const addTasks = (newTasks) => {
    tasks.value = [...tasks.value, ...newTasks]
  }

  /**
   * 设置任务列表
   */
  const setTasks = (newTasks) => {
    tasks.value = newTasks
  }

  // 监听页面可见性变化
  if (typeof document !== 'undefined') {
    document.addEventListener('visibilitychange', () => {
      if (document.hidden) {
        // 页面隐藏时，不进行重连
        console.log('页面隐藏，暂停重连')
      } else {
        // 页面可见时，如果未连接则尝试连接
        console.log('页面可见，检查连接状态')
        if (!connected.value && shouldReconnect.value) {
          // 如果当前没有连接且允许重连，则尝试连接
          if (!reconnectTimer) {
            reconnectDelay.value = 3000 // 重置延迟
            scheduleReconnect()
          }
        }
      }
    })
  }

  return {
    ws,
    connected,
    tasks,
    connect,
    disconnect,
    updateTask,
    resetTasks,
    addTasks,
    setTasks
  }
})

