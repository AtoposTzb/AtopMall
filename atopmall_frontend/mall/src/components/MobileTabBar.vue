<template>
  <div class="mobile-tab-bar">
    <router-link to="/" class="tab-item" :class="{ active: $route.path === '/' }">
      <el-icon :size="22"><HomeFilled /></el-icon>
      <span>首页</span>
    </router-link>
    <router-link to="/goods" class="tab-item" :class="{ active: $route.path.startsWith('/goods') }">
      <el-icon :size="22"><Grid /></el-icon>
      <span>分类</span>
    </router-link>
    <router-link to="/cart" class="tab-item" :class="{ active: $route.path === '/cart' }">
      <el-badge :value="cartStore.totalCount" :hidden="cartStore.totalCount === 0" :max="99">
        <el-icon :size="22"><ShoppingCart /></el-icon>
      </el-badge>
      <span>购物车</span>
    </router-link>
    <div class="tab-item" :class="{ active: $route.path.startsWith('/user') }" @click="handleUserClick">
      <el-icon :size="22"><User /></el-icon>
      <span>{{ userStore.isAuthenticated ? '我的' : '登录' }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useCartStore } from '@/store/cart'
import { useAuthModal } from '@/composables/useAuthModal'

const router = useRouter()
const userStore = useUserStore()
const cartStore = useCartStore()
const { open } = useAuthModal()

const handleUserClick = () => {
  if (userStore.isAuthenticated) {
    router.push('/user')
  } else {
    open('login')
  }
}
</script>

<style lang="scss" scoped>
.mobile-tab-bar {
  display: none;
}

@media (max-width: $bp-mobile) {
  .mobile-tab-bar {
    display: flex;
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    z-index: 1000;
    background: #fff;
    border-top: 1px solid $border-base;
    padding: 4px 0 env(safe-area-inset-bottom, 6px);
    justify-content: space-around;
    align-items: center;
    box-shadow: 0 -2px 10px rgba(0, 0, 0, 0.06);

    .tab-item {
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 2px;
      padding: 4px 12px;
      color: $text-secondary;
      text-decoration: none;
      cursor: pointer;
      transition: color 0.2s;
      font-size: 10px;
      min-width: 56px;

      &:hover {
        color: $primary-color;
      }

      &.active {
        color: $primary-color;
      }

      :deep(.el-badge__content) {
        font-size: 10px;
      }
    }
  }
}
</style>