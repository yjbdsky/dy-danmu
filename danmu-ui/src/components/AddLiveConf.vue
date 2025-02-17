<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { createLiveConf } from '../api/live-conf'
import type { CreateLiveConfRequest } from '../types/models/live_conf'

const props = defineProps<{
  modelValue: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

watch(() => props.modelValue, (newVal) => {
  if (!newVal) {
    form.value = {
      room_display_id: '',
      url: '',
      name: '',
      enable: true
    }
  }
})

const form = ref<CreateLiveConfRequest>({
  room_display_id: '',
  url: '',
  name: '',
  enable: true
})

const isLoading = ref(false)

async function handleSubmit() {
  try {
    isLoading.value = true
    await createLiveConf(form.value)
    ElMessage.success('添加成功')
    emit('success')
    emit('update:modelValue', false)
  } catch (error) {
    console.error('[LiveConf] Create error:', error)
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <el-dialog
    title="添加直播配置"
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    width="90%"
    :max-width="500"
  >
    <el-form :model="form" label-width="100px">
      <el-form-item label="房间ID">
        <el-input v-model="form.room_display_id" placeholder="请输入房间ID" />
      </el-form-item>
      
      <el-form-item label="名称">
        <el-input v-model="form.name" placeholder="请输入名称" />
      </el-form-item>
      
      <el-form-item label="URL">
        <el-input v-model="form.url" placeholder="请输入URL" />
      </el-form-item>
      
      <el-form-item label="状态">
        <el-switch v-model="form.enable" />
      </el-form-item>
    </el-form>

    <template #footer>
      <el-button @click="$emit('update:modelValue', false)">取消</el-button>
      <el-button 
        type="primary" 
        :loading="isLoading"
        @click="handleSubmit"
      >
        保存
      </el-button>
    </template>
  </el-dialog>
</template> 