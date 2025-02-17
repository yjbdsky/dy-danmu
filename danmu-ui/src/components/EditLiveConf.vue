<script setup lang="ts">
import { ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { updateLiveConf } from '../api/live-conf'
import type { UpdateLiveConfRequest, LiveConf } from '../types/models/live_conf'

const props = defineProps<{
  modelValue: boolean
  liveConf: LiveConf | null
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'success'): void
}>()

const form = ref<UpdateLiveConfRequest>({
  id: '',
  room_display_id: '',
  url: '',
  name: '',
  enable: true
})

// 监听 liveConf 变化，更新表单数据
watch(() => props.liveConf, (newVal) => {
  if (newVal) {
    form.value = {
      id: newVal.id,
      room_display_id: newVal.room_display_id,
      url: newVal.url,
      name: newVal.name,
      enable: newVal.enable
    }
  }
}, { immediate: true })

const isLoading = ref(false)

async function handleSubmit() {
  try {
    // 检查表单数据是否有变化
    if (props.liveConf && 
        props.liveConf.room_display_id === form.value.room_display_id &&
        props.liveConf.url === form.value.url &&
        props.liveConf.name === form.value.name &&
        props.liveConf.enable === form.value.enable) {
      ElMessage.info('数据未发生变化')
      emit('update:modelValue', false)
      return
    }

    isLoading.value = true
    await updateLiveConf(form.value)
    ElMessage.success('更新成功')
    emit('success')
    emit('update:modelValue', false)
  } catch (error) {
    console.error('[LiveConf] Update error:', error)
  } finally {
    isLoading.value = false
  }
}
</script>

<template>
  <el-dialog
    title="编辑直播配置"
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