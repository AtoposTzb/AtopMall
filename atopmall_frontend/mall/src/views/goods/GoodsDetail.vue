<template>
  <div class="goods-detail-page" v-loading="loading">
    <div class="container" v-if="goods">
      <!-- 面包屑导航 -->
      <div class="breadcrumb-wrapper">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ path: '/goods' }">全部商品</el-breadcrumb-item>
          <el-breadcrumb-item>{{ goods.name }}</el-breadcrumb-item>
        </el-breadcrumb>
      </div>

      <!-- 商品主体 -->
      <div class="product-card card">
        <div class="product-body">
          <!-- 左侧：图片区 -->
          <div class="gallery" @mouseenter="pauseAutoPlay" @mouseleave="resumeAutoPlay">
            <div class="main-image-box">
              <transition name="fade" mode="out-in">
                <img :src="currentImage" :key="currentImage" class="main-image" />
              </transition>
              <div class="badges">
                <span v-if="goods.is_hot" class="badge hot">热销</span>
                <span v-if="goods.is_new" class="badge new">新品</span>
              </div>
              <div class="nav-arrows" v-if="goods.images?.length && goods.images.length > 1">
                <button class="nav-arrow left" @click.stop="prevImage">
                  <el-icon><ArrowLeft /></el-icon>
                </button>
                <button class="nav-arrow right" @click.stop="nextImage">
                  <el-icon><ArrowRight /></el-icon>
                </button>
              </div>
              <div class="image-dots" v-if="goods.images?.length && goods.images.length > 1">
                <span
                  v-for="(img, idx) in goods.images"
                  :key="idx"
                  class="dot"
                  :class="{ active: currentImage === img }"
                  @click.stop="selectImage(idx)"
                ></span>
              </div>
            </div>
            <div class="thumb-strip" v-if="goods.images?.length">
              <div
                v-for="(img, index) in goods.images"
                :key="index"
                class="thumb-item"
                :class="{ active: currentImage === img }"
                @click="selectImage(index)"
              >
                <img :src="img" class="thumb-img" />
              </div>
            </div>
          </div>

          <!-- 右侧：信息区 -->
          <div class="info">
            <h1 class="product-name">{{ goods.name }}</h1>
            <p class="product-brief" v-if="goods.goods_brief">{{ goods.goods_brief }}</p>

            <!-- 价格卡片 -->
            <div class="price-card">
              <div class="price-row">
                <span class="price-label">售价</span>
                <span class="price-currency">¥</span>
                <span class="price-number">{{ goods.shop_price }}</span>
              </div>
            </div>

            <!-- 服务标签 -->
            <div class="service-row">
              <span v-if="goods.ship_free" class="service-tag">
                <el-icon><Van /></el-icon> 包邮
              </span>
              <span class="service-tag">
                <el-icon><CircleCheck /></el-icon> 正品保障
              </span>
              <span class="service-tag">
                <el-icon><Clock /></el-icon> 48小时发货
              </span>
            </div>

            <div class="divider"></div>

            <!-- 规格信息 -->
            <div class="spec-row">
              <div class="spec-item" v-if="goods.brand?.name">
                <span class="spec-label">品牌</span>
                <span class="spec-value">{{ goods.brand.name }}</span>
              </div>
              <div class="spec-item" v-if="goods.ctegory?.name">
                <span class="spec-label">分类</span>
                <span class="spec-value">{{ goods.ctegory.name }}</span>
              </div>
            </div>

            <div class="divider"></div>

            <!-- 数量选择 -->
            <div class="quantity-row">
              <span class="quantity-label">数量</span>
              <el-input-number
                v-model="quantity"
                :min="1"
                :max="99"
                size="large"
                class="quantity-input"
              />
              <span class="quantity-hint" v-if="quantity > 1">
                小计：<strong>¥{{ (Number(goods.shop_price) * quantity).toFixed(2) }}</strong>
              </span>
            </div>

            <!-- 操作按钮 -->
            <div class="action-row">
              <button
                class="btn btn-cart"
                @click="handleAddCart"
                :disabled="addCartLoading"
              >
                <el-icon v-if="!addCartLoading"><ShoppingCart /></el-icon>
                <span v-if="addCartLoading" class="loading-spinner"></span>
                {{ addCartLoading ? '添加中...' : '加入购物车' }}
              </button>
              <button class="btn btn-buy" @click="handleBuyNow">
                立即购买
              </button>
              <button
                class="btn btn-fav"
                :class="{ favorited: isFaved }"
                @click="handleToggleFav"
              >
                <el-icon><Star /></el-icon>
                <span>{{ isFaved ? '已收藏' : '收藏' }}</span>
              </button>
            </div>

            <!-- 购物车入口 -->
            <div class="cart-entry" v-if="cartStore.totalCount > 0">
              <router-link to="/cart" class="cart-entry-link">
                <el-icon><ShoppingCart /></el-icon>
                <span>购物车 ({{ cartStore.totalCount }} 件)</span>
                <el-icon class="arrow"><ArrowRight /></el-icon>
              </router-link>
            </div>
          </div>
        </div>
      </div>

      <!-- 商品详情 -->
      <div class="desc-card card">
        <div class="desc-header">
          <h2>商品详情</h2>
          <div class="desc-divider"></div>
        </div>
        <div class="desc-content" v-if="goods.desc">
          <div class="desc-text" v-html="goods.desc"></div>
        </div>
        <div class="desc-images" v-if="goods.desc_images?.length">
          <img
            v-for="(img, index) in goods.desc_images"
            :key="index"
            :src="img"
            class="desc-img"
            loading="lazy"
          />
        </div>
        <div class="desc-empty" v-if="!goods.desc && !goods.desc_images?.length">
          <el-icon><Document /></el-icon>
          <p>暂无详情描述</p>
        </div>
      </div>
    </div>

    <!-- 回到顶部 -->
    <div class="back-to-top" v-show="showBackTop" @click="scrollToTop">
      <el-icon :size="22"><ArrowUp /></el-icon>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { getGoodsDetail, type GoodsItem } from "@/api/goods";
