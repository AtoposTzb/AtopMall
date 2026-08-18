<template>
  <div class="user-center-page">

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
            <el-menu-item index="/user/message">
              <el-icon><ChatLineSquare /></el-icon>
              <span>我的留言</span>
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


const userStore = useUserStore()
</script>

<style lang="scss" scoped>
.user-center-page {
  min-height: 100vh;
  background: linear-gradient(180deg, #f0f5ff 0%, #f5f7fa 30%);
}

.container {
  padding: 28px 0;
}

.user-center-layout {
  display: flex;
  gap: 24px;
  min-height: 600px;
}

.sidebar {
  width: 240px;
  flex-shrink: 0;
  background: #fff;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);

  .user-info {
    text-align: center;
    padding: 28px 24px 24px;
    background: linear-gradient(135deg, $primary-color 0%, #66b1ff 100%);
    position: relative;

    &::after {
      content: '';
      position: absolute;
      bottom: 0;
      left: 0;
      right: 0;
      height: 20px;
      background: #fff;
      border-radius: 50% 50% 0 0;
    }

    :deep(.el-avatar) {
      margin-bottom: 12px;
      background: rgba(255, 255, 255, 0.2);
      border: 3px solid rgba(255, 255, 255, 0.5);
      box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
      transition: transform 0.3s;

      &:hover {
        transform: scale(1.05);
      }
    }

    h3 {
      font-size: 16px;
      font-weight: 600;
      color: #fff;
      margin: 0;
      padding-bottom: 8px;
    }
  }

  .user-menu {
    border-right: none;
    padding: 8px 0;

    :deep(.el-menu-item) {
      height: 48px;
      line-height: 48px;
      margin: 2px 8px;
      border-radius: 8px;
      font-size: 14px;
      color: $text-primary;
      transition: all 0.2s;

      &:hover {
        background: #f0f5ff;
        color: $primary-color;
      }

      &.is-active {
        background: linear-gradient(135deg, #ecf5ff 0%, #f0f5ff 100%);
        color: $primary-color;
        font-weight: 600;

        &::before {
          content: '';
          position: absolute;
          left: 0;
          top: 50%;
          transform: translateY(-50%);
          width: 3px;
          height: 20px;
          background: $primary-color;
          border-radius: 0 3px 3px 0;
        }
      }

      .el-icon {
        margin-right: 8px;
        font-size: 18px;
      }
    }
  }
}

.content {
  flex: 1;
  background: #fff;
  border-radius: 12px;
  padding: 28px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  min-height: 500px;
}
</style>