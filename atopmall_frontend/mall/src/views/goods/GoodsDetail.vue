<template>
  <div class="goods-detail-page" v-loading="loading">
    <AppHeader />

    <div class="container" v-if="goods">
      <!-- 面包屑导航 -->
      <div class="breadcrumb-wrapper">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item :to="{ path: '/goods' }"
            >全部商品</el-breadcrumb-item
          >
          <el-breadcrumb-item>{{ goods.name }}</el-breadcrumb-item>
        </el-breadcrumb>
      </div>

      <div class="detail-card card">
        <div class="detail-main flex">
          <!-- 商品图片 -->
          <div class="detail-images">
            <div class="main-image-wrapper">
              <img :src="currentImage" class="main-image" />
              <div class="image-badges">
                <span v-if="goods.is_hot" class="badge badge-hot">热销</span>
                <span v-if="goods.is_new" class="badge badge-new">新品</span>
              </div>
            </div>
            <div class="thumbnail-list flex" v-if="goods.images?.length">
              <img
                v-for="(img, index) in goods.images"
                :key="index"
                :src="img"
                class="thumbnail"
                :class="{ active: currentImage === img }"
                @click="currentImage = img"
              />
            </div>
          </div>

          <!-- 商品信息 -->
          <div class="detail-info">
            <h1 class="goods-name">{{ goods.name }}</h1>
            <p class="goods-brief">{{ goods.goods_brief }}</p>

            <div class="goods-price">
              <div class="price-main">
                <span class="price-symbol">¥</span>
                <span class="price-value">{{ goods.shop_price }}</span>
              </div>
            </div>

            <div class="goods-meta">
              <div class="meta-item" v-if="goods.brand?.name">
                <span class="meta-label">品牌：</span>
                <span class="meta-value">{{ goods.brand.name }}</span>
              </div>
              <div class="meta-item" v-if="goods.ctegory?.name">
                <span class="meta-label">分类：</span>
                <span class="meta-value">{{ goods.ctegory.name }}</span>
              </div>
              
              <div class="meta-item tags">
                <span v-if="goods.ship_free" class="service-tag">
                  <el-icon><Van /></el-icon>
                  包邮
                </span>
              </div>
            </div>

            <!-- 数量选择 -->
            <div class="quantity-section">
              <span class="quantity-label">数量：</span>
              <el-input-number
                v-model="quantity"
                :min="1"
                :max="99"
                size="large"
              />
            </div>

            <!-- 操作按钮 -->
            <div class="action-buttons">
              <el-button
                type="primary"
                size="large"
                @click="handleAddCart"
                :loading="addCartLoading"
                class="btn-cart"
              >
                <el-icon><ShoppingCart /></el-icon>
                加入购物车
              </el-button>
              <el-button
                type="danger"
                size="large"
                @click="handleBuyNow"
                class="btn-buy"
              >
                立即购买
              </el-button>
              <el-button
                :type="isFaved ? 'warning' : 'default'"
                size="large"
                @click="handleToggleFav"
                class="btn-fav"
              >
                <el-icon><Star /></el-icon>
                {{ isFaved ? "取消收藏" : "收藏" }}
              </el-button>
            </div>

            <!-- 购物车入口 -->
            <div class="cart-entry" v-if="cartStore.totalCount > 0">
              <router-link to="/cart" class="cart-entry-link">
                <el-icon><ShoppingCart /></el-icon>
                <span>购物车中有 <strong>{{ cartStore.totalCount }}</strong> 件商品</span>
                <el-icon><ArrowRight /></el-icon>
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
        <div
          class="desc-empty"
          v-if="!goods.desc && !goods.desc_images?.length"
        >
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
const quantity = ref(1);
const loading = ref(false);
const addCartLoading = ref(false);
const isFaved = ref(false);
const showBackTop = ref(false);

const goodsId = Number(route.params.id);

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
  // 检查购物车中是否已有该商品，没有则添加
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
    currentImage.value =
      goods.value.front_image || goods.value.images?.[0] || "";
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
});
</script>

<style lang="scss" scoped>
.breadcrumb-wrapper {
  margin: 20px 0;
}

