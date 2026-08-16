<template>
  <div class="category-nav" @mouseleave="activeCategory = null">
    <!-- 左侧分类按钮 -->
    <div class="category-trigger" @mouseenter="handleTriggerEnter">
      <el-icon><Menu /></el-icon>
      <span>全部商品分类</span>
      <el-icon class="arrow"><ArrowDown /></el-icon>
    </div>

    <!-- Mega Menu 面板 -->
    <div class="mega-menu" v-show="activeCategory !== null">
      <!-- 左侧一级分类列表 -->
      <div class="mega-menu-left">
        <div
          v-for="cat in categories"
          :key="cat.id"
          class="mega-menu-item"
          :class="{ active: activeCategory === cat.id }"
          @mouseenter="activeCategory = cat.id"
        >
          <span>{{ cat.name }}</span>
          <el-icon><ArrowRight /></el-icon>
        </div>
      </div>

      <!-- 右侧二三级分类展示 -->
      <div class="mega-menu-right" v-if="currentCategory">
        <div v-for="sub in currentCategory.sub_category" :key="sub.id" class="sub-category-row">
          <span class="sub-category-name">{{ sub.name }}</span>
          <div class="sub-category-items">
            <router-link
              v-for="item in sub.sub_category"
              :key="item.id"
              :to="{ path: '/goods', query: { ctg: item.id } }"
              class="sub-category-link"
              @click="activeCategory = null"
            >
              {{ item.name }}
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { getCategoryList, type CategoryItem } from '@/api/category'

const categories = ref<CategoryItem[]>([])
const activeCategory = ref<number | null>(null)
const loaded = ref(false)

const currentCategory = computed(() => {
  return categories.value.find(cat => cat.id === activeCategory.value)
})

const loadCategories = async () => {
  if (loaded.value) return
  try {
    const res = await getCategoryList()
    categories.value = res as unknown as CategoryItem[]
    loaded.value = true
  } catch (error) {
    console.error('加载分类失败', error)
  }
}

const handleTriggerEnter = async () => {
  await loadCategories()
  if (categories.value.length > 0) {
    activeCategory.value = categories.value[0].id
  }
}
</script>

<style lang="scss" scoped>
.category-nav {
  position: relative;
  display: inline-block;
}

.category-trigger {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: $primary-color;
  color: #fff;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  border-radius: 4px 4px 0 0;
  transition: background 0.2s;

  &:hover {
    background: #2c85e6;
  }

  .arrow {
    margin-left: 4px;
    transition: transform 0.2s;
  }

  &:hover .arrow {
    transform: rotate(180deg);
  }
}

.mega-menu {
  position: absolute;
  top: 100%;
  left: 0;
  display: flex;
  background: #fff;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
  border-radius: 0 0 4px 4px;
  z-index: 1000;
  min-height: 300px;
}

.mega-menu-left {
  width: 200px;
  background: #333;
  padding: 8px 0;

  .mega-menu-item {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 20px;
    color: #fff;
    font-size: 14px;
    cursor: pointer;
    transition: background 0.2s;

    &:hover,
    &.active {
      background: $primary-color;
    }
  }
}

.mega-menu-right {
  flex: 1;
  min-width: 500px;
  padding: 20px 24px;
  display: flex;
  flex-direction: column;
  gap: 16px;

  .sub-category-row {
    display: flex;
    align-items: flex-start;
    gap: 16px;

    .sub-category-name {
      font-weight: 600;
      font-size: 14px;
      color: $text-primary;
      min-width: 80px;
      padding-top: 2px;
    }

    .sub-category-items {
      display: flex;
      flex-wrap: wrap;
      gap: 8px 16px;
      flex: 1;
    }

    .sub-category-link {
      font-size: 13px;
      color: $text-regular;
      text-decoration: none;
      transition: color 0.2s;

      &:hover {
        color: $primary-color;
      }
    }
  }
}
</style>