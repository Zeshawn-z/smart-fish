<template>
  <div class="admin-sidebar">
    <!-- 品牌区域 -->
    <div class="sidebar-brand">
      <span class="brand-text">管理后台</span>
    </div>

    <!-- 导航菜单 -->
    <el-scrollbar class="sidebar-menu-scroll">
      <el-menu
        :default-active="active"
        :default-openeds="['data', 'ops', 'community']"
        @select="(key: string) => $emit('update:active', key)"
      >
        <el-sub-menu index="data">
          <template #title>
            <el-icon><Folder /></el-icon>
            <span>数据管理</span>
          </template>
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
        </el-sub-menu>

        <el-sub-menu index="ops">
          <template #title>
            <el-icon><Operation /></el-icon>
            <span>运营管理</span>
          </template>
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
        </el-sub-menu>

        <el-sub-menu index="community">
          <template #title>
            <el-icon><ChatDotRound /></el-icon>
            <span>社区管理</span>
          </template>
          <el-menu-item index="posts">
            <el-icon><Document /></el-icon>
            <span>帖子管理</span>
          </el-menu-item>
          <el-menu-item index="comments">
            <el-icon><ChatLineSquare /></el-icon>
            <span>评论管理</span>
          </el-menu-item>
          <el-menu-item index="iot_devices">
            <el-icon><Cpu /></el-icon>
            <span>IoT设备</span>
          </el-menu-item>
          <el-menu-item index="fishing_records">
            <el-icon><DataLine /></el-icon>
            <span>垂钓记录</span>
          </el-menu-item>
          <el-menu-item index="fish_caught">
            <el-icon><Trophy /></el-icon>
            <span>渔获管理</span>
          </el-menu-item>
          <el-menu-item index="suggestions">
            <el-icon><Compass /></el-icon>
            <span>垂钓建议</span>
          </el-menu-item>
        </el-sub-menu>
      </el-menu>
    </el-scrollbar>

    <!-- 底部用户信息 -->
    <div class="sidebar-footer">
      <div class="sidebar-user">
        <el-avatar :size="30" :icon="UserFilled" class="sidebar-avatar" />
        <div class="sidebar-user-info">
          <span class="sidebar-username">{{ username }}</span>
          <span class="sidebar-role">{{ roleName }}</span>
        </div>
      </div>
      <el-tooltip content="退出登录" placement="top">
        <el-button :icon="SwitchButton" circle size="small" type="danger" plain @click="$emit('logout')" />
      </el-tooltip>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  Ship, MapLocation, Location, Monitor, Connection, Bell,
  Notification, User, UserFilled, SwitchButton, Folder, Operation,
  ChatDotRound, Document, Cpu, DataLine, ChatLineSquare, Trophy, Compass
} from '@element-plus/icons-vue'

defineProps<{
  active: string
  showUsers: boolean
  username: string
  roleName: string
}>()

defineEmits<{
  (e: 'update:active', value: string): void
  (e: 'logout'): void
}>()
</script>

<style scoped>
.admin-sidebar {
  width: 220px;
  display: flex;
  flex-direction: column;
  background: #fff;
  border-right: 1px solid #e4e7ed;
  flex-shrink: 0;
}

/* 品牌区域 */
.sidebar-brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px;
  border-bottom: 1px solid #e4e7ed;
}

.brand-text {
  font-size: 15px;
  font-weight: 700;
  color: #303133;
  letter-spacing: 1px;
}

/* 菜单滚动区域 */
.sidebar-menu-scroll {
  flex: 1;
  overflow: hidden;
}

.sidebar-menu-scroll :deep(.el-menu) {
  border-right: none;
}

.sidebar-menu-scroll :deep(.el-sub-menu__title) {
  font-size: 13px;
  font-weight: 600;
}

.sidebar-menu-scroll :deep(.el-menu-item) {
  height: 42px;
  line-height: 42px;
  font-size: 14px;
}

/* 底部用户区 */
.sidebar-footer {
  padding: 12px 16px;
  border-top: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.sidebar-user {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.sidebar-avatar {
  background: #409eff;
  flex-shrink: 0;
}

.sidebar-user-info {
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.sidebar-username {
  font-size: 13px;
  font-weight: 600;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.sidebar-role {
  font-size: 11px;
  color: #909399;
}
</style>
