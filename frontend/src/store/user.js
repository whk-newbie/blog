import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const username = ref(localStorage.getItem('username') || '')
  const userId = ref(localStorage.getItem('userId') || '')
  const isDefaultPassword = ref(localStorage.getItem('isDefaultPassword') === 'true')

  const setToken = (newToken) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUserInfo = (info) => {
    username.value = info.username
    userId.value = info.id
    isDefaultPassword.value = info.is_default_password || false

    localStorage.setItem('username', info.username)
    localStorage.setItem('userId', info.id)
    localStorage.setItem('isDefaultPassword', info.is_default_password || false)
  }

  const clearUserInfo = () => {
    token.value = ''
    username.value = ''
    userId.value = ''
    isDefaultPassword.value = false

    localStorage.removeItem('token')
    localStorage.removeItem('username')
    localStorage.removeItem('userId')
    localStorage.removeItem('isDefaultPassword')
  }

  const login = (loginData) => {
    setToken(loginData.token)
    setUserInfo(loginData.user)
  }

  const logout = () => {
    clearUserInfo()
  }

  const isLoggedIn = () => {
    return !!token.value
  }

  return {
    token,
    username,
    userId,
    isDefaultPassword,
    setToken,
    setUserInfo,
    clearUserInfo,
    login,
    logout,
    isLoggedIn
  }
})

