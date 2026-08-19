<template>
  <div class="order-list-page">
    <div class="container">
      <!-- 面包屑导航 -->
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item>我的订单</el-breadcrumb-item>
        </el-breadcrumb>
      </div>

      <div class="page-header">
        <el-button @click="$router.back()" class="back-btn">
          <el-icon><ArrowLeft /></el-icon>
          返回
        </el-button>
        <h2 class="page-title">我的订单</h2>
      </div>

      <div v-loading="loading">
        <div v-for="order in orders" :key="order.id" class="card order-card">
          <div class="order-header flex-between">
            <div class="order-sn">
              <span class="label">订单号：</span>
              <span class="value">{{ order.order_sn }}</span>
            </div>
            <el-tag :type="getStatusType(order.status)" effect="dark">
              {{ getStatusText(order.status) }}
            </el-tag>
          </div>
          
          <div class="order-body" @click="$router.push(`/order/${order.id}`)">
            <div class="order-info">
              <p class="info-item">
                <el-icon><User /></el-icon>
                <span>{{ order.name }} {{ order.mobile }}</span>
              </p>
              <p class="info-item">
                <el-icon><Location /></el-icon>
                <span>{{ order.address }}</span>
              </p>
            </div>
            
            <div class="order-total">
              <span class="label">合计：</span>
              <span class="price">{{ order.total }}</span>
            </div>
            
            <div class="order-time">
              <el-icon><Clock /></el-icon>
              <span>{{ formatTime(order.add_time) }}</span>
            </div>
          </div>
          
          <div class="order-footer">
            <el-button type="primary" size="small" @click="$router.push(`/order/${order.id}`)">
              查看详情
            </el-button>
          </div>
        </div>
      </div>

      <div v-if="!loading && orders.length === 0" class="empty">
        <el-icon class="empty-icon"><Document /></el-icon>
        <p>暂无订单</p>
        <el-button type="primary" @click="$router.push('/goods')">去购物</el-button>
      </div>

      <div class="pagination" v-if="total > 0">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="loadOrders"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getOrderList, type OrderItem } from '@/api/order'

const orders = ref<OrderItem[]>([])
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = 10

const getStatusText = (status: string): string => {
  const statusMap: Record<string, string> = {
    'paying': '待支付',
    'WAIT_BUYER_PAY': '待支付',
    'TRADE_SUCCESS': '已支付',
    'TRADE_CLOSED': '已关闭',
    'TRADE_FINISHED': '已完成',
    'pending': '待支付',
    'paid': '已支付',
    'shipped': '已发货',
    'received': '已收货',
    'cancelled': '已取消',
    'refunded': '已退款'
  }
  return statusMap[status] || status
}

const getStatusType = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  if (status === 'paying' || status === 'WAIT_BUYER_PAY' || status === 'pending') return 'warning'
  if (status === 'TRADE_SUCCESS' || status === 'paid' || status === 'received' || status === 'TRADE_FINISHED') return 'success'
  if (status === 'shipped') return 'primary'
  if (status === 'TRADE_CLOSED' || status === 'cancelled') return 'info'
  return 'info'
}

const formatTime = (time: string): string => {
  if (!time) return ''
  const date = new Date(time)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

const loadOrders = async () => {
  loading.value = true
  try {
    const res = await getOrderList({ p: currentPage.value, pnum: pageSize })
    orders.value = ((res as any).data || []).sort((a: OrderItem, b: OrderItem) =>
      new Date(b.add_time).getTime() - new Date(a.add_time).getTime()
    )
    total.value = (res as any).total || 0
  } catch (error) {
    console.error('加载订单失败', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadOrders()
})
</script>

<style lang="scss" scoped>
.breadcrumb {
  padding: 16px 0 8px;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;

  .back-btn {
    flex-shrink: 0;
  }
}

.page-title {
  font-size: 22px;
  margin: 0;
}

.order-card {
  margin-bottom: 16px;
  transition: all 0.3s ease;

  &:hover {
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.12);
    transform: translateY(-2px);
  }

  .order-header {
    padding-bottom: 12px;
    border-bottom: 1px solid $border-lighter;
    margin-bottom: 12px;

    .order-sn {
      .label {
        color: $text-secondary;
        font-size: 14px;
      }
      .value {
        color: $text-primary;
        font-weight: 500;
      }
    }
  }

  .order-body {
    cursor: pointer;
    padding: 12px 0;

    .order-info {
      margin-bottom: 12px;

      .info-item {
        display: flex;
        align-items: center;
        gap: 8px;
        color: $text-regular;
        font-size: 14px;
        margin-bottom: 6px;

        .el-icon {
          color: $text-secondary;
        }
      }
    }

    .order-total {
      margin-bottom: 8px;

      .label {
        color: $text-secondary;
        font-size: 14px;
      }

      .price {
        font-size: 20px;
        font-weight: bold;
        color: $danger-color;
      }
    }

    .order-time {
      display: flex;
      align-items: center;
      gap: 6px;
      color: $text-secondary;
      font-size: 13px;
    }
  }

  .order-footer {
    padding-top: 12px;
    border-top: 1px solid $border-lighter;
    display: flex;
    justify-content: flex-end;
  }
}

// ==================== 响应式 ====================
@media (max-width: $bp-mobile) {
  .page-title {
    font-size: 18px;
  }

  .order-card .order-header {
    flex-direction: column;
    gap: 8px;
    align-items: flex-start;
  }

  .order-body .order-info .info-item span {
    font-size: 13px;
  }
}
</style>