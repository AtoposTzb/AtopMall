import { ref } from 'vue'

export type AuthMode = 'login' | 'register'

const visible = ref(false)
const mode = ref<AuthMode>('login')

export function useAuthModal() {
  const open = (m?: AuthMode) => {
    mode.value = m || 'login'
    visible.value = true
  }

  const close = () => {
    visible.value = false
  }

  const switchMode = (m: AuthMode) => {
    mode.value = m
  }

  return { visible, mode, open, close, switchMode }
}