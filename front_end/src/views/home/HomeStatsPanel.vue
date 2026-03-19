<template>
  <el-card class="stats-card">
    <div v-loading="loading">
      <el-skeleton :rows="1" animated :loading="loading">
        <template #default>
          <div v-if="stats.length > 0">
            <el-row :gutter="20">
              <el-col
                v-for="stat in stats"
                :key="stat.key"
                :xs="8" :sm="6" :md="4" :lg="3"
              >
                <el-statistic :title="stat.label" :value="stat.value" class="stat-item">
                  <template #suffix>
                    <el-icon class="stat-icon" :style="{ color: stat.color }">
                      <component :is="stat.icon" />
                    </el-icon>
                  </template>
                </el-statistic>
              </el-col>
            </el-row>
          </div>
          <div v-else class="no-data-message">
            <el-empty description="暂无统计数据" />
          </div>
        </template>
      </el-skeleton>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { SystemSummary } from '@/types'
import {
  Location, MapLocation, Monitor, Connection, UserFilled,
  Warning, Bell, CircleCheckFilled, Avatar
} from '@element-plus/icons-vue'

const props = defineProps<{
  summary: SystemSummary | null
  loading: boolean
}>()

const STATS_LABELS: Record<string, { label: string; icon: any; color: string }> = {
  open_spots: { label: '开放水域', icon: MapLocation, color: '#409eff' },
  total_spots: { label: '水域总数', icon: Location, color: '#67c23a' },
  online_devices: { label: '在线设备', icon: Monitor, color: '#e6a23c' },
  total_devices: { label: '设备总数', icon: Connection, color: '#909399' },
  online_gateways: { label: '在线网关', icon: CircleCheckFilled, color: '#17a2b8' },
  total_gateways: { label: '网关总数', icon: Connection, color: '#909399' },
  total_users: { label: '注册用户', icon: Avatar, color: '#6610f2' },
  total_fishing_count: { label: '当前垂钓人数', icon: UserFilled, color: '#67c23a' },
  active_reminders: { label: '活跃提醒', icon: Warning, color: '#f56c6c' },
  recent_notices: { label: '近期通知', icon: Bell, color: '#17a2b8' },
}

const stats = computed(() => {
  if (!props.summary) return []
  const keys: (keyof SystemSummary)[] = [
    'total_fishing_count', 'open_spots', 'online_devices',
    'total_users', 'active_reminders', 'recent_notices',
    'online_gateways', 'total_spots', 'total_devices'
  ]
  return keys.map(key => ({
    key,
    value: props.summary![key],
    ...STATS_LABELS[key]
  })).filter(s => s.label)
})
</script>

<style scoped>
.stats-card {
  margin-bottom: 30px;
  border-radius: 12px !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08) !important;
  transition: all 0.3s;
}

.stats-card:hover {
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.12) !important;
}

.stat-item {
  padding: 16px;
  position: relative;
  overflow: hidden;
}

.stat-item::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  width: 100%;
  height: 3px;
  background: linear-gradient(90deg, #409eff, #36b5ff);
  transform: scaleX(0);
  transform-origin: left;
  transition: transform 0.3s ease;
}

.stat-item:hover::after {
  transform: scaleX(1);
}

.stat-item :deep(.el-statistic__content) {
  font-size: 28px !important;
  font-weight: 600;
  background: linear-gradient(45deg, #409eff, #36b5ff);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-item :deep(.el-statistic__title) {
  font-size: 14px;
  color: #888;
  letter-spacing: 1px;
}

.stat-icon {
  margin-left: 5px;
  font-size: 18px;
  color: #409eff;
}

.no-data-message {
  padding: 30px 0;
  text-align: center;
}

@media (max-width: 768px) {
  .stats-card { margin-bottom: 15px; }
  .stat-item :deep(.el-statistic__content) { font-size: 20px !important; }
  .stat-item :deep(.el-statistic__title) { font-size: 12px; }
}
</style>
