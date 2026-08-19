<template>
  <header class="app-header">
    <div class="container header-main">
      <!-- 移动端汉堡菜单按钮 -->
      <div class="hamburger show-mobile" @click="mobileMenuOpen = !mobileMenuOpen">
        <el-icon :size="24"><Menu /></el-icon>
      </div>

      <!-- Logo -->
      <router-link to="/" class="logo">
        <el-icon :size="28"><Shop /></el-icon>
        <span class="logo-text">AtopMall</span>
      </router-link>

      <!-- 搜索框 -->
      <div class="header-search hide-mobile">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索商品"
          size="large"
          class="search-input"
          @keyup.enter="handleSearch"
        >
          <template #suffix>
            <div class="search-suffix" @click="handleSearch">
              <el-icon :size="20"><Search /></el-icon>
            </div>
          </template>
        </el-input>
      </div>

      <!-- 右侧导航 -->
      <div class="header-right">
        <template v-if="userStore.isAuthenticated">
          <router-link to="/user" class="header-link hide-mobile">{{ userStore.userName }}</router-link>
          <el-button size="small" class="hide-mobile" @click="handleLogout">退出</el-button>
        </template>
        <template v-else>
          <span class="header-link hide-mobile" @click="openLogin">登录</span>
          <span class="header-link hide-mobile" @click="openRegister">注册</span>
        </template>
        <router-link to="/cart" class="header-link cart-link">
          <el-badge :value="cartStore.totalCount" :hidden="cartStore.totalCount === 0">
            <el-icon :size="20"><ShoppingCart /></el-icon>
          </el-badge>
          <span class="hide-mobile">购物车</span>
        </router-link>
      </div>
    </div>

    <!-- 横向分类标签栏（PC端） -->
    <div class="container header-tab-bar hide-mobile">
      <div class="tab-bar">
        <!-- 全部商品分类（级联菜单） -->
        <div
          ref="categoryTriggerRef"
          class="tab-item tab-category-all"
          :class="{ active: showCatDropdown }"
          @mouseenter="openCatDropdown"
          @mouseleave="closeCatDropdown"
        >
          <el-icon :size="14"><Menu /></el-icon>
          <span>全部商品分类</span>
        </div>

        <!-- 全部 -->
        <router-link
          to="/goods"
          class="tab-item"
          :class="{ active: activeTab === 0 }"
          @click="activeTab = 0"
        >
          全部
        </router-link>

        <!-- Tab分类 -->
        <router-link
          v-for="tab in tabCategories"
          :key="tab.id"
          :to="{ path: '/goods', query: { ctg: tab.id } }"
          class="tab-item"
          :class="{ active: activeTab === tab.id }"
          @click="activeTab = tab.id"
        >
          {{ tab.name }}
        </router-link>
      </div>
    </div>

    <!-- 移动端展开菜单 -->
    <transition name="slide-down">
      <div class="mobile-menu" v-if="mobileMenuOpen">
        <div class="mobile-search">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索商品"
            size="large"
            @keyup.enter="handleMobileSearch"
          >
            <template #suffix>
              <div class="search-suffix" @click="handleMobileSearch">
                <el-icon :size="20"><Search /></el-icon>
              </div>
            </template>
          </el-input>
        </div>
        <div class="mobile-nav">
          <router-link to="/goods" class="mobile-nav-item" @click="mobileMenuOpen = false">全部商品</router-link>
          <router-link
            v-for="tab in tabCategories"
            :key="tab.id"
            :to="{ path: '/goods', query: { ctg: tab.id } }"
            class="mobile-nav-item"
            @click="mobileMenuOpen = false"
          >
            {{ tab.name }}
          </router-link>
        </div>
        <div class="mobile-user" v-if="!userStore.isAuthenticated">
          <el-button type="primary" size="large" @click="openLogin; mobileMenuOpen = false">登录</el-button>
          <el-button size="large" @click="openRegister; mobileMenuOpen = false">注册</el-button>
        </div>
        <div class="mobile-user" v-else>
          <router-link to="/user" class="mobile-nav-item" @click="mobileMenuOpen = false">个人中心</router-link>
          <el-button size="small" @click="handleLogout">退出</el-button>
        </div>
      </div>
    </transition>
  </header>

  <!-- Teleport: 级联下拉菜单渲染到 body，彻底脱离所有父容器堆叠上下文 -->
  <Teleport to="body">
    <div
      v-show="showCatDropdown"
      class="cascade-dropdown"
      :style="dropdownStyle"
      @mouseenter="cancelCloseTimer"
      @mouseleave="closeCatDropdown"
    >
      <div class="cascade-left">
        <div
          v-for="cat in allCategories"
          :key="cat.id"
          class="cascade-l1-item"
          :class="{ hover: hoveredL1 === cat.id }"
          @mouseenter="hoveredL1 = cat.id"
        >
          <router-link :to="{ path: '/goods', query: { ctg: cat.id } }" class="l1-name">
            {{ cat.name }}
          </router-link>
          <el-icon v-if="cat.sub_category?.length" :size="12"><ArrowRight /></el-icon>
        </div>
        <div v-if="catLoading" class="cascade-loading">
          <el-icon class="is-loading"><Loading /></el-icon>
          <span>加载中...</span>
        </div>
        <div v-else-if="catError" class="cascade-error">
          <span>分类加载失败</span>
          <el-button type="primary" link size="small" @click="loadCategories">重试</el-button>
        </div>
      </div>
      <div class="cascade-right" v-if="hoveredL1 && getSubCategories(hoveredL1).length > 0">
        <div class="sub-grid">
          <div v-for="sub in getSubCategories(hoveredL1)" :key="sub.id" class="sub-column">
            <router-link :to="{ path: '/goods', query: { ctg: sub.id } }" class="sub-title">
              {{ sub.name }}
            </router-link>
            <div class="sub-links" v-if="sub.sub_category?.length">
              <router-link
                v-for="item in sub.sub_category"
                :key="item.id"
                :to="{ path: '/goods', query: { ctg: item.id } }"
                class="sub-link"
              >
                {{ item.name }}
              </router-link>
            </div>
          </div>
        </div>
      </div>
      <div class="cascade-right cascade-right-empty" v-else-if="hoveredL1">
        <span class="empty-tip">该分类暂无子分类</span>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useCartStore } from '@/store/cart'
