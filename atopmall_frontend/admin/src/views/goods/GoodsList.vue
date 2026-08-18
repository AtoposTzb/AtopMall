<template>
  <div class="goods-list-page">
    <div class="page-header flex-between">
      <h2>商品管理 <span class="total-count">（共 {{ goodsTotal }} 件）</span></h2>
      <el-button type="primary" @click="$router.push('/goods/create')">新增商品</el-button>
    </div>

    <!-- 搜索栏 -->
    <el-card class="mb-md">
      <el-form :inline="true" :model="searchParams">
        <el-form-item label="关键词">
          <el-input v-model="searchParams.q" placeholder="商品名称" clearable @keyup.enter="loadGoods" />
        </el-form-item>
        <el-form-item label="分类">
          <el-tree-select
            v-model="searchParams.c"
            :data="categoryTree"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            placeholder="全部"
            check-strictly
            clearable
            filterable
            style="width: 180px"
          />
        </el-form-item>
        <el-form-item label="品牌">
          <el-select v-model="searchParams.b" placeholder="全部" clearable filterable style="width: 160px">
            <el-option v-for="b in brandOptions" :key="b.id" :label="b.name" :value="b.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="价格">
          <el-input-number v-model="searchParams.pmin" :min="0" placeholder="最低" size="small" controls-position="right" />
          <span class="mx-sm">-</span>
          <el-input-number v-model="searchParams.pmax" :min="0" placeholder="最高" size="small" controls-position="right" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadGoods">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
        <el-form-item>
          <el-radio-group v-model="onsaleFilter" @change="onOnsaleChange" size="small">
            <el-radio-button value="0">全部</el-radio-button>
            <el-radio-button value="1">已上架</el-radio-button>
            <el-radio-button value="2">未上架</el-radio-button>
          </el-radio-group>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 商品列表 -->
    <el-card>
      <el-table :data="goodsList" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="商品图片" width="100">
          <template #default="{ row }">
            <el-image :src="row.front_image" style="width: 60px; height: 60px;" fit="cover" />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="商品名称" min-width="180" show-overflow-tooltip />
        <el-table-column label="分类" width="120">
          <template #default="{ row }">
            {{ row.ctegory?.name }}
          </template>
        </el-table-column>
        <el-table-column label="品牌" width="120">
          <template #default="{ row }">
            {{ row.brand?.name }}
          </template>
        </el-table-column>
        <el-table-column label="售价" width="100">
          <template #default="{ row }">
            <span class="price">{{ formatPrice(row.shop_price) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="库存" width="80">
          <template #default="{ row }">
            {{ stocksMap[row.id] ?? '...' }}
          </template>
        </el-table-column>
        <el-table-column label="上架" width="80">
          <template #default="{ row }">
            <el-switch :model-value="row.on_sale" size="small" @change="(val: boolean) => toggleStatus(row, { sale: val, new: row.is_new, hot: row.is_hot })" />
          </template>
        </el-table-column>
        <el-table-column label="新品" width="80">
          <template #default="{ row }">
            <el-switch :model-value="row.is_new" size="small" @change="(val: boolean) => toggleStatus(row, { sale: row.on_sale, new: val, hot: row.is_hot })" />
          </template>
        </el-table-column>
        <el-table-column label="热销" width="80">
          <template #default="{ row }">
            <el-switch :model-value="row.is_hot" size="small" @change="(val: boolean) => toggleStatus(row, { sale: row.on_sale, new: row.is_new, hot: val })" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push(`/goods/edit/${row.id}`)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination mt-md">
        <el-pagination
          v-model:current-page="searchParams.p"
          :page-size="searchParams.pnum"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadGoods"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getGoodsList, deleteGoods, updateGoodsStatus, getGoodsStock, type GoodsItem } from '@/api/goods'
import { getCategoryList, getBrandList, type CategoryItem, type BrandItem } from '@/api/category'

interface CategoryTreeNode {
  id: number
  name: string
  level: number
  parent: number | null
  children: CategoryTreeNode[]
}

const goodsList = ref<GoodsItem[]>([])
const loading = ref(false)
const total = ref(0)
const categoryTree = ref<CategoryTreeNode[]>([])
const brandOptions = ref<BrandItem[]>([])
const goodsTotal = ref(0)
const onsaleFilter = ref('0')
const stocksMap = ref<Record<number, number>>({})

const searchParams = reactive({
  p: 1,
  pnum: 10,
  q: '',
  c: undefined as number | undefined,
  b: undefined as number | undefined,
  pmin: undefined as number | undefined,
  pmax: undefined as number | undefined,
  onsale: undefined as number | undefined
})

const buildCategoryTree = (list: CategoryItem[]): CategoryTreeNode[] => {
  return list.map(item => ({
    id: item.id,
    name: item.name,
    level: item.level,
    parent: item.parent,
    children: item.sub_category ? buildCategoryTree(item.sub_category) : []
  }))
}

const formatPrice = (val: any): string => {
  const num = Number(String(val).replace(/[¥￥]/g, ''))
  return isNaN(num) ? '0.00' : num.toFixed(2)
}

const loadFilters = async () => {
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

const loadGoods = async () => {
  loading.value = true
  try {
    const params = { ...searchParams }
    if (params.c === undefined) delete params.c
    if (params.b === undefined) delete params.b
    if (params.onsale === undefined) delete params.onsale
    const res = await getGoodsList(params) as any
    goodsList.value = res.data || []
    total.value = res.total || 0
    goodsTotal.value = res.total || 0
    fetchStocks(res.data || [])
  } finally {
    loading.value = false
  }
}

const fetchStocks = async (list: GoodsItem[]) => {
  const results = await Promise.allSettled(
    list.map(item => getGoodsStock(item.id))
  )
  const map: Record<number, number> = {}
  results.forEach((r, i) => {
    if (r.status === 'fulfilled') {
      map[list[i].id] = (r.value as any).num ?? 0
    }
  })
  stocksMap.value = map
}

const resetSearch = () => {
  searchParams.q = ''
  searchParams.c = undefined
  searchParams.b = undefined
  searchParams.pmin = undefined
  searchParams.pmax = undefined
  searchParams.p = 1
  onsaleFilter.value = '0'
  searchParams.onsale = 0
  loadGoods()
}

const onOnsaleChange = (val: string) => {
  searchParams.onsale = Number(val)
  searchParams.p = 1
  loadGoods()
}

const toggleStatus = async (row: GoodsItem, data: { sale: boolean; new: boolean; hot: boolean }) => {
  try {
    await updateGoodsStatus(row.id, data)
    row.on_sale = data.sale
    row.is_new = data.new
    row.is_hot = data.hot
    ElMessage.success('更新成功')
  } catch {
    await loadGoods()
  }
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定删除该商品吗？', '提示', { type: 'warning' })
    await deleteGoods(id)
    ElMessage.success('删除成功')
    await loadGoods()
  } catch {
    /* 用户取消 */
  }
}

onMounted(() => {
  loadFilters()
  loadGoods()
})
</script>

<style lang="scss" scoped>
.page-header {
  margin-bottom: 20px;
}

.mb-md {
  margin-bottom: 16px;
}

.mt-md {
  margin-top: 16px;
}

.mx-sm {
  margin: 0 8px;
}

.price {
  color: #e4393c;
  font-weight: 500;
}

.total-count {
  color: $text-secondary;
  font-size: 14px;
  font-weight: normal;
}
</style>