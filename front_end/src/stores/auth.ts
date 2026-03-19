import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User } from '@/types'
import { AuthService } from '@/services'

const TOKEN_REFRESH_INTERVAL = 50 * 60 * 1000 // 50分钟
const TOKEN_MAX_AGE = 3 * 24 * 60 * 60 * 1000 // 3天

export const useAuthStore = defineStore('auth', () => {
  // ===== State =====
  const user = ref<User | null>(null)
  const accessToken = ref<string | null>(localStorage.getItem('access_token'))
  const refreshToken = ref<string | null>(localStorage.getItem('refresh_token'))
  const isLoading = ref(false)
  let refreshTimer: ReturnType<typeof setInterval> | null = null

  // ===== Getters =====
  const isLoggedIn = computed(() => !!accessToken.value)
  const isAdmin = computed(() => user.value?.role === 'admin')
  const isStaff = computed(() => user.value?.role === 'staff' || user.value?.role === 'admin')
  const username = computed(() => user.value?.username || '')

  // ===== Actions =====
  async function login(usernameInput: string, password: string) {
    isLoading.value = true
    try {
      const data = await AuthService.login({ username: usernameInput, password })
      setTokens(data.access_token, data.refresh_token)
      user.value = data.user
      startAutoRefresh()
      return data.user
    } finally {
      isLoading.value = false
    }
  }

  async function register(usernameInput: string, password: string, phone?: string, email?: string) {
    isLoading.value = true
    try {
      const data = await AuthService.register({
        username: usernameInput,
        password,
        phone,
        email
      })
      return data
    } finally {
      isLoading.value = false
    }
  }

  async function fetchUser() {
    if (!accessToken.value) return null
    try {
      user.value = await AuthService.getMe()
      return user.value
    } catch {
      logout()
      return null
    }
  }

  async function updateProfile(data: { phone?: string; email?: string }) {
    user.value = await AuthService.updateMe(data)
    return user.value
  }

  async function updatePassword(oldPassword: string, newPassword: string) {
    return AuthService.updatePassword(oldPassword, newPassword)
  }

  function setTokens(access: string, refresh: string) {
    accessToken.value = access
    refreshToken.value = refresh
    localStorage.setItem('access_token', access)
    localStorage.setItem('refresh_token', refresh)
    localStorage.setItem('token_time', Date.now().toString())
  }

  function logout() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('token_time')
    stopAutoRefresh()
  }

  function startAutoRefresh() {
    stopAutoRefresh()
    refreshTimer = setInterval(async () => {
      if (!refreshToken.value) return
      try {
        const data = await AuthService.refreshToken(refreshToken.value)
        accessToken.value = data.access_token
        localStorage.setItem('access_token', data.access_token)
      } catch {
        logout()
      }
    }, TOKEN_REFRESH_INTERVAL)
  }

  function stopAutoRefresh() {
    if (refreshTimer) {
      clearInterval(refreshTimer)
      refreshTimer = null
    }
  }

  // 初始化：检查 token 是否过期
  async function initialize() {
    const tokenTime = localStorage.getItem('token_time')
    if (tokenTime) {
      const elapsed = Date.now() - parseInt(tokenTime)
      if (elapsed > TOKEN_MAX_AGE) {
        logout()
        return
      }
    }

    if (accessToken.value) {
      await fetchUser()
      if (user.value) {
        startAutoRefresh()
      }
    }
  }

  return {
    // State
    user,
    accessToken,
    refreshToken,
    isLoading,
    // Getters
    isLoggedIn,
    isAdmin,
    isStaff,
    username,
    // Actions
    login,
    register,
    fetchUser,
    updateProfile,
    updatePassword,
    logout,
    initialize
  }
})
