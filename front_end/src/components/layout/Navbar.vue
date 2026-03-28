<template>
  <el-header class="navbar" :class="{ scrolled: isScrolled }">
    <div class="navbar-container">
      <!-- 品牌 -->
      <div class="navbar-brand" @click="router.push('/')">
        <div class="brand-icon-wrap">
          <el-icon :size="24" color="#fff"><Ship /></el-icon>
        </div>
        <span class="brand-text">智钓蓝海</span>
      </div>

      <!-- 桌面端菜单 -->
      <el-menu
        v-if="!isMobile"
        :default-active="activeMenu"
        mode="horizontal"
        :ellipsis="false"
        class="navbar-menu"
        router
      >
        <el-menu-item index="/">
          <el-icon><HomeFilled /></el-icon> 首页
        </el-menu-item>
        <el-menu-item index="/spots">
          <el-icon><MapLocation /></el-icon> 垂钓水域
        </el-menu-item>
        <el-menu-item index="/reminders">
          <el-icon><Bell /></el-icon> 信息中心
        </el-menu-item>
        <el-menu-item index="/community">
          <el-icon><ChatDotRound /></el-icon> 钓友社区
        </el-menu-item>
        <el-menu-item index="/agent">
          <el-icon><Monitor /></el-icon> 智能助手
        </el-menu-item>
        <el-menu-item index="/dashboard">
          <el-icon><DataLine /></el-icon> 数据大屏
        </el-menu-item>
        <el-menu-item index="/profile" v-if="authStore.isLoggedIn">
          <el-icon><User /></el-icon> 个人中心
        </el-menu-item>
        <el-menu-item index="/admin" v-if="authStore.isStaff">
          <el-icon><Setting /></el-icon> 管理后台
        </el-menu-item>
      </el-menu>

      <!-- 桌面端右侧 -->
      <div v-if="!isMobile" class="navbar-right">
        <!-- 心知天气小组件 -->
        <div id="tp-weather-widget"></div>

        <template v-if="authStore.isLoggedIn">
          <el-dropdown trigger="click">
            <span class="user-info">
              <el-avatar :size="32" :src="myAvatar.src" :style="myAvatar.hasAvatar ? {} : myAvatar.style" class="user-avatar">
                {{ myAvatar.hasAvatar ? '' : myAvatar.letter }}
              </el-avatar>
              <span class="username">{{ authStore.username }}</span>
              <el-icon class="dropdown-arrow"><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="router.push('/profile')">
                  <el-icon><User /></el-icon> 个人中心
                </el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">
                  <el-icon><SwitchButton /></el-icon> 退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
        <template v-else>
          <el-button type="primary" round @click="router.push('/auth')">登录 / 注册</el-button>
        </template>
      </div>

      <!-- 移动端：汉堡按钮 -->
      <div v-if="isMobile" class="mobile-actions">
        <el-icon :size="20" @click="router.push('/reminders')"><Bell /></el-icon>
        <el-icon :size="24" class="hamburger" @click="drawerVisible = true">
          <Menu />
        </el-icon>
      </div>
    </div>

    <!-- 移动端侧边栏 -->
    <el-drawer
      v-model="drawerVisible"
      direction="rtl"
      size="260px"
      :show-close="false"
      class="mobile-drawer"
    >
      <template #header>
        <div class="drawer-header">
          <div class="brand-icon-wrap small">
            <el-icon :size="18" color="#fff"><Ship /></el-icon>
          </div>
          <span class="drawer-brand">智钓蓝海</span>
        </div>
      </template>

      <div class="drawer-menu">
        <div class="drawer-item" :class="{ active: activeMenu === '/' }" @click="navigateMobile('/')">
          <el-icon><HomeFilled /></el-icon> <span>首页</span>
        </div>
        <div class="drawer-item" :class="{ active: activeMenu === '/spots' }" @click="navigateMobile('/spots')">
          <el-icon><MapLocation /></el-icon> <span>垂钓水域</span>
        </div>
        <div class="drawer-item" :class="{ active: activeMenu === '/reminders' }" @click="navigateMobile('/reminders')">
          <el-icon><Bell /></el-icon> <span>信息中心</span>
        </div>
        <div class="drawer-item" :class="{ active: activeMenu === '/community' }" @click="navigateMobile('/community')">
          <el-icon><ChatDotRound /></el-icon> <span>钓友社区</span>
        </div>
        <div class="drawer-item" :class="{ active: activeMenu === '/agent' }" @click="navigateMobile('/agent')">
          <el-icon><Monitor /></el-icon> <span>智能助手</span>
        </div>
        <div class="drawer-item" :class="{ active: activeMenu === '/dashboard' }" @click="navigateMobile('/dashboard')">
          <el-icon><DataLine /></el-icon> <span>数据大屏</span>
        </div>
        <div v-if="authStore.isLoggedIn" class="drawer-item" :class="{ active: activeMenu === '/profile' }" @click="navigateMobile('/profile')">
          <el-icon><User /></el-icon> <span>个人中心</span>
        </div>
        <div v-if="authStore.isStaff" class="drawer-item" :class="{ active: activeMenu === '/admin' }" @click="navigateMobile('/admin')">
          <el-icon><Setting /></el-icon> <span>管理后台</span>
        </div>
      </div>

      <div class="drawer-footer">
        <template v-if="authStore.isLoggedIn">
          <div class="drawer-user">
            <el-avatar :size="36" :src="myAvatar.src" :style="myAvatar.hasAvatar ? {} : myAvatar.style">
              {{ myAvatar.hasAvatar ? '' : myAvatar.letter }}
            </el-avatar>
            <div class="drawer-user-info">
              <span class="drawer-username">{{ authStore.username }}</span>
              <span class="drawer-role">{{ authStore.isAdmin ? '管理员' : authStore.isStaff ? '工作人员' : '用户' }}</span>
            </div>
          </div>
          <div class="drawer-actions">
            <el-button size="small" @click="navigateMobile('/profile')">个人中心</el-button>
            <el-button size="small" type="danger" plain @click="handleLogout(); drawerVisible = false">退出</el-button>
          </div>
        </template>
        <template v-else>
          <el-button type="primary" round style="width:100%" @click="navigateMobile('/auth')">登录 / 注册</el-button>
        </template>
      </div>
    </el-drawer>
  </el-header>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useAvatar } from '@/composables/useAvatar'
