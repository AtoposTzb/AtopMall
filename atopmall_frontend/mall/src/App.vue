<template>
  <div id="app">
    <router-view />
    <AuthModal />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useUserStore } from '@/store/user'
import { useCartStore } from '@/store/cart'
import AuthModal from '@/components/AuthModal.vue'

const userStore = useUserStore()
const cartStore = useCartStore()

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