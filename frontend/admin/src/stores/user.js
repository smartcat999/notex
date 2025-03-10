import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, logout, getProfile } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('admin_token') || '')
  const user = ref(null)
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)

  async function loginUser(credentials) {
    try {
      loading.value = true
      const response = await login(credentials)
      token.value = response.data.token
      localStorage.setItem('admin_token', response.data.token)
      await fetchUserProfile()
      return response
    } finally {
      loading.value = false
    }
  }

  async function logoutUser() {
    try {
      loading.value = true
      await logout()
      token.value = ''
      user.value = null
      localStorage.removeItem('admin_token')
    } finally {
      loading.value = false
    }
  }

  async function fetchUserProfile() {
    try {
      loading.value = true
      const response = await getProfile()
      user.value = response.data
    } finally {
      loading.value = false
    }
  }

  return {
    token,
    user,
    loading,
    isAuthenticated,
    loginUser,
    logoutUser,
    fetchUserProfile,
  }
}) 