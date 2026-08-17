<template>
  <div class="not-found-page">
    <div class="nf-container">
      <div class="nf-visual">
        <div class="nf-code">404</div>
        <div class="nf-icon">
          <el-icon :size="80"><Warning /></el-icon>
        </div>
      </div>
      <h2 class="nf-title">页面未找到</h2>
      <p class="nf-desc">抱歉，您访问的页面不存在或已被移除</p>
      <div class="nf-actions">
        <el-button type="primary" size="large" round @click="$router.push('/')">
          <el-icon><HomeFilled /></el-icon>
          返回首页
        </el-button>
        <el-button size="large" round @click="$router.go(-1)">
          <el-icon><ArrowLeft /></el-icon>
          返回上页
        </el-button>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.not-found-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f0f5ff 0%, #e8f0fe 50%, #f5f7fa 100%);
}

.nf-container {
  text-align: center;
  padding: 60px 40px;
  animation: fadeInUp 0.6s ease;
}

.nf-visual {
  position: relative;
  display: inline-block;
  margin-bottom: 24px;
}

.nf-code {
  font-size: 140px;
  font-weight: 800;
  line-height: 1;
  background: linear-gradient(135deg, $primary-color 0%, #a0cfff 60%, $primary-color 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  opacity: 0.15;
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  user-select: none;
}

.nf-icon {
  color: $primary-color;
  position: relative;
  z-index: 1;
  animation: float 3s ease-in-out infinite;
}

.nf-title {
  font-size: 24px;
  font-weight: 600;
  color: $text-primary;
  margin: 0 0 12px;
}

.nf-desc {
  font-size: 15px;
  color: $text-secondary;
  margin: 0 0 32px;
}

.nf-actions {
  display: flex;
  gap: 12px;
  justify-content: center;
  flex-wrap: wrap;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-12px);
  }
}
</style>