import { useCartStore } from "@/store/cart";
import { useUserStore } from "@/store/user";
import { useAuthModal } from "@/composables/useAuthModal";
import { addFav, deleteFav, getFavStatus } from "@/api/favorite";

const route = useRoute();
const router = useRouter();
const cartStore = useCartStore();
const userStore = useUserStore();
const { open } = useAuthModal();

const goods = ref<GoodsItem | null>(null);
const currentImage = ref("");
const imageIndex = ref(0);
const quantity = ref(1);
const loading = ref(false);
const addCartLoading = ref(false);
const isFaved = ref(false);
const showBackTop = ref(false);
let autoPlayTimer: ReturnType<typeof setInterval> | null = null;
let isPaused = false;

const goodsId = Number(route.params.id);

const imageList = () => goods.value?.images || [];

const selectImage = (index: number) => {
  const list = imageList();
  if (list.length > 0) {
    imageIndex.value = index;
    currentImage.value = list[index];
  }
};

const nextImage = () => {
  const list = imageList();
  if (list.length <= 1) return;
  imageIndex.value = (imageIndex.value + 1) % list.length;
  currentImage.value = list[imageIndex.value];
};

const prevImage = () => {
  const list = imageList();
  if (list.length <= 1) return;
  imageIndex.value = (imageIndex.value - 1 + list.length) % list.length;
  currentImage.value = list[imageIndex.value];
};

const startAutoPlay = () => {
  stopAutoPlay();
  autoPlayTimer = setInterval(() => {
    if (!isPaused) {
      nextImage();
    }
  }, 3000);
};

const stopAutoPlay = () => {
  if (autoPlayTimer) {
    clearInterval(autoPlayTimer);
    autoPlayTimer = null;
  }
};

const pauseAutoPlay = () => {
  isPaused = true;
};

const resumeAutoPlay = () => {
  isPaused = false;
};

const handleScroll = () => {
  showBackTop.value = window.scrollY > 400;
};

const scrollToTop = () => {
  window.scrollTo({ top: 0, behavior: "smooth" });
};

const handleAddCart = async () => {
  if (!userStore.isAuthenticated) {
    open("login");
    return;
  }
  addCartLoading.value = true;
  try {
    await cartStore.addCartItem(goodsId, quantity.value);
  } finally {
    addCartLoading.value = false;
  }
};

