<template>
  <div class="goods-list-page">

    <div class="container goods-container">
      <div class="page-header">
        <h2 class="page-title">{{ searchParams.q ? `搜索结果：${searchParams.q}` : '全部商品' }}</h2>
        <p class="page-subtitle">发现优质好物，尽在 AtopMall</p>
      </div>

      <div class="search-bar">
        <el-input
          v-model="searchParams.q"
          placeholder="搜索商品名称、描述..."
          size="large"
          clearable
          @keyup.enter="handleSearch"
          class="search-input"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
          <template #suffix>
            <div class="goods-search-suffix" @click="handleSearch">
              <el-icon :size="20"><Search /></el-icon>
            </div>
          </template>
        </el-input>
      </div>

      <div class="goods-layout">
        <div class="sidebar">
          <div class="sidebar-section">
            <div class="section-title">商品分类</div>
            <div class="category-tree">
              <div
                v-for="cat in categories"
                :key="cat.id"
                class="category-item"
              >
                <div
                  class="category-title"
                  :class="{ active: searchParams.c === cat.id || activeCategory === cat.id }"
                  @click="handleCategoryClick(cat.id)"
                >
                  <el-icon class="expand-icon"><ArrowDown v-if="expandedCategories.includes(cat.id)" /><ArrowRight v-else /></el-icon>
                  <span class="cat-name">{{ cat.name }}</span>
                  <span v-if="cat.sub_category?.length" class="cat-count">{{ cat.sub_category.length }}</span>
                </div>
                <div v-if="expandedCategories.includes(cat.id) && cat.sub_category?.length" class="sub-category-list">
                  <div
                    v-for="sub in cat.sub_category"
                    :key="sub.id"
                    class="sub-category-item"
                  >
                    <div
                      class="sub-category-title"
                      :class="{ active: searchParams.c === sub.id || activeSubCategory === sub.id }"
                      @click="handleSubCategoryClick(sub.id)"
                    >
                      <el-icon v-if="sub.sub_category?.length" class="expand-icon"><ArrowDown v-if="expandedSubCategories.includes(sub.id)" /><ArrowRight v-else /></el-icon>
                      <span class="sub-name" :class="{ 'no-children': !sub.sub_category?.length }">{{ sub.name }}</span>
                    </div>
                    <div v-if="expandedSubCategories.includes(sub.id) && sub.sub_category?.length" class="third-category-list">
                      <div
                        v-for="item in sub.sub_category"
                        :key="item.id"
                        class="third-category-link"
                        :class="{ active: searchParams.c === item.id }"
                        @click="selectCategory(item.id)"
                      >
                        {{ item.name }}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="sidebar-section">
            <div class="section-title">价格区间</div>
            <div class="price-section">
              <div class="price-inputs">
                <el-input-number
                  v-model="searchParams.pmin"
                  :min="0"
                  placeholder="最低价"
                  size="default"
                  controls-position="right"
                  class="price-input"
                />
                <span class="price-separator">-</span>
                <el-input-number
                  v-model="searchParams.pmax"
                  :min="0"
                  placeholder="最高价"
                  size="default"
                  controls-position="right"
                  class="price-input"
                />
              </div>
              <el-button type="primary" size="small" @click="handleSearch">确定</el-button>
            </div>
          </div>

          <div class="sidebar-section">
            <div class="section-title">快捷筛选</div>
            <div class="filter-options">
              <el-checkbox v-model="hotChecked" @change="handleFilterChange" border>
                <el-icon><TrendCharts /></el-icon>
                热销商品
              </el-checkbox>
              <el-checkbox v-model="newChecked" @change="handleFilterChange" border>
                <el-icon><Star /></el-icon>
                新品上市
              </el-checkbox>
            </div>
          </div>
        </div>

        <div class="main-content">
          <div class="filter-bar">
          <div class="filter-left">
            <el-tag
              :type="!searchParams.c ? 'primary' : 'info'"
              class="filter-tag"
              effect="dark"
              @click="
                searchParams.c = undefined;
                handleSearch();
              "
            >
              全部
            </el-tag>
            <el-tag
              v-for="cat in flatCategories.slice(0, 6)"
              :key="cat.id"
              :type="searchParams.c === cat.id ? 'primary' : 'info'"
              class="filter-tag"
              effect="dark"
              @click="
                searchParams.c = cat.id;
                handleSearch();
              "
            >
              {{ cat.name }}
            </el-tag>
          </div>
            <div class="filter-right">
              <span class="result-count">共 {{ total }} 件商品</span>
            </div>
          </div>

          <div class="goods-grid" v-loading="loading">
            <div
              v-for="goods in goodsList"
              :key="goods.id"
              class="goods-card"
              @click="goDetail(goods.id)"
            >
              <div class="goods-image-wrapper">
                <img :src="goods.front_image" :alt="goods.name" class="goods-img" />
                <div class="goods-tags">
                  <span v-if="goods.is_new" class="tag tag-new">新品</span>
                  <span v-if="goods.is_hot" class="tag tag-hot">热销</span>
                </div>
              </div>
              <div class="goods-info">
                <h3 class="goods-name text-ellipsis-2">{{ goods.name }}</h3>
                <p class="goods-brief text-ellipsis">{{ goods.goods_brief }}</p>
                <div class="goods-footer">
                  <span class="price">{{ goods.shop_price }}</span>
                </div>
              </div>
            </div>
          </div>

          <div v-if="!loading && goodsList.length === 0" class="empty">
            <template v-if="loadError">
              <el-icon class="empty-icon"><WarningFilled /></el-icon>
              <p>加载失败</p>
              <el-button type="primary" @click="loadGoods">重新加载</el-button>
            </template>
            <template v-else>
              <el-icon class="empty-icon"><GoodsFilled /></el-icon>
              <p>暂无商品</p>
              <el-button type="primary" @click="handleSearch">重新搜索</el-button>
            </template>
          </div>

          <div class="pagination" v-if="total > 0">
            <el-pagination
              v-model:current-page="searchParams.p"
              :page-size="searchParams.pnum"
              :total="total"
              layout="prev, pager, next, jumper, total"
              @current-change="loadGoods"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from "vue";
