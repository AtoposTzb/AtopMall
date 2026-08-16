import { defineStore } from "pinia";
import { ref, computed } from "vue";
import {
  getCartList,
  addToCart,
  updateCartItem,
  deleteCartItem,
  type CartItem,
} from "@/api/cart";
import { ElMessage } from "element-plus";

export const useCartStore = defineStore("cart", () => {
  const cartList = ref<CartItem[]>([]);
  const loading = ref(false);
  const initialized = ref(false); // 是否已完成首次加载

  const totalCount = computed(() => {
    return cartList.value.reduce((sum, item) => sum + item.nums, 0);
  });

  const checkedCount = computed(() => {
    return cartList.value
      .filter((item) => item.checked)
      .reduce((sum, item) => sum + item.nums, 0);
  });

  const checkedTotalPrice = computed(() => {
    return cartList.value
      .filter((item) => item.checked)
      .reduce((sum, item) => sum + item.goods_price * item.nums, 0)
      .toFixed(2);
  });

  const isAllChecked = computed(() => {
    return (
      cartList.value.length > 0 && cartList.value.every((item) => item.checked)
    );
  });

  // 内部：仅刷新列表，不自动勾选
  const refreshCartList = async () => {
    try {
      const res = await getCartList();
      const data = (res as any)?.data;
      cartList.value = data && Array.isArray(data) ? data : [];
    } catch (error) {
      // 静默失败
    }
  };

  // 首次加载：自动勾选所有商品，仅执行一次
  const loadCartList = async () => {
    loading.value = true;
    try {
      const res = await getCartList();
      const data = (res as any)?.data;
      if (data && Array.isArray(data) && data.length > 0) {
        cartList.value = data;
        if (!initialized.value) {
          initialized.value = true;
          // 仅首次加载时自动勾选所有商品
          const uncheckedItems = cartList.value.filter((item) => !item.checked);
          if (uncheckedItems.length > 0) {
            await Promise.all(
              uncheckedItems.map((item) =>
                updateCartItem(item.goods_id, {
                  nums: item.nums,
                  checked: true,
                }).catch(() => {}),
              ),
            );
            await refreshCartList();
          }
        }
      } else {
        cartList.value = [];
        initialized.value = true;
      }
    } catch (error) {
      ElMessage.error("获取购物车失败");
      cartList.value = [];
    } finally {
      loading.value = false;
    }
  };

  // 添加商品到购物车
  const addCartItem = async (goodsId: number, nums: number) => {
    try {
      await addToCart({ goods: goodsId, nums });
      ElMessage.success("添加成功");
      await refreshCartList();
    } catch (error) {
      ElMessage.error("添加失败");
    }
  };

  // 更新购物车商品
  const updateCart = async (
    goodsId: number,
    data: { nums?: number; checked?: boolean },
  ) => {
    const currentItem = cartList.value.find(
      (item) => item.goods_id === goodsId,
    );
    // 始终携带 checked：若调用方显式传入则使用，否则保留当前状态
    const payload: { nums: number; checked: boolean } = {
      nums: data.nums ?? currentItem?.nums ?? 1,
      checked: data.checked !== undefined ? data.checked : (currentItem?.checked ?? true),
    };
    try {
      await updateCartItem(goodsId, payload);
      await refreshCartList();
    } catch (error) {
      ElMessage.error("更新失败");
    }
  };

  // 删除购物车商品
  const removeCartItem = async (goodsId: number) => {
    try {
      await deleteCartItem(goodsId);
      ElMessage.success("删除成功");
      await refreshCartList();
    } catch (error) {
      ElMessage.error("删除失败");
    }
  };

  // 全选/取消全选
  const toggleAllChecked = async (checked: boolean) => {
    const promises = cartList.value.map((item) =>
      updateCartItem(item.goods_id, { nums: item.nums, checked }),
    );
    await Promise.all(promises);
    await refreshCartList();
  };

  // 批量删除
  const batchDelete = async () => {
    const checkedItems = cartList.value.filter((item) => item.checked);
    const promises = checkedItems.map((item) => deleteCartItem(item.goods_id));
    await Promise.all(promises);
    ElMessage.success("删除成功");
    await refreshCartList();
  };

  return {
    cartList,
    loading,
    totalCount,
    checkedCount,
    checkedTotalPrice,
    isAllChecked,
    loadCartList,
    addCartItem,
    updateCart,
    removeCartItem,
    toggleAllChecked,
    batchDelete,
  };
});