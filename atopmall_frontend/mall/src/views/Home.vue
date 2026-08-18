<template>
  <div class="home-page">

    <div class="container" v-if="pageError">
      <div class="page-error-bar">
        <div class="page-error-content">
          <el-icon><WarningFilled /></el-icon>
          <span>首页数据加载失败，请检查网络后重试</span>
        </div>
        <el-button type="primary" size="small" @click="retryAll">重新加载</el-button>
      </div>
    </div>

    <!-- ==================== 大 Banner 区域 ==================== -->
    <div class="hero-wrapper">
      <!-- 轮播图背景层 -->
      <div class="hero-carousel-bg">
        <el-carousel
          v-if="bannerList.length > 0"
          height="520px"
          trigger="click"
          :interval="5000"
          arrow="hover"
          :autoplay="true"
        >
          <el-carousel-item v-for="(banner, idx) in bannerList" :key="banner.id">
            <a :href="banner.url || '#'" class="carousel-link" target="_blank">
              <img :src="banner.image" :alt="'Banner ' + (idx + 1)" class="carousel-img" />
            </a>
          </el-carousel-item>
        </el-carousel>
        <div v-else class="carousel-placeholder">
          <el-icon :size="48"><PictureFilled /></el-icon>
          <p>暂无轮播图，请前往管理后台添加</p>
        </div>
      </div>

      <!-- 浮于轮播图右侧的用户信息卡片 -->
      <div class="hero-user-card">
        <template v-if="userStore.isAuthenticated">
          <router-link to="/user" class="huc-avatar-link">
            <div class="huc-avatar">
              <el-icon :size="24"><UserFilled /></el-icon>
            </div>
            <span class="huc-name">{{ userStore.userName }}</span>
          </router-link>
        </template>
        <template v-else>
          <div class="huc-avatar-guest" @click="openLogin">
            <el-icon :size="24"><UserFilled /></el-icon>
          </div>
          <span class="huc-welcome">Hi，欢迎</span>
          <span class="huc-sub-welcome">来到AtopMall</span>
        </template>
        <div class="huc-auth-row" v-if="!userStore.isAuthenticated">
          <span class="huc-auth-link" @click="openLogin">登录</span>
          <span class="huc-divider">|</span>
          <span class="huc-auth-link" @click="openRegister">注册</span>
        </div>
        <div class="huc-actions">
          <a
            v-if="userStore.isAuthenticated"
            class="huc-action-btn"
            @click="$router.push('/user/orders')"
          >
            <el-icon :size="15"><Document /></el-icon>
            <span>我的订单</span>
          </a>
          <a v-else class="huc-action-btn" @click="openLogin">
            <el-icon :size="15"><Document /></el-icon>
            <span>我的订单</span>
          </a>
          <a
            v-if="userStore.isAuthenticated"
            class="huc-action-btn"
            @click="$router.push('/user/favorite')"
          >
            <el-icon :size="15"><Star /></el-icon>
            <span>我的收藏</span>
          </a>
          <a v-else class="huc-action-btn" @click="openLogin">
            <el-icon :size="15"><Star /></el-icon>
            <span>我的收藏</span>
          </a>
          <a
            v-if="userStore.isAuthenticated"
            class="huc-action-btn"
            @click="$router.push('/user/address')"
          >
            <el-icon :size="15"><Location /></el-icon>
            <span>地址管理</span>
          </a>
          <a v-else class="huc-action-btn" @click="openLogin">
            <el-icon :size="15"><Location /></el-icon>
            <span>地址管理</span>
          </a>
          <a
            v-if="userStore.isAuthenticated"
            class="huc-action-btn huc-action-cart"
            @click="$router.push('/user/cart')"
          >
            <el-icon :size="15"><ShoppingCart /></el-icon>
            <span>购物车</span>
            <el-badge :value="cartStore.totalCount" :hidden="cartStore.totalCount === 0" class="huc-cart-badge" />
          </a>
          <a v-else class="huc-action-btn huc-action-cart" @click="openLogin">
            <el-icon :size="15"><ShoppingCart /></el-icon>
            <span>购物车</span>
            <el-badge :value="cartStore.totalCount" :hidden="cartStore.totalCount === 0" class="huc-cart-badge" />
          </a>
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
              <span class="price">{{ goods.shop_price }}</span>
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
              <span class="price">{{ goods.shop_price }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部 -->
    <footer class="app-footer">
      <div class="container">
        <div class="features-row">
          <div class="feature-item">
            <el-icon :size="28"><Shield /></el-icon>
            <div>
              <h3>品质保障</h3>
              <p>正品行货 品质护航</p>
            </div>
          </div>
          <div class="feature-item">
            <el-icon :size="28"><Truck /></el-icon>
            <div>
              <h3>极速物流</h3>
              <p>多仓直发 极速送达</p>
            </div>
          </div>
          <div class="feature-item">
            <el-icon :size="28"><CircleCheck /></el-icon>
            <div>
              <h3>售后无忧</h3>
              <p>7天无理由退换货</p>
            </div>
          </div>
          <div class="feature-item">
            <el-icon :size="28"><Headset /></el-icon>
            <div>
              <h3>贴心服务</h3>
              <p>7x24小时客服在线</p>
            </div>
          </div>
        </div>
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
import { getBannerList, type BannerItem } from '@/api/category'
import { useUserStore } from '@/store/user'
import { useCartStore } from '@/store/cart'
import { useAuthModal } from '@/composables/useAuthModal'


const userStore = useUserStore()
const cartStore = useCartStore()
const { open } = useAuthModal()

const openLogin = () => open('login')
const openRegister = () => open('register')

const newGoods = ref<GoodsItem[]>([])
const hotGoods = ref<GoodsItem[]>([])
const loading = ref(false)
const pageError = ref(false)
const bannerList = ref<BannerItem[]>([])

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
}

