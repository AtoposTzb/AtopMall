<template>
  <el-dialog
    v-model="modalVisible"
    :title="mode === 'login' ? '用户登录' : '用户注册'"
    width="420px"
    :close-on-click-modal="false"
    :close-on-press-escape="true"
    destroy-on-close
    @close="handleClose"
  >
    <!-- 登录表单 -->
    <template v-if="mode === 'login'">
      <el-form ref="loginFormRef" :model="loginForm" :rules="loginRules" label-width="0" @keyup.enter="handleLogin">
        <el-form-item prop="mobile">
          <el-input v-model="loginForm.mobile" placeholder="请输入手机号" size="large" />
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            size="large"
            show-password
          />
        </el-form-item>
        <el-form-item prop="captcha">
          <div class="captcha-row">
            <el-input v-model="loginForm.captcha" placeholder="请输入验证码" size="large" />
            <img
              v-if="captchaImage"
              :src="captchaImage"
              class="captcha-img"
              @click="loadCaptcha"
              title="点击刷新验证码"
            />
          </div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" class="submit-btn" :loading="loginLoading" @click="handleLogin">
            登 录
          </el-button>
        </el-form-item>
      </el-form>
      <div class="form-footer">
        <span>还没有账号？</span>
        <el-button type="primary" link @click="switchMode('register')">立即注册</el-button>
      </div>
    </template>

    <!-- 注册表单 -->
    <template v-else>
      <el-form ref="registerFormRef" :model="registerForm" :rules="registerRules" label-width="0">
        <el-form-item prop="mobile">
          <el-input v-model="registerForm.mobile" placeholder="请输入手机号" size="large" />
        </el-form-item>
        <el-form-item prop="email">
          <el-input v-model="registerForm.email" placeholder="请输入邮箱" size="large" />
        </el-form-item>
        <el-form-item prop="code">
          <div class="captcha-row">
            <el-input v-model="registerForm.code" placeholder="验证码" size="large" />
            <el-button type="primary" size="large" :disabled="countdown > 0" @click="handleSendCode" class="code-btn">
              {{ countdown > 0 ? `${countdown}s` : '发送验证码' }}
            </el-button>
          </div>
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="registerForm.password" type="password" placeholder="请输入密码" size="large" show-password />
        </el-form-item>
        <el-form-item prop="confirmPassword">
          <el-input v-model="registerForm.confirmPassword" type="password" placeholder="请确认密码" size="large" show-password />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" size="large" class="submit-btn" :loading="registerLoading" @click="handleRegister">
            注 册
          </el-button>
        </el-form-item>
      </el-form>
      <div class="form-footer">
        <span>已有账号？</span>
        <el-button type="primary" link @click="switchMode('login')">去登录</el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, type FormInstance } from 'element-plus'
import { login, register, getCaptcha, sendEmailCode } from '@/api/user'
import { useUserStore } from '@/store/user'
import { useAuthModal, type AuthMode } from '@/composables/useAuthModal'

const userStore = useUserStore()
const { visible, mode, close, switchMode: switchAuthMode } = useAuthModal()

const modalVisible = computed({
  get: () => visible.value,
  set: (val: boolean) => { if (!val) close() }
})

const switchMode = (m: AuthMode) => {
  switchAuthMode(m)
}

const handleClose = () => {
  close()
}

// ========== 登录表单 ==========
const loginFormRef = ref<FormInstance>()
const loginLoading = ref(false)
const captchaImage = ref('')
const captchaId = ref('')

const loginForm = reactive({
  mobile: '',
  password: '',
  captcha: '',
  captcha_id: ''
})

const loginRules = {
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
  const valid = await loginFormRef.value?.validate().catch(() => false)
  if (!valid) return

  loginLoading.value = true
  try {
    loginForm.captcha_id = captchaId.value
    const res = await login(loginForm)
    userStore.setToken(res.token)
    await userStore.initUser()
    ElMessage.success('登录成功')
    close()
  } catch (error) {
    loadCaptcha()
  } finally {
    loginLoading.value = false
  }
}

// ========== 注册表单 ==========
const registerFormRef = ref<FormInstance>()
const registerLoading = ref(false)
const countdown = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

const registerForm = reactive({
  mobile: '',
  email: '',
  code: '',
  password: '',
  confirmPassword: ''
})

const validateConfirmPassword = (_rule: any, value: string, callback: any) => {
  if (value !== registerForm.password) {
    callback(new Error('两次密码输入不一致'))
  } else {
    callback()
  }
}

const registerRules = {
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
  if (!registerForm.email) {
    ElMessage.warning('请先输入邮箱')
    return
  }
  try {
    await sendEmailCode({ email: registerForm.email, type: 1 })
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
  const valid = await registerFormRef.value?.validate().catch(() => false)
  if (!valid) return

  registerLoading.value = true
  try {
    const res = await register({
      mobile: registerForm.mobile,
      password: registerForm.password,
      email: registerForm.email,
      code: registerForm.code
    })
    userStore.setToken(res.token)
    ElMessage.success('注册成功')
    close()
  } catch (error) {
    // 错误已在拦截器中处理
  } finally {
    registerLoading.value = false
  }
}

onMounted(() => {
  loadCaptcha()
})
</script>

<style lang="scss" scoped>
.captcha-row {
  display: flex;
  gap: 8px;
  width: 100%;

  .captcha-img {
    height: 40px;
    cursor: pointer;
    border-radius: 4px;
    flex-shrink: 0;
  }

  .code-btn {
    flex-shrink: 0;
    min-width: 110px;
  }
}

.submit-btn {
  width: 100%;
}

.form-footer {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #909399;
}
</style>