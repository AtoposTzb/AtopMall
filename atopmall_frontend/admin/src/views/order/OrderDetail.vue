<template>
  <div class="order-detail-page" v-loading="loading">
    <div class="page-header flex-between">
      <h2>订单详情</h2>
      <el-button @click="$router.push('/order')">返回订单列表</el-button>
    </div>

    <el-card class="mb-md" v-if="order">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="订单号">{{ order.order_sn }}</el-descriptions-item>
        <el-descriptions-item label="订单状态">
          <el-tag :type="getStatusType(order.status)" size="large">{{ getStatusText(order.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="收货人">{{ order.name }}</el-descriptions-item>
        <el-descriptions-item label="联系电话">{{ order.mobile }}</el-descriptions-item>
        <el-descriptions-item label="收货地址" :span="2">{{ order.address }}</el-descriptions-item>
        <el-descriptions-item label="订单金额">
          <span class="price" style="font-size: 20px;">{{ order.total }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="支付方式">{{ order.pay_type || '未支付' }}</el-descriptions-item>
        <el-descriptions-item label="创建时间" :span="2">{{ order.add_time }}</el-descriptions-item>
        <el-descriptions-item label="留言" v-if="order.post">{{ order.post }}</el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card class="mb-md" v-if="order?.goods">
      <h3>商品清单</h3>
      <el-table :data="order.goods" style="width: 100%">
        <el-table-column label="商品图片" width="100">
          <template #default="{ row }">
            <el-image :src="row.image" style="width: 60px; height: 60px;" fit="cover" />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="商品名称" />
        <el-table-column label="单价" width="120">
          <template #default="{ row }">
            <span class="price">{{ row.price }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="nums" label="数量" width="100" />
        <el-table-column label="小计" width="120">
          <template #default="{ row }">
            <span class="price">{{ (row.price * row.nums).toFixed(2) }}</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <div class="empty" v-if="!order && !loading">
      <el-icon class="empty-icon"><Document /></el-icon>
      <p>订单不存在</p>
      <el-button type="primary" @click="$router.push('/order')">返回订单列表</el-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getOrderDetail, type OrderItem } from '@/api/order'

const route = useRoute()
const order = ref<OrderItem | null>(null)
const loading = ref(false)

const orderId = Number(route.params.id)

const getStatusText = (status: string): string => {
  const statusMap: Record<string, string> = {
    'paying': '待支付',
    'WAIT_BUYER_PAY': '待支付',
    'TRADE_SUCCESS': '已支付',
    'TRADE_CLOSED': '已关闭',
    'TRADE_FINISHED': '已完成'
  }
  return statusMap[status] || status
}

const getStatusType = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  if (status === 'paying' || status === 'WAIT_BUYER_PAY') return 'warning'
  if (status === 'TRADE_SUCCESS') return 'success'
  if (status === 'TRADE_FINISHED') return 'success'
  if (status === 'TRADE_CLOSED') return 'info'
  return 'info'
}

onMounted(async () => {
  loading.value = true
  try {
    const res = await getOrderDetail(orderId)
    order.value = res as unknown as OrderItem
  } finally {
    loading.value = false
  }
})
</script>

<style lang="scss" scoped>
.page-header {
  margin-bottom: 20px;
}

h2, h3 {
  margin-bottom: 16px;
}

.mb-md {
  margin-bottom: 16px;
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