onMounted(() => {
  loadData()
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
  margin: 0 0 20px 0;
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

// ==================== Banner 大区域 ====================
.hero-wrapper {
  position: relative;
  z-index: 0;
  height: 520px;
  overflow: visible;
}

// 轮播图背景层
.hero-carousel-bg {
  position: absolute;
  inset: 0;
  z-index: 0;
  overflow: hidden;

  .carousel-img {
    width: 100%;
    height: 100%;
    object-fit: cover;
    object-position: center;
    display: block;
  }

  .carousel-link {
    display: block;
    width: 100%;
    height: 100%;
  }

  .carousel-placeholder {
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    background: linear-gradient(135deg, #f5f7fa 0%, #ebeef5 100%);
    color: #909399;
    gap: 14px;
    font-size: 14px;
  }

  :deep(.el-carousel),
  :deep(.el-carousel__container) {
    height: 100% !important;
  }

  :deep(.el-carousel__indicators) {
    bottom: 20px;
    left: 50%;
    transform: translateX(-50%);

    .el-carousel__indicator {
      padding: 4px;

      .el-carousel__button {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: rgba(255, 255, 255, 0.5);
        opacity: 1;
        transition: all 0.35s ease;
      }

      &.is-active .el-carousel__button {
        width: 24px;
        border-radius: 4px;
        background: #fff;
        box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2);
      }
    }
  }

  :deep(.el-carousel__arrow) {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background: rgba(0, 0, 0, 0.12);
    color: #fff;
    font-size: 14px;
    transition: all 0.3s;
    opacity: 0;
    z-index: 5;

    &:hover {
      background: rgba(0, 0, 0, 0.25);
    }
  }

  &:hover :deep(.el-carousel__arrow) {
    opacity: 1;
  }
}

// ======= 浮于轮播图右侧的用户信息卡片 =======
.hero-user-card {
  position: absolute;
  right: 6%;
  top: 50%; 
  transform: translateY(-50%); 
  z-index: 5; 
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 20px;          /* 增大内部元素间距 */
  padding: 28px 24px; /* 上下、左右内边距，横向留白加宽 */
  width: 280px;       /* 280px，解决过窄问题 */
  background: rgba(255, 255, 255, 0.94);
  backdrop-filter: blur(12px);
  border-radius: 18px; 
  box-shadow: 0 8px 36px rgba(0, 0, 0, 0.13); 
  border: 1px solid rgba(255, 255, 255, 0.7);
  max-width: 90vw;    /* 小屏幕防溢出，响应式兜底 */
  box-sizing: border-box;
}

.huc-avatar-link {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  flex-shrink: 0;

  &:hover .huc-name {
    color: $primary-color;
  }
}

.huc-avatar {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, $primary-color, #66b1ff);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.45);
  transition: transform 0.2s;

  &:hover {
    transform: scale(1.08);
  }
}

.huc-name {
  font-size: 15px;
  font-weight: 600;
  color: $text-primary;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  transition: color 0.2s;
}

.huc-avatar-guest {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  background: linear-gradient(135deg, $primary-color, #66b1ff);
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  flex-shrink: 0;
  cursor: pointer;
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.45);
  transition: transform 0.2s;

  &:hover {
    transform: scale(1.08);
  }
}

.huc-welcome {
  font-size: 14px;
  font-weight: 600;
  color: $text-primary;
  white-space: nowrap;
  margin-top: 2px;
}

.huc-sub-welcome {
  font-size: 12px;
  color: $text-secondary;
  white-space: nowrap;
}

.huc-auth-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.huc-auth-link {
  font-size: 15px;
  color: $primary-color;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s;
  white-space: nowrap;

  &:hover {
    color: #398ee5;
    text-decoration: underline;
  }
}

.huc-divider {
  color: #dcdfe6;
  font-size: 15px;
  user-select: none;
}

.huc-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
  width: 100%;
}

.huc-action-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 9px 6px;
  background: #f5f7fa;
  border-radius: 22px;
  color: $text-regular;
  text-decoration: none;
  font-size: 12px;
  white-space: nowrap;
  cursor: pointer;
  transition: all 0.25s;
  border: 1px solid transparent;
  font-weight: 500;

  &:hover {
    color: $primary-color;
    background: #ecf5ff;
    border-color: rgba(64, 158, 255, 0.25);
    box-shadow: 0 2px 8px rgba(64, 158, 255, 0.12);
  }

  &.huc-action-cart {
    position: relative;
  }

  .huc-cart-badge {
    position: absolute;
    top: 1px;
    right: 8px;
  }
}

