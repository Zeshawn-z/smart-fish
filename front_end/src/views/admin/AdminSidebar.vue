<template>
  <el-aside width="200px" class="admin-aside">
    <el-menu :default-active="active" @select="(key: string) => $emit('update:active', key)">
      <el-menu-item index="regions">
        <el-icon><MapLocation /></el-icon>
        <span>区域管理</span>
      </el-menu-item>
      <el-menu-item index="spots">
        <el-icon><Location /></el-icon>
        <span>水域管理</span>
      </el-menu-item>
      <el-menu-item index="devices">
        <el-icon><Monitor /></el-icon>
        <span>设备管理</span>
      </el-menu-item>
      <el-menu-item index="gateways">
        <el-icon><Connection /></el-icon>
        <span>网关管理</span>
      </el-menu-item>
      <el-menu-item index="reminders">
        <el-icon><Bell /></el-icon>
        <span>提醒管理</span>
      </el-menu-item>
      <el-menu-item index="notices">
        <el-icon><Notification /></el-icon>
        <span>通知管理</span>
      </el-menu-item>
      <el-menu-item index="users" v-if="showUsers">
        <el-icon><User /></el-icon>
        <span>用户管理</span>
      </el-menu-item>
    </el-menu>
  </el-aside>
</template>

<script setup lang="ts">
import { MapLocation, Location, Monitor, Connection, Bell, Notification, User } from '@element-plus/icons-vue'

defineProps<{
  active: string
  showUsers: boolean
}>()

defineEmits<{
  (e: 'update:active', value: string): void
}>()
</script>

<style scoped>
.admin-aside {
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
  margin-right: 16px;
}
</style>
