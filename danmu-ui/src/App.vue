<script setup lang="ts">
import TheNavbar from './components/TheNavbar.vue'
import MobileNavbar from './components/mobile/MobileNavbar.vue'
import { useRoute } from 'vue-router'
import { useDevice } from './composables/useDevice'
import { computed } from 'vue'

const route = useRoute()
const { isMobile } = useDevice()
const showNavbar = computed(() => route.path !== '/login')
console.log('Current route:', route.path)  // 添加控制台日志
</script>

<template>
  <MobileNavbar v-if="showNavbar && isMobile" />
  <TheNavbar v-else-if="showNavbar" />
  <div :class="showNavbar ? 'pt-14 sm:pt-16' : ''">
    <router-view />
  </div>
</template>


