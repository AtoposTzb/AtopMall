<template>
  <div class="home-page">
    <AppHeader />

    <div class="container" v-if="pageError">
      <div class="page-error-bar">
        <div class="page-error-content">
          <el-icon><WarningFilled /></el-icon>
          <span>首页数据加载失败，请检查网络后重试</span>
        </div>
        <el-button type="primary" size="small" @click="retryAll">重新加载</el-button>
      </div>
    </div>

    <div class="hero-section">
      <div class="container hero-layout">
        <!-- 左侧全部分类 -->
        <div class="sidebar" @mouseleave="handleCategoryLeave">
          <div class="sidebar-title">
            <el-icon><Menu /></el-icon>
            <span>全部商品分类</span>
          </div>
          <div class="sidebar-list" v-loading="catLoading">
            <template v-if="catError">
              <div class="sidebar-error">
                <el-icon><WarningFilled /></el-icon>
                <span>分类加载失败</span>
                <el-button type="primary" link size="small" @click="loadCategories">点击重试</el-button>
              </div>
            </template>
            <template v-else-if="!catLoading && categories.length === 0">
              <div class="sidebar-empty">
                <span>暂无分类</span>
              </div>
            </template>
            <template v-else>
              <div
                v-for="cat in categories"
                :key="cat.id"
                class="sidebar-item"
                :class="{ active: activeCategory === cat.id }"
                @mouseenter="handleCategoryEnter(cat.id)"
              >
                <span class="item-name">{{ cat.name }}</span>
                <el-icon><ArrowRight /></el-icon>
                <div class="category-dropdown" v-show="activeCategory === cat.id && cat.sub_category?.length">
                  <div class="dropdown-inner">
                    <div v-for="sub in cat.sub_category" :key="sub.id" class="dropdown-column">
                      <router-link :to="{ path: '/goods', query: { ctg: sub.id } }" class="column-title">
                        {{ sub.name }}
                      </router-link>
                      <div class="column-links">
                        <router-link
                          v-for="item in sub.sub_category"
                          :key="item.id"
                          :to="{ path: '/goods', query: { ctg: item.id } }"
                          class="link-item"
                        >
                          {{ item.name }}
                        </router-link>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </template>
          </div>
        </div>

        <!-- 中间轮播 -->
        <div class="center-panel">
          <el-carousel height="420px" indicator-position="outside" arrow="always" :interval="4000">
            <el-carousel-item v-for="banner in bannerList" :key="banner.id">
              <a :href="banner.url || '#'" class="carousel-link" target="_blank">
                <img :src="banner.image" :alt="'Banner ' + banner.index" class="carousel-img" />
              </a>
            </el-carousel-item>
          </el-carousel>
        </div>

        <!-- 右侧用户信息 + 快捷入口 -->
        <div class="right-panel">
          <div class="user-entry">
            <template v-if="userStore.isAuthenticated">
              <div class="user-avatar">
                <el-icon :size="30"><UserFilled /></el-icon>
              </div>
              <p class="user-greeting">Hi，{{ userStore.userName || '用户' }}</p>
            </template>
            <template v-else>
              <div class="user-avatar">
                <el-icon :size="30"><UserFilled /></el-icon>
              </div>
              <p class="user-greeting">Hi，欢迎光临</p>
              <div class="user-btns">
                <span class="btn-login" @click="openLogin">登录</span>
                <span class="btn-register" @click="openRegister">注册</span>
              </div>
            </template>
          </div>

          <div class="quick-links">
            <div class="quick-item" @click="$router.push('/order')">
              <el-icon :size="20"><Document /></el-icon>
              <span>我的订单</span>
            </div>
            <div class="quick-item" @click="$router.push('/cart')">
              <el-icon :size="20"><ShoppingCart /></el-icon>
              <span>购物车</span>
            </div>
            <div class="quick-item" @click="$router.push('/user')">
              <el-icon :size="20"><Star /></el-icon>
              <span>我的收藏</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 中部服务保障 -->
    <div class="features-section">
      <div class="container">
        <div class="features-grid">
          <div class="feature-item">
            <el-icon :size="32"><Shield /></el-icon>
            <div>
              <h3>品质保障</h3>
              <p>正品行货 品质护航</p>
            </div>
          </div>
          <div class="feature-item">
            <el-icon :size="32"><Truck /></el-icon>
            <div>
              <h3>极速物流</h3>
              <p>多仓直发 极速送达</p>
            </div>
          </div>
          <div class="feature-item">
            <el-icon :size="32"><CircleCheck /></el-icon>
            <div>
              <h3>售后无忧</h3>
              <p>7天无理由退换货</p>
            </div>
          </div>
          <div class="feature-item">
            <el-icon :size="32"><Headset /></el-icon>
            <div>
              <h3>贴心服务</h3>
              <p>7x24小时客服在线</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 下方功能入口 -->
    <div class="container">
      <div class="entry-grid">
        <div class="entry-card" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%)" @click="$router.push('/goods?ishot=1')">
          <div class="entry-icon">
            <el-icon :size="28"><TrendCharts /></el-icon>
          </div>
          <div class="entry-text">
            <div class="entry-title">热销爆品</div>
            <div class="entry-desc">大家都在买</div>
          </div>
        </div>
        <div class="entry-card" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)" @click="$router.push('/goods?isnew=1')">
          <div class="entry-icon">
            <el-icon :size="28"><Present /></el-icon>
          </div>
          <div class="entry-text">
            <div class="entry-title">新品首发</div>
            <div class="entry-desc">抢先体验</div>
          </div>
        </div>
        <div class="entry-card" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%)" @click="$router.push('/goods?istab=1')">
          <div class="entry-icon">
            <el-icon :size="28"><Goods /></el-icon>
          </div>
          <div class="entry-text">
            <div class="entry-title">推荐好物</div>
            <div class="entry-desc">精选品质</div>
          </div>
        </div>
        <div class="entry-card" style="background: linear-gradient(135deg, #fa709a 0%, #fee140 100%)" @click="$router.push('/goods')">
          <div class="entry-icon">
            <el-icon :size="28"><Discount /></el-icon>
          </div>
          <div class="entry-text">
            <div class="entry-title">限时特惠</div>
            <div class="entry-desc">超值抢购</div>
          </div>
        </div>
        <div class="entry-card" style="background: linear-gradient(135deg, #a18cd1 0%, #fbc2eb 100%)" @click="$router.push('/goods')">
          <div class="entry-icon">
            <el-icon :size="28"><More /></el-icon>
          </div>
          <div class="entry-text">
            <div class="entry-title">全部分类</div>
            <div class="entry-desc">海量商品</div>
          </div>
        </div>
      </div>
    </div>

    <!-- 新品上市 -->
    <div class="container goods-section" v-loading="loading">
      <div class="section-header">
        <h2 class="section-title">新品上市</h2>
        <router-link to="/goods?isnew=1" class="view-all">
          查看全部 <el-icon><ArrowRight /></el-icon>
        </router-link>
      </div>
      <div class="goods-grid">
        <div v-if="!loading && newGoods.length === 0" class="goods-empty">
          <el-icon :size="40"><GoodsFilled /></el-icon>
          <span>暂无新品，敬请期待</span>
        </div>
        <div
          v-for="goods in newGoods"
          :key="goods.id"
          class="goods-card"
          @click="$router.push(`/goods/${goods.id}`)"
        >
          <div class="goods-image-wrapper">
            <img :src="goods.front_image" :alt="goods.name" class="goods-img" />
            <div class="goods-tags">
              <span v-if="goods.is_new" class="tag tag-new">新品</span>
              <span v-if="goods.is_hot" class="tag tag-hot">热销</span>
            </div>
          </div>
          <div class="goods-info">
            <h3 class="goods-name">{{ goods.name }}</h3>
            <p class="goods-brief">{{ goods.goods_brief }}</p>
            <div class="goods-footer">
              <span class="price">¥{{ goods.shop_price }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 热门推荐 -->
    <div class="container goods-section">
      <div class="section-header">
        <h2 class="section-title">热销推荐</h2>
        <router-link to="/goods?ishot=1" class="view-all">
          查看全部 <el-icon><ArrowRight /></el-icon>
        </router-link>
      </div>
      <div class="goods-grid">
        <div v-if="!loading && hotGoods.length === 0" class="goods-empty">
          <el-icon :size="40"><GoodsFilled /></el-icon>
          <span>暂无热销商品</span>
        </div>
        <div
          v-for="goods in hotGoods"
          :key="goods.id"
          class="goods-card"
          @click="$router.push(`/goods/${goods.id}`)"
        >
          <div class="goods-image-wrapper">
            <img :src="goods.front_image" :alt="goods.name" class="goods-img" />
            <div class="goods-tags">
              <span v-if="goods.is_new" class="tag tag-new">新品</span>
              <span v-if="goods.is_hot" class="tag tag-hot">热销</span>
            </div>
          </div>
          <div class="goods-info">
            <h3 class="goods-name">{{ goods.name }}</h3>
            <p class="goods-brief">{{ goods.goods_brief }}</p>
            <div class="goods-footer">
              <span class="price">¥{{ goods.shop_price }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部 -->
    <footer class="app-footer">
      <div class="container">
        <div class="footer-content">
          <div class="footer-section">
            <h4>购物指南</h4>
            <ul>
              <li><a href="#">购物流程</a></li>
              <li><a href="#">会员介绍</a></li>
              <li><a href="#">常见问题</a></li>
            </ul>
          </div>
          <div class="footer-section">
            <h4>配送方式</h4>
            <ul>
              <li><a href="#">上门自提</a></li>
              <li><a href="#">快递运输</a></li>
              <li><a href="#">物流配送</a></li>
            </ul>
          </div>
          <div class="footer-section">
            <h4>售后服务</h4>
            <ul>
              <li><a href="#">退换货政策</a></li>
              <li><a href="#">退款说明</a></li>
              <li><a href="#">联系客服</a></li>
            </ul>
          </div>
          <div class="footer-section">
            <h4>联系我们</h4>
            <p>客服热线：400-888-8888</p>
            <p>工作时间：9:00-21:00</p>
          </div>
        </div>
        <div class="footer-bottom">
          <p>&copy; 2026 AtopMall. All rights reserved.</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getGoodsList, type GoodsItem } from '@/api/goods'
import { getCategoryList, getBannerList, type CategoryItem, type BannerItem } from '@/api/category'
import { useUserStore } from '@/store/user'
import { useAuthModal } from '@/composables/useAuthModal'
import AppHeader from '@/components/AppHeader.vue'

const userStore = useUserStore()
const { open } = useAuthModal()

const openLogin = () => open('login')
const openRegister = () => open('register')

const newGoods = ref<GoodsItem[]>([])
const hotGoods = ref<GoodsItem[]>([])
const categories = ref<CategoryItem[]>([])
const activeCategory = ref<number | null>(null)
const loading = ref(false)
const catLoading = ref(false)
const catError = ref(false)
const pageError = ref(false)
const bannerList = ref<BannerItem[]>([])

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

const loadData = async () => {
  loading.value = true
  pageError.value = false
  try {
    const [newRes, hotRes, bannerRes] = await Promise.all([
      getGoodsList({ isnew: 1, p: 1, pnum: 8 }).catch(() => ({ data: [] } as any)),
      getGoodsList({ ishot: 1, p: 1, pnum: 8 }).catch(() => ({ data: [] } as any)),
      getBannerList().catch(() => [] as BannerItem[]),
    ])
    newGoods.value = (newRes as any).data || []
    hotGoods.value = (hotRes as any).data || []
    bannerList.value = Array.isArray(bannerRes) ? bannerRes : []
    if (bannerList.value.length > 0) {
      bannerList.value.sort((a, b) => a.index - b.index)
    }
  } catch (error) {
    console.error('加载首页数据失败', error)
    pageError.value = true
  } finally {
    loading.value = false
  }
}

const retryAll = () => {
  loadData()
  loadCategories()
}

const handleCategoryEnter = (catId: number) => {
  activeCategory.value = catId
}

const handleCategoryLeave = () => {
  activeCategory.value = null
}

onMounted(() => {
  loadData()
  loadCategories()
})
</script>

<style lang="scss" scoped>
.home-page {
  background: $bg-page;
  min-height: 100vh;
}

// ======= 页面错误栏 =======
.page-error-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 20px;
  margin: 12px 0;
  background: #fef0f0;
  border: 1px solid #fde2e2;
  border-radius: 8px;

  .page-error-content {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #f56c6c;
    font-size: 14px;

    .el-icon {
      font-size: 18px;
    }
  }
}

