import { ref, onUnmounted } from 'vue'

export function useDraggable() {
  const offsetY = ref(0)
  let dragging = false
  let startY = 0
  let startOffsetY = 0

  const onMouseDown = (e: MouseEvent) => {
    dragging = true
    startY = e.clientY
    startOffsetY = offsetY.value
    document.addEventListener('mousemove', onMouseMove)
    document.addEventListener('mouseup', onMouseUp)
    document.body.style.userSelect = 'none'
    e.preventDefault()
  }

  const onMouseMove = (e: MouseEvent) => {
    if (!dragging) return
    const delta = e.clientY - startY
    offsetY.value = startOffsetY + delta
  }

  const onMouseUp = () => {
    dragging = false
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
    document.body.style.userSelect = ''
  }

  onUnmounted(() => {
    document.removeEventListener('mousemove', onMouseMove)
    document.removeEventListener('mouseup', onMouseUp)
    document.body.style.userSelect = ''
  })

  return { offsetY, onMouseDown }
}