// ======= 商品区 =======
.goods-section {
  padding: 40px 0;

  .goods-grid {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 18px;
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
    }
  }

  .goods-card {
    background: #fff;
    border-radius: 10px;
    overflow: hidden;
    cursor: pointer;
    transition: all 0.3s cubic-bezier(0.25, 0.46, 0.45, 0.94);
    box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);

    &:hover {
      transform: translateY(-6px);
      box-shadow: 0 12px 28px rgba(0, 0, 0, 0.12);
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
        transform: scale(1.08);
      }

      .goods-tags {
        position: absolute;
        top: 10px;
        left: 10px;
        display: flex;
        gap: 6px;

        .tag {
          padding: 3px 8px;
          border-radius: 4px;
          font-size: 11px;
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
      padding: 14px;

      .goods-name {
        font-size: 14px;
        font-weight: 500;
        margin-bottom: 6px;
        line-height: 1.4;
        color: $text-primary;
        overflow: hidden;
        text-overflow: ellipsis;
        display: -webkit-box;
        -webkit-line-clamp: 2;
        -webkit-box-orient: vertical;
      }

      .goods-brief {
        font-size: 12px;
        color: $text-secondary;
        margin-bottom: 10px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .goods-footer {
        display: flex;
        align-items: baseline;
        gap: 8px;

        .price {
          font-size: 20px;
          font-weight: bold;
          color: $danger-color;
        }

        .market-price {
          font-size: 12px;
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
  margin-bottom: 24px;
  padding-bottom: 14px;
  border-bottom: 1px solid $border-lighter;

  .section-title {
    font-size: 22px;
    font-weight: bold;
    color: $text-primary;
    position: relative;
    padding-left: 14px;

    &::before {
      content: '';
      position: absolute;
      left: 0;
      top: 50%;
      transform: translateY(-50%);
      width: 4px;
      height: 20px;
      background: $primary-color;
      border-radius: 2px;
    }
  }

  .view-all {
    display: flex;
    align-items: center;
    gap: 4px;
    color: $text-secondary;
    text-decoration: none;
    font-size: 13px;
    transition: all 0.2s;

    &:hover {
      color: $primary-color;
      transform: translateX(3px);
    }
  }
}

// ======= 底部 =======
.app-footer {
  background: #2c3e50;
  color: #ecf0f1;
  padding: 48px 0 24px;
  margin-top: 60px;

  .features-row {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 40px;
    padding: 24px 0;
    margin-bottom: 36px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
  }

  .feature-item {
    display: flex;
    align-items: center;
    gap: 12px;
    transition: all 0.3s;

    &:hover {
      transform: translateY(-2px);

      .el-icon {
        opacity: 1;
      }
    }

    .el-icon {
      color: #409eff;
      flex-shrink: 0;
      opacity: 0.8;
      transition: all 0.3s;
    }

    h3 {
      font-size: 14px;
      margin-bottom: 2px;
      color: #ecf0f1;
      font-weight: 600;
      white-space: nowrap;
    }

    p {
      font-size: 12px;
      color: rgba(255, 255, 255, 0.55);
      margin: 0;
      white-space: nowrap;
    }
  }

  .footer-content {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 40px;
    margin-bottom: 36px;
  }

  .footer-section {
    h4 {
      font-size: 16px;
      font-weight: 600;
      margin-bottom: 18px;
      color: #fff;
      position: relative;
      padding-bottom: 10px;

      &::after {
        content: '';
        position: absolute;
        bottom: 0;
        left: 0;
        width: 28px;
        height: 2px;
        background: $primary-color;
        border-radius: 1px;
      }
    }

    p {
      font-size: 13px;
      line-height: 2;
      color: #bdc3c7;
    }

    ul {
      list-style: none;
      padding: 0;

      li {
        margin-bottom: 10px;

        a {
          color: #bdc3c7;
          text-decoration: none;
          font-size: 13px;
          transition: color 0.2s;

          &:hover {
            color: $primary-color;
          }
        }
      }
    }
  }

  .footer-bottom {
    text-align: center;
    padding-top: 20px;
    border-top: 1px solid rgba(255, 255, 255, 0.08);

    p {
      font-size: 12px;
      color: #7f8c8d;
    }
  }
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes carouselFadeIn {
  from { opacity: 0.6; }
  to { opacity: 1; }
}
</style>