import {
  Ship, User, SwitchButton, HomeFilled, MapLocation,
  Bell, Setting, ArrowDown, Menu, ChatDotRound, DataLine, Monitor
} from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

/** 当前用户头像信息（响应式） */
const myAvatar = computed(() => useAvatar(authStore.user?.avatar, authStore.user?.username))

const activeMenu = computed(() => (route.meta.activeMenu as string) || route.path)
const isMobile = ref(false)
const isScrolled = ref(false)
const drawerVisible = ref(false)

function handleLogout() {
  authStore.logout()
  router.push('/')
}

function navigateMobile(path: string) {
  drawerVisible.value = false
  router.push(path)
}

function checkScreen() {
  isMobile.value = window.innerWidth < 768
}

function onAppScroll(e: Event) {
  const scrollTop = (e as CustomEvent).detail as number
  isScrolled.value = scrollTop > 10
}

onMounted(() => {
  checkScreen()
  window.addEventListener('resize', checkScreen)
  document.addEventListener('app-scroll', onAppScroll)

  // 桌面端加载心知天气小组件
  if (!isMobile.value) {
    const w = window as any
    w['SeniverseWeatherWidgetObject'] = 'SeniverseWeatherWidget'
    w['SeniverseWeatherWidget'] = w['SeniverseWeatherWidget'] || function (...args: any[]) {
      (w['SeniverseWeatherWidget'].q = w['SeniverseWeatherWidget'].q || []).push(args)
    }
    w['SeniverseWeatherWidget'].l = +new Date()
    // 直接注入脚本（SPA 中 load 事件早已触发，不能用 addEventListener('load')）
    const script = document.createElement('script')
    script.src = '//cdn.sencdn.com/widget2/static/js/bundle.js?t=' + parseInt((new Date().getTime() / 100000000).toString(), 10)
    script.charset = 'utf-8'
    script.async = true
    document.head.appendChild(script)
    w['SeniverseWeatherWidget']('show', {
      flavor: 'slim',
      location: 'YB1UX38K6DY1',
      geolocation: true,
      language: 'zh-Hans',
      unit: 'c',
      theme: 'auto',
      token: '8237a328-a343-42d7-9de5-075fb91ab4b0',
      hover: 'enabled',
      container: 'tp-weather-widget'
    })
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', checkScreen)
  document.removeEventListener('app-scroll', onAppScroll)
})
</script>

<style scoped>
.navbar {
  height: 60px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  position: sticky;
  top: 0;
  z-index: 1000;
  transition: all 0.3s;
  padding: 0;
}

.navbar.scrolled {
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.1);
  backdrop-filter: blur(8px);
  background: rgba(255, 255, 255, 0.95);
}

