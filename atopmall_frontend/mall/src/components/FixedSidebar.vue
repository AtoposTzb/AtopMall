<template>
  <div class="fixed-sidebar" :style="{ transform: `translateY(calc(-50% + ${offsetY}px))` }">
    <div class="fsb-drag-handle" @mousedown="onMouseDown">
      <span class="drag-line"></span>
      <span class="drag-line"></span>
      <span class="drag-line"></span>
    </div>
    <div class="fsb-item" @click="$router.push('/user/cart')">
      <div class="fsb-icon">
        <el-badge :value="cartStore.totalCount" :hidden="cartStore.totalCount === 0">
          <el-icon :size="20"><ShoppingCart /></el-icon>
        </el-badge>
      </div>
      <span class="fsb-label">购物车</span>
    </div>
    <div class="fsb-item" @click="userStore.isAuthenticated ? $router.push('/user') : openLogin()">
      <div class="fsb-icon">
        <el-icon :size="20"><User /></el-icon>
      </div>
      <span class="fsb-label">我的</span>
    </div>
    <div class="fsb-item" @click="userStore.isAuthenticated ? $router.push('/user/orders') : openLogin()">
      <div class="fsb-icon">
        <el-icon :size="20"><Document /></el-icon>
      </div>
      <span class="fsb-label">订单</span>
    </div>
    <div class="fsb-item" @click="userStore.isAuthenticated ? $router.push('/user/favorite') : openLogin()">
      <div class="fsb-icon">
        <el-icon :size="20"><Star /></el-icon>
      </div>
      <span class="fsb-label">收藏</span>
    </div>
    <div class="fsb-item" @click="scrollToTop">
      <div class="fsb-icon">
        <el-icon :size="20"><Top /></el-icon>
      </div>
      <span class="fsb-label">顶部</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from '@/store/user'
import { useCartStore } from '@/store/cart'
import { useAuthModal } from '@/composables/useAuthModal'
import { useDraggable } from '@/composables/useDraggable'

const userStore = useUserStore()
const cartStore = useCartStore()
const { open } = useAuthModal()
const { offsetY, onMouseDown } = useDraggable()

const openLogin = () => open('login')

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}
</script>

<style lang="scss" scoped>
.fixed-sidebar {
  position: fixed;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 999;
  display: flex;
  flex-direction: column;
  gap: 2px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  padding: 0 0 6px 0;
  user-select: none;
  transition: box-shadow 0.2s;

  .fsb-drag-handle {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 3px;
    padding: 8px 0;
    cursor: grab;
    border-radius: 10px 10px 0 0;
    transition: background 0.2s;

    &:hover {
      background: #f5f7fa;
    }

    &:active {
      cursor: grabbing;
      background: #ecf5ff;
    }

    .drag-line {
      display: block;
      width: 18px;
      height: 2px;
      background: #c0c4cc;
      border-radius: 1px;
      transition: all 0.2s;
    }

    &:hover .drag-line {
      background: #909399;
    }
  }

  .fsb-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    padding: 12px 10px;
    cursor: pointer;
    transition: all 0.2s;
    color: $text-regular;
    width: 64px;

    &:hover {
      color: $primary-color;
      background: #ecf5ff;
    }

    .fsb-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 24px;
    }

    .fsb-label {
      font-size: 11px;
      white-space: nowrap;
    }
  }
}

@media (max-width: $bp-mobile) {
  .fixed-sidebar {
    display: none;
  }
}
</style>