<template>
  <div class="brand-list-page">
    <div class="page-header flex-between">
      <h2>品牌管理 <span class="total-count">（共 {{ brandTotal }} 个）</span></h2>
    </div>

    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="品牌列表" name="brand">
        <el-card>
          <div class="card-header flex-between mb-md">
            <el-input
              v-model="searchKeyword"
              placeholder="搜索品牌名称"
              clearable
              style="width: 240px"
              @input="onSearchInput"
            />
            <el-button type="primary" @click="handleAdd">新增品牌</el-button>
          </div>
          <el-table :data="displayBrands" v-loading="loading" style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column label="品牌Logo" width="120">
              <template #default="{ row }">
                <el-image :src="row.logo" style="width: 80px; height: 40px;" fit="contain" />
              </template>
            </el-table-column>
            <el-table-column prop="name" label="品牌名称" min-width="200" />
            <el-table-column label="操作" width="280">
              <template #default="{ row }">
                <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
                <el-button type="danger" link @click="handleDelete(row.id)">删除</el-button>
                <el-button type="primary" link @click="handleViewGoods(row)">查看商品</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination mt-md">
            <el-pagination
              v-model:current-page="currentPage"
              :page-size="pageSize"
              :total="filteredTotal"
              layout="total, prev, pager, next"
              @current-change="onPageChange"
            />
          </div>
        </el-card>
      </el-tab-pane>

      <el-tab-pane label="品牌分类绑定" name="categoryBrand">
        <el-card>
          <div class="card-header flex-between mb-md">
            <span></span>
            <el-button type="primary" @click="handleAddCategoryBrand">新增绑定</el-button>
          </div>
          <el-table :data="categoryBrands" v-loading="cbLoading" style="width: 100%">
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column label="分类名称" min-width="150">
              <template #default="{ row }">
                {{ row.category?.name }}
              </template>
            </el-table-column>
            <el-table-column label="品牌名称" min-width="150">
              <template #default="{ row }">
                {{ row.brand?.name }}
              </template>
            </el-table-column>
            <el-table-column label="品牌Logo" width="120">
              <template #default="{ row }">
                <el-image :src="row.brand?.logo" style="width: 80px; height: 40px;" fit="contain" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="200">
              <template #default="{ row }">
                <el-button type="primary" link @click="handleEditCategoryBrand(row)">编辑</el-button>
                <el-button type="danger" link @click="handleDeleteCategoryBrand(row.id)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-tab-pane>
    </el-tabs>

    <!-- 新增/编辑品牌弹窗 -->
    <el-dialog v-model="showDialog" :title="editingId ? '编辑品牌' : '新增品牌'" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="品牌名称">
          <el-input v-model="form.name" placeholder="请输入品牌名称" />
        </el-form-item>
        <el-form-item label="品牌Logo">
          <ImageUpload v-model="form.logo" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" :disabled="!form.name || !form.logo" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>

    <!-- 查看品牌商品弹窗 -->
    <el-dialog v-model="showGoodsDialog" :title="`${goodsBrandName} 的商品列表`" width="800px">
      <el-table :data="goodsList" v-loading="goodsLoading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column label="商品主图" width="80">
          <template #default="{ row }">
            <el-image :src="row.front_image" style="width: 50px; height: 50px;" fit="cover" />
          </template>
        </el-table-column>
        <el-table-column prop="name" label="商品名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="shop_price" label="价格" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.on_sale ? 'success' : 'info'" size="small">
              {{ row.on_sale ? '上架中' : '已下架' }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="showGoodsDialog = false">关闭</el-button>
      </template>
    </el-dialog>

    <!-- 新增/编辑品牌分类绑定弹窗 -->
    <el-dialog v-model="showCbDialog" :title="editingCbId ? '编辑绑定' : '新增绑定'" width="500px">
      <el-form :model="cbForm" label-width="80px">
        <el-form-item label="分类">
          <el-tree-select
            v-model="cbForm.category_id"
            :data="categoryTree"
            :props="{ label: 'name', value: 'id', children: 'children' }"
            placeholder="请选择分类"
            check-strictly
            filterable
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="品牌">
          <el-select
            v-model="cbForm.brand_id"
            placeholder="请选择品牌"
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="item in brandOptions"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCbDialog = false">取消</el-button>
        <el-button type="primary" :disabled="!cbForm.category_id || !cbForm.brand_id" @click="handleSaveCategoryBrand">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  getBrandList,
  createBrand,
  updateBrand,
  deleteBrand,
  getCategoryList,
  getCategoryBrandList,
  createCategoryBrand,
  updateCategoryBrand,
  deleteCategoryBrand,
  type BrandItem,
  type CategoryItem,
  type CategoryBrandItem
} from '@/api/category'
import { getGoodsList, type GoodsItem } from '@/api/goods'
import ImageUpload from '@/components/ImageUpload.vue'

interface CategoryTreeNode {
  id: number
  name: string
  level: number
  parent: number | null
  children: CategoryTreeNode[]
}

const activeTab = ref('brand')
const loading = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = 10
const showDialog = ref(false)
const editingId = ref<number | null>(null)

