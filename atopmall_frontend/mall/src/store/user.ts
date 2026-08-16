import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getUserDetail, type UserInfo } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  // 状态
  const token = ref<string>(localStorage.getItem('token') || '')
  const userInfo = ref<UserInfo | null>(null)

  // 计算属性
  const isAuthenticated = computed(() => !!token.value)
  const userName = computed(() => userInfo.value?.name || '')

  // 初始化用户信息
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

  // 设置 Token
  const setToken = (newToken: string) => {
    token.value = newToken
    localStorage.setItem('token', newToken)
  }

  // 设置用户信息
  const setUserInfo = (info: UserInfo) => {
    userInfo.value = info
  }

  // 退出登录
  const logout = () => {
    token.value = ''
    userInfo.value = null
    localStorage.removeItem('token')
  }

  return {
    token,
    userInfo,
    isAuthenticated,
    userName,
    initUser,
    setToken,
    setUserInfo,
    logout
  }
})
