<template>
  <div class="banner-list-page">
    <div class="page-header flex-between">
      <h2>轮播图管理</h2>
      <el-button type="primary" @click="handleAdd">新增轮播图</el-button>
    </div>

    <el-card>
      <el-table :data="banners" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="图片" width="200">
          <template #default="{ row }">
            <el-image :src="row.image" style="width: 160px; height: 80px;" fit="cover" />
          </template>
        </el-table-column>
        <el-table-column prop="url" label="跳转链接" min-width="200" />
        <el-table-column prop="index" label="排序" width="80" />
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="showDialog" :title="editingId ? '编辑轮播图' : '新增轮播图'" width="520px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="图片">
          <ImageUpload v-model="form.image" />
        </el-form-item>
        <el-form-item label="跳转方式">
          <el-radio-group v-model="linkType" @change="onLinkTypeChange">
            <el-radio value="product">商品详情页</el-radio>
            <el-radio value="custom">自定义链接</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="linkType === 'product'" label="选择商品">
          <el-select
            v-model="selectedProductId"
            filterable
            remote
            reserve-keyword
            placeholder="输入商品名称搜索"
            :remote-method="searchProducts"
            :loading="productSearching"
            clearable
            style="width: 100%"
          >
            <el-option
              v-for="item in productOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-if="linkType === 'custom'" label="跳转链接">
          <el-input v-model="form.url" placeholder="请输入链接，如 https://example.com" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="form.index" :min="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getBannerList, createBanner, updateBanner, deleteBanner, type BannerItem } from '@/api/category'
import { getGoodsList, type GoodsItem } from '@/api/goods'
import ImageUpload from '@/components/ImageUpload.vue'

const banners = ref<BannerItem[]>([])
const loading = ref(false)
const showDialog = ref(false)
const editingId = ref<number | null>(null)

const form = reactive({
  image: '',
  url: '',
  index: 0
})

const linkType = ref<'product' | 'custom'>('product')
const selectedProductId = ref<number | undefined>()
const productOptions = ref<GoodsItem[]>([])
const productSearching = ref(false)

const loadBanners = async () => {
  loading.value = true
  try {
    const res = await getBannerList()
    banners.value = res as unknown as BannerItem[]
  } finally {
    loading.value = false
  }
}

const searchProducts = async (query: string) => {
  if (!query) {
    productOptions.value = []
    return
  }
  productSearching.value = true
  try {
    const res = await getGoodsList({ q: query, p: 1, pnum: 20 })
    productOptions.value = (res as any).data || []
  } finally {
    productSearching.value = false
  }
}

const onLinkTypeChange = () => {
  if (linkType.value === 'product') {
    form.url = ''
  } else {
    selectedProductId.value = undefined
  }
}

const parseUrlMode = (url: string) => {
  if (!url) return
  const match = url.match(/^\/goods\/(\d+)$/)
  if (match) {
    linkType.value = 'product'
    selectedProductId.value = Number(match[1])
  } else {
    linkType.value = 'custom'
    form.url = url
  }
}

const handleAdd = () => {
  editingId.value = null
  form.image = ''
  form.url = ''
  form.index = 0
  linkType.value = 'product'
  selectedProductId.value = undefined
  productOptions.value = []
  showDialog.value = true
}

const handleEdit = (row: BannerItem) => {
  editingId.value = row.id
  form.image = row.image
  form.url = ''
  form.index = row.index
  parseUrlMode(row.url)
  showDialog.value = true
}

const handleSave = async () => {
  if (linkType.value === 'product') {
    if (!selectedProductId.value) {
      ElMessage.warning('请选择商品')
      return
    }
    form.url = `/goods/${selectedProductId.value}`
  } else {
    if (!form.url) {
      ElMessage.warning('请输入跳转链接')
      return
    }
  }

  try {
    if (editingId.value) {
      await updateBanner(editingId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createBanner(form)
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    await loadBanners()
  } catch (error) {
    // 错误已在拦截器中处理
  }
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定删除该轮播图吗？', '提示', { type: 'warning' })
    await deleteBanner(id)
    ElMessage.success('删除成功')
    await loadBanners()
  } catch (error) {
    // 用户取消或错误
  }
}

onMounted(() => {
  loadBanners()
})
</script>

<style lang="scss" scoped>
.page-header {
  margin-bottom: 20px;
}
</style>