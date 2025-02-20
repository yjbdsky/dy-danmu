<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useWindowSize } from '@vueuse/core'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getAuth, updateUser } from '../../api/auth'
import type { Auth } from '../../types/models/auth'

const { width } = useWindowSize()
const isMobile = computed(() => width.value < 640)
const router = useRouter()
const loading = ref(false)
const userInfo = ref<Auth | null>(null)

const form = reactive({
  name: '',
  email: '',
  password: '',
  confirmPassword: ''
})

// 表单验证规则
const rules = {
  name: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { min: 6, message: '密码长度至少6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { 
      validator: (_: any, value: string, callback: Function) => {
        if (form.password && value !== form.password) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 获取用户信息
async function fetchUserInfo() {
  loading.value = true
  try {
    const res = await getAuth()
    userInfo.value = res.data.data
    // 填充表单
    form.name = userInfo.value.name
    form.email = userInfo.value.email
  } catch (error) {
    ElMessage.error('获取用户信息失败')
  } finally {
    loading.value = false
  }
}

// 保存修改
async function handleSave() {
  try {
    await ElMessageBox.confirm('确认保存修改？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    loading.value = true
    await updateUser({
      name: form.name,
      email: form.email,
      password: form.password || undefined // 如果密码为空则不修改
    })
    
    ElMessage.success('保存成功')
    router.back()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('保存失败')
    }
  } finally {
    loading.value = false
  }
}

// 取消修改
function handleCancel() {
  router.back()
}

fetchUserInfo()
</script>

<template>
  <div class="min-h-screen bg-gray-50 py-4 sm:py-8 px-4 sm:px-6">
    <div class="max-w-2xl mx-auto bg-white rounded-lg shadow-sm p-4 sm:p-6">
      <h2 class="text-xl font-medium mb-6">个人资料</h2>
      
      <el-form 
        v-if="userInfo"
        :model="form"
        :rules="rules"
        :label-width="isMobile ? '70px' : '100px'"
        :label-position="isMobile ? 'top' : 'right'"
        class="space-y-4"
      >
        <el-form-item label="用户名" prop="name">
          <el-input 
            v-model="form.name"
            class="!max-w-md"
          />
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input 
            v-model="form.email"
            class="!max-w-md"
          />
        </el-form-item>

        <el-form-item label="角色">
          <el-tag>{{ userInfo.role }}</el-tag>
        </el-form-item>

        <el-form-item label="新密码" prop="password">
          <el-input 
            v-model="form.password"
            type="password"
            class="!max-w-md"
            placeholder="不修改请留空"
          />
        </el-form-item>

        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input 
            v-model="form.confirmPassword"
            type="password"
            class="!max-w-md"
            placeholder="不修改请留空"
          />
        </el-form-item>

        <el-form-item>
          <div class="flex gap-4">
            <el-button 
              type="primary"
              :loading="loading"
              @click="handleSave"
            >
              保存
            </el-button>
            <el-button @click="handleCancel">取消</el-button>
          </div>
        </el-form-item>
      </el-form>

      <div v-else class="h-40 flex items-center justify-center">
        <el-skeleton :rows="4" />
      </div>
    </div>
  </div>
</template> 