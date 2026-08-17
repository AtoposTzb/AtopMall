<template>
  <div class="register-page">
    <div class="register-wrapper">
      <div class="register-card">
        <div class="card-header">
          <router-link to="/" class="logo-link">
            <el-icon :size="32"><Shop /></el-icon>
            <span class="logo-text">AtopMall</span>
          </router-link>
          <p class="subtitle">创建您的账号，开启购物之旅</p>
        </div>
        <el-form ref="formRef" :model="form" :rules="rules" label-width="0" class="register-form">
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
          <el-form-item prop="email">
            <el-input
              v-model="form.email"
              placeholder="请输入邮箱"
              size="large"
              class="custom-input"
            >
              <template #prefix>
                <el-icon><Message /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="code">
            <div class="code-row">
              <el-input
                v-model="form.code"
                placeholder="验证码"
                size="large"
                class="custom-input code-input"
              >
                <template #prefix>
                  <el-icon><Key /></el-icon>
                </template>
              </el-input>
              <el-button
                type="primary"
                size="large"
                :disabled="countdown > 0"
                @click="handleSendCode"
                class="code-btn"
              >
                {{ countdown > 0 ? `${countdown}s` : '发送验证码' }}
              </el-button>
            </div>
          </el-form-item>
          <el-form-item prop="password">
            <el-input
              v-model="form.password"
              type="password"
              placeholder="请输入密码（至少6位）"
              size="large"
              show-password
              class="custom-input"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item prop="confirmPassword">
            <el-input
              v-model="form.confirmPassword"
              type="password"
              placeholder="请确认密码"
              size="large"
              show-password
              class="custom-input"
            >
              <template #prefix>
                <el-icon><Lock /></el-icon>
              </template>
            </el-input>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" size="large" class="submit-btn" :loading="loading" @click="handleRegister">
              {{ loading ? '注册中...' : '注 册' }}
            </el-button>
          </el-form-item>
        </el-form>
        <div class="card-footer">
          <router-link to="/login" class="footer-link">已有账号？去登录</router-link>
          <router-link to="/" class="footer-link">返回首页</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance } from 'element-plus'
import { register, sendEmailCode } from '@/api/user'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const loading = ref(false)
const countdown = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

const form = reactive({
  mobile: '',
  email: '',
  code: '',
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (value !== form.password) {
    callback(new Error('两次密码输入不一致'))
  } else {
    callback()
  }
}

const rules = {
  mobile: [{ required: true, message: '请输入手机号', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  code: [{ required: true, message: '请输入验证码', trigger: 'blur' }],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
}

const handleSendCode = async () => {
  if (!form.email) {
    ElMessage.warning('请先输入邮箱')
    return
  }
  try {
    await sendEmailCode({ email: form.email, type: 1 })
    ElMessage.success('验证码已发送')
    countdown.value = 60
    timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0 && timer) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error) {
    // 错误已在拦截器中处理
  }
}

const handleRegister = async () => {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return

  loading.value = true
  try {
    const res = await register({
      mobile: form.mobile,
      password: form.password,
      email: form.email,
      code: form.code
    })
    userStore.setToken(res.token)
    ElMessage.success('注册成功')
    router.push('/')
  } catch (error) {
    // 错误已在拦截器中处理
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.register-page {
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

.register-wrapper {
  position: relative;
  z-index: 1;
}

.register-card {
  width: 420px;
  background: #fff;
  border-radius: 16px;
  padding: 40px 40px 32px;
  box-shadow: 0 20px 60px rgba(64, 158, 255, 0.12), 0 4px 16px rgba(0, 0, 0, 0.06);
  animation: cardSlideUp 0.5s ease;

  .card-header {
    text-align: center;
    margin-bottom: 32px;

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

  .register-form {
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

    .code-row {
      display: flex;
      gap: 10px;
      width: 100%;

      .code-input {
        flex: 1;
      }

      .code-btn {
        height: 40px;
        border-radius: 10px;
        white-space: nowrap;
        flex-shrink: 0;
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
</style>

<style lang="scss">
.register-page .el-form-item__error {
  padding-left: 4px;
}
</style>