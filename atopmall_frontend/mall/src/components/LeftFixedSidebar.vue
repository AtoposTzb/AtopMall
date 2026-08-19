<template>
  <div class="left-fixed-sidebar" :style="{ transform: `translateY(calc(-50% + ${offsetY}px))` }">
    <div class="lfs-drag-handle" @mousedown="onMouseDown">
      <span class="drag-line"></span>
      <span class="drag-line"></span>
      <span class="drag-line"></span>
    </div>
    <div class="lfs-item" @click="$router.push('/goods?ishot=1')">
      <div class="lfs-icon">
        <el-icon :size="20"><TrendCharts /></el-icon>
      </div>
      <span class="lfs-label">热销</span>
    </div>
    <div class="lfs-item" @click="$router.push('/goods?isnew=1')">
      <div class="lfs-icon">
        <el-icon :size="20"><Present /></el-icon>
      </div>
      <span class="lfs-label">新品</span>
    </div>
    <div class="lfs-item" @click="$router.push('/goods?istab=1')">
      <div class="lfs-icon">
        <el-icon :size="20"><Goods /></el-icon>
      </div>
      <span class="lfs-label">推荐</span>
    </div>
    <div class="lfs-item" @click="$router.push('/goods')">
      <div class="lfs-icon">
        <el-icon :size="20"><Discount /></el-icon>
      </div>
      <span class="lfs-label">特惠</span>
    </div>
    <div class="lfs-item" @click="$router.push('/goods')">
      <div class="lfs-icon">
        <el-icon :size="20"><More /></el-icon>
      </div>
      <span class="lfs-label">分类</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useDraggable } from '@/composables/useDraggable'

const { offsetY, onMouseDown } = useDraggable()
</script>

<style lang="scss" scoped>
.left-fixed-sidebar {
  position: fixed;
  left: 12px;
  top: 50%;
  transform: translateY(-50%);
  z-index: 999;
  display: flex;
  flex-direction: column;
  gap: 2px;
  background: #fff;
  border-radius: 10px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  padding: 0 0 6px 0;
  user-select: none;
  transition: box-shadow 0.2s;

  .lfs-drag-handle {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    gap: 3px;
    padding: 8px 0;
    cursor: grab;
    border-radius: 10px 10px 0 0;
    transition: background 0.2s;

    &:hover {
      background: #f5f7fa;
    }

    &:active {
      cursor: grabbing;
      background: #ecf5ff;
    }

    .drag-line {
      display: block;
      width: 18px;
      height: 2px;
      background: #c0c4cc;
      border-radius: 1px;
      transition: all 0.2s;
    }

    &:hover .drag-line {
      background: #909399;
    }
  }

  .lfs-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 4px;
    padding: 12px 10px;
    cursor: pointer;
    transition: all 0.2s;
    color: $text-regular;
    width: 64px;

    &:hover {
      color: $primary-color;
      background: #ecf5ff;
    }

    .lfs-icon {
      display: flex;
      align-items: center;
      justify-content: center;
      height: 24px;
    }

    .lfs-label {
      font-size: 11px;
      white-space: nowrap;
    }
  }
}

@media (max-width: $bp-mobile) {
  .left-fixed-sidebar {
    display: none;
  }
}
</style>