// ======= Hero 区域 =======
.hero-section {
  padding: 20px 0;
  background: #fff;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.04);
}

.hero-layout {
  display: flex;
  gap: 16px;
  align-items: stretch;
}

// ======= 左侧分类 =======
.sidebar {
  width: 240px;
  flex-shrink: 0;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.06);
  overflow: visible;

  .sidebar-title {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 14px 18px;
    background: $primary-color;
    color: #fff;
    font-size: 15px;
    font-weight: 600;
  }

  .sidebar-list {
    padding: 8px 0;
    min-height: 200px;
  }

  .sidebar-item {
    position: relative;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 18px;
    cursor: pointer;
    font-size: 14px;
    color: $text-regular;
    transition: all 0.2s;

    &:hover,
    &.active {
      color: $primary-color;
      background: #ecf5ff;
    }

    .item-name {
      flex: 1;
    }
  }

  .sidebar-error,
  .sidebar-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 8px;
    padding: 40px 20px;
    color: #909399;
    font-size: 14px;
  }

  .category-dropdown {
    position: absolute;
    left: 100%;
    top: 0;
    width: 640px;
    min-height: 320px;
    background: #fff;
    border-radius: 0 10px 10px 10px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
    z-index: 100;
    padding: 20px;

    .dropdown-inner {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 20px;
    }

    .dropdown-column {
      .column-title {
        display: block;
        font-size: 14px;
        font-weight: 600;
        color: $text-primary;
        margin-bottom: 10px;
        padding-bottom: 8px;
        border-bottom: 1px solid $border-lighter;
        text-decoration: none;

        &:hover {
          color: $primary-color;
        }
      }

      .column-links {
        display: flex;
        flex-wrap: wrap;
        gap: 6px;

        .link-item {
          font-size: 12px;
          color: $text-secondary;
          padding: 3px 8px;
          border-radius: 4px;
          text-decoration: none;
          transition: all 0.2s;

          &:hover {
            color: $primary-color;
            background: #ecf5ff;
          }
        }
      }
    }
  }
}

