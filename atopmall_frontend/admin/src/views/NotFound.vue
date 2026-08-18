<template>
  <div class="not-found-page flex-center">
    <div class="content">
      <h1>404</h1>
      <p>抱歉，您访问的页面不存在</p>
      <el-button type="primary" @click="$router.push('/')">返回首页</el-button>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.not-found-page {
  min-height: 100vh;

  .content {
    text-align: center;

    h1 {
      font-size: 80px;
      color: $primary-color;
      margin-bottom: 16px;
    }

    p {
      font-size: 18px;
      color: $text-secondary;
      margin-bottom: 24px;
    }
  }
}
</style>
