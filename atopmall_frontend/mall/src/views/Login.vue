<template>
  <div class="login-page flex-center">
    <el-card class="login-card">
      <h2 class="title">用户登录</h2>
      <el-form ref="formRef" :model="form" :rules="rules" label-width="0">
        <el-form-item prop="mobile">
          <el-input v-model="form.mobile" placeholder="请输入手机号" prefix-icon="Phone" size="large" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="Lock"
            size="large"
            show-password
          />
        </el-form-item>
        <el-form-item prop="captcha">
          <div class="captcha-input flex">
            <el-input v-model="form.captcha" placeholder="请输入验证码" size="large" />
            <img
              v-if="captchaImage"
              :src="captchaImage"
              class="captcha-img"
              @click="loadCaptcha"
              title="点击刷新验证码"
            />
            <el-button v-else type="primary" size="large" @click="loadCaptcha">获取验证码</el-button>
          </div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" class="login-btn" :loading="loading" @click="handleLogin">
            登 录
          </el-button>
        </el-form-item>
      </el-form>
      <div class="footer flex-between">
        <el-button type="primary" link @click="goHome">返回首页</el-button>
        <router-link to="/register">注册新账号</router-link>
      </div>
    </el-card>
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

const goHome = () => {
  router.push('/')
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
  background: $bg-page;
}

.login-card {
  width: 400px;

  .title {
    text-align: center;
    margin-bottom: 24px;
    color: $text-primary;
  }

  .captcha-input {
    width: 100%;
    gap: 8px;

    .captcha-img {
      height: 40px;
      cursor: pointer;
      border-radius: 4px;
    }
  }

  .login-btn {
    width: 100%;
  }

  .footer {
    margin-top: 12px;
  }
}
</style>