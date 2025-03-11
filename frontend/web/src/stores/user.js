import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, register, logout, getProfile } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)

  async function loginUser(credentials) {
    try {
      loading.value = true
      const response = await login(credentials)
      if (response && response.token) {
        token.value = response.token
        localStorage.setItem('token', response.token)
        await fetchUserProfile()
        return response
      }
      throw new Error('登录失败：未获取到有效的 token')
    } finally {
      loading.value = false
    }
  }

  async function registerUser(userData) {
    try {
      loading.value = true
      const response = await register(userData)
      return response
    } finally {
      loading.value = false
    }
  }

  async function logoutUser() {
    try {
      loading.value = true
      await logout()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      token.value = ''
      user.value = null
      localStorage.removeItem('token')
      loading.value = false
    }
  }

  async function fetchUserProfile() {
    try {
      loading.value = true
      const response = await getProfile()
      if (response && response.user) {
        user.value = response.user
      }
    } finally {
      loading.value = false
    }
  }

  function setUser(userData) {
    user.value = userData
  }

  return {
    token,
    user,
    loading,
    isAuthenticated,
    loginUser,
    registerUser,
    logoutUser,
    fetchUserProfile,
    setUser,
  }
}) 