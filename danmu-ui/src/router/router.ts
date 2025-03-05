import { createRouter, createWebHistory } from 'vue-router';
import LoginView from '../views/login.vue'  //
import { useAuthStore } from '../stores/auth';
import HomeView from '../views/home.vue';
import LiveConfDetail from '../views/live-conf/LiveConfDetail.vue'
import Profile from '../views/profile/Profile.vue'
import Users from '../views/users/Users.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',  // 添加根路由
      redirect: '/home'  // 重定向到登录页
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,  // 使用直接导入的组件
      meta: { requiresAuth: false }
    },
    {
      path: '/home',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: true }
    },
    {
      path: '/users',
      name: 'users',
      component: Users,
      meta: { requiresAuth: false, requiresAdmin: true }
    },
    {
      path: '/profile',
      name: 'profile',
      component: Profile,
      meta: { requiresAuth: true }
    },
    {
      path: '/live-conf/:id',
      name: 'live-conf-detail',
      component: LiveConfDetail,
      meta: { requiresAuth: true }
    }
    // ... 其他路由
  ]
});

// 路由守卫
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore();
  const isAdmin = authStore.auth?.role === 'admin';
  
  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    // 需要认证但未登录，重定向到登录页
    next({ name: 'login' });
  } else if (to.meta.requiresAdmin && !isAdmin) {
    // 需要管理员权限但用户不是管理员
    next({ name: 'home' });
  } else if (to.name === 'login' && authStore.isAuthenticated) {
    // 已登录用户访问登录页，重定向到首页
    next({ name: 'home' });
  } else {
    next();
  }
});

export default router;