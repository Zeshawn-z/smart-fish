import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import { nextTick } from 'vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: () => import('@/views/HomePage.vue'),
      meta: { title: '首页' }
    },
    {
      path: '/spots',
      name: 'FishingSpots',
      component: () => import('@/views/FishingSpotsView.vue'),
      meta: { title: '垂钓水域' }
    },
    {
      path: '/reminders',
      name: 'Reminders',
      component: () => import('@/views/ReminderView.vue'),
      meta: { title: '信息中心' }
    },
    {
      path: '/community',
      name: 'Community',
      component: () => import('@/views/CommunityView.vue'),
      meta: { title: '钓友社区' }
    },
    {
      path: '/community/:id',
      name: 'PostDetail',
      component: () => import('@/views/PostDetailView.vue'),
      meta: { title: '帖子详情', activeMenu: '/community' }
    },
    {
      path: '/auth',
      name: 'Auth',
      component: () => import('@/views/AuthView.vue'),
      meta: { title: '登录' }
    },
    {
      path: '/profile',
      name: 'Profile',
      component: () => import('@/views/UserView.vue'),
      meta: { title: '个人中心', requiresAuth: true }
    },
    {
      path: '/admin',
      name: 'Admin',
      component: () => import('@/views/AdminView.vue'),
      meta: { title: '管理后台', requiresAuth: true, requiresStaff: true, hideFooter: true }
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('@/views/DashboardView.vue'),
      meta: { title: '数据大屏', hideFooter: true, hideNavbar: true, fullscreen: true }
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'NotFound',
      component: () => import('@/views/NotFound.vue'),
      meta: { title: '404' }
    }
  ]
})

// 路由守卫
router.beforeEach((to, _from, next) => {
  // 更新页面标题
  document.title = `${to.meta.title || '智钓蓝海'} - 智钓蓝海信息平台`

  const authStore = useAuthStore()

  // 需要认证的页面
  if (to.meta.requiresAuth && !authStore.isLoggedIn) {
    ElMessage.warning('请先登录后再访问')
    next({ name: 'Auth', query: { redirect: to.fullPath } })
    return
  }

  // 需要 staff 权限的页面
  if (to.meta.requiresStaff && !authStore.isStaff) {
    next({ name: 'Home' })
    return
  }

  // 已登录用户访问登录页，重定向到首页
  if (to.name === 'Auth' && authStore.isLoggedIn) {
    next({ name: 'Home' })
    return
  }

  next()
})

// 路由切换后自动滚动到顶部（el-scrollbar 接管了滚动，需要手动操作 DOM）
router.afterEach(() => {
  nextTick(() => {
    const scrollbar = document.querySelector('.app-scrollbar .el-scrollbar__wrap') as HTMLElement | null
    if (scrollbar) {
      scrollbar.scrollTo({ top: 0, behavior: 'smooth' })
    }
  })
})

export default router