.navbar-container {
  display: flex;
  align-items: center;
  max-width: 1400px;
  width: 100%;
  height: 100%;
  margin: 0 auto;
  padding: 0 24px;
}

/* 品牌 */
.navbar-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  margin-right: 32px;
  flex-shrink: 0;
  transition: transform 0.2s;
}

.navbar-brand:hover {
  transform: scale(1.02);
}

.brand-icon-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #409eff, #007bff);
  border-radius: 10px;
  box-shadow: 0 3px 10px rgba(64, 158, 255, 0.3);
  transition: box-shadow 0.3s;
}

.brand-icon-wrap.small {
  width: 30px;
  height: 30px;
  border-radius: 8px;
}

.navbar-brand:hover .brand-icon-wrap {
  box-shadow: 0 4px 14px rgba(64, 158, 255, 0.45);
}

.brand-text {
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(135deg, #2c6fbb, #409eff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  letter-spacing: 1px;
}

/* 菜单 */
.navbar-menu {
  flex: 1;
  border-bottom: none !important;
}

.navbar-menu :deep(.el-menu-item) {
  font-weight: 500;
  transition: all 0.2s;
}

.navbar-menu :deep(.el-menu-item.is-active) {
  font-weight: 600;
}

/* 右侧 */
.navbar-right {
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

#tp-weather-widget {
  display: flex;
  align-items: center;
  margin-right: 12px;
}

/* 心知天气 SDK 动态注入的元素需要穿透 scoped */
#tp-weather-widget :deep(*) {
  font-size: 13px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 4px 8px;
  border-radius: 20px;
  transition: background 0.2s;
}

.user-info:hover {
  background: #f5f7fa;
}

.user-avatar {
  flex-shrink: 0;
}

.username {
  font-size: 14px;
  color: #606266;
  font-weight: 500;
}

.dropdown-arrow {
  color: #c0c4cc;
  font-size: 12px;
  transition: transform 0.2s;
}

.badge-item {
  margin-left: 4px;
}

/* 移动端 */
.mobile-actions {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-left: auto;
}

.hamburger {
  cursor: pointer;
  color: #606266;
  transition: color 0.2s;
}

.hamburger:hover {
  color: #409eff;
}

.mobile-badge :deep(.el-badge__content) {
  top: -4px;
  right: -4px;
}

/* 侧边栏 */
.drawer-header {
  display: flex;
  align-items: center;
  gap: 10px;
}

.drawer-brand {
  font-size: 16px;
  font-weight: 700;
  color: #303133;
}

.drawer-menu {
  padding: 8px 0;
}

.drawer-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 8px;
  margin: 2px 0;
  cursor: pointer;
  font-size: 15px;
  color: #606266;
  transition: all 0.2s;
}

.drawer-item:hover {
  background: #ecf5ff;
  color: #409eff;
}

.drawer-item.active {
  background: #409eff;
  color: #fff;
}

.drawer-badge {
  margin-left: auto;
}

.drawer-footer {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 16px;
  border-top: 1px solid #ebeef5;
  background: #fff;
}

.drawer-user {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 12px;
}

.drawer-user-info {
  display: flex;
  flex-direction: column;
}

.drawer-username {
  font-weight: 600;
  font-size: 14px;
  color: #303133;
}

.drawer-role {
  font-size: 12px;
  color: #909399;
}

.drawer-actions {
  display: flex;
  gap: 8px;
}

@media (max-width: 768px) {
  .navbar-container {
    padding: 0 16px;
  }
}
</style>
