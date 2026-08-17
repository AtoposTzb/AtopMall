<template>
  <div class="user-favorite">
    <h3 class="mb-md">我的收藏</h3>

    <div v-loading="loading" class="fav-grid">
      <div
        v-for="item in favorites"
        :key="item.id"
        class="fav-item card"
        @click="goToDetail(item)"
      >
        <div class="fav-image">
          <img
            v-if="item.front_image"
            :src="item.front_image"
            :alt="item.name"
          />
          <div v-else class="image-placeholder">
            <el-icon><Picture /></el-icon>
          </div>
        </div>
        <div class="fav-content">
          <h4 class="fav-name">{{ item.name }}</h4>
          <p v-if="item.goods_brief" class="fav-brief">
            {{ item.goods_brief }}
          </p>
          <div class="fav-footer">
            <span class="price">¥{{ item.shop_price }}</span>
            <el-button
              type="danger"
              size="small"
              @click.stop="handleRemove(item.goods_id || item.id)"
            >
              取消收藏
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="!loading && favorites.length === 0" class="empty">
      <el-icon class="empty-icon"><Star /></el-icon>
      <p>暂无收藏</p>
      <el-button type="primary" @click="$router.push('/goods')"
        >去逛逛</el-button
      >
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import { getFavList, deleteFav, type FavItem } from "@/api/favorite";

const router = useRouter();
const favorites = ref<FavItem[]>([]);
const loading = ref(false);

const loadFavorites = async () => {
  loading.value = true;
  try {
    const res = await getFavList();
    favorites.value = (res as any).data || [];
  } finally {
    loading.value = false;
  }
};

const goToDetail = (item: FavItem) => {
  const goodsId = item.goods_id || item.id;
  router.push(`/goods/${goodsId}`);
};

const handleRemove = async (goodsId: number) => {
  try {
    await deleteFav(goodsId);
    ElMessage.success("已取消收藏");
    await loadFavorites();
  } catch (error) {
    // 错误已在拦截器中处理
  }
};

onMounted(() => {
  loadFavorites();
});
</script>

<style lang="scss" scoped>
.user-favorite {
  h3 {
    font-size: 20px;
    font-weight: 600;
    color: $text-primary;
    margin-bottom: 28px;
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

.fav-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.fav-item {
  display: flex;
  padding: 16px;
  cursor: pointer;
  transition: all 0.3s;
  border-radius: 12px;
  border: 1px solid #f0f0f0;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
    border-color: transparent;
  }
}

.fav-image {
  width: 120px;
  height: 120px;
  flex-shrink: 0;
  margin-right: 16px;
  border-radius: 8px;
  overflow: hidden;
  background: #f5f7fa;

  img {
    width: 100%;
    height: 100%;
    object-fit: cover;
  }

  .image-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #c0c4cc;
    font-size: 32px;
  }
}

.fav-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  min-width: 0;

  .fav-name {
    font-size: 15px;
    font-weight: 500;
    margin-bottom: 8px;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
  }

  .fav-brief {
    font-size: 13px;
    color: $text-secondary;
    margin-bottom: 8px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .fav-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;

    .price {
      font-size: 18px;
      font-weight: bold;
      color: $danger-color;
    }
  }
}
</style>