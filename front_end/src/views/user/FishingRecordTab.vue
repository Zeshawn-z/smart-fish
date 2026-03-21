<template>
  <div class="record-tab">
    <div class="tab-header">
      <h3>我的垂钓记录</h3>
      <el-button type="primary" size="small" :icon="Plus" @click="showCreateDialog = true">新增记录</el-button>
    </div>

    <div v-loading="recordStore.isLoading" class="record-list">
      <el-empty v-if="recordStore.records.length === 0 && !recordStore.isLoading" description="暂无垂钓记录" />

      <div v-for="record in recordStore.records" :key="record.id" class="record-card">
        <div class="record-header">
          <span class="record-id">#{{ record.id }}</span>
          <span class="record-time">{{ formatTime(record.start_time) }} - {{ formatTime(record.end_time) }}</span>
        </div>

        <div class="record-info">
          <div class="info-item">
            <el-icon><Location /></el-icon>
            <span>{{ record.latitude.toFixed(4) }}, {{ record.longitude.toFixed(4) }}</span>
          </div>
          <div v-if="record.device_id" class="info-item">
            <el-icon><Monitor /></el-icon>
            <span>设备: {{ record.device_id }}</span>
          </div>
        </div>

        <!-- 渔获列表 -->
        <div v-if="(record.fish_caught || record.caught || []).length > 0" class="caught-section">
          <span class="section-label">渔获 ({{ (record.fish_caught || record.caught || []).length }})</span>
          <div class="caught-list">
            <div v-for="fish in (record.fish_caught || record.caught || [])" :key="fish.id || fish.fish_id" class="caught-item">
              <div class="fish-info">
                <el-tag size="small">{{ fish.fish_type || '未知鱼种' }}</el-tag>
                <span class="fish-weight">{{ fish.weight }}kg</span>
                <span v-if="fish.bait_type" class="fish-bait">饵料: {{ fish.bait_type }}</span>
                <span v-if="fish.fishing_depth" class="fish-depth">深度: {{ fish.fishing_depth }}m</span>
              </div>
              <el-image v-if="fish.image_url" :src="fish.image_url" fit="cover" :preview-src-list="[fish.image_url]" class="fish-thumb" />
            </div>
          </div>
        </div>

        <!-- 设备传感器数据 -->
        <SensorDataCard v-if="record.device_data" :data="record.device_data" />
      </div>
    </div>

    <CreateRecordDialog v-model="showCreateDialog" @created="refreshRecords" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, Location, Monitor } from '@element-plus/icons-vue'
import { useFishingRecordStore } from '@/stores/community'
import { useAuthStore } from '@/stores/auth'

import SensorDataCard from './SensorDataCard.vue'
import CreateRecordDialog from './CreateRecordDialog.vue'

const authStore = useAuthStore()
const recordStore = useFishingRecordStore()
const showCreateDialog = ref(false)

function formatTime(t: string) {
  if (!t) return '--'
  return t.replace('T', ' ').slice(0, 16)
}

function refreshRecords() {
  if (authStore.user?.id) {
    recordStore.fetchRecords(authStore.user.id)
  }
}

onMounted(() => refreshRecords())
</script>

<style scoped>
.tab-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.tab-header h3 { margin: 0; font-size: 16px; font-weight: 600; color: #1d2129; }

.record-list { display: flex; flex-direction: column; gap: 12px; }
.record-card { background: #fafbfc; border-radius: 8px; border: 1px solid #f0f0f0; padding: 16px; }

.record-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
.record-id { font-weight: 700; color: #5b8ff9; font-size: 14px; }
.record-time { font-size: 12.5px; color: #86909c; }

.record-info { display: flex; gap: 20px; margin-bottom: 12px; font-size: 13px; color: #4e5969; }
.info-item { display: flex; align-items: center; gap: 4px; }

.caught-section { margin-top: 12px; padding-top: 12px; border-top: 1px dashed #e5e6eb; }
.section-label { font-size: 13px; font-weight: 600; color: #1d2129; display: block; margin-bottom: 8px; }
.caught-list { display: flex; flex-direction: column; gap: 8px; }
.caught-item { display: flex; justify-content: space-between; align-items: center; padding: 8px 10px; background: #fff; border-radius: 6px; border: 1px solid #f0f0f0; }
.fish-info { display: flex; align-items: center; gap: 10px; font-size: 13px; color: #4e5969; flex-wrap: wrap; }
.fish-weight { font-weight: 600; color: #1d2129; }
.fish-thumb { width: 50px; height: 50px; border-radius: 6px; flex-shrink: 0; }
</style>