import { useRouter, useRoute } from "vue-router";
import { getGoodsList, type GoodsItem } from "@/api/goods";
import { getCategoryList, type CategoryItem } from "@/api/category";


const router = useRouter();
const route = useRoute();

const goodsList = ref<GoodsItem[]>([]);
const categories = ref<CategoryItem[]>([]);
const loading = ref(false);
const loadError = ref(false);
const total = ref(0);
const hotChecked = ref(false);
const newChecked = ref(false);
const activeCategory = ref<number | null>(null);
const activeSubCategory = ref<number | null>(null);
const expandedCategories = ref<number[]>([]);
const expandedSubCategories = ref<number[]>([]);

const searchParams = reactive({
  p: 1,
  pnum: 12,
  q: "",
  c: undefined as number | undefined,
  pmin: undefined as number | undefined,
  pmax: undefined as number | undefined,
  ishot: undefined as number | undefined,
  isnew: undefined as number | undefined,
  istab: undefined as number | undefined,
});

const flatCategories = computed(() => {
  const result: CategoryItem[] = [];
  const flatten = (items: CategoryItem[]) => {
    items.forEach((item) => {
      result.push(item);
      if (item.sub_category?.length) {
        flatten(item.sub_category);
      }
    });
  };
  flatten(categories.value);
  return result;
});

const toggleCategory = (id: number) => {
  if (expandedCategories.value.includes(id)) {
    expandedCategories.value = expandedCategories.value.filter(catId => catId !== id);
  } else {
    expandedCategories.value.push(id);
  }
  activeCategory.value = id;
};

const toggleSubCategory = (id: number) => {
  if (expandedSubCategories.value.includes(id)) {
    expandedSubCategories.value = expandedSubCategories.value.filter(subId => subId !== id);
  } else {
    expandedSubCategories.value.push(id);
  }
  activeSubCategory.value = id;
};

// 点击一级分类：展开/折叠，同时设为筛选条件
const handleCategoryClick = (id: number) => {
  toggleCategory(id)
  searchParams.c = id
  handleSearch()
}

// 点击二级分类：展开/折叠（如有子分类），同时设为筛选条件
const handleSubCategoryClick = (id: number) => {
  toggleSubCategory(id)
  searchParams.c = id
  handleSearch()
}

// 点击三级分类：设为筛选条件
const selectCategory = (id: number) => {
  searchParams.c = id
  handleSearch()
}

const goDetail = (id: number) => {
  router.push(`/goods/${id}`);
};

const loadGoods = async () => {
  const params: any = {
    p: searchParams.p,
    pnum: searchParams.pnum,
  };
  if (searchParams.q) params.q = searchParams.q;
  if (searchParams.c !== undefined) params.c = searchParams.c;
  if (searchParams.pmin !== undefined) params.pmin = searchParams.pmin;
  if (searchParams.pmax !== undefined) params.pmax = searchParams.pmax;
  if (searchParams.ishot !== undefined) params.ishot = searchParams.ishot;
  if (searchParams.isnew !== undefined) params.isnew = searchParams.isnew;
  if (searchParams.istab !== undefined) params.istab = searchParams.istab;

  loading.value = true;
  loadError.value = false;
  try {
    const res = await getGoodsList(params);
    const allGoods = (res as any).data || []
    const onSaleGoods = allGoods.filter((g: GoodsItem) => g.on_sale)
    goodsList.value = onSaleGoods
    total.value = onSaleGoods.length
  } catch (error) {
    console.error("加载商品列表失败", error);
    loadError.value = true;
  } finally {
    loading.value = false;
  }
};