import { useAuthModal } from '@/composables/useAuthModal'
import { getCategoryList, type CategoryItem } from '@/api/category'

const router = useRouter()
const userStore = useUserStore()
const cartStore = useCartStore()
const { open } = useAuthModal()
const searchKeyword = ref('')

const openLogin = () => open('login')
const openRegister = () => open('register')

const handleSearch = () => {
  if (searchKeyword.value) {
    mobileMenuOpen.value = false
    router.push({ path: '/goods', query: { q: searchKeyword.value } })
  }
}

const handleMobileSearch = () => {
  handleSearch()
}

const handleLogout = () => {
  userStore.logout()
  router.push('/')
}

// ==================== 分类标签栏 ====================
const categories = ref<CategoryItem[]>([])
const catLoading = ref(false)
const catError = ref(false)
const activeTab = ref(0)
const showCatDropdown = ref(false)
const hoveredL1 = ref<number | null>(null)
const mobileMenuOpen = ref(false)
let closeTimer: ReturnType<typeof setTimeout> | null = null

const categoryTriggerRef = ref<HTMLElement | null>(null)
const dropdownStyle = ref<Record<string, string>>({})

const updateDropdownPosition = () => {
  if (!categoryTriggerRef.value) return
  const rect = categoryTriggerRef.value.getBoundingClientRect()
  dropdownStyle.value = {
    position: 'fixed',
    top: rect.bottom + 4 + 'px',
    left: rect.left + 'px',
  }
}

const flatTabCategories = (list: CategoryItem[]): CategoryItem[] => {
  const result: CategoryItem[] = []
  for (const item of list) {
    if (item.isTab === true) result.push(item)
    if (item.sub_category?.length) {
      result.push(...flatTabCategories(item.sub_category))
    }
  }
  return result
}

const tabCategories = computed(() => {
  return flatTabCategories(categories.value)
})

const allCategories = computed(() => {
  return categories.value
})

const loadCategories = async () => {
  catLoading.value = true
  catError.value = false
  try {
    const res = await getCategoryList()
    categories.value = (res as unknown as CategoryItem[]) || []
  } catch {
    catError.value = true
    categories.value = []
  } finally {
    catLoading.value = false
  }
}

const getSubCategories = (catId: number): CategoryItem[] => {
  const cat = categories.value.find(c => c.id === catId)
  return cat?.sub_category || []
}

const openCatDropdown = async () => {
  if (closeTimer) {
    clearTimeout(closeTimer)
    closeTimer = null
  }
  showCatDropdown.value = true
  await nextTick()
  updateDropdownPosition()
}

const closeCatDropdown = () => {
  closeTimer = setTimeout(() => {
    showCatDropdown.value = false
    hoveredL1.value = null
  }, 150)
}

const cancelCloseTimer = () => {
  if (closeTimer) {
    clearTimeout(closeTimer)
    closeTimer = null
  }
}

const handleResize = () => {
  if (showCatDropdown.value) {
    updateDropdownPosition()
  }
}

