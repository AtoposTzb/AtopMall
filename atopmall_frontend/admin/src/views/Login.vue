<template>
  <div class="login-page flex-center">
    <div class="login-wrapper">
      <div class="login-brand">
        <div class="brand-icon">
          <el-icon :size="48"><ShoppingBag /></el-icon>
        </div>
        <h1>AtopMall</h1>
        <p>管理后台</p>
      </div>
      <el-card class="login-card" shadow="always">
        <h2 class="title">管理员登录</h2>
        <el-form ref="formRef" :model="form" :rules="rules" label-width="0">
          <el-form-item prop="mobile">
            <el-input
              v-model="form.mobile"
              placeholder="请输入手机号"
              size="large"
              :prefix-icon="Phone"
            />
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              :prefix-icon="Lock"
              show-password
              @keyup.enter="handleLogin"
            />
          </el-form-item>
          <el-form-item prop="captcha">
            <div class="captcha-row">
              <el-input
                v-model="form.captcha"
                placeholder="请输入验证码"
                size="large"
                :prefix-icon="Key"
                maxlength="5"
                class="captcha-input"
              />
              <div class="captcha-img" @click="refreshCaptcha" title="点击刷新验证码">
                <img v-if="captchaImg" :src="captchaImg" alt="验证码" />
                <span v-else class="captcha-loading">加载中...</span>
              </div>
            </div>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="large" class="login-btn" :loading="loading" @click="handleLogin">
              登 录
            </el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance } from 'element-plus'
import { Phone, Lock, Key } from '@element-plus/icons-vue'
import { login, getCaptcha } from '@/api/user'
import { useUserStore } from '@/store/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)
const captchaImg = ref('')
const captchaId = ref('')

const form = reactive({
  mobile: '',
  password: '',
  captcha: ''
})

const rules = {
  mobile: [
    { required: true, message: '请输入手机号', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入正确的手机号', trigger: 'blur' }
  ],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }]
}

const refreshCaptcha = async () => {
  try {
    const res = await getCaptcha()
    captchaId.value = res.id
    captchaImg.value = res.picBase64
  } catch (error) {
    // 错误已在拦截器中处理
  }
}

const handleLogin = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  if (!captchaId.value) {
    ElMessage.warning('请先获取验证码')
    refreshCaptcha()
    return
  }

  loading.value = true
  try {
    const res = await login({
      mobile: form.mobile,
      password: form.password,
      captcha: form.captcha,
      captcha_id: captchaId.value
    })
    userStore.setToken(res.token)
    if (!userStore.isAdmin) {
      ElMessage.error('您不是管理员，无法登录管理后台')
      userStore.logout()
      return
    }
    await userStore.initUser()
    ElMessage.success('登录成功')
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (error) {
    // 登录失败时刷新验证码
    refreshCaptcha()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  refreshCaptcha()
})
</script>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 40%, #0f3460 100%);
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: -50%;
    left: -50%;
    width: 200%;
    height: 200%;
    background: radial-gradient(circle at 30% 50%, rgba(64, 158, 255, 0.08) 0%, transparent 50%),
                radial-gradient(circle at 70% 50%, rgba(103, 194, 58, 0.06) 0%, transparent 50%);
    animation: bgFloat 20s ease-in-out infinite;
  }

  @keyframes bgFloat {
    0%, 100% { transform: translate(0, 0); }
    50% { transform: translate(-2%, -2%); }
  }
}

.login-wrapper {
  position: relative;
  z-index: 1;
}

.login-brand {
  text-align: center;
  margin-bottom: 32px;
  color: #fff;

  .brand-icon {
    width: 80px;
    height: 80px;
    margin: 0 auto 16px;
    background: rgba(255, 255, 255, 0.15);
    border-radius: 20px;
    display: flex;
    align-items: center;
    justify-content: center;
    backdrop-filter: blur(10px);
    color: #409eff;
  }

  h1 {
    font-size: 28px;
    font-weight: 700;
    margin-bottom: 4px;
    letter-spacing: 2px;
  }

  p {
    font-size: 14px;
    opacity: 0.7;
    letter-spacing: 4px;
  }
}

.login-card {
  width: 400px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.98);

  .title {
    text-align: center;
    margin-bottom: 28px;
    color: $text-primary;
    font-size: 20px;
    font-weight: 600;
  }

  .login-btn {
    width: 100%;
    height: 44px;
    font-size: 16px;
    letter-spacing: 4px;
  }
}

.captcha-row {
  display: flex;
  gap: 12px;
  align-items: center;

  .captcha-input {
    flex: 1;
  }

  .captcha-img {
    width: 120px;
    height: 40px;
    border-radius: 4px;
    overflow: hidden;
    cursor: pointer;
    flex-shrink: 0;
    border: 1px solid $border-base;
    background: #f5f7fa;
    display: flex;
    align-items: center;
    justify-content: center;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    .captcha-loading {
      font-size: 12px;
      color: $text-secondary;
    }
  }
}
</style>