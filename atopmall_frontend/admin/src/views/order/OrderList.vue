<template>
  <div class="order-list-page">
    <h2>订单管理</h2>

    <el-card class="mb-md">
      <el-form :inline="true">
        <el-form-item label="订单状态">
          <el-select v-model="filterStatus" placeholder="全部" clearable @change="loadOrders">
            <el-option label="待支付" value="paying" />
            <el-option label="待支付" value="WAIT_BUYER_PAY" />
            <el-option label="已支付" value="TRADE_SUCCESS" />
            <el-option label="已关闭" value="TRADE_CLOSED" />
            <el-option label="已完成" value="TRADE_FINISHED" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadOrders">查询</el-button>
          <el-button @click="filterStatus = ''; loadOrders()">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card>
      <el-table :data="orders" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="order_sn" label="订单号" width="200" />
        <el-table-column prop="name" label="收货人" width="120" />
        <el-table-column prop="mobile" label="联系电话" width="140" />
        <el-table-column label="订单金额" width="120">
          <template #default="{ row }">
            <span class="price" style="font-size: 15px;">{{ row.total }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="pay_type" label="支付方式" width="100" />
        <el-table-column prop="add_time" label="下单时间" width="180" />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push(`/order/${row.id}`)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination mt-md">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadOrders"
        />
      </div>
    </el-card>
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
const filterStatus = ref('')

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

const loadOrders = async () => {
  loading.value = true
  try {
    const params: any = { p: currentPage.value, pnum: pageSize }
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    const res = await getOrderList(params)
    orders.value = (res as any).data || []
    total.value = (res as any).total || 0
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadOrders()
})
</script>

<style lang="scss" scoped>
h2 {
  margin-bottom: 20px;
}

.mb-md {
  margin-bottom: 16px;
}

.mt-md {
  margin-top: 16px;
}
</style>