import { ref, onMounted, onUnmounted } from 'vue'

export function useDevice() {
  const isMobile = ref(false)
  let cleanup: (() => void) | null = null
  
  function checkDevice() {
    isMobile.value = window.innerWidth <= 768
  }
  
  onMounted(() => {
    checkDevice()
    window.addEventListener('resize', checkDevice)
    cleanup = () => window.removeEventListener('resize', checkDevice)
  })
  
  if (cleanup) {
    onUnmounted(cleanup)
  }
  
  return {
    isMobile
  }
} 