<template>
  <div class="dashboard-page">
    <h2 class="page-title">控制台</h2>

    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <p class="stat-label">商品总数</p>
              <p class="stat-value">{{ stats.goodsCount }}</p>
            </div>
            <el-icon class="stat-icon" style="color: #409eff;"><Goods /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <p class="stat-label">品牌总数</p>
              <p class="stat-value">{{ stats.brandCount }}</p>
            </div>
            <el-icon class="stat-icon" style="color: #9b59b6;"><CollectionTag /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <p class="stat-label">订单总数</p>
              <p class="stat-value">{{ stats.orderCount }}</p>
            </div>
            <el-icon class="stat-icon" style="color: #67c23a;"><Document /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <p class="stat-label">用户总数</p>
              <p class="stat-value">{{ stats.userCount }}</p>
            </div>
            <el-icon class="stat-icon" style="color: #e6a23c;"><User /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <p class="stat-label">今日销售额</p>
              <p class="stat-value">¥{{ stats.todaySales }}</p>
            </div>
            <el-icon class="stat-icon" style="color: #f56c6c;"><Money /></el-icon>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-lg">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>最新订单</span>
          </template>
          <el-table :data="recentOrders" style="width: 100%" v-if="recentOrders.length">
            <el-table-column prop="order_sn" label="订单号" width="180" />
            <el-table-column prop="name" label="收货人" />
            <el-table-column label="金额" width="100">
              <template #default="{ row }">
                <span class="price">{{ row.total }}</span>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getOrderStatusType(row.status)" size="small">{{ getOrderStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-else description="暂无订单" />
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>热销商品</span>
          </template>
          <el-table :data="hotGoods" style="width: 100%" v-if="hotGoods.length">
            <el-table-column prop="name" label="商品名称" show-overflow-tooltip />
            <el-table-column label="价格" width="100">
              <template #default="{ row }">
                <span class="price">{{ row.shop_price }}</span>
              </template>
            </el-table-column>
            <el-table-column label="热销" width="80">
              <template #default="{ row }">
                <el-tag v-if="row.is_hot" type="danger" size="small">热销</el-tag>
              </template>
            </el-table-column>
          </el-table>
          <el-empty v-else description="暂无热销商品" />
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getGoodsList } from '@/api/goods'
import { getBrandList } from '@/api/category'
import { getOrderList } from '@/api/order'
import { getUserList } from '@/api/user'

const stats = reactive({
  goodsCount: 0,
  brandCount: 0,
  orderCount: 0,
  userCount: 0,
  todaySales: '0.00'
})

const recentOrders = ref<any[]>([])
const hotGoods = ref<any[]>([])

const getOrderStatusText = (status: string): string => {
  const map: Record<string, string> = {
    'paying': '待支付',
    'WAIT_BUYER_PAY': '待支付',
    'TRADE_SUCCESS': '已支付',
    'TRADE_CLOSED': '已关闭',
    'TRADE_FINISHED': '已完成'
  }
  return map[status] || status
}

const getOrderStatusType = (status: string): string => {
  if (status === 'paying' || status === 'WAIT_BUYER_PAY') return 'warning'
  if (status === 'TRADE_SUCCESS' || status === 'TRADE_FINISHED') return 'success'
  if (status === 'TRADE_CLOSED') return 'info'
  return 'info'
}

const loadStats = async () => {
  try {
    const [goodsRes, brandRes, orderRes, userRes] = await Promise.all([
      getGoodsList({ p: 1, pnum: 1 }),
      getBrandList(1, 1),
      getOrderList({ p: 1, pnum: 5 }),
      getUserList({ pn: 1, psize: 1 })
    ])
    stats.goodsCount = (goodsRes as any).total || 0
    stats.brandCount = (brandRes as any).total || 0
    stats.orderCount = (orderRes as any).total || 0
    stats.userCount = Array.isArray(userRes) ? userRes.length : (userRes as any)?.total || 0
    recentOrders.value = (orderRes as any).data || []

    // 计算今日销售额
    const todayOrders = recentOrders.value.filter((o: any) => {
      if (!o.add_time || o.status === 'TRADE_CLOSED') return false
      const today = new Date().toLocaleDateString()
      const orderDate = new Date(o.add_time).toLocaleDateString()
      return orderDate === today
    })
    const totalSales = todayOrders.reduce((sum: number, o: any) => sum + (o.total || 0), 0)
    stats.todaySales = totalSales.toFixed(2)
  } catch (error) {
    console.error('加载统计数据失败', error)
  }
}

const loadHotGoods = async () => {
  try {
    const res = await getGoodsList({ ishot: 1, p: 1, pnum: 5 })
    hotGoods.value = (res as any).data || []
  } catch (error) {
    console.error('加载热销商品失败', error)
  }
}

onMounted(() => {
  loadStats()
  loadHotGoods()
})
</script>

<style lang="scss" scoped>
.page-title {
  margin-bottom: 20px;
}

.stat-card {
  .stat-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .stat-label {
    color: $text-secondary;
    font-size: 14px;
    margin-bottom: 8px;
  }

  .stat-value {
    font-size: 28px;
    font-weight: bold;
    color: $text-primary;
  }

  .stat-icon {
    font-size: 48px;
    opacity: 0.8;
  }
}

.mt-lg {
  margin-top: 20px;
}
</style>