const loadCategories = async () => {
  try {
    const res = await getCategoryList();
    categories.value = res as unknown as CategoryItem[];
  } catch (error) {
    console.error("加载分类失败", error);
  }
};

const handleSearch = () => {
  searchParams.p = 1;
  loadGoods();
};

const handleFilterChange = () => {
  searchParams.ishot = hotChecked.value ? 1 : undefined;
  searchParams.isnew = newChecked.value ? 1 : undefined;
  handleSearch();
};

const syncFromRoute = () => {
  if (route.query.ctg) {
    searchParams.c = Number(route.query.ctg)
  } else if (route.query.c) {
    searchParams.c = Number(route.query.c)
  } else {
    searchParams.c = undefined
  }
  if (route.query.ishot) {
    searchParams.ishot = Number(route.query.ishot)
    hotChecked.value = true
  } else {
    searchParams.ishot = undefined
    hotChecked.value = false
  }
  if (route.query.isnew) {
    searchParams.isnew = Number(route.query.isnew)
    newChecked.value = true
  } else {
    searchParams.isnew = undefined
    newChecked.value = false
  }
  if (route.query.istab) {
    searchParams.istab = Number(route.query.istab)
  } else {
    searchParams.istab = undefined
  }
  if (route.query.skey) {
    searchParams.q = String(route.query.skey)
  } else if (route.query.q) {
    searchParams.q = String(route.query.q)
  } else {
    searchParams.q = ''
  }
  searchParams.p = 1
  activeCategory.value = searchParams.c ?? null
  activeSubCategory.value = null
  loadGoods()
}

onMounted(() => {
  syncFromRoute()
  loadCategories()
})

watch(
  () => route.query,
  () => {
    syncFromRoute()
  }
)
</script>

<style lang="scss" scoped>
.page-header {
  text-align: center;
  padding: 32px 0 24px;

  .page-title {
    font-size: 32px;
    font-weight: bold;
    margin-bottom: 8px;
    color: $text-primary;
  }

  .page-subtitle {
    font-size: 16px;
    color: $text-secondary;
  }
}

.search-bar {
  margin-bottom: 24px;

  .search-input {
    :deep(.el-input__wrapper) {
      border-radius: 24px;
      box-shadow: 0 2px 12px rgba(0, 0, 0, 0.06);
      transition: all 0.3s;
      padding-right: 0;

      &:hover {
        box-shadow: 0 4px 16px rgba(64, 158, 255, 0.12);
      }

      &.is-focus {
        box-shadow: 0 4px 20px rgba(64, 158, 255, 0.2);
      }
    }
  }
}

.goods-layout {
  display: flex;
  gap: 24px;
}

.sidebar {
  width: 220px;
  flex-shrink: 0;

  .sidebar-section {
    background: #fff;
    border-radius: 8px;
    padding: 16px;
    margin-bottom: 16px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

    .section-title {
      font-weight: 600;
      font-size: 15px;
      color: $text-primary;
      padding-bottom: 12px;
      margin-bottom: 12px;
      border-bottom: 2px solid $primary-color;
    }
  }
}

