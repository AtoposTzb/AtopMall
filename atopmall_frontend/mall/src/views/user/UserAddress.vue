<template>
  <div class="user-address">
    <div class="flex-between mb-md">
      <h3>地址管理</h3>
      <el-button type="primary" @click="handleAdd">新增地址</el-button>
    </div>

    <div v-loading="loading">
      <div v-for="addr in addresses" :key="addr.id" class="address-item card flex-between">
        <div>
          <p><strong>{{ addr.signer_name }}</strong> {{ addr.signer_mobile }}</p>
          <p class="text-secondary">{{ addr.province }} {{ addr.city }} {{ addr.district }} {{ addr.address }}</p>
        </div>
        <div class="flex">
          <el-button type="primary" link @click="handleEdit(addr)">编辑</el-button>
          <el-button type="danger" link @click="handleDelete(addr.id)">删除</el-button>
        </div>
      </div>

      <div v-if="!loading && addresses.length === 0" class="empty">
        <el-icon class="empty-icon"><Location /></el-icon>
        <p>暂无收货地址</p>
      </div>
    </div>

    <!-- 新增/编辑地址弹窗 -->
    <el-dialog 
      v-model="showDialog" 
      :title="editingId ? '编辑地址' : '新增地址'" 
      width="500px"
      @close="handleClose"
    >
      <el-form 
        ref="formRef"
        :model="form" 
        :rules="rules"
        label-width="80px"
      >
        <el-form-item label="收货人" prop="signer_name">
          <el-input v-model="form.signer_name" placeholder="请输入收货人姓名" />
        </el-form-item>
        <el-form-item label="手机号" prop="signer_mobile">
          <el-input v-model="form.signer_mobile" placeholder="请输入手机号" />
        </el-form-item>
        <el-form-item label="省份" prop="province">
          <el-input v-model="form.province" placeholder="请输入省份" />
        </el-form-item>
        <el-form-item label="城市" prop="city">
          <el-input v-model="form.city" placeholder="请输入城市" />
        </el-form-item>
        <el-form-item label="地区" prop="district">
          <el-input v-model="form.district" placeholder="请输入地区" />
        </el-form-item>
        <el-form-item label="详细地址" prop="address">
          <el-input 
            v-model="form.address" 
            type="textarea" 
            :rows="3"
            placeholder="请输入详细地址"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="saving">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { getAddressList, createAddress, updateAddress, deleteAddress, type AddressItem } from '@/api/address'

const addresses = ref<AddressItem[]>([])
const loading = ref(false)
const saving = ref(false)
const showDialog = ref(false)
const editingId = ref<number | null>(null)
const formRef = ref<FormInstance>()

const form = reactive({
  signer_name: '',
  signer_mobile: '',
  province: '',
  city: '',
  district: '',
  address: ''
})

const rules: FormRules = {
  signer_name: [
    { required: true, message: '请输入收货人姓名', trigger: 'blur' },
    { min: 2, max: 20, message: '长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  signer_mobile: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  province: [
    { required: true, message: '请输入省份', trigger: 'blur' }
  ],
  city: [
    { required: true, message: '请输入城市', trigger: 'blur' }
  ],
  district: [
    { required: true, message: '请输入地区', trigger: 'blur' }
  ],
  address: [
    { required: true, message: '请输入详细地址', trigger: 'blur' },
    { min: 5, max: 100, message: '长度在 5 到 100 个字符', trigger: 'blur' }
  ]
}

const loadAddresses = async () => {
  loading.value = true
  try {
    const res = await getAddressList()
    addresses.value = (res as any).data || []
  } finally {
    loading.value = false
  }
}

const resetForm = () => {
  form.signer_name = ''
  form.signer_mobile = ''
  form.province = ''
  form.city = ''
  form.district = ''
  form.address = ''
  editingId.value = null
  formRef.value?.clearValidate()
}

const handleAdd = () => {
  resetForm()
  showDialog.value = true
}

const handleEdit = (addr: AddressItem) => {
  editingId.value = addr.id
  form.signer_name = addr.signer_name
  form.signer_mobile = addr.signer_mobile
  form.province = addr.province
  form.city = addr.city
  form.district = addr.district
  form.address = addr.address
  showDialog.value = true
}

const handleClose = () => {
  resetForm()
}

const handleSave = async () => {
  if (!formRef.value) return
  
  const valid = await formRef.value.validate().catch(() => false)
  if (!valid) return

  saving.value = true
  try {
    if (editingId.value) {
      await updateAddress(editingId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createAddress(form)
      ElMessage.success('添加成功')
    }
    showDialog.value = false
    resetForm()
    await loadAddresses()
  } catch (error) {
    // 错误已在拦截器中处理
  } finally {
    saving.value = false
  }
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定删除该地址吗？', '提示', { 
      type: 'warning',
      confirmButtonText: '确定',
      cancelButtonText: '取消'
    })
    await deleteAddress(id)
    ElMessage.success('删除成功')
    await loadAddresses()
  } catch (error) {
    // 用户取消或错误
  }
}

onMounted(() => {
  loadAddresses()
})
</script>

<style lang="scss" scoped>
.user-address {
  h3 {
    font-size: 20px;
    font-weight: 600;
    color: $text-primary;
    margin: 0;
    padding-bottom: 16px;
    border-bottom: 2px solid #f0f5ff;
    position: relative;

    &::after {
      content: '';
      position: absolute;
      bottom: -2px;
      left: 0;
      width: 60px;
      height: 2px;
      background: $primary-color;
      border-radius: 1px;
    }
  }
}

.address-item {
  margin-bottom: 12px;
  padding: 16px;
  transition: all 0.3s;
  border-radius: 10px;
  border: 1px solid #f0f0f0;

  &:hover {
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
    border-color: transparent;
  }
}
</style>