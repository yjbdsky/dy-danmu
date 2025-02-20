<script setup lang="ts">
import TheNavbar from './components/TheNavbar.vue'
import MobileNavbar from './components/mobile/MobileNavbar.vue'
import { useRoute } from 'vue-router'
import { useDevice } from './composables/useDevice'
import { computed, onMounted } from 'vue'

const route = useRoute()
const { isMobile } = useDevice()
const showNavbar = computed(() => route.path !== '/login')
onMounted(() => {
  console.log('BASE_URL:', import.meta.env.VITE_BASE_URL)
  console.log('API_URL:', import.meta.env.VITE_API_URL)
  console.log('所有环境变量：', import.meta.env)
})
</script>

<template>
  <MobileNavbar v-if="showNavbar && isMobile" />
  <TheNavbar v-else-if="showNavbar" />
  <div :class="showNavbar ? 'pt-14 sm:pt-16' : ''">
    <router-view />
  </div>
</template>


