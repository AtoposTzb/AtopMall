<template>
  <div class="goods-edit-page">
    <h2>{{ isEdit ? '编辑商品' : '新增商品' }}</h2>

    <el-card v-loading="loading">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px">
        <el-form-item label="商品名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入商品名称" />
        </el-form-item>
        <el-form-item label="商品货号" prop="goods_sn">
          <el-input v-model="form.goods_sn" placeholder="请输入商品货号" />
        </el-form-item>
        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="市场价格" prop="market_price">
              <el-input-number v-model="form.market_price" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="本店价格" prop="shop_price">
              <el-input-number v-model="form.shop_price" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="库存" prop="stocks">
              <el-input-number v-model="form.stocks" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="商品分类" prop="category">
              <el-tree-select
                v-model="form.category"
                :data="categoryTree"
                :props="{ label: 'name', value: 'id', children: 'children' }"
                placeholder="请选择分类"
                check-strictly
                filterable
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="品牌" prop="brand">
              <el-select v-model="form.brand" placeholder="请选择品牌" filterable style="width: 100%">
                <el-option v-for="b in brandOptions" :key="b.id" :label="b.name" :value="b.id" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="商品简介" prop="goods_brief">
          <el-input v-model="form.goods_brief" type="textarea" :rows="3" placeholder="请输入商品简介" />
        </el-form-item>
        <el-form-item label="封面图">
          <ImageUpload v-model="form.front_image" />
        </el-form-item>
        <el-form-item label="轮播图">
          <MultiImageUpload v-model="form.images" />
        </el-form-item>
        <el-form-item label="详情图">
          <MultiImageUpload v-model="form.desc_images" />
        </el-form-item>
        <el-form-item label="是否包邮">
          <el-switch v-model="form.ship_free" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">
            {{ isEdit ? '保存修改' : '创建商品' }}
          </el-button>
          <el-button @click="$router.back()">返回</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance } from 'element-plus'
import { getGoodsDetail, createGoods, updateGoods } from '@/api/goods'
import { getCategoryList, getBrandList, type CategoryItem, type BrandItem } from '@/api/category'
import ImageUpload from '@/components/ImageUpload.vue'
import MultiImageUpload from '@/components/MultiImageUpload.vue'

interface CategoryTreeNode {
  id: number
  name: string
  level: number
  parent: number | null
  children: CategoryTreeNode[]
}

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const submitting = ref(false)

const isEdit = computed(() => !!route.params.id)
const goodsId = Number(route.params.id)

const categoryTree = ref<CategoryTreeNode[]>([])
const brandOptions = ref<BrandItem[]>([])

const form = reactive({
  name: '',
  goods_sn: '',
  market_price: 0,
  shop_price: 0,
  stocks: 0,
  category: undefined as number | undefined,
  brand: undefined as number | undefined,
  goods_brief: '',
  front_image: '',
  images: [] as string[],
  desc_images: [] as string[],
  ship_free: false
})

const rules = {
  name: [{ required: true, message: '请输入商品名称', trigger: 'blur' }],
  goods_sn: [{ required: true, message: '请输入商品货号', trigger: 'blur' }],
  market_price: [{ required: true, message: '请输入市场价格', trigger: 'blur' }],
  shop_price: [{ required: true, message: '请输入本店价格', trigger: 'blur' }],
  stocks: [{ required: true, message: '请输入库存', trigger: 'blur' }],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  brand: [{ required: true, message: '请选择品牌', trigger: 'change' }]
}

const buildCategoryTree = (list: CategoryItem[]): CategoryTreeNode[] => {
  return list.map(item => ({
    id: item.id,
    name: item.name,
    level: item.level,
    parent: item.parent,
    children: item.sub_category ? buildCategoryTree(item.sub_category) : []
  }))
}

const handleSubmit = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  submitting.value = true
  try {
    if (isEdit.value) {
      await updateGoods(goodsId, form)
      ElMessage.success('更新成功')
    } else {
      await createGoods(form)
      ElMessage.success('创建成功')
    }
    router.push('/goods')
  } catch {
    /* 拦截器已处理 */
  } finally {
    submitting.value = false
  }
}

const loadGoodsDetail = async () => {
  if (!isEdit.value) return
  loading.value = true
  try {
    const res = await getGoodsDetail(goodsId) as any
    Object.assign(form, {
      name: res.name,
      goods_sn: res.goods_sn,
      market_price: res.market_price,
      shop_price: res.shop_price,
      stocks: res.stocks || 0,
      category: res.ctegory?.id,
      brand: res.brand?.id,
      goods_brief: res.goods_brief,
      front_image: res.front_image,
      images: res.images || [],
      desc_images: res.desc_images || [],
      ship_free: res.ship_free
    })
  } finally {
    loading.value = false
  }
}

const loadOptions = async () => {
  try {
    const [catRes, brandRes] = await Promise.all([
      getCategoryList() as any,
      getBrandList(1, 1000) as any
    ])
    categoryTree.value = buildCategoryTree(catRes || [])
    brandOptions.value = brandRes.data || []
  } catch {
    /* ignore */
  }
}

onMounted(() => {
  loadOptions()
  loadGoodsDetail()
})
</script>

<style lang="scss" scoped>
h2 {
  margin-bottom: 20px;
}
</style>