<script setup lang="ts">
import { ref, onMounted } from 'vue'
import dayjs from 'dayjs'

const props = defineProps<{
  modelValue: [number, number]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: [number, number]): void
}>()

const startTime = ref(dayjs().startOf('day').valueOf())
const endTime = ref(dayjs().endOf('day').valueOf())

const activePreset = ref('today')

const presets = [
  { label: '今日', value: 'today' },
  { label: '昨日', value: 'yesterday' },
  { label: '本周', value: 'week' },
  { label: '本月', value: 'month' }
]

function updateTimeRange(preset: string) {
  activePreset.value = preset
  const now = dayjs()
  
  switch (preset) {
    case 'today':
      startTime.value = now.startOf('day').valueOf()
      endTime.value = now.endOf('day').valueOf()
      break
    case 'yesterday':
      startTime.value = now.subtract(1, 'day').startOf('day').valueOf()
      endTime.value = now.subtract(1, 'day').endOf('day').valueOf()
      break
    case 'week':
      startTime.value = now.startOf('week').valueOf()
      endTime.value = now.endOf('week').valueOf()
      break
    case 'month':
      startTime.value = now.startOf('month').valueOf()
      endTime.value = now.endOf('month').valueOf()
      break
  }
  
  emit('update:modelValue', [startTime.value, endTime.value])
}

function handleTimeChange() {
  activePreset.value = ''
  emit('update:modelValue', [startTime.value, endTime.value])
}

onMounted(() => {
  updateTimeRange('today')  // 默认选择今日
})
</script>

<template>
  <div class="space-y-2">
    <!-- 快捷按钮 -->
    <div class="flex gap-2">
      <el-button
        v-for="preset in presets"
        :key="preset.value"
        :type="activePreset === preset.value ? 'primary' : 'default'"
        @click="updateTimeRange(preset.value)"
      >
        {{ preset.label }}
      </el-button>
    </div>
    
    <!-- 时间选择器 -->
    <div class="flex items-center gap-2">
      <el-tag 
        type="success" 
        effect="plain"
        size="large"
        class="!w-[70px] !text-center"
      >
        开始
      </el-tag>
      <el-date-picker
        v-model="startTime"
        type="datetime"
        placeholder="请选择开始时间"
        value-format="x"
        @change="handleTimeChange"
      />
    </div>
    
    <div class="flex items-center gap-2">
      <el-tag 
        type="warning" 
        effect="plain"
        size="large"
        class="!w-[70px] !text-center"
      >
        结束
      </el-tag>
      <el-date-picker
        v-model="endTime"
        type="datetime"
        placeholder="请选择结束时间"
        value-format="x"
        @change="handleTimeChange"
      />
    </div>
  </div>
</template>

<style scoped>
:deep(.el-range-editor) {
  --el-border-color: #dcdfe6;
  --el-border-color-hover: #c0c4cc;
  background-color: white;
}

:deep(.el-button) {
  --el-button-hover-bg-color: #f5f7fa;
  --el-button-active-bg-color: #e4e7ed;
}

:deep(.el-button.el-button--primary) {
  --el-button-hover-bg-color: #409eff;
  --el-button-active-bg-color: #337ecc;
}

:deep(.el-range-editor.el-input__wrapper) {
  width: 100%;
}

@media (min-width: 768px) {
  :deep(.el-range-editor.el-input__wrapper) {
    width: 420px;
  }
}

:deep(.el-picker__popper.el-popper) {
  max-width: 90vw;
}
</style> 