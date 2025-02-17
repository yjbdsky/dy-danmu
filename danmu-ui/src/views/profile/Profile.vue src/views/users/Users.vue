<template>
  <!-- No changes to template section -->
</template>

<script setup>
import { useWindowSize } from '@vueuse/core'
import { computed } from 'vue'

const { width } = useWindowSize()
const isMobile = computed(() => width.value < 640)
</script> 