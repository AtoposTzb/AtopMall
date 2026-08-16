<template>
  <div class="register-page flex-center">
    <el-card class="register-card">
      <h2 class="title">用户注册</h2>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="0">
        <el-form-item prop="mobile">
          <el-input v-model="form.mobile" placeholder="请输入手机号" prefix-icon="Phone" size="large" />
        </el-form-item>
        <el-form-item prop="email">
          <el-input v-model="form.email" placeholder="请输入邮箱" prefix-icon="Message" size="large" />
        </el-form-item>
        <el-form-item prop="code">
          <div class="code-input flex">
            <el-input v-model="form.code" placeholder="验证码" size="large" />
            <el-button type="primary" size="large" :disabled="countdown > 0" @click="handleSendCode">
              {{ countdown > 0 ? `${countdown}s` : '发送验证码' }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="form.password" type="password" placeholder="请输入密码" prefix-icon="Lock" size="large" show-password />
        </el-form-item>
        <el-form-item prop="confirmPassword">
          <el-input v-model="form.confirmPassword" type="password" placeholder="请确认密码" prefix-icon="Lock" size="large" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" class="register-btn" :loading="loading" @click="handleRegister">
            注 册
          </el-button>
        </el-form-item>
      </el-form>
      <div class="footer flex-between">
        <el-button type="primary" link @click="goHome">返回首页</el-button>
        <router-link to="/login">已有账号？去登录</router-link>
      </div>
    </el-card>
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

const goHome = () => {
  router.push('/')
}
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
  background: $bg-page;
}

.register-card {
  width: 400px;

  .title {
    text-align: center;
    margin-bottom: 24px;
  }

  .code-input {
    width: 100%;
    gap: 8px;
  }

  .register-btn {
    width: 100%;
  }

  .footer {
    margin-top: 12px;
  }
}
</style>