.category-tree {
  .category-item {
    margin-bottom: 2px;

    .category-title {
      display: flex;
      align-items: center;
      gap: 4px;
      padding: 10px 8px;
      cursor: pointer;
      border-radius: 4px;
      font-size: 14px;
      font-weight: 500;
      color: $text-primary;
      transition: all 0.2s;

      &:hover {
        background: #ecf5ff;
        color: $primary-color;
      }

      &.active {
        background: $primary-color;
        color: #fff;

        .cat-count {
          background: rgba(255, 255, 255, 0.3);
        }
      }

      .expand-icon {
        font-size: 12px;
        flex-shrink: 0;
      }

      .cat-name {
        flex: 1;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .cat-count {
        flex-shrink: 0;
        font-size: 11px;
        background: #f0f2f5;
        color: $text-secondary;
        padding: 1px 6px;
        border-radius: 10px;
        min-width: 18px;
        text-align: center;
      }
    }
  }

  .sub-category-list {
    padding-left: 16px;

    .sub-category-item {
      margin-bottom: 1px;

      .sub-category-title {
        display: flex;
        align-items: center;
        gap: 4px;
        padding: 8px 8px;
        cursor: pointer;
        border-radius: 4px;
        font-size: 13px;
        color: $text-regular;
        transition: all 0.2s;

        &:hover {
          background: #f5f7fa;
          color: $primary-color;
        }

        &.active {
          background: #ecf5ff;
          color: $primary-color;
          font-weight: 500;
        }

        .expand-icon {
          font-size: 11px;
          flex-shrink: 0;
        }

        .sub-name {
          flex: 1;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;

          &.no-children {
            padding-left: 15px; // 无子分类时对齐图标位置
          }
        }
      }
    }
  }

  .third-category-list {
    padding-left: 20px;
    display: flex;
    flex-direction: column;
    gap: 1px;

    .third-category-link {
      padding: 7px 12px;
      font-size: 13px;
      color: $text-secondary;
      cursor: pointer;
      border-radius: 4px;
      transition: all 0.2s;
      position: relative;

      &::before {
        content: '';
        position: absolute;
        left: 0;
        top: 50%;
        transform: translateY(-50%);
        width: 4px;
        height: 4px;
        border-radius: 50%;
        background: #d0d5dd;
        transition: background 0.2s;
      }

      &:hover {
        background: #f0f5ff;
        color: $primary-color;

        &::before {
          background: $primary-color;
        }
      }

      &.active {
        background: #ecf5ff;
        color: $primary-color;
        font-weight: 500;

        &::before {
          background: $primary-color;
        }
      }
    }
  }
}

.price-section {
  display: flex;
  flex-direction: column;
  gap: 12px;

  .price-inputs {
    display: flex;
    align-items: center;
    gap: 8px;

    .price-input {
      flex: 1;
    }

    .price-separator {
      color: $text-secondary;
      font-size: 16px;
    }
  }
}

.filter-options {
  display: flex;
  flex-direction: column;
  gap: 12px;

  .el-checkbox {
    display: flex;
    align-items: center;
    gap: 4px;
  }
}

.main-content {
  flex: 1;
}

.filter-bar {
  background: #fff;
  border-radius: 8px;
  padding: 16px 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  display: flex;
  justify-content: space-between;
  align-items: center;

  .filter-left {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;

    .filter-tag {
      cursor: pointer;
      transition: all 0.2s;

      &:hover {
        transform: translateY(-1px);
      }
    }
  }

  .filter-right {
    .result-count {
      font-size: 14px;
      color: $text-secondary;
    }
  }
}

.goods-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 20px;
  min-height: 200px;
}

.goods-card {
  background: #fff;
  border-radius: 8px;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.3s;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

  &:hover {
    transform: translateY(-6px);
    box-shadow: 0 12px 24px rgba(0, 0, 0, 0.12);
  }

  .goods-image-wrapper {
    position: relative;
    overflow: hidden;

    .goods-img {
      width: 100%;
      height: 220px;
      object-fit: cover;
      transition: transform 0.3s;
    }

    &:hover .goods-img {
      transform: scale(1.08);
    }

    .goods-tags {
      position: absolute;
      top: 12px;
      left: 12px;
      display: flex;
      gap: 6px;

      .tag {
        padding: 4px 10px;
        border-radius: 4px;
        font-size: 12px;
        font-weight: 600;
        color: #fff;

        &.tag-new {
          background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
        }

        &.tag-hot {
          background: linear-gradient(135deg, #ff6b6b 0%, #ee5a6f 100%);
        }
      }
    }
  }

  .goods-info {
    padding: 16px;

    .goods-name {
      font-size: 15px;
      font-weight: 500;
      margin-bottom: 8px;
      line-height: 1.4;
      color: $text-primary;
    }

    .goods-brief {
      font-size: 13px;
      color: $text-secondary;
      margin-bottom: 12px;
    }

    .goods-footer {
      display: flex;
      align-items: baseline;
      gap: 8px;

      .price {
        font-size: 22px;
        font-weight: bold;
        color: $danger-color;
      }

      .market-price {
        font-size: 13px;
        color: $text-secondary;
        text-decoration: line-through;
      }
    }
  }
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

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 40px;
}
</style>

<style lang="scss">
.goods-search-suffix {
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #409eff;
  padding: 0 14px;
  transition: all 0.2s;
  border-radius: 0 24px 24px 0;

  &:hover {
    background: rgba(64, 158, 255, 0.1);
    color: #337ecc;
  }

  &:active {
    background: rgba(64, 158, 255, 0.18);
  }
}
</style>