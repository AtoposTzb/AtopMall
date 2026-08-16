<template>
  <div class="order-detail-page">
    <AppHeader />

    <div class="container" v-loading="loading">
      <!-- 面包屑导航 -->
      <div class="breadcrumb">
        <el-breadcrumb separator="/">
          <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
          <el-breadcrumb-item>
            <a href="javascript:;" @click="goBack">我的订单</a>
          </el-breadcrumb-item>
          <el-breadcrumb-item>订单详情</el-breadcrumb-item>
        </el-breadcrumb>
      </div>

      <div v-if="order">
        <h2 class="page-title">订单详情</h2>

        <div class="card">
          <div class="card-header">
            <h3>订单信息</h3>
            <el-tag
              :type="getStatusType(order.status)"
              size="large"
              effect="dark"
            >
              {{ getStatusText(order.status) }}
            </el-tag>
          </div>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="订单号">{{
              order.order_sn
            }}</el-descriptions-item>
            <el-descriptions-item label="下单时间">{{
              order.add_time
            }}</el-descriptions-item>
            <el-descriptions-item label="收货人">{{
              order.name
            }}</el-descriptions-item>
            <el-descriptions-item label="联系电话">{{
              order.mobile
            }}</el-descriptions-item>
            <el-descriptions-item label="收货地址" :span="2">{{
              order.address
            }}</el-descriptions-item>
            <el-descriptions-item label="订单金额">
              <span class="price" style="font-size: 20px; font-weight: bold"
                >¥{{ order.total }}</span
              >
            </el-descriptions-item>
            <el-descriptions-item label="支付方式">{{
              order.pay_type || "未支付"
            }}</el-descriptions-item>
            <el-descriptions-item label="留言" v-if="order.post">{{
              order.post
            }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <div class="card">
          <h3>商品清单</h3>
          <div
            v-for="goods in order.goods"
            :key="goods.id"
            class="goods-item flex"
          >
            <img :src="goods.image" class="goods-img" />
            <div class="goods-info">
              <p>{{ goods.name }}</p>
              <p class="text-secondary">
                <span class="price">¥{{ goods.price }}</span> × {{ goods.nums }}
              </p>
            </div>
          </div>
        </div>

        <!-- 待支付：显示支付倒计时和支付方式 -->
        <div class="card" v-if="isPayable">
          <h3>完成支付</h3>

          <div class="countdown-bar" :class="{ urgent: remainingSeconds <= 120 }" v-if="!isExpired">
            <div class="countdown-left">
              <el-icon class="countdown-icon" :class="{ pulse: remainingSeconds <= 120 }"><Clock /></el-icon>
              <span class="countdown-label">支付剩余</span>
              <span class="countdown-time">{{ countdownText }}</span>
            </div>
            <div class="countdown-progress">
              <div class="progress-track">
                <div class="progress-fill" :style="{ width: (remainingSeconds / COUNTDOWN_SECONDS * 100) + '%' }"></div>
              </div>
            </div>
          </div>
          <div class="countdown-bar expired" v-else>
            <div class="expired-content">
              <el-icon><WarningFilled /></el-icon>
              <span>支付已超时，订单可能已关闭</span>
            </div>
            <el-button
              type="primary"
              size="small"
              @click="loadOrder"
              >刷新状态</el-button
            >
          </div>

          <template v-if="!isExpired">
            <div class="pay-section">
              <div class="pay-methods">
                <div class="pay-method-label">选择支付方式</div>
                <div class="pay-method-cards">
                  <div
                    class="pay-card"
                    :class="{ active: payMethod === 'alipay' }"
                    @click="payMethod = 'alipay'"
                  >
                    <div class="pay-card-icon alipay-icon">
                      <span class="icon-text">支</span>
                    </div>
                    <div class="pay-card-info">
                      <div class="pay-card-name">支付宝</div>
                      <div class="pay-card-desc">推荐使用，安全便捷</div>
                    </div>
                    <div class="pay-card-check" v-if="payMethod === 'alipay'">
                      <el-icon><Check /></el-icon>
                    </div>
                  </div>
                  <div
                    class="pay-card"
                    :class="{ active: payMethod === 'wechat' }"
                    @click="payMethod = 'wechat'"
                  >
                    <div class="pay-card-icon wechat-icon">
                      <span class="icon-text">微</span>
                    </div>
                    <div class="pay-card-info">
                      <div class="pay-card-name">微信支付</div>
                      <div class="pay-card-desc">暂不支持，请使用支付宝</div>
                    </div>
                    <div class="pay-card-check" v-if="payMethod === 'wechat'">
                      <el-icon><Check /></el-icon>
                    </div>
                  </div>
                  <div
                    class="pay-card"
                    :class="{ active: payMethod === 'mock' }"
                    @click="payMethod = 'mock'"
                  >
                    <div class="pay-card-icon mock-icon">
                      <span class="icon-text">演</span>
                    </div>
                    <div class="pay-card-info">
                      <div class="pay-card-name">模拟支付</div>
                      <div class="pay-card-desc">演示环境使用</div>
                    </div>
                    <div class="pay-card-check" v-if="payMethod === 'mock'">
                      <el-icon><Check /></el-icon>
                    </div>
                  </div>
                </div>
              </div>

              <div class="pay-message">
                <div class="pay-message-label">留言备注</div>
                <el-input
                  v-model="postMessage"
                  type="textarea"
                  :rows="2"
                  placeholder="选填：配送时间、包装要求等"
                  maxlength="200"
                  show-word-limit
                />
              </div>

              <div class="pay-actions">
                <el-button @click="goBack">返回订单列表</el-button>
                <el-button
                  type="danger"
                  size="large"
                  @click="handlePay"
                  :loading="paying"
                  class="pay-btn"
                >
                  <el-icon><Wallet /></el-icon>
                  去支付 ¥{{ order.total }}
                </el-button>
              </div>
            </div>
          </template>

          <div class="pay-actions" v-else>
            <el-button @click="goBack">返回订单列表</el-button>
          </div>
        </div>

        <!-- 已支付/已发货等 -->
        <div class="card action-card" v-if="!isPayable">
          <el-button type="primary" @click="goBack">返回订单列表</el-button>
          <el-button @click="$router.push('/')">返回首页</el-button>
        </div>
      </div>

      <div v-else-if="!loading" class="empty">
        <el-icon class="empty-icon"><Document /></el-icon>
        <p>订单不存在</p>
        <el-button type="primary" @click="goBack">返回订单列表</el-button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { getOrderDetail, type OrderItem } from "@/api/order";
import AppHeader from "@/components/AppHeader.vue";

const route = useRoute();
const router = useRouter();
const order = ref<OrderItem | null>(null);
const loading = ref(false);
const paying = ref(false);
const payMethod = ref("alipay");
const postMessage = ref("");

// 支付倒计时（与后端 RocketMQ delay_level=14 对应，10分钟）
const COUNTDOWN_SECONDS = 600;
const remainingSeconds = ref(COUNTDOWN_SECONDS);
const isExpired = ref(false);
let timer: ReturnType<typeof setInterval> | null = null;

const countdownText = computed(() => {
  const m = Math.floor(remainingSeconds.value / 60);
  const s = remainingSeconds.value % 60;
  return `${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`;
});

function startCountdown() {
  stopCountdown();
  remainingSeconds.value = COUNTDOWN_SECONDS;
  isExpired.value = false;

  timer = setInterval(() => {
    remainingSeconds.value--;
    if (remainingSeconds.value <= 0) {
      isExpired.value = true;
      stopCountdown();
      loadOrder();
    }
  }, 1000);
}

function stopCountdown() {
  if (timer) {
    clearInterval(timer);
    timer = null;
  }
}

const orderId = Number(route.params.id);

// 后端订单状态值：paying(默认), TRADE_SUCCESS, TRADE_CLOSED, WAIT_BUYER_PAY, TRADE_FINISHED
const isPayable = computed(() => {
  const s = order.value?.status;
  return s === "paying" || s === "WAIT_BUYER_PAY";
});

const getStatusText = (status: string): string => {
  const statusMap: Record<string, string> = {
    paying: "待支付",
    WAIT_BUYER_PAY: "待支付",
    TRADE_SUCCESS: "已支付",
    TRADE_CLOSED: "已关闭",
    TRADE_FINISHED: "已完成",
    pending: "待支付",
    paid: "已支付",
    shipped: "已发货",
    received: "已收货",
    cancelled: "已取消",
    refunded: "已退款",
  };
  return statusMap[status] || status;
};

const getStatusType = (
  status: string,
): "primary" | "success" | "warning" | "info" | "danger" => {
  if (
    status === "paying" ||
    status === "WAIT_BUYER_PAY" ||
    status === "pending"
  )
    return "warning";
  if (
    status === "TRADE_SUCCESS" ||
    status === "paid" ||
    status === "received" ||
    status === "TRADE_FINISHED"
  )
    return "success";
  if (status === "shipped") return "primary";
  if (status === "TRADE_CLOSED" || status === "cancelled") return "info";
  return "info";
};

// 返回上一页（用户中心 > 订单列表），如果没有上一页则回首页
const goBack = () => {
  if (window.history.length > 1) {
    router.back();
  } else {
    router.push("/");
  }
};

const handlePay = () => {
  if (payMethod.value === "mock") {
    // 模拟支付成功
    paying.value = true;
    setTimeout(() => {
      paying.value = false;
      ElMessage.success("模拟支付成功！");
      router.push("/");
    }, 1500);
    return;
  }
  // 真实支付：使用后端返回的支付宝支付链接
  if (payMethod.value === "alipay" && order.value?.alipay_url) {
    window.open(order.value.alipay_url);
  } else {
    ElMessage.info("当前选择微信支付，请使用支付宝支付或选择模拟支付");
  }
};

async function loadOrder() {
  loading.value = true;
  try {
    const res = await getOrderDetail(orderId);
    order.value = res as unknown as OrderItem;
    if (order.value && isPayable.value) {
      startCountdown();
    }
  } catch (error) {
    console.error("加载订单详情失败", error);
  } finally {
    loading.value = false;
  }
}

onMounted(() => {
  loadOrder();
});

onUnmounted(() => {
  stopCountdown();
});
</script>

<style lang="scss" scoped>
.breadcrumb {
  padding: 20px 0 0;

  a {
    color: $text-regular;
    text-decoration: none;
    cursor: pointer;

    &:hover {
      color: $primary-color;
    }
  }
}

.page-title {
  font-size: 24px;
  font-weight: bold;
  margin: 20px 0;
  color: $text-primary;
}

.card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }

  h3 {
    margin-bottom: 16px;
    font-size: 16px;
    color: $text-primary;
  }
}

