import { useAuthStore } from '../stores/auth'
import { computed } from 'vue'
import { House, Setting } from '@element-plus/icons-vue'
import type { Component } from 'vue'

interface NavItem {
  name: string
  path: string
  icon: Component
  requireAdmin?: boolean
}

export function useNavbar() {
  const authStore = useAuthStore()
  
  const isAdmin = computed(() => authStore.auth?.role === 'admin')
  const authName = computed(() => authStore.auth?.name || '')

  const navigation: NavItem[] = [
    { name: '首页', path: '/home', icon: House },
    { name: '用户管理', path: '/users', icon: Setting, requireAdmin: true },
  ]

  function handleLogout() {
    authStore.logout()
  }

  return {
    isAdmin,
    authName,
    navigation,
    handleLogout
  }
} 