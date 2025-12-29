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
  const maxReconnectAttempts = 5
  const reconnectDelay = ref(3000) // 初始重连延迟3秒
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
        reconnectDelay.value = 3000
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
        
        // 自动重连
        if (reconnectAttempts.value < maxReconnectAttempts) {
          scheduleReconnect()
        } else {
          console.error('WebSocket重连次数已达上限')
        }
      }
    } catch (error) {
      console.error('WebSocket连接失败:', error)
      connected.value = false
      scheduleReconnect()
    }
  }

  /**
   * 断开WebSocket连接
   */
  const disconnect = () => {
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
   * 安排重连
   */
  const scheduleReconnect = () => {
    if (reconnectTimer) {
      return
    }

    reconnectAttempts.value++
    reconnectTimer = setTimeout(() => {
      reconnectTimer = null
      console.log(`尝试重连 (${reconnectAttempts.value}/${maxReconnectAttempts})...`)
      connect()
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

