<template>
  <div id="app">
    <AppHeader v-if="showHeader" />
    <router-view />
    <AuthModal />
    <!-- 全局固定左侧悬浮侧边栏 -->
    <LeftFixedSidebar />
    <!-- 全局固定右侧悬浮侧边栏 -->
    <FixedSidebar />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useCartStore } from '@/store/cart'
import AppHeader from '@/components/AppHeader.vue'
import LeftFixedSidebar from '@/components/LeftFixedSidebar.vue'
import AuthModal from '@/components/AuthModal.vue'
import FixedSidebar from '@/components/FixedSidebar.vue'

const route = useRoute()
const userStore = useUserStore()
const cartStore = useCartStore()

const showHeader = computed(() => {
  return !['Login', 'Register', 'NotFound'].includes(String(route.name))
})

onMounted(async () => {
  try {
    await userStore.initUser()
    if (userStore.isAuthenticated) {
      cartStore.loadCartList()
    }
  } catch (error) {
    console.error('App 初始化失败', error)
  }
})
</script>

<style lang="scss">
#app {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}
</style>