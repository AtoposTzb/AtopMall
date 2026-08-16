<template>
  <div class="checkout-page">
    <AppHeader />

    <div class="container">
      <div class="page-header">
        <h2 class="page-title">确认订单</h2>
        <p class="page-subtitle">请核对订单信息并完成支付</p>
      </div>

      <!-- 收货地址 -->
      <div class="section-card">
        <div class="section-header">
          <div class="section-title">
            <el-icon>
              <Location />
            </el-icon>
            <span>收货地址</span>
          </div>
          <el-button type="primary" link @click="showAddressDialog = true">
            <el-icon>
              <Plus />
            </el-icon>
            新增地址
          </el-button>
        </div>

        <div v-if="addresses.length" class="address-list">
          <div v-for="addr in addresses" :key="addr.id" class="address-item"
            :class="{ selected: selectedAddress?.id === addr.id }" @click="selectedAddress = addr">
            <div class="address-header">
              <strong class="signer-name">{{ addr.signer_name }}</strong>
              <span class="signer-mobile">{{ addr.signer_mobile }}</span>
              <el-icon v-if="selectedAddress?.id === addr.id" class="check-icon">
                <Check />
              </el-icon>
            </div>
            <p class="address-detail">{{ addr.province }} {{ addr.city }} {{ addr.district }} {{ addr.address }}</p>
          </div>
        </div>

        <div v-else class="empty-address">
          <el-icon class="empty-icon">
            <Location />
          </el-icon>
          <p>暂无收货地址</p>
          <el-button type="primary" @click="showAddressDialog = true">添加地址</el-button>
        </div>
      </div>

      <!-- 商品清单 -->
      <div class="section-card">
        <div class="section-header">
          <div class="section-title">
            <el-icon>
              <GoodsFilled />
            </el-icon>
            <span>商品清单</span>
          </div>
          <span class="goods-count">共 {{ checkedItems.length }} 件商品</span>
        </div>

        <div class="goods-list">
          <div v-for="item in checkedItems" :key="item.id" class="goods-item">
            <div class="goods-image">
              <img :src="item.goods_image" :alt="item.goods_name" />
            </div>
            <div class="goods-info">
              <p class="goods-name">{{ item.goods_name }}</p>
              <div class="goods-meta">
                <span class="goods-price price">¥{{ item.goods_price }}</span>
                <span class="goods-nums">× {{ item.nums }}</span>
              </div>
            </div>
            <div class="goods-total">
              <span class="price">¥{{ (item.goods_price * item.nums).toFixed(2) }}</span>
            </div>
          </div>
        </div>

        <div class="order-summary">
          <div class="summary-row">
            <span>商品总额</span>
            <span class="price">¥{{ cartStore.checkedTotalPrice }}</span>
          </div>
          <div class="summary-row">
            <span>运费</span>
            <span class="free-shipping">免运费</span>
          </div>
          <div class="summary-row total-row">
            <span>应付金额</span>
            <span class="total-price price">¥{{ cartStore.checkedTotalPrice }}</span>
          </div>
        </div>
      </div>

      <!-- 提交订单 -->
      <div class="submit-section">
        <div class="submit-info">
          <div class="info-item">
            <span class="label">收货人：</span>
            <span class="value">{{ selectedAddress?.signer_name || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="label">联系电话：</span>
            <span class="value">{{ selectedAddress?.signer_mobile || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="label">收货地址：</span>
            <span class="value">{{ selectedAddress ?
              `${selectedAddress.province}${selectedAddress.city}${selectedAddress.district}${selectedAddress.address}`
              : '-' }}</span>
          </div>
        </div>
        <div class="submit-action">
          <div class="amount-wrapper">
            <span class="amount-label">应付金额：</span>
            <span class="amount-value price">¥{{ cartStore.checkedTotalPrice }}</span>
          </div>
          <el-button type="danger" size="large" :loading="submitting" @click="handleSubmitOrder"
            :disabled="!selectedAddress || checkedItems.length === 0" class="submit-btn">
            <el-icon v-if="!submitting">
              <Check />
            </el-icon>
            提交订单
          </el-button>
        </div>
      </div>

      <!-- 新增地址弹窗 -->
      <el-dialog v-model="showAddressDialog" title="添加收货地址" width="500px" :close-on-click-modal="false">
        <el-form ref="addressFormRef" :model="addressForm" :rules="addressRules" label-width="100px">
          <el-form-item label="收货人" prop="signer_name">
            <el-input v-model="addressForm.signer_name" placeholder="请输入收货人姓名" />
          </el-form-item>
          <el-form-item label="手机号" prop="signer_mobile">
            <el-input v-model="addressForm.signer_mobile" placeholder="请输入手机号" />
          </el-form-item>
          <el-form-item label="省份" prop="province">
            <el-input v-model="addressForm.province" placeholder="请输入省份" />
          </el-form-item>
          <el-form-item label="城市" prop="city">
            <el-input v-model="addressForm.city" placeholder="请输入城市" />
          </el-form-item>
          <el-form-item label="地区" prop="district">
            <el-input v-model="addressForm.district" placeholder="请输入地区" />
          </el-form-item>
          <el-form-item label="详细地址" prop="address">
            <el-input v-model="addressForm.address" type="textarea" :rows="3" placeholder="请输入详细地址" />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="showAddressDialog = false">取消</el-button>
          <el-button type="primary" @click="handleSaveAddress" :loading="savingAddress">保存</el-button>
        </template>
      </el-dialog>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useCartStore } from '@/store/cart'
import { getAddressList, createAddress, type AddressItem } from '@/api/address'
import { createOrder } from '@/api/order'
import AppHeader from '@/components/AppHeader.vue'

const router = useRouter()
const cartStore = useCartStore()

const addresses = ref<AddressItem[]>([])
const selectedAddress = ref<AddressItem | null>(null)
const showAddressDialog = ref(false)
const submitting = ref(false)
const savingAddress = ref(false)
const addressFormRef = ref<FormInstance>()

const addressForm = reactive({
  signer_name: '',
  signer_mobile: '',
  province: '',
  city: '',
  district: '',
  address: ''
})

const addressRules: FormRules = {
  signer_name: [
    { required: true, message: '请输入收货人姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  signer_mobile: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  province: [{ required: true, message: '请输入省份', trigger: 'blur' }],
  city: [{ required: true, message: '请输入城市', trigger: 'blur' }],
  district: [{ required: true, message: '请输入地区', trigger: 'blur' }],
  address: [
    { required: true, message: '请输入详细地址', trigger: 'blur' },
    { min: 5, max: 100, message: '长度在 5 到 100 个字符', trigger: 'blur' }
  ]
}

const checkedItems = computed(() => cartStore.cartList.filter(item => item.checked))

const loadAddresses = async () => {
  try {
    const res = await getAddressList()
    addresses.value = (res as any).data || []
    if (addresses.value.length && !selectedAddress.value) {
      selectedAddress.value = addresses.value[0]
    }
  } catch (error) {
    console.error('加载地址失败', error)
  }
}

const handleSaveAddress = async () => {
  if (!addressFormRef.value) return

  const valid = await addressFormRef.value.validate().catch(() => false)
  if (!valid) return

  savingAddress.value = true
  try {
    await createAddress(addressForm)
    ElMessage.success('添加成功')
    showAddressDialog.value = false
    addressForm.signer_name = ''
    addressForm.signer_mobile = ''
    addressForm.province = ''
    addressForm.city = ''
    addressForm.district = ''
    addressForm.address = ''
    addressFormRef.value?.resetFields()
    await loadAddresses()
  } catch (error) {
    // 错误已在拦截器中处理
  } finally {
    savingAddress.value = false
  }
}

const handleSubmitOrder = async () => {
  if (!selectedAddress.value) {
    ElMessage.warning('请选择收货地址')
    return
  }
  if (checkedItems.value.length === 0) {
    ElMessage.warning('请选择商品')
    return
  }

  submitting.value = true
  try {
    const addr = selectedAddress.value
    const res = await createOrder({
      address: `${addr.province}${addr.city}${addr.district}${addr.address}`,
      name: addr.signer_name,
      mobile: addr.signer_mobile,
      post: ''
    })
    ElMessage.success('订单创建成功，请选择支付方式并完成支付')
    await cartStore.loadCartList()
    router.push(`/order/${res.id}`)
  } catch (error) {
    // 错误已在拦截器中处理
  } finally {
    submitting.value = false
  }
}

onMounted(async () => {
  await cartStore.loadCartList()
  loadAddresses()
})
</script>

<style lang="scss" scoped>
.page-header {
  padding: 32px 0 24px;

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

.section-card {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    padding-bottom: 16px;
    border-bottom: 1px solid $border-lighter;

    .section-title {
      display: flex;
      align-items: center;
      gap: 8px;
      font-size: 18px;
      font-weight: 600;
      color: $text-primary;

      .el-icon {
        font-size: 20px;
        color: $primary-color;
      }
    }

    .goods-count {
      font-size: 14px;
      color: $text-secondary;
    }
  }
}

.address-list {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.address-item {
  padding: 16px;
  border: 2px solid $border-light;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;

  &:hover {
    border-color: $primary-color;
    box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
  }

  &.selected {
    border-color: $primary-color;
    background: linear-gradient(135deg, #ecf5ff 0%, #f0f9ff 100%);

    .check-icon {
      position: absolute;
      top: 12px;
      right: 12px;
      font-size: 20px;
      color: $primary-color;
    }
  }

  .address-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;

    .signer-name {
      font-size: 16px;
      color: $text-primary;
    }

    .signer-mobile {
      font-size: 14px;
      color: $text-secondary;
    }
  }

  .address-detail {
    font-size: 14px;
    color: $text-regular;
    line-height: 1.5;
  }
}

.empty-address {
  text-align: center;
  padding: 40px 0;
  color: $text-secondary;

  .empty-icon {
    font-size: 48px;
    margin-bottom: 12px;
    opacity: 0.3;
  }

  p {
    margin-bottom: 16px;
    font-size: 15px;
  }
}

.goods-list {
  .goods-item {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px 0;
    border-bottom: 1px solid $border-lighter;

    &:last-child {
      border-bottom: none;
    }

    .goods-image {
      width: 80px;
      height: 80px;
      flex-shrink: 0;
      border-radius: 8px;
      overflow: hidden;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);

      img {
        width: 100%;
        height: 100%;
        object-fit: cover;
      }
    }

    .goods-info {
      flex: 1;

      .goods-name {
        font-size: 15px;
        color: $text-primary;
        margin-bottom: 8px;
        line-height: 1.4;
      }

      .goods-meta {
        display: flex;
        align-items: center;
        gap: 12px;

        .goods-price {
          font-size: 16px;
          font-weight: 600;
        }

        .goods-nums {
          font-size: 14px;
          color: $text-secondary;
        }
      }
    }

    .goods-total {
      .price {
        font-size: 18px;
        font-weight: 600;
      }
    }
  }
}

.order-summary {
  margin-top: 20px;
  padding-top: 20px;
  border-top: 2px solid $border-lighter;

  .summary-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 8px 0;
    font-size: 15px;

    span:first-child {
      color: $text-secondary;
    }

    .free-shipping {
      color: $success-color;
      font-weight: 500;
    }

    &.total-row {
      padding-top: 16px;
      margin-top: 8px;
      border-top: 1px dashed $border-light;

      span:first-child {
        font-size: 16px;
        font-weight: 600;
        color: $text-primary;
      }

      .total-price {
        font-size: 28px;
        font-weight: bold;
      }
    }
  }
}

.submit-section {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  position: sticky;
  bottom: 20px;

  .submit-info {
    display: flex;
    gap: 32px;
    padding-bottom: 20px;
    margin-bottom: 20px;
    border-bottom: 1px solid $border-lighter;

    .info-item {
      display: flex;
      align-items: center;
      gap: 8px;

      .label {
        color: $text-secondary;
        font-size: 14px;
      }

      .value {
        color: $text-primary;
        font-size: 14px;
        font-weight: 500;
      }
    }
  }

  .submit-action {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .amount-wrapper {
      display: flex;
      align-items: baseline;
      gap: 8px;

      .amount-label {
        font-size: 16px;
        color: $text-primary;
        font-weight: 600;
      }

      .amount-value {
        font-size: 32px;
        font-weight: bold;
      }
    }

    .submit-btn {
      padding: 14px 48px;
      font-size: 16px;
      font-weight: 600;
      display: flex;
      align-items: center;
      gap: 6px;
    }
  }
}
</style>