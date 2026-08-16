<template>
  <div class="user-center-page">
    <AppHeader />
    
    <div class="container">
      <div class="user-center-layout">
        <!-- 侧边栏 -->
        <aside class="sidebar">
          <div class="user-info">
            <el-avatar :size="60" icon="User" />
            <h3>{{ userStore.userName }}</h3>
          </div>
          
          <el-menu :default-active="$route.path" router class="user-menu">
            <el-menu-item index="/user">
              <el-icon><User /></el-icon>
              <span>个人信息</span>
            </el-menu-item>
            <el-menu-item index="/user/address">
              <el-icon><Location /></el-icon>
              <span>地址管理</span>
            </el-menu-item>
            <el-menu-item index="/user/favorite">
              <el-icon><Star /></el-icon>
              <span>我的收藏</span>
            </el-menu-item>
            <el-menu-item index="/user/cart">
              <el-icon><ShoppingCart /></el-icon>
              <span>购物车</span>
            </el-menu-item>
            <el-menu-item index="/user/orders">
              <el-icon><Document /></el-icon>
              <span>我的订单</span>
            </el-menu-item>
          </el-menu>
        </aside>

        <!-- 内容区 -->
        <main class="content">
          <router-view />
        </main>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useUserStore } from '@/store/user'
import AppHeader from '@/components/AppHeader.vue'

const userStore = useUserStore()
</script>

<style lang="scss" scoped>
.user-center-page {
  min-height: 100vh;
  background: #f5f7fa;
}

.container {
  padding: 24px 0;
}

.user-center-layout {
  display: flex;
  gap: 24px;
  min-height: 600px;
}

.sidebar {
  width: 220px;
  flex-shrink: 0;
  background: #fff;
  border-radius: 8px;
  padding: 24px 0;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

  .user-info {
    text-align: center;
    padding: 0 24px 24px;
    border-bottom: 1px solid #f0f0f0;
    margin-bottom: 16px;

    .el-avatar {
      margin-bottom: 12px;
      background: $primary-color;
    }

    h3 {
      font-size: 16px;
      color: $text-primary;
      margin: 0;
    }
  }

  .user-menu {
    border-right: none;

    .el-menu-item {
      height: 48px;
      line-height: 48px;

      .el-icon {
        margin-right: 8px;
      }
    }
  }
}

.content {
  flex: 1;
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}
</style>
