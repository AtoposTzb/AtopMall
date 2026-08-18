import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getUserDetail, type UserInfo } from '@/api/user'
import { isAdminUser } from '@/utils/jwt'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => isAdminUser(token.value))
  const userName = computed(() => userInfo.value?.name || '')

  const initUser = async () => {
    if (token.value && !userInfo.value) {
      try {
        const res = await getUserDetail()
        userInfo.value = res
      } catch (error) {
        console.error('获取用户信息失败', error)
      }
    }
  }

  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  const setUserInfo = (info: UserInfo) => {
    userInfo.value = info
  }

  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  return {
    token,
    userInfo,
    isAuthenticated,
    isAdmin,
    userName,
    initUser,
    setToken,
    setUserInfo,
    logout
  }
})