<!-- views/Home.vue -->
<script setup lang="ts">
import { ref, watch } from 'vue'
import UserMultiSelect from './UserMultiSelect.vue'
import TimeRangePicker from './TimeRangePicker.vue'
import ToUserSelect from './ToUserSelect.vue'

const selectedUserIds = ref<number[]>([])
const timeRange = ref<[number, number]>([Date.now(), Date.now()])  // 添加时间范围状态
const selectedToUserIds = ref<number[]>([])

// 监听时间范围变化
watch(timeRange, (newVal) => {
  console.log('TimeRange changed:', {
    start: new Date(newVal[0]).toLocaleString(),
    end: new Date(newVal[1]).toLocaleString(),
    raw: newVal
  })
})

// 监听用户选择变化
watch(selectedUserIds, (newVal) => {
  console.log('Selected users changed:', newVal)
})

// 监听接收用户选择变化
watch(selectedToUserIds, (newVal) => {
  console.log('Selected to-users changed:', newVal)
})
</script>

<template>
  <div class="space-y-8">  <!-- 增加组件之间的间距 -->
    <TimeRangePicker v-model="timeRange" />
    <UserMultiSelect v-model="selectedUserIds" />
    <ToUserSelect 
      v-model="selectedToUserIds"
      room-display-id="949735304546"
    />
  </div>
</template>