// ======= 中间轮播 =======
.center-panel {
  flex: 1;
  min-width: 0;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.06);

  .carousel-link {
    display: block;
    width: 100%;
    height: 100%;
  }

  .carousel-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }
}

// ======= 右侧面板 =======
.right-panel {
  width: 240px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.user-entry {
  background: #fff;
  border-radius: 10px;
  padding: 24px 16px;
  text-align: center;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.06);

  .user-avatar {
    width: 60px;
    height: 60px;
    margin: 0 auto 14px;
    border-radius: 50%;
    background: linear-gradient(135deg, #ecf5ff 0%, #d9ecff 100%);
    display: flex;
    align-items: center;
    justify-content: center;
    color: $primary-color;
    box-shadow: 0 2px 8px rgba(64, 158, 255, 0.15);
  }

  .user-greeting {
    font-size: 14px;
    color: $text-primary;
    font-weight: 500;
  }

  .user-btns {
    display: flex;
    gap: 12px;
    justify-content: center;

    a,
    span {
      display: inline-block;
      padding: 6px 20px;
      border-radius: 6px;
      font-size: 13px;
      text-decoration: none;
      cursor: pointer;
      transition: all 0.2s;
    }

    .btn-login {
      background: $primary-color;
      color: #fff;

      &:hover {
        opacity: 0.85;
      }
    }

    .btn-register {
      background: #f5f7fa;
      color: $text-primary;

      &:hover {
        background: #ecf5ff;
        color: $primary-color;
      }
    }
  }
}

.quick-links {
  background: #fff;
  border-radius: 10px;
  padding: 14px 16px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.06);
  display: flex;
  justify-content: space-around;

  .quick-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 6px;
    cursor: pointer;
    color: $text-regular;
    font-size: 12px;
    transition: all 0.3s;
    padding: 4px 8px;
    border-radius: 8px;

    &:hover {
      color: $primary-color;
      background: #ecf5ff;
      transform: translateY(-2px);
    }
  }
}