const showGoodsDialog = ref(false)
const goodsLoading = ref(false)
const goodsList = ref<GoodsItem[]>([])
const goodsBrandName = ref('')

const allBrands = ref<BrandItem[]>([])
const searchKeyword = ref('')
const brandTotal = ref(0)
let searchTimer: ReturnType<typeof setTimeout> | null = null

const filteredBrands = computed(() => {
  if (!searchKeyword.value) return allBrands.value
  const kw = searchKeyword.value.toLowerCase()
  return allBrands.value.filter(b => b.name.toLowerCase().includes(kw))
})

const filteredTotal = computed(() => filteredBrands.value.length)

const displayBrands = computed(() => {
  const start = (currentPage.value - 1) * pageSize
  return filteredBrands.value.slice(start, start + pageSize)
})

const form = reactive({
  name: '',
  logo: ''
})

const categoryBrands = ref<CategoryBrandItem[]>([])
const cbLoading = ref(false)
const showCbDialog = ref(false)
const editingCbId = ref<number | null>(null)
const categoryTree = ref<CategoryTreeNode[]>([])
const brandOptions = ref<BrandItem[]>([])

const cbForm = reactive({
  category_id: 0,
  brand_id: 0
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

const loadBrands = async () => {
  loading.value = true
  try {
    const res = await getBrandList(1, 1000) as any
    allBrands.value = res.data || []
    total.value = res.total || 0
    brandTotal.value = res.total || 0
    currentPage.value = 1
  } finally {
    loading.value = false
  }
}

const onSearchInput = () => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(() => {
    currentPage.value = 1
  }, 300)
}

const onPageChange = () => {
  // computed 自动响应 currentPage 变化
}

const handleAdd = () => {
  editingId.value = null
  form.name = ''
  form.logo = ''
  showDialog.value = true
}

const handleEdit = (row: BrandItem) => {
  editingId.value = row.id
  form.name = row.name
  form.logo = row.logo
  showDialog.value = true
}

const handleSave = async () => {
  try {
    if (editingId.value) {
      await updateBrand(editingId.value, form)
      ElMessage.success('更新成功')
    } else {
      await createBrand(form)
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    await loadBrands()
  } catch (error) {
    // 错误已在拦截器中处理
  }
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定删除该品牌吗？', '提示', { type: 'warning' })
    await deleteBrand(id)
    ElMessage.success('删除成功')
    await loadBrands()
  } catch (error) {
    // 用户取消或错误
  }
}

const handleViewGoods = async (row: BrandItem) => {
  goodsBrandName.value = row.name
  showGoodsDialog.value = true
  goodsLoading.value = true
  try {
    const res = await getGoodsList({ b: row.id, p: 1, pnum: 1000 }) as any
    goodsList.value = res.data || []
  } finally {
    goodsLoading.value = false
  }
}

const loadCategoryBrands = async () => {
  cbLoading.value = true
  try {
    const res = await getCategoryBrandList() as any
    categoryBrands.value = res.data || []
  } finally {
    cbLoading.value = false
  }
}

const loadOptions = async () => {
  const [catRes, brandRes] = await Promise.all([
    getCategoryList() as any,
    getBrandList(1, 1000) as any
  ])
  categoryTree.value = buildCategoryTree(catRes || [])
  brandOptions.value = brandRes.data || []
}

const handleAddCategoryBrand = () => {
  editingCbId.value = null
  cbForm.category_id = 0
  cbForm.brand_id = 0
  showCbDialog.value = true
  loadOptions()
}

const handleEditCategoryBrand = (row: CategoryBrandItem) => {
  editingCbId.value = row.id
  cbForm.category_id = row.category?.id || 0
  cbForm.brand_id = row.brand?.id || 0
  showCbDialog.value = true
  loadOptions()
}

const handleSaveCategoryBrand = async () => {
  try {
    if (editingCbId.value) {
      await updateCategoryBrand(editingCbId.value, cbForm)
      ElMessage.success('更新成功')
    } else {
      await createCategoryBrand(cbForm)
      ElMessage.success('创建成功')
    }
    showCbDialog.value = false
    await loadCategoryBrands()
  } catch (error) {
    // 错误已在拦截器中处理
  }
}

const handleDeleteCategoryBrand = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定删除该绑定吗？', '提示', { type: 'warning' })
    await deleteCategoryBrand(id)
    ElMessage.success('删除成功')
    await loadCategoryBrands()
  } catch (error) {
    // 用户取消或错误
  }
}

const handleTabChange = (name: string) => {
  if (name === 'categoryBrand') {
    loadCategoryBrands()
    if (categoryTree.value.length === 0) {
      loadOptions()
    }
  }
}

onMounted(() => {
  loadBrands()
})
</script>

<style lang="scss" scoped>
.page-header {
  margin-bottom: 20px;
}

.mt-md {
  margin-top: 16px;
}

.mb-md {
  margin-bottom: 16px;
}

.card-header {
  font-weight: 500;
}

.total-count {
  color: $text-secondary;
  font-size: 14px;
  font-weight: normal;
}
</style>