onMounted(() => {
  loadCategories()
  window.addEventListener('resize', handleResize)
  window.addEventListener('scroll', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  window.removeEventListener('scroll', handleResize)
})
</script>

<style lang="scss" scoped>
.app-header {
  background: #fff;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  position: sticky;
  top: 0;
  z-index: 1000;
  padding: 16px 0 0 0;

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
    flex-shrink: 0;

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
        border-radius: 24px;
        box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08);
        transition: all 0.3s;
        padding-right: 0;

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
    flex-shrink: 0;

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

  // ==================== 横向分类标签栏 ====================
  .header-tab-bar {
    margin-top: 12px;
    border-top: 1px solid #f0f0f0;
    padding: 0;
  }

  .tab-bar {
    display: flex;
    align-items: center;
    gap: 2px;
    overflow-x: auto;
    white-space: nowrap;
    padding: 0 4px;

    &::-webkit-scrollbar {
      height: 0;
    }
  }

  .tab-item {
    position: relative;
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 12px 22px;
    font-size: 14px;
    color: $text-regular;
    text-decoration: none;
    cursor: pointer;
    transition: all 0.25s;
    font-weight: 400;
    border-radius: 6px;
    margin: 0 2px;

    &:hover {
      color: $primary-color;
      background: #ecf5ff;
    }

    &.active {
      color: $primary-color;
      font-weight: 600;
      background: #ecf5ff;

      &::after {
        content: '';
        position: absolute;
        bottom: 0;
        left: 50%;
        transform: translateX(-50%);
        width: 28px;
        height: 3px;
        border-radius: 3px;
        background: $primary-color;
      }
    }
  }

  .tab-category-all {
    position: relative;
  }
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

// 移动端汉堡菜单样式
.hamburger {
  cursor: pointer;
  padding: 8px;
  color: $text-primary;
  margin-right: 8px;
}

// 移动端展开菜单
.mobile-menu {
  background: #fff;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  padding: 16px;

  .mobile-search {
    margin-bottom: 16px;
  }

  .mobile-nav {
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-bottom: 16px;

    .mobile-nav-item {
      display: block;
      padding: 12px 16px;
      font-size: 15px;
      color: $text-primary;
      text-decoration: none;
      border-radius: 8px;
      transition: all 0.2s;

      &:hover {
        background: #ecf5ff;
        color: $primary-color;
      }
    }
  }

  .mobile-user {
    display: flex;
    gap: 12px;
    align-items: center;
    padding-top: 12px;
    border-top: 1px solid #f0f0f0;
  }
}

// 移动端滑入动画
.slide-down-enter-active {
  transition: all 0.3s ease;
}
.slide-down-leave-active {
  transition: all 0.2s ease;
}
.slide-down-enter-from {
  opacity: 0;
  transform: translateY(-10px);
}
.slide-down-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

// 响应式：移动端
@media (max-width: $bp-mobile) {
  .app-header {
    padding: 10px 0 0 0;

    .header-main {
      .logo {
        font-size: 20px;
        flex-shrink: 0;
      }

      .header-right {
        gap: 8px;
        margin-left: auto;

        .cart-link {
          padding: 6px 10px;
        }
      }
    }
  }
}
</style>

<!-- 级联下拉菜单样式（非 scoped，因为 Teleport 到 body 后脱离组件作用域） -->
<style lang="scss">
// 搜索框后缀图标
.search-suffix {
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #409eff;
  padding: 0 12px;
  transition: all 0.2s;
  border-radius: 0 24px 24px 0;

  &:hover {
    background: rgba(64, 158, 255, 0.1);
    color: #337ecc;
  }

  &:active {
    background: rgba(64, 158, 255, 0.18);
  }
}

.cascade-dropdown {
  display: flex;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.15);
  z-index: 9999;
  min-height: 380px;
  animation: fadeIn 0.2s ease;
  overflow: hidden;

  .cascade-left {
    width: 200px;
    flex-shrink: 0;
    background: #f8fafc;
    padding: 8px 0;
    border-right: 1px solid #ebeef5;
    max-height: 460px;
    overflow-y: auto;
  }

  .cascade-l1-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 20px;
    cursor: pointer;
    font-size: 13px;
    color: #303133;
    transition: all 0.15s;
    margin: 2px 8px;
    border-radius: 6px;

    &:hover,
    &.hover {
      color: #409eff;
      background: #ecf5ff;
    }

    .l1-name {
      flex: 1;
      color: inherit;
      text-decoration: none;
    }

    .el-icon {
      color: #c0c4cc;
      transition: all 0.15s;
    }

    &:hover .el-icon {
      color: #409eff;
    }
  }

  .cascade-right {
    width: 520px;
    flex-shrink: 0;
    padding: 20px;
    max-height: 460px;
    overflow-y: auto;

    .sub-grid {
      display: grid;
      grid-template-columns: repeat(2, 1fr);
      gap: 20px 24px;
    }

    .sub-column {
      .sub-title {
        display: block;
        font-size: 14px;
        font-weight: 600;
        color: #303133;
        margin-bottom: 10px;
        padding-bottom: 8px;
        border-bottom: 1px solid #ebeef5;
        text-decoration: none;

        &:hover {
          color: #409eff;
        }
      }

      .sub-links {
        display: flex;
        flex-wrap: wrap;
        gap: 6px;

        .sub-link {
          font-size: 12px;
          color: #909399;
          padding: 3px 8px;
          border-radius: 4px;
          text-decoration: none;
          transition: all 0.2s;

          &:hover {
            color: #409eff;
            background: #ecf5ff;
          }
        }
      }
    }
  }

  .cascade-right-empty {
    display: flex;
    align-items: center;
    justify-content: center;

    .empty-tip {
      color: #c0c4cc;
      font-size: 14px;
    }
  }

  .cascade-loading,
  .cascade-error {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 40px 20px;
    color: #909399;
    font-size: 13px;
  }
}

// 移动端隐藏级联下拉菜单
@media (max-width: $bp-mobile) {
  .cascade-dropdown {
    display: none !important;
  }
}
</style>