<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
    <div class="max-w-sm w-full bg-white p-6 rounded-lg shadow-sm border border-gray-100">
      <!-- Logo and Title -->
      <div class="text-center">
        <h2 class="mt-4 text-xl font-medium text-gray-900">
          danmu-nu
        </h2>
      </div>

      <!-- Login Form -->
      <el-form @submit.prevent="handleLogin" class="mt-6 space-y-4">
        <el-form-item>
          <el-input
            v-model="loginForm.email"
            type="text"
            class="rounded-md !h-10"
            placeholder="请输入邮箱"
          />
        </el-form-item>

        <el-form-item>
          <el-input
            v-model="loginForm.password"
            type="password"
            class="!h-10"
            placeholder="请输入密码"
          />
        </el-form-item>
        <el-button
          type="primary"
          :loading="isLoading"
          class="w-full !h-10"
          @click="handleLogin"
        >
          登录
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import { useAuthStore } from '../stores/auth';
import type { LoginRequest } from '../types/models/auth';

const router = useRouter();
const authStore = useAuthStore();
const loginForm = ref<LoginRequest>({
  email: '',
  password: ''
});
const isLoading = ref(false);

async function handleLogin() {
  if (!loginForm.value.email || !loginForm.value.password) {
    ElMessage.warning('请输入邮箱和密码');
    return;
  }

  try {
    isLoading.value = true;
    await authStore.login(loginForm.value);
    ElMessage.success('登录成功');
    await router.push('/home');
  } catch (error: any) {
    console.error('登录失败:', error);
    ElMessage.error(error?.response?.data?.msg || '登录失败，请稍后重试');
  } finally {
    isLoading.value = false;
  }
}
</script>

