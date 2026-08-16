<template>
  <header class="app-header">
    <div class="container flex-between">
      <!-- Logo -->
      <router-link to="/" class="logo">
        <el-icon :size="28"><Shop /></el-icon>
        <span class="logo-text">AtopMall</span>
      </router-link>

      <!-- 搜索框 -->
      <div class="header-search">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索商品"
          size="large"
          class="search-input"
          @keyup.enter="handleSearch"
        >
          <template #append>
            <el-button @click="handleSearch">
              <el-icon><Search /></el-icon>
            </el-button>
          </template>
        </el-input>
      </div>

      <!-- 右侧导航 -->
      <div class="header-right">
        <template v-if="userStore.isAuthenticated">
          <router-link to="/user" class="header-link">{{ userStore.userName }}</router-link>
          <el-button size="small" @click="handleLogout">退出</el-button>
        </template>
        <template v-else>
          <span class="header-link" @click="openLogin">登录</span>
          <span class="header-link" @click="openRegister">注册</span>
        </template>
        <router-link to="/cart" class="header-link cart-link">
          <el-badge :value="cartStore.totalCount" :hidden="cartStore.totalCount === 0">
            <el-icon :size="20"><ShoppingCart /></el-icon>
          </el-badge>
          <span>购物车</span>
        </router-link>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useCartStore } from '@/store/cart'
import { useAuthModal } from '@/composables/useAuthModal'

const router = useRouter()
const userStore = useUserStore()
const cartStore = useCartStore()
const { open } = useAuthModal()
const searchKeyword = ref('')

const openLogin = () => open('login')
const openRegister = () => open('register')

const handleSearch = () => {
  if (searchKeyword.value) {
    router.push({ path: '/goods', query: { q: searchKeyword.value } })
  }
}

const handleLogout = () => {
  userStore.logout()
  router.push('/')
}
</script>

<style lang="scss" scoped>
.app-header {
  background: #fff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  position: sticky;
  top: 0;
  z-index: 1000;
  padding: 16px 0;

  .container {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .logo {
    display: flex;
    align-items: center;
    gap: 8px;
    color: $primary-color;
    font-size: 24px;
    font-weight: bold;
    text-decoration: none;

    .logo-text {
      background: linear-gradient(135deg, $primary-color 0%, #66b1ff 100%);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      background-clip: text;
    }
  }

  .header-search {
    flex: 1;
    max-width: 420px;
    margin: 0 40px;

    .search-input {
      :deep(.el-input__wrapper) {
        border-radius: 20px;
        box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
        transition: all 0.3s;

        &:hover {
          box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
        }

        &.is-focus {
          box-shadow: 0 4px 16px rgba(64, 158, 255, 0.25);
        }
      }
    }
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 20px;

    .header-link {
      color: $text-primary;
      text-decoration: none;
      font-size: 14px;
      cursor: pointer;
      transition: color 0.2s;

      &:hover {
        color: $primary-color;
      }

      &.cart-link {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 8px 16px;
        background: #f5f7fa;
        border-radius: 20px;
        transition: all 0.3s;

        &:hover {
          background: #ecf5ff;
          color: $primary-color;
        }
      }
    }
  }
}
</style>