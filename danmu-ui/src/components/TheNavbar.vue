<script setup lang="ts">
import { useNavbar } from '../composables/useNavbar'
import { User } from '@element-plus/icons-vue'

const { handleLogout, isAdmin, authName, navigation } = useNavbar()
</script>

<template>
  <div class="fixed top-0 left-0 right-0 bg-indigo-600 h-14 px-4 z-50">
    <div class="max-w-7xl mx-auto h-full">
      <div class="flex justify-between items-center h-full">
        <!-- 左侧导航菜单 -->
        <div class="flex items-center space-x-2">
          <router-link 
            v-for="item in navigation"
            :key="item.path"
            :to="item.path"
            v-show="!item.requireAdmin || isAdmin"
            class="flex items-center px-4 py-2 bg-indigo-700 rounded-md text-white hover:bg-indigo-800 transition-colors"
          >
            <el-icon class="mr-1">
              <component :is="item.icon" ></component>
            </el-icon>
            <span>{{ item.name }}</span>
          </router-link>
        </div>

        <!-- 右侧用户菜单 -->
        <el-dropdown trigger="click">
          <div class="flex items-center px-4 py-2 bg-indigo-700 rounded-md text-white hover:bg-indigo-800 cursor-pointer transition-colors">
            <el-icon class="mr-1"><User /></el-icon>
            <span>{{ authName }}</span>
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>
                <router-link to="/profile" class="text-gray-700">个人资料</router-link>
              </el-dropdown-item>
              <el-dropdown-item v-if="isAdmin">
                <router-link to="/users" class="text-gray-700">用户管理</router-link>
              </el-dropdown-item>
              <el-dropdown-item divided>
                <span class="text-red-600 block w-full" @click="handleLogout">登出</span>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
  </div>
</template>

<style scoped>
:deep(.el-menu) {
  --el-menu-bg-color: transparent;
  --el-menu-text-color: white;
  --el-menu-hover-bg-color: rgba(255, 255, 255, 0.1);
}

:deep(.el-menu-item) {
  --el-menu-active-color: white;
}
</style> 