// ======= 服务保障 =======
.features-section {
  background: #fff;
  padding: 32px 0;
  margin-top: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.04);

  .features-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 0;
  }

  .feature-item {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 14px;
    padding: 18px 24px;
    transition: all 0.3s;
    border-right: 1px solid #f0f0f0;

    &:last-child {
      border-right: none;
    }

    &:hover {
      background: #f8fafc;
      transform: translateY(-2px);

      .el-icon {
        transform: scale(1.1);
      }
    }

    .el-icon {
      color: $primary-color;
      flex-shrink: 0;
      transition: transform 0.3s;
    }

    h3 {
      font-size: 16px;
      margin-bottom: 4px;
      color: $text-primary;
    }

    p {
      font-size: 13px;
      color: $text-secondary;
      margin: 0;
    }
  }
}

// ======= 入口卡片 =======
.entry-grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 16px;
  padding: 28px 0;
}

.entry-card {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 20px 20px;
  border-radius: 10px;
  color: #fff;
  cursor: pointer;
  transition: all 0.35s cubic-bezier(0.25, 0.46, 0.45, 0.94);

  &:hover {
    transform: translateY(-6px);
    box-shadow: 0 12px 28px rgba(0, 0, 0, 0.2);
  }

  .entry-icon {
    flex-shrink: 0;
    display: flex;
    align-items: center;
  }

  .entry-text {
    flex: 1;
  }

  .entry-title {
    font-size: 17px;
    font-weight: 600;
    margin-bottom: 4px;
  }

  .entry-desc {
    font-size: 13px;
    opacity: 0.85;
  }
}