const handleBuyNow = async () => {
  if (!userStore.isAuthenticated) {
    open("login");
    return;
  }
  const existing = cartStore.cartList.find(
    (item) => item.goods_id === goodsId
  );
  if (!existing) {
    await cartStore.addCartItem(goodsId, quantity.value);
  } else {
    ElMessage.info("该商品已在购物车中");
  }
  router.push("/cart");
};

const handleToggleFav = async () => {
  if (!userStore.isAuthenticated) {
    open("login");
    return;
  }
  try {
    if (isFaved.value) {
      await deleteFav(goodsId);
      isFaved.value = false;
      ElMessage.success("已取消收藏");
    } else {
      await addFav(goodsId);
      isFaved.value = true;
      ElMessage.success("收藏成功");
    }
  } catch (error) {
    // 错误已在拦截器中处理
  }
};

const checkFavStatus = async () => {
  if (!userStore.isAuthenticated) return;
  try {
    await getFavStatus(goodsId);
    isFaved.value = true;
  } catch {
    isFaved.value = false;
  }
};

onMounted(async () => {
  window.addEventListener("scroll", handleScroll);
  loading.value = true;
  try {
    const goodsRes = await getGoodsDetail(goodsId);
    goods.value = goodsRes as unknown as GoodsItem;
    const list = imageList();
    if (list.length > 0) {
      currentImage.value = goods.value.front_image || list[0];
      imageIndex.value = list.findIndex((img: string) => img === currentImage.value);
      if (imageIndex.value < 0) imageIndex.value = 0;
    }
    startAutoPlay();
  } catch (error) {
    console.error("加载商品详情失败", error);
    ElMessage.error("商品不存在或已下架");
  } finally {
    loading.value = false;
  }

  checkFavStatus();
});

onUnmounted(() => {
  window.removeEventListener("scroll", handleScroll);
  stopAutoPlay();
});
</script>

<style lang="scss" scoped>
// ========== 面包屑 ==========
.breadcrumb-wrapper {
  margin: 20px 0;
}

// ========== 商品主体卡片 ==========
.product-card {
  padding: 32px;

  .product-body {
    display: flex;
    gap: 48px;
  }
}

