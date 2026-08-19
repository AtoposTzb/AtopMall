<template>
  <div class="login-page">
    <div class="login-wrapper">
      <div class="login-card">
        <div class="card-header">
          <router-link to="/" class="logo-link">
            <el-icon :size="32"><Shop /></el-icon>
            <span class="logo-text">AtopMall</span>
          </router-link>
          <p class="subtitle">欢迎回来，请登录您的账号</p>
        </div>
        <el-form ref="formRef" :model="form" :rules="rules" label-width="0" class="login-form">
          <el-form-item prop="mobile">
            <el-input
              v-model="form.mobile"
              placeholder="请输入手机号"
              size="large"
              class="custom-input"
            >
              <template #prefix>
                <el-icon><Phone /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              show-password
              class="custom-input"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="captcha">
            <div class="captcha-row">
              <el-input
                v-model="form.captcha"
                placeholder="验证码"
                size="large"
                class="custom-input captcha-input"
              >
                <template #prefix>
                  <el-icon><Key /></el-icon>
                </template>
              </el-input>
              <img
                v-if="captchaImage"
                :src="captchaImage"
                class="captcha-img"
                @click="loadCaptcha"
                title="点击刷新验证码"
              />
              <el-button v-else type="primary" size="large" @click="loadCaptcha" class="captcha-btn">获取验证码</el-button>
            </div>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="large" class="submit-btn" :loading="loading" @click="handleLogin">
              {{ loading ? '登录中...' : '登 录' }}
            </el-button>
          </el-form-item>
        </el-form>
        <div class="card-footer">
          <router-link to="/register" class="footer-link">还没有账号？立即注册</router-link>
          <router-link to="/" class="footer-link home-link">返回首页</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance } from 'element-plus'
import { login, getCaptcha } from '@/api/user'
import { useUserStore } from '@/store/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)
const captchaImage = ref('')
const captchaId = ref('')

const form = reactive({
  mobile: '',
  password: '',
  captcha: '',
  captcha_id: ''
})

const rules = {
  mobile: [{ required: true, message: '请输入手机号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }],
  captcha: [{ required: true, message: '请输入验证码', trigger: 'blur' }]
}

const loadCaptcha = async () => {
  try {
    const res = await getCaptcha()
    captchaImage.value = res.picBase64
    captchaId.value = res.id
  } catch (error) {
    ElMessage.error('获取验证码失败')
  }
}

const handleLogin = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    form.captcha_id = captchaId.value
    const res = await login(form)
    userStore.setToken(res.token)
    await userStore.initUser()
    ElMessage.success('登录成功')
    const redirect = (route.query.redirect as string) || '/'
    router.push(redirect)
  } catch (error) {
    loadCaptcha()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadCaptcha()
})
</script>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #e8f0fe 0%, #d4e4fc 30%, #c5d9f8 60%, #ecf5ff 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: -120px;
    right: -120px;
    width: 500px;
    height: 500px;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(64, 158, 255, 0.08) 0%, transparent 70%);
  }

  &::after {
    content: '';
    position: absolute;
    bottom: -80px;
    left: -80px;
    width: 360px;
    height: 360px;
    border-radius: 50%;
    background: radial-gradient(circle, rgba(64, 158, 255, 0.06) 0%, transparent 70%);
  }
}

.login-wrapper {
  position: relative;
  z-index: 1;
}

.login-card {
  width: 420px;
  background: #fff;
  border-radius: 16px;
  padding: 44px 40px 36px;
  box-shadow: 0 20px 60px rgba(64, 158, 255, 0.12), 0 4px 16px rgba(0, 0, 0, 0.06);
  animation: cardSlideUp 0.5s ease;

  .card-header {
    text-align: center;
    margin-bottom: 36px;

    .logo-link {
      display: inline-flex;
      align-items: center;
      gap: 10px;
      text-decoration: none;
      color: $primary-color;
      margin-bottom: 12px;

      .logo-text {
        font-size: 28px;
        font-weight: 700;
        background: linear-gradient(135deg, $primary-color 0%, #66b1ff 100%);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        background-clip: text;
      }
    }

    .subtitle {
      font-size: 14px;
      color: $text-secondary;
      margin: 0;
    }
  }

  .login-form {
    .custom-input {
      :deep(.el-input__wrapper) {
        border-radius: 10px;
        box-shadow: 0 0 0 1px $border-light inset;
        transition: all 0.3s;
        padding: 0 16px;

        &:hover {
          box-shadow: 0 0 0 1px $primary-color inset;
        }

        &.is-focus {
          box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.25) inset;
        }
      }
    }

    .captcha-row {
      display: flex;
      gap: 10px;
      width: 100%;

      .captcha-input {
        flex: 1;
      }

      .captcha-img {
        height: 40px;
        width: 110px;
        cursor: pointer;
        border-radius: 10px;
        border: 1px solid $border-light;
        object-fit: cover;
        transition: all 0.2s;

        &:hover {
          border-color: $primary-color;
          box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.15);
        }
      }

      .captcha-btn {
        height: 40px;
        border-radius: 10px;
        white-space: nowrap;
      }
    }

    .submit-btn {
      width: 100%;
      height: 46px;
      border-radius: 10px;
      font-size: 16px;
      font-weight: 600;
      letter-spacing: 4px;
      background: linear-gradient(135deg, $primary-color 0%, #66b1ff 100%);
      border: none;
      transition: all 0.3s;
      margin-top: 4px;

      &:hover:not(:disabled) {
        transform: translateY(-1px);
        box-shadow: 0 8px 24px rgba(64, 158, 255, 0.35);
      }

      &:active:not(:disabled) {
        transform: translateY(0);
      }
    }
  }

  .card-footer {
    display: flex;
    justify-content: space-between;
    margin-top: 20px;

    .footer-link {
      font-size: 13px;
      color: $text-secondary;
      text-decoration: none;
      transition: color 0.2s;

      &:hover {
        color: $primary-color;
      }
    }
  }
}

@keyframes cardSlideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: $bp-mobile) {
  .login-card {
    width: 90%;
    padding: 32px 24px;
  }
}
</style>

<style lang="scss">
.login-page .el-form-item__error {
  padding-left: 4px;
}
</style>