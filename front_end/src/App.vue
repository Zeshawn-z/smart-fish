<template>
  <el-scrollbar class="app-scrollbar" @scroll="onScroll">
    <div id="app-wrapper">
      <Navbar />
      <main class="main-content" :class="{ 'admin-page': route.meta.hideFooter }">
        <router-view />
      </main>
      <Footer v-if="!route.meta.hideFooter" />
    </div>
  </el-scrollbar>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import Navbar from '@/components/layout/Navbar.vue'
import Footer from '@/components/layout/Footer.vue'

const route = useRoute()

/** 将 el-scrollbar 的滚动位置广播给 Navbar 等组件 */
function onScroll({ scrollTop }: { scrollTop: number }) {
  document.dispatchEvent(new CustomEvent('app-scroll', { detail: scrollTop }))
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

/* 禁用浏览器原生滚动条，由 el-scrollbar 接管 */
html, body, #app {
  height: 100%;
  overflow: hidden;
}

.app-scrollbar {
  height: 100vh;
}

#app-wrapper {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #f5f7fa;
}

.main-content {
  flex: 1;
  padding: 20px 24px;
  max-width: 1400px;
  width: 100%;
  margin: 0 auto;
}

/* 后台管理页面：全宽、无 padding */
.main-content.admin-page {
  padding: 0;
  max-width: none;
  margin: 0;
}
</style>