.goods-item {
  align-items: center;
  padding: 12px 0;
  border-bottom: 1px solid $border-lighter;
  gap: 12px;

  &:last-child {
    border-bottom: none;
  }

  .goods-img {
    width: 80px;
    height: 80px;
    object-fit: cover;
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  }

  .goods-info {
    p {
      margin-bottom: 4px;
      font-size: 15px;
      color: $text-primary;
    }

    .text-secondary {
      font-size: 14px;
      color: $text-secondary;

      .price {
        color: $danger-color;
        font-weight: 600;
      }
    }
  }
}

.countdown-bar {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 16px 20px;
  background: linear-gradient(135deg, #fef0f0 0%, #fff5f5 100%);
  border: 1px solid #fde2e2;
  border-radius: 10px;
  margin-bottom: 24px;
  font-size: 14px;
  color: #f56c6c;
  transition: all 0.3s;

  &.urgent {
    background: linear-gradient(135deg, #fef0f0 0%, #fde2e2 100%);
    border-color: #fab6b6;
    animation: urgentPulse 2s ease-in-out infinite;
  }

  .countdown-left {
    display: flex;
    align-items: center;
    gap: 8px;
  }

  .countdown-icon {
    font-size: 18px;

    &.pulse {
      animation: pulse 1s ease-in-out infinite;
    }
  }

  .countdown-label {
    color: #606266;
    font-size: 13px;
  }

  .countdown-time {
    font-size: 22px;
    font-weight: bold;
    font-variant-numeric: tabular-nums;
    color: #f56c6c;
    letter-spacing: 1px;
  }

  .countdown-progress {
    .progress-track {
      width: 100%;
      height: 4px;
      background: #fde2e2;
      border-radius: 2px;
      overflow: hidden;

      .progress-fill {
        height: 100%;
        background: linear-gradient(90deg, #f56c6c 0%, #e6a23c 100%);
        border-radius: 2px;
        transition: width 1s linear;
      }
    }
  }

  &.expired {
    background: #f5f7fa;
    border-color: #e4e7ed;
    color: #909399;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;

    .expired-content {
      display: flex;
      align-items: center;
      gap: 8px;
    }
  }
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.4; }
}

@keyframes urgentPulse {
  0%, 100% { border-color: #fab6b6; }
  50% { border-color: #f56c6c; }
}

.pay-section {
  .pay-methods {
    margin-bottom: 24px;

    .pay-method-label {
      font-size: 15px;
      font-weight: 600;
      color: $text-primary;
      margin-bottom: 14px;
    }

    .pay-method-cards {
      display: flex;
      flex-direction: column;
      gap: 10px;
    }

    .pay-card {
      display: flex;
      align-items: center;
      gap: 14px;
      padding: 16px 18px;
      border: 2px solid $border-light;
      border-radius: 10px;
      cursor: pointer;
      transition: all 0.3s;
      position: relative;
      background: #fff;

      &:hover {
        border-color: #b3d8ff;
        box-shadow: 0 2px 12px rgba(64, 158, 255, 0.08);
        transform: translateX(4px);
      }

      &.active {
        border-color: $primary-color;
        background: linear-gradient(135deg, #ecf5ff 0%, #f0f9ff 100%);
        box-shadow: 0 2px 12px rgba(64, 158, 255, 0.12);

        .pay-card-check {
          opacity: 1;
          transform: scale(1);
        }
      }

      .pay-card-icon {
        width: 44px;
        height: 44px;
        border-radius: 10px;
        display: flex;
        align-items: center;
        justify-content: center;
        flex-shrink: 0;

        .icon-text {
          font-size: 18px;
          font-weight: bold;
          color: #fff;
        }

        &.alipay-icon {
          background: linear-gradient(135deg, #1677ff 0%, #4096ff 100%);
        }

        &.wechat-icon {
          background: linear-gradient(135deg, #07c160 0%, #2dc97a 100%);
        }

        &.mock-icon {
          background: linear-gradient(135deg, #909399 0%, #b0b3b8 100%);
        }
      }

      .pay-card-info {
        flex: 1;

        .pay-card-name {
          font-size: 15px;
          font-weight: 600;
          color: $text-primary;
          margin-bottom: 2px;
        }

        .pay-card-desc {
          font-size: 12px;
          color: $text-secondary;
        }
      }

      .pay-card-check {
        width: 24px;
        height: 24px;
        border-radius: 50%;
        background: $primary-color;
        display: flex;
        align-items: center;
        justify-content: center;
        opacity: 0;
        transform: scale(0.6);
        transition: all 0.3s;

        .el-icon {
          color: #fff;
          font-size: 14px;
        }
      }
    }
  }

  .pay-message {
    margin-bottom: 24px;

    .pay-message-label {
      font-size: 15px;
      font-weight: 600;
      color: $text-primary;
      margin-bottom: 12px;
    }
  }

  .pay-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;

    .pay-btn {
      min-width: 160px;
      font-size: 16px;
      font-weight: 600;
      letter-spacing: 1px;
      border-radius: 8px;
      transition: all 0.3s;

      &:hover {
        transform: translateY(-1px);
        box-shadow: 0 4px 14px rgba(245, 108, 108, 0.35);
      }
    }
  }
}

.card {
  > .pay-actions {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }
}

.action-card {
  display: flex;
  justify-content: center;
  gap: 16px;
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