<template>
  <el-card class="dashboard-card">
    <template #header>
      <div class="card-header">
        <el-icon class="card-icon notice-icon"><ChatLineRound /></el-icon>
        <span class="card-title">近期通知</span>
      </div>
    </template>
    <div v-loading="loading">
      <el-skeleton :rows="2" animated :loading="loading">
        <template #default>
          <div v-if="notices.length > 0">
            <el-alert
              v-for="n in notices"
              :key="n.id"
              type="info"
              show-icon
              class="animated-alert mb-10"
              :closable="false"
            >
              <template #icon>
                <el-icon><Bell /></el-icon>
              </template>
              <template #default>
                <div class="alert-content">
                  <div class="notice-inline">
                    <span class="notice-title-inline">{{ n.title }}</span>
                    <span class="notice-time-inline">{{ formatTime(n.timestamp) }}</span>
                  </div>
                  <div class="notice-excerpt">{{ n.content }}</div>
                </div>
              </template>
            </el-alert>
          </div>
          <div v-else class="no-data-text">暂无重要通知</div>
        </template>
      </el-skeleton>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { Notice } from '@/types'
import { ChatLineRound, Bell } from '@element-plus/icons-vue'

defineProps<{
  notices: Notice[]
  loading: boolean
}>()

function formatTime(ts: string): string {
  const d = new Date(ts)
  const now = new Date()
  const diff = now.getTime() - d.getTime()
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`
  return d.toLocaleDateString('zh-CN')
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

.notice-icon { background-color: #6610f2; }

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

.alert-content { display: flex; flex-direction: column; gap: 4px; }

.notice-inline {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.notice-title-inline {
  font-weight: 600;
  font-size: 13px;
  color: #303133;
}

.notice-time-inline {
  font-size: 11px;
  color: #c0c4cc;
}

.notice-excerpt {
  font-size: 12px;
  color: #909399;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.mb-10 { margin-bottom: 10px; }

.no-data-text {
  padding: 20px;
  text-align: center;
  color: #909399;
  font-size: 14px;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
