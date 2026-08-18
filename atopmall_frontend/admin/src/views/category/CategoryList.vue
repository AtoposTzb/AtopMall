<template>
  <div class="category-list-page">
    <div class="page-header flex-between">
      <h2>分类管理</h2>
      <el-button type="primary" @click="handleAdd">新增分类</el-button>
    </div>

    <el-card>
      <el-table
        :key="tableKey"
        :data="categories"
        v-loading="loading"
        row-key="id"
        :tree-props="{ children: 'sub_category' }"
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" label="分类名称" min-width="200">
          <template #default="{ row }">
            <span :style="{ paddingLeft: (row.level - 1) * 24 + 'px' }">
              <el-icon v-if="row.level > 1" class="sub-icon"><ArrowRight /></el-icon>
              <span class="level-badge" :class="`level-${row.level}`">L{{ row.level }}</span>
              {{ row.name }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="是否Tab" width="100">
          <template #default="{ row }">
            <el-tag :type="row.isTab ? 'success' : 'info'" size="small">
              {{ row.isTab ? '是' : '否' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
            <el-button
              v-if="row.level < 3"
              type="success"
              link
              @click="handleAddChild(row)"
            >
              添加子分类
            </el-button>
            <el-button type="danger" link @click="handleDelete(row.id)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 新增/编辑弹窗 -->
    <el-dialog v-model="showDialog" :title="dialogTitle" width="500px">
      <el-form :model="form" label-width="80px">
        <el-form-item label="分类名称">
          <el-input v-model="form.name" placeholder="请输入分类名称" />
        </el-form-item>
        <el-form-item label="父分类">
          <el-select
            v-model="form.parent"
            placeholder="无（顶级分类）"
            :disabled="!!editingId"
            clearable
            @change="onParentChange"
          >
            <el-option :value="0" label="无（顶级分类）" />
            <el-option
              v-for="cat in allParentOptions"
              :key="cat.id"
              :label="cat.name"
              :value="cat.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="分类级别">
          <el-tag>{{ form.level }}</el-tag>
        </el-form-item>
        <el-form-item label="是否Tab">
          <el-switch
            :model-value="isTabValue"
            :active-value="true"
            :inactive-value="false"
            @change="isTabValue = $event"
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
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getCategoryList, createCategory, updateCategory, deleteCategory, type CategoryItem } from '@/api/category'

const categories = ref<CategoryItem[]>([])
const loading = ref(false)
const saving = ref(false)
const showDialog = ref(false)
const editingId = ref<number | null>(null)
const tableKey = ref(0)
const isTabValue = ref(false)

const form = reactive({
  name: '',
  parent: 0,
  level: 1
})

const dialogTitle = computed(() => {
  if (editingId.value) return '编辑分类'
  if (form.parent > 0) return '添加子分类'
  return '新增分类'
})

const loadCategories = async () => {
  loading.value = true
  try {
    const res = await getCategoryList()
    categories.value = (res as any) || []
  } finally {
    loading.value = false
  }
}

const handleAdd = () => {
  editingId.value = null
  form.name = ''
  form.parent = 0
  form.level = 1
  isTabValue.value = false
  showDialog.value = true
}

const handleAddChild = (row: CategoryItem) => {
  editingId.value = null
  form.name = ''
  form.parent = row.id
  form.level = calcLevel(row.id)
  isTabValue.value = false
  showDialog.value = true
}

const handleEdit = (row: CategoryItem) => {
  editingId.value = row.id
  form.name = row.name
  form.parent = row.parent || 0
  form.level = row.level
  isTabValue.value = row.isTab ?? false
  showDialog.value = true
}

// 递归查找分类（树形结构）
const findCategoryById = (list: CategoryItem[], id: number): CategoryItem | null => {
  for (const item of list) {
    if (item.id === id) return item
    if (item.sub_category?.length) {
      const found = findCategoryById(item.sub_category, id)
      if (found) return found
    }
  }
  return null
}

// 扁平化树形结构（用于父分类下拉选择）
const flattenTree = (list: CategoryItem[]): CategoryItem[] => {
  const result: CategoryItem[] = []
  for (const item of list) {
    result.push(item)
    if (item.sub_category?.length) {
      result.push(...flattenTree(item.sub_category))
    }
  }
  return result
}

// 父分类可选列表（所有分类，排除自身，且 level < 3）
const allParentOptions = computed(() => {
  const flat = flattenTree(categories.value)
  return flat.filter(c => c.id !== editingId.value && c.level < 3)
})

// 根据父分类自动计算级别
const calcLevel = (parentId: number): number => {
  if (parentId === 0) return 1
  const parent = findCategoryById(categories.value, parentId)
  return parent ? parent.level + 1 : 1
}

const onParentChange = (val: number | string | null) => {
  const parentId = typeof val === 'number' ? val : 0
  form.level = calcLevel(parentId || 0)
}

const handleSave = async () => {
  saving.value = true
  try {
    if (editingId.value) {
      await updateCategory(editingId.value, { name: form.name, is_tab: isTabValue.value })
      const target = findCategoryById(categories.value, editingId.value)
      if (target) {
        target.name = form.name
        target.isTab = isTabValue.value
      }
      ElMessage.success('更新成功')
    } else {
      const parentId = form.parent || 0
      const postData = {
        name: form.name.trim(),
        parent: parentId,
        level: calcLevel(parentId),
        is_tab: isTabValue.value
      }
      await createCategory(postData)
      ElMessage.success('创建成功')
    }
    showDialog.value = false
    await loadCategories()
    tableKey.value++
  } catch {
    /* 拦截器已统一处理错误提示 */
  } finally {
    saving.value = false
  }
}

const handleDelete = async (id: number) => {
  try {
    await ElMessageBox.confirm('确定删除该分类吗？删除后子分类将变为顶级分类', '提示', { type: 'warning' })
    await deleteCategory(id)
    ElMessage.success('删除成功')
    await loadCategories()
    tableKey.value++
  } catch (error) {
    // 用户取消或错误
  }
}

onMounted(() => {
  loadCategories()
})
</script>

<style lang="scss" scoped>
.page-header {
  margin-bottom: 20px;
}

.sub-icon {
  margin-right: 4px;
  font-size: 12px;
  color: $text-secondary;
}

.level-badge {
  display: inline-block;
  font-size: 11px;
  padding: 1px 6px;
  border-radius: 3px;
  margin-right: 6px;
  font-weight: 500;

  &.level-1 {
    background: #ecf5ff;
    color: #409eff;
  }

  &.level-2 {
    background: #f0f9eb;
    color: #67c23a;
  }

  &.level-3 {
    background: #fdf6ec;
    color: #e6a23c;
  }
}
</style>