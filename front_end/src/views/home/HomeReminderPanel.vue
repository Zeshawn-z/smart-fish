<template>
  <el-card class="dashboard-card">
    <template #header>
      <div class="card-header">
        <el-icon class="card-icon alert-icon"><Warning /></el-icon>
        <span class="card-title">公开提醒</span>
      </div>
    </template>
    <div v-loading="loading">
      <el-skeleton :rows="2" animated :loading="loading">
        <template #default>
          <div v-if="reminders.length > 0">
            <el-alert
              v-for="r in reminders"
              :key="r.id"
              :type="getAlertType(r.level)"
              show-icon
              class="animated-alert mb-10"
              :closable="false"
            >
              <template #icon>
                <el-icon><component :is="getReminderIcon(r.reminder_type)" /></el-icon>
              </template>
              <template #default>
                <div class="alert-content">
                  <span class="alert-message">{{ r.message }}</span>
                  <router-link :to="`/reminders?id=${r.id}`" class="alert-link">查看详情</router-link>
                </div>
              </template>
            </el-alert>
          </div>
          <div v-else class="no-data-text">暂无安全提醒</div>
        </template>
      </el-skeleton>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { Reminder } from '@/types'
import { Warning, UserFilled, Odometer, Bell } from '@element-plus/icons-vue'

defineProps<{
  reminders: Reminder[]
  loading: boolean
}>()

const getReminderIcon = (type: string) => {
  switch (type) {
    case 'weather': return Warning
    case 'capacity': return UserFilled
    case 'environment': return Odometer
    default: return Bell
  }
}

const getAlertType = (level: number) => {
  switch (level) {
    case 3: return 'error'
    case 2: return 'warning'
    case 1: return 'info'
    default: return 'success'
  }
}
</script>

<style scoped>
.dashboard-card {
  margin-bottom: 25px;
  border-radius: 12px !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08) !important;
  border: none;
  transition: all 0.3s;
}

.dashboard-card:hover {
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1) !important;
}

.dashboard-card :deep(.el-card__header) {
  padding: 18px 24px;
  background: linear-gradient(45deg, #fafafa, #f6f9ff) !important;
  border-bottom: 1px solid #e4e7ed;
}

.card-header { display: flex; align-items: center; }

.card-icon {
  margin-right: 8px;
  font-size: 20px;
  padding: 3px;
  border-radius: 8px;
  color: white;
}

.alert-icon { background-color: #dc3545; }

.card-title {
  font-size: 17px !important;
  font-weight: 600;
  color: #333;
  letter-spacing: 0.5px;
}

.animated-alert {
  transition: all 0.3s ease;
  animation: fadeIn 0.5s ease-in-out;
}

.animated-alert:hover { transform: translateX(5px); }

.alert-content { display: flex; align-items: flex-start; gap: 8px; }

.alert-message { font-size: 13px; flex: 1; line-height: 1.5; }

.alert-link {
  margin-left: 10px;
  font-size: 12px;
  color: #409EFF;
  text-decoration: none;
  white-space: nowrap;
}

.alert-link:hover { text-decoration: underline; }

.mb-10 { margin-bottom: 10px; }

.no-data-text {
  padding: 20px;
  text-align: center;
  color: #909399;
  font-size: 14px;
}

.el-alert--error {
  background-color: #fff0f0 !important;
  border: 1px solid rgba(245, 108, 108, 0.3);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