.detail-card {
  .detail-main {
    gap: 48px;
  }

  .detail-images {
    width: 480px;
    flex-shrink: 0;

    .main-image-wrapper {
      position: relative;
      border-radius: 12px;
      overflow: hidden;
      box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);

      .main-image {
        width: 480px;
        height: 480px;
        object-fit: cover;
        display: block;
      }

      .image-badges {
        position: absolute;
        top: 16px;
        left: 16px;
        display: flex;
        gap: 8px;

        .badge {
          padding: 6px 12px;
          border-radius: 6px;
          font-size: 13px;
          font-weight: 600;
          color: #fff;
          box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);

          &.badge-hot {
            background: linear-gradient(135deg, #ff6b6b 0%, #ee5a6f 100%);
          }

          &.badge-new {
            background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
          }
        }
      }
    }

    .thumbnail-list {
      margin-top: 16px;
      gap: 12px;
    }

    .thumbnail {
      width: 72px;
      height: 72px;
      object-fit: cover;
      border: 3px solid transparent;
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.3s;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      }

      &.active {
        border-color: $primary-color;
      }
    }
  }

  .detail-info {
    flex: 1;

    .goods-name {
      font-size: 28px;
      font-weight: 600;
      margin-bottom: 12px;
      line-height: 1.3;
      color: $text-primary;
    }

    .goods-brief {
      color: $text-secondary;
      font-size: 15px;
      margin-bottom: 24px;
      line-height: 1.6;
    }

    .goods-price {
      background: linear-gradient(135deg, #fff5f5 0%, #ffe8e8 100%);
      padding: 24px;
      border-radius: 12px;
      margin-bottom: 28px;
      display: flex;
      align-items: center;
      gap: 16px;
      flex-wrap: wrap;

      .price-main {
        display: flex;
        align-items: baseline;
        color: $danger-color;

        .price-symbol {
          font-size: 20px;
          font-weight: 600;
        }

        .price-value {
          font-size: 36px;
          font-weight: 700;
          margin-left: 4px;
        }
      }

      .market-price {
        color: $text-secondary;
        font-size: 14px;

        del {
          opacity: 0.7;
        }
      }

      .discount-tag {
        background: $danger-color;
        color: #fff;
        padding: 4px 10px;
        border-radius: 6px;
        font-size: 13px;
        font-weight: 600;
      }
    }

    .goods-meta {
      margin-bottom: 28px;

      .meta-item {
        display: flex;
        align-items: center;
        margin-bottom: 16px;
        font-size: 15px;

        .meta-label {
          color: $text-secondary;
          min-width: 60px;
        }

        .meta-value {
          color: $text-primary;
          font-weight: 500;

          &.out-of-stock {
            color: $danger-color;
            font-weight: 600;
          }
        }

        &.tags {
          gap: 8px;

          .service-tag {
            display: inline-flex;
            align-items: center;
            gap: 4px;
            padding: 6px 12px;
            background: #f0f9ff;
            color: $primary-color;
            border-radius: 6px;
            font-size: 13px;
            font-weight: 500;
          }
        }
      }
    }

    .quantity-section {
      display: flex;
      align-items: center;
      margin: 28px 0;
      gap: 16px;

      .quantity-label {
        color: $text-secondary;
        font-size: 15px;
      }

      .stock-hint {
        display: inline-flex;
        align-items: center;
        gap: 4px;
        color: $warning-color;
        font-size: 13px;
        font-weight: 500;
        padding: 4px 10px;
        background: #fff7e6;
        border-radius: 6px;
      }
    }

    .action-buttons {
      display: flex;
      gap: 16px;
      margin-top: 32px;

      .btn-cart,
      .btn-buy {
        flex: 1;
        height: 48px;
        font-size: 16px;
        font-weight: 600;
      }

      .btn-fav {
        min-width: 120px;
        height: 48px;
        font-size: 15px;
      }
    }
  }
}

.desc-card {
  margin-top: 32px;

  .desc-header {
    margin-bottom: 24px;

    h2 {
      font-size: 22px;
      font-weight: 600;
      margin-bottom: 12px;
      color: $text-primary;
    }

    .desc-divider {
      height: 2px;
      background: linear-gradient(90deg, $primary-color 0%, transparent 100%);
      border-radius: 1px;
    }
  }

  .desc-content {
    margin-bottom: 24px;

    .desc-text {
      font-size: 15px;
      line-height: 1.8;
      color: $text-regular;

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
    color: $text-secondary;

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

// ======= 购物车入口 =======
.cart-entry {
  margin-top: 20px;

  .cart-entry-link {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 12px 20px;
    background: #f0f9ff;
    border: 1px solid #d0e8ff;
    border-radius: 10px;
    color: $primary-color;
    text-decoration: none;
    font-size: 14px;
    transition: all 0.2s;

    &:hover {
      background: #dceeff;
      border-color: $primary-color;
    }

    strong {
      font-size: 16px;
      margin: 0 2px;
    }
  }
}

// ======= 回到顶部 =======
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
</style>