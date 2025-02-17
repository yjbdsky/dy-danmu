<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useWindowSize } from '@vueuse/core'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { listUsers, deleteUser, resetPassword, register } from '../../api/auth'
import type { Auth, RegisterRequest } from '../../types/models/auth'
import { Role } from '../../types/models/auth'

const { width } = useWindowSize()
const isMobile = computed(() => width.value < 640)
const loading = ref(false)
const userList = ref<Auth[]>([])
const dialogVisible = ref(false)

const form = ref<RegisterRequest>({
  name: '',
  email: '',
  password: '',
  role: Role.User
})

// 表单验证规则
const rules = {
  name: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ],
  password: [
    { 
      validator: (rule: any, value: string, callback: Function) => {
        if (value && (value.length < 6 || value.length > 50)) {
          callback(new Error('密码长度需在6-50位之间'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ],
  role: [{ required: true, message: '请选择角色', trigger: 'change' }]
}

async function fetchUsers() {
  loading.value = true
  try {
    const res = await listUsers()
    userList.value = res.data.data
  } catch (error) {
    ElMessage.error('获取用户列表失败')
  } finally {
    loading.value = false
  }
}

async function handleDelete(id: string) {
  try {
    await ElMessageBox.confirm('确认删除该用户？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await deleteUser(id)
    ElMessage.success('删除成功')
    fetchUsers()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

async function handleResetPassword(id: string) {
  try {
    await ElMessageBox.confirm('确认重置该用户密码？', '警告', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await resetPassword(id)
    ElMessage.success('密码重置成功')
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('重置密码失败')
    }
  }
}

async function handleCreate() {
  try {
    await register(form.value)
    ElMessage.success('创建成功')
    dialogVisible.value = false
    fetchUsers()
    // 重置表单
    form.value = {
      name: '',
      email: '',
      password: '',
      role: Role.User
    }
  } catch (error) {
    ElMessage.error('创建失败')
  }
}

onMounted(fetchUsers)
</script>

<template>
  <div class="p-2 sm:p-4">
    <div class="max-w-5xl mx-auto">
      <!-- 用户列表 -->
      <el-table 
        v-loading="loading"
        :data="userList" 
        class="w-full"
        :size="isMobile ? 'small' : 'default'"
      >
        <el-table-column label="用户名" prop="name" min-width="100" />
        <el-table-column 
          label="邮箱" 
          prop="email" 
          min-width="140"
          :show-overflow-tooltip="isMobile" 
        />
        <el-table-column label="角色" prop="role" width="100">
          <template #default="{ row }">
            <el-tag :type="row.role === 'admin' ? 'danger' : 'info'">
              {{ row.role }}
            </el-tag>
          </template>
        </el-table-column>
        
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <div class="flex gap-2">
              <el-button 
                type="danger" 
                link
                v-if="row.role !== 'admin'"
                @click="handleDelete(row.id)"
              >
                删除
              </el-button>
              <el-button 
                type="warning" 
                link
                v-if="row.role !== 'admin'"
                @click="handleResetPassword(row.id)"
              >
                重置密码
              </el-button>
              <span v-if="row.role === 'admin'" class="text-gray-400 text-sm">
                无可用操作
              </span>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <!-- 创建用户按钮 -->
      <div class="fixed right-8 bottom-8">
        <el-button
          type="primary"
          circle
          size="large"
          @click="dialogVisible = true"
        >
          <el-icon><Plus /></el-icon>
        </el-button>
      </div>
    </div>

    <!-- 创建用户对话框 -->
    <el-dialog
      v-model="dialogVisible"
      title="创建用户"
      :width="isMobile ? '90%' : '500px'"
    >
      <el-form 
        ref="formRef"
        :model="form"
        :rules="rules"
        :label-width="isMobile ? '70px' : '80px'"
        :label-position="isMobile ? 'top' : 'right'"
      >
        <el-form-item label="用户名" prop="name">
          <el-input v-model="form.name" />
        </el-form-item>
        
        <el-form-item label="邮箱" prop="email">
          <el-input v-model="form.email" />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input 
            v-model="form.password" 
            type="password"
            placeholder="不填则使用默认密码"
          />
        </el-form-item>
        
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role">
            <el-option label="管理员" :value="Role.Admin" />
            <el-option label="普通用户" :value="Role.User" />
          </el-select>
        </el-form-item>
      </el-form>
      
      <template #footer>
        <div class="flex justify-end gap-2">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleCreate">确认</el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template> 