// ========== 左侧：图片区 ==========
.gallery {
  width: 480px;
  flex-shrink: 0;

  .main-image-box {
    position: relative;
    width: 480px;
    height: 480px;
    border-radius: 12px;
    overflow: hidden;
    background: #f9fafb;
    border: 1px solid #f0f0f0;
    cursor: pointer;

    .main-image {
      width: 100%;
      height: 100%;
      object-fit: cover;
      position: absolute;
      top: 0;
      left: 0;
    }

    .badges {
      position: absolute;
      top: 16px;
      left: 16px;
      display: flex;
      gap: 8px;
      z-index: 2;

      .badge {
        padding: 5px 14px;
        border-radius: 20px;
        font-size: 12px;
        font-weight: 600;
        color: #fff;
        letter-spacing: 0.5px;

        &.hot {
          background: linear-gradient(135deg, #ff6b6b, #ee5a6f);
        }

        &.new {
          background: linear-gradient(135deg, #4facfe, #00f2fe);
        }
      }
    }

    .nav-arrows {
      position: absolute;
      top: 50%;
      left: 0;
      right: 0;
      transform: translateY(-50%);
      display: flex;
      justify-content: space-between;
      padding: 0 12px;
      pointer-events: none;
      z-index: 3;
      opacity: 0;
      transition: opacity 0.3s;

      .nav-arrow {
        width: 36px;
        height: 36px;
        border-radius: 50%;
        background: rgba(255, 255, 255, 0.85);
        border: none;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        pointer-events: auto;
        color: #333;
        font-size: 16px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
        transition: all 0.2s;

        &:hover {
          background: #fff;
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        }
      }
    }

    &:hover .nav-arrows {
      opacity: 1;
    }

    .image-dots {
      position: absolute;
      bottom: 16px;
      left: 50%;
      transform: translateX(-50%);
      display: flex;
      gap: 8px;
      z-index: 2;

      .dot {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: rgba(255, 255, 255, 0.5);
        cursor: pointer;
        transition: all 0.3s;

        &.active {
          background: #fff;
          width: 24px;
          border-radius: 4px;
        }
      }
    }
  }

  .thumb-strip {
    display: flex;
    gap: 10px;
    margin-top: 16px;

    .thumb-item {
      width: 72px;
      height: 72px;
      border-radius: 8px;
      overflow: hidden;
      border: 2px solid transparent;
      cursor: pointer;
      transition: all 0.25s ease;
      box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
      }

      &.active {
        border-color: $primary-color;
        box-shadow: 0 0 0 2px rgba($primary-color, 0.2);
      }

      .thumb-img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }
  }
}

// ========== 右侧：信息区 ==========
.info {
  flex: 1;
  min-width: 0;

  .product-name {
    font-size: 24px;
    font-weight: 700;
    line-height: 1.4;
    color: #1a1a1a;
    margin-bottom: 8px;
  }

  .product-brief {
    font-size: 14px;
    color: #999;
    line-height: 1.6;
    margin-bottom: 20px;
  }

  // 价格卡片
  .price-card {
    background: linear-gradient(135deg, #fff5f5 0%, #ffe8e8 100%);
    border-radius: 12px;
    padding: 20px 24px;
    margin-bottom: 16px;

    .price-row {
      display: flex;
      align-items: baseline;
      gap: 6px;

      .price-label {
        font-size: 13px;
        color: #999;
        margin-right: 8px;
      }

      .price-currency {
        font-size: 22px;
        font-weight: 700;
        color: #e4393c;
      }

      .price-number {
        font-size: 38px;
        font-weight: 800;
        color: #e4393c;
        line-height: 1;
        letter-spacing: -1px;
      }
    }
  }

  // 服务标签
  .service-row {
    display: flex;
    gap: 12px;
    flex-wrap: wrap;
    margin-bottom: 16px;

    .service-tag {
      display: inline-flex;
      align-items: center;
      gap: 4px;
      font-size: 12px;
      color: #666;
      padding: 4px 10px;
      background: #f8f9fa;
      border-radius: 4px;
      border: 1px solid #eee;
    }
  }

  .divider {
    height: 1px;
    background: #f0f0f0;
    margin: 16px 0;
  }

  // 规格信息
  .spec-row {
    display: flex;
    flex-direction: column;
    gap: 12px;

    .spec-item {
      display: flex;
      align-items: center;
      font-size: 14px;

      .spec-label {
        color: #999;
        width: 48px;
        flex-shrink: 0;
      }

      .spec-value {
        color: #333;
        font-weight: 500;
      }
    }
  }

  // 数量选择
  .quantity-row {
    display: flex;
    align-items: center;
    gap: 16px;
    margin-bottom: 28px;

    .quantity-label {
      font-size: 14px;
      color: #999;
    }

    .quantity-input {
      :deep(.el-input-number__decrease),
      :deep(.el-input-number__increase) {
        border-radius: 0;
      }

      :deep(.el-input__wrapper) {
        border-radius: 0;
      }
    }

    .quantity-hint {
      font-size: 14px;
      color: #666;

      strong {
        color: #e4393c;
      }
    }
  }

  // 操作按钮
  .action-row {
    display: flex;
    gap: 12px;
    margin-bottom: 16px;

    .btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 8px;
      height: 52px;
      border: none;
      border-radius: 10px;
      font-size: 16px;
      font-weight: 600;
      cursor: pointer;
      transition: all 0.2s ease;
      padding: 0 24px;

      &:active {
        transform: scale(0.98);
      }
    }

    .btn-cart {
      flex: 1;
      background: linear-gradient(135deg, #ff9a3c, #ff6f3c);
      color: #fff;

      &:hover {
        box-shadow: 0 4px 16px rgba(255, 111, 60, 0.35);
        transform: translateY(-1px);
      }

      &:disabled {
        opacity: 0.7;
        cursor: not-allowed;
        transform: none;
      }

      .loading-spinner {
        width: 14px;
        height: 14px;
        border: 2px solid rgba(255, 255, 255, 0.3);
        border-top-color: #fff;
        border-radius: 50%;
        animation: spin 0.6s linear infinite;
      }
    }

    .btn-buy {
      flex: 1;
      background: linear-gradient(135deg, #e4393c, #c0392b);
      color: #fff;

      &:hover {
        box-shadow: 0 4px 16px rgba(228, 57, 60, 0.35);
        transform: translateY(-1px);
      }
    }

    .btn-fav {
      min-width: 120px;
      background: #fff;
      color: #666;
      border: 2px solid #e0e0e0;
      flex-direction: column;
      gap: 2px;

      .el-icon {
        font-size: 18px;
      }

      span {
        font-size: 12px;
      }

      &:hover {
        border-color: #ffc069;
        color: #ff9a3c;
      }

      &.favorited {
        border-color: #ffc069;
        color: #ff9a3c;
        background: #fff8f0;
      }
    }
  }

  // 购物车入口
  .cart-entry {
    .cart-entry-link {
      display: flex;
      align-items: center;
      gap: 8px;
      padding: 10px 16px;
      background: #f5f8ff;
      border: 1px solid #e0e8f5;
      border-radius: 8px;
      color: $primary-color;
      text-decoration: none;
      font-size: 13px;
      transition: all 0.2s;

      &:hover {
        background: #e8f0ff;
        border-color: $primary-color;
      }

      .arrow {
        margin-left: auto;
      }
    }
  }
}

// ========== 商品详情卡片 ==========
.desc-card {
  margin-top: 24px;

  .desc-header {
    margin-bottom: 24px;

    h2 {
      font-size: 20px;
      font-weight: 600;
      margin-bottom: 12px;
      color: #1a1a1a;
    }

    .desc-divider {
      height: 2px;
      background: linear-gradient(90deg, $primary-color 0%, transparent 100%);
      border-radius: 1px;
    }
  }

  .desc-content {
    .desc-text {
      font-size: 15px;
      line-height: 1.8;
      color: #444;

      :deep(p) {
        margin-bottom: 16px;
      }

      :deep(img) {
        max-width: 100%;
        height: auto;
        border-radius: 8px;
        margin: 16px 0;
      }
    }
  }

  .desc-images {
    .desc-img {
      width: 100%;
      margin-bottom: 8px;
      border-radius: 8px;
      display: block;
    }
  }

  .desc-empty {
    text-align: center;
    padding: 60px 0;
    color: #999;

    .el-icon {
      font-size: 48px;
      margin-bottom: 12px;
      opacity: 0.3;
    }

    p {
      font-size: 15px;
    }
  }
}

// ========== 回到顶部 ==========
.back-to-top {
  position: fixed;
  right: 40px;
  bottom: 80px;
  width: 44px;
  height: 44px;
  background: $primary-color;
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  box-shadow: 0 4px 16px rgba(64, 158, 255, 0.35);
  transition: all 0.3s;
  z-index: 500;

  &:hover {
    transform: translateY(-4px);
    box-shadow: 0 8px 24px rgba(64, 158, 255, 0.5);
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

// ========== 图片切换过渡 ==========
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.35s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

// ==================== 响应式 ====================
@media (max-width: $bp-mobile) {
  .product-card {
    padding: 16px;

    .product-body {
      flex-direction: column;
      gap: 24px;
    }
  }

  .gallery {
    width: 100%;

    .main-image-box {
      width: 100%;
      height: 0;
      padding-bottom: 100%;
      position: relative;

      .main-image {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
      }
    }

    .thumb-strip {
      .thumb-item {
        width: 56px;
        height: 56px;
      }
    }
  }

  .info {
    .product-name {
      font-size: 18px;
    }

    .price-card .price-row {
      .price-number {
        font-size: 28px;
      }
    }

    .action-row {
      flex-wrap: wrap;

      .btn {
        height: 46px;
        font-size: 14px;
        padding: 0 16px;
      }

      .btn-fav {
        min-width: 100%;
        flex-direction: row;
        gap: 8px;
      }
    }
  }

  .back-to-top {
    right: 16px;
    bottom: 60px;
    width: 40px;
    height: 40px;
  }
}
</style>