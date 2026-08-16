<template>
  <div class="cart-page">
    <AppHeader />
    
    <div class="container">
      <div class="page-header">
        <div class="header-left">
          <h2 class="page-title">我的购物车</h2>
          <p class="page-subtitle">共 {{ cartStore.totalCount }} 件商品</p>
        </div>
        <el-button type="primary" plain @click="$router.push('/goods')">
          <el-icon><ArrowLeft /></el-icon>
          继续购物
        </el-button>
      </div>

      <div class="cart-wrapper" v-loading="cartStore.loading">
        <div class="cart-header">
          <div class="header-check">
            <el-checkbox
              :model-value="cartStore.isAllChecked"
              @change="cartStore.toggleAllChecked($event as boolean)"
            >
              全选
            </el-checkbox>
          </div>
          <div class="header-goods">商品信息</div>
          <div class="header-price">单价</div>
          <div class="header-quantity">数量</div>
          <div class="header-total">小计</div>
          <div class="header-action">操作</div>
        </div>

        <div v-if="cartStore.cartList.length" class="cart-items">
          <div v-for="item in cartStore.cartList" :key="item.id" class="cart-item">
            <div class="item-check">
              <el-checkbox
                :model-value="item.checked"
                @change="cartStore.updateCart(item.goods_id, { checked: !item.checked })"
              />
            </div>
            <div class="item-goods">
              <img :src="item.goods_image" class="goods-img" />
              <div class="goods-info">
                <span class="goods-name">{{ item.goods_name }}</span>
              </div>
            </div>
            <div class="item-price">
              <span class="price">¥{{ item.goods_price }}</span>
            </div>
            <div class="item-quantity">
              <el-input-number
                :model-value="item.nums"
                :min="1"
                :max="99"
                size="default"
                @change="(val: number) => cartStore.updateCart(item.goods_id, { nums: val })"
              />
            </div>
            <div class="item-total">
              <span class="price">¥{{ (item.goods_price * item.nums).toFixed(2) }}</span>
            </div>
            <div class="item-action">
              <el-button type="danger" link @click="cartStore.removeCartItem(item.goods_id)">
                <el-icon><Delete /></el-icon>
                删除
              </el-button>
            </div>
          </div>
        </div>

        <div v-else class="empty">
          <el-icon class="empty-icon"><ShoppingCart /></el-icon>
          <p>购物车是空的</p>
          <el-button type="primary" @click="$router.push('/goods')">去购物</el-button>
        </div>
      </div>

      <div class="cart-footer" v-if="cartStore.cartList.length">
        <div class="footer-left">
          <el-button type="danger" link @click="cartStore.batchDelete">
            <el-icon><Delete /></el-icon>
            删除选中
          </el-button>
        </div>
        <div class="footer-right">
          <div class="summary">
            <span>已选 <strong>{{ cartStore.checkedCount }}</strong> 件商品</span>
            <div class="total-wrapper">
              <span class="total-label">合计：</span>
              <span class="total-price price">¥{{ cartStore.checkedTotalPrice }}</span>
            </div>
          </div>
          <el-button 
            type="danger" 
            size="large" 
            :disabled="cartStore.checkedCount === 0" 
            @click="$router.push('/checkout')"
            class="checkout-btn"
          >
            去结算
          </el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useCartStore } from '@/store/cart'
import AppHeader from '@/components/AppHeader.vue'

const cartStore = useCartStore()

onMounted(() => {
  cartStore.loadCartList()
})
</script>

<style lang="scss" scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 32px 0 24px;

  .header-left {
    .page-title {
      font-size: 28px;
      font-weight: bold;
      margin-bottom: 8px;
      color: $text-primary;
    }

    .page-subtitle {
      font-size: 15px;
      color: $text-secondary;
    }
  }
}

.cart-wrapper {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
}

.cart-header {
  display: grid;
  grid-template-columns: 60px 3fr 1fr 1fr 1fr 80px;
  gap: 16px;
  align-items: center;
  padding: 16px 0;
  border-bottom: 2px solid $border-light;
  font-weight: 600;
  color: $text-secondary;
  font-size: 14px;

  .header-goods {
    grid-column: 2;
  }

  .header-price,
  .header-quantity,
  .header-total,
  .header-action {
    text-align: center;
  }
}

.cart-items {
  .cart-item {
    display: grid;
    grid-template-columns: 60px 3fr 1fr 1fr 1fr 80px;
    gap: 16px;
    align-items: center;
    padding: 20px 0;
    border-bottom: 1px solid $border-lighter;
    transition: background 0.2s;

    &:hover {
      background: #fafafa;
    }

    &:last-child {
      border-bottom: none;
    }

    .item-check {
      display: flex;
      justify-content: center;
    }

    .item-goods {
      display: flex;
      align-items: center;
      gap: 16px;

      .goods-img {
        width: 100px;
        height: 100px;
        object-fit: cover;
        border-radius: 8px;
        box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
      }

      .goods-info {
        flex: 1;

        .goods-name {
          font-size: 15px;
          color: $text-primary;
          line-height: 1.5;
        }
      }
    }

    .item-price,
    .item-total {
      text-align: center;

      .price {
        font-size: 16px;
        font-weight: 600;
        color: $danger-color;
      }
    }

    .item-quantity {
      display: flex;
      justify-content: center;
    }

    .item-action {
      display: flex;
      justify-content: center;

      .el-button {
        display: flex;
        align-items: center;
        gap: 4px;
      }
    }
  }
}

.cart-footer {
  background: #fff;
  border-radius: 12px;
  padding: 20px 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  bottom: 20px;

  .footer-left {
    .el-button {
      display: flex;
      align-items: center;
      gap: 4px;
    }
  }

  .footer-right {
    display: flex;
    align-items: center;
    gap: 24px;

    .summary {
      display: flex;
      flex-direction: column;
      align-items: flex-end;
      gap: 8px;

      span {
        font-size: 14px;
        color: $text-secondary;

        strong {
          color: $danger-color;
          font-size: 16px;
        }
      }

      .total-wrapper {
        display: flex;
        align-items: baseline;
        gap: 8px;

        .total-label {
          font-size: 15px;
          color: $text-primary;
        }

        .total-price {
          font-size: 28px;
          font-weight: bold;
        }
      }
    }

    .checkout-btn {
      padding: 12px 40px;
      font-size: 16px;
      font-weight: 600;
    }
  }
}

.empty {
  text-align: center;
  padding: 80px 0;
  color: $text-secondary;

  .empty-icon {
    font-size: 64px;
    margin-bottom: 16px;
    opacity: 0.3;
  }

  p {
    font-size: 16px;
    margin-bottom: 24px;
  }
}
</style>