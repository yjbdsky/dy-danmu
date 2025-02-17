<script setup lang="ts">
import { useNavbar } from '../../composables/useNavbar'
import { User } from '@element-plus/icons-vue'

const { handleLogout, isAdmin, authName, navigation } = useNavbar()
</script>

<template>
  <div class="fixed top-0 left-0 right-0 bg-indigo-600 h-12 px-3 z-50">
    <div class="h-full">
      <div class="flex justify-between items-center h-full">
        <!-- 左侧导航菜单 -->
        <div class="flex items-center space-x-1">
          <router-link 
            v-for="item in navigation"
            :key="item.path"
            :to="item.path"
            v-show="!item.requireAdmin || isAdmin"
            class="flex items-center px-3 py-1.5 bg-indigo-700/50 rounded-md text-white text-sm active:bg-indigo-800"
          >
            <el-icon>
             <component :is="item.icon" />
            </el-icon>
            <span>{{ item.name }}</span>
          </router-link>
        </div>

        <!-- 右侧用户菜单 -->
        <el-dropdown trigger="click" placement="bottom-end">
          <div class="flex items-center px-3 py-1.5 bg-indigo-700/50 rounded-md text-white text-sm active:bg-indigo-800">
            <el-icon class="mr-1"><User /></el-icon>
            <span class="max-w-[100px] truncate">{{ authName }}</span>
          </div>

          <template #dropdown>
            <el-dropdown-menu class="!min-w-[120px]">
              <el-dropdown-item>
                <router-link to="/profile" class="flex items-center justify-center text-gray-700">
                  个人资料
                </router-link>
              </el-dropdown-item>
              <el-dropdown-item v-if="isAdmin">
                <router-link to="/users" class="flex items-center justify-center text-gray-700">
                  用户管理
                </router-link>
              </el-dropdown-item>
              <el-dropdown-item divided>
                <div 
                  class="flex items-center justify-center text-red-600 w-full" 
                  @click="handleLogout"
                >
                  登出
                </div>
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
  </div>
</template>

<style scoped>
:deep(.el-dropdown-menu) {
  --el-dropdown-menuItem-hover-fill: rgb(249 250 251);
  --el-dropdown-menuItem-hover-color: currentColor;
}

:deep(.el-dropdown-menu__item) {
  padding: 8px 12px;
  justify-content: center;
}
</style>