// ======= 商品区 =======
.goods-section {
  padding: 42px 0;

  .goods-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 20px;
    min-height: 200px;

    .goods-empty {
      grid-column: 1 / -1;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 10px;
      padding: 60px 20px;
      color: #909399;
      font-size: 14px;

      .el-icon {
        font-size: 40px;
      }
    }
  }

  .goods-card {
    background: #fff;
    border-radius: 10px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.35s cubic-bezier(0.25, 0.46, 0.45, 0.94);
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.06);

    &:hover {
      transform: translateY(-8px);
      box-shadow: 0 16px 32px rgba(0, 0, 0, 0.12);
    }

    .goods-image-wrapper {
      position: relative;
      overflow: hidden;
      aspect-ratio: 1 / 1;

      .goods-img {
        width: 100%;
        height: 100%;
        object-fit: cover;
        transition: transform 0.4s;
      }

      &:hover .goods-img {
        transform: scale(1.1);
      }

      .goods-tags {
        position: absolute;
        top: 12px;
        left: 12px;
        display: flex;
        gap: 6px;

        .tag {
          padding: 4px 10px;
          border-radius: 4px;
          font-size: 12px;
          font-weight: 600;
          color: #fff;

          &.tag-new {
            background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
          }

          &.tag-hot {
            background: linear-gradient(135deg, #ff6b6b 0%, #ee5a6f 100%);
          }
        }
      }
    }

    .goods-info {
      padding: 16px;

      .goods-name {
        font-size: 15px;
        font-weight: 500;
        margin-bottom: 8px;
        line-height: 1.4;
        color: $text-primary;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
      }

      .goods-brief {
        font-size: 13px;
        color: $text-secondary;
        margin-bottom: 12px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .goods-footer {
        display: flex;
        align-items: baseline;
        gap: 8px;

        .price {
          font-size: 22px;
          font-weight: bold;
          color: $danger-color;
        }

        .market-price {
          font-size: 13px;
          color: $text-secondary;
          text-decoration: line-through;
        }
      }
    }
  }
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 28px;
  padding-bottom: 16px;
  border-bottom: 1px solid $border-lighter;

  .section-title {
    font-size: 24px;
    font-weight: bold;
    color: $text-primary;
    position: relative;
    padding-left: 16px;

    &::before {
      content: '';
      position: absolute;
      left: 0;
      top: 50%;
      transform: translateY(-50%);
      width: 4px;
      height: 22px;
      background: $primary-color;
      border-radius: 2px;
    }
  }

  .view-all {
    color: $text-secondary;
    text-decoration: none;
    font-size: 14px;
    transition: all 0.2s;

    &:hover {
      color: $primary-color;
      transform: translateX(4px);
    }
  }
}

// ======= 底部 =======
.app-footer {
  background: #2c3e50;
  color: #ecf0f1;
  padding: 40px 0 20px;
  margin-top: 60px;

  .footer-content {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 32px;
    margin-bottom: 32px;
  }

  .footer-section {
    h4 {
      font-size: 16px;
      font-weight: bold;
      margin-bottom: 16px;
      color: #fff;
    }

    p {
      font-size: 14px;
      line-height: 1.8;
      color: #bdc3c7;
    }

    ul {
      list-style: none;
      padding: 0;

      li {
        margin-bottom: 8px;

        a {
          color: #bdc3c7;
          text-decoration: none;
          font-size: 14px;

          &:hover {
            color: #fff;
          }
        }
      }
    }
  }

  .footer-bottom {
    text-align: center;
    padding-top: 20px;
    border-top: 1px solid #34495e;

    p {
      font-size: 13px;
      color: #95a5a6;
    }
  }
}
</style>