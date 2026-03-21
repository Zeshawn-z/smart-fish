<template>
  <div class="panel spot-detail">
    <div class="panel-title">水域概况</div>
    <div v-if="spot" class="detail-content">
      <div class="detail-header">
        <div class="spot-status-dot" :style="{ background: statusColor }"></div>
        <div class="spot-name">{{ spot.name }}</div>
        <span class="spot-type">{{ WATER_TYPE_MAP[spot.water_type] }} · {{ SPOT_STATUS_MAP[spot.status] }}</span>
      </div>

      <div class="detail-metrics">
        <div class="metric">
          <span class="metric-val">{{ spot.capacity }}</span>
          <span class="metric-label">容量</span>
        </div>
        <div class="metric" v-if="spot.bound_device">
          <span class="metric-val" :style="{ color: spot.bound_device.status === 'online' ? '#00ff88' : '#ff6b6b' }">
            {{ spot.bound_device.status === 'online' ? '在线' : '离线' }}
          </span>
          <span class="metric-label">设备</span>
        </div>
        <div class="metric" v-if="spot.bound_device?.water_temp">
          <span class="metric-val">{{ spot.bound_device.water_temp.toFixed(1) }}°</span>
          <span class="metric-label">水温</span>
        </div>
        <div class="metric" v-if="spot.bound_device?.fishing_count !== undefined">
          <span class="metric-val accent">{{ spot.bound_device.fishing_count }}</span>
          <span class="metric-label">在钓</span>
        </div>
      </div>

      <div class="detail-desc" v-if="spot.description" :title="spot.description">{{ spot.description }}</div>

      <div class="detail-region" v-if="region">
        📍 {{ region.province }} · {{ region.city }}
      </div>
    </div>
    <div v-else class="empty-hint">
      <div class="empty-icon">🗺️</div>
      <div>在地图上点击水域查看详情</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { FishingSpot, Region } from '@/types'
import { WATER_TYPE_MAP, SPOT_STATUS_MAP } from '@/types'

const props = defineProps<{
  spot: FishingSpot | null
  regions: Region[]
}>()

const region = computed(() => {
  if (!props.spot) return null
  return props.regions.find(r => r.id === props.spot!.region_id) ?? null
})

const statusColor = computed(() => {
  const s = props.spot?.status
  return s === 'open' ? '#00ff88' : s === 'closed' ? '#ff6b6b' : '#ffd93d'
})
</script>

<style scoped>
.spot-detail {
  padding: 10px;
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.detail-content {
  flex: 1;
  overflow: hidden;
  margin-top: 4px;
}

.detail-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.spot-status-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
  box-shadow: 0 0 6px currentColor;
}

.spot-name {
  font-size: 14px;
  font-weight: 600;
  color: #e8f0ff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  min-width: 0;
}

.spot-type {
  font-size: 10px;
  color: #5a7a9a;
  white-space: nowrap;
  flex-shrink: 0;
  margin-left: auto;
}

.detail-metrics {
  display: flex;
  gap: 6px;
  margin-bottom: 8px;
}

.metric {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 5px 4px;
  border-radius: 5px;
  background: rgba(0, 212, 255, 0.04);
  border: 1px solid rgba(0, 212, 255, 0.08);
  min-width: 0;
}

.metric-val {
  font-size: 14px;
  font-weight: 700;
  color: #00d4ff;
  font-family: 'DIN Alternate', 'Courier New', monospace;
  line-height: 1.2;
}

.metric-val.accent {
  color: #00ff88;
}

.metric-label {
  font-size: 9px;
  color: #5a7a9a;
  margin-top: 1px;
}

.detail-desc {
  font-size: 11px;
  color: #8ba3c7;
  line-height: 1.4;
  margin-bottom: 6px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
  cursor: default;
}

.detail-region {
  font-size: 10px;
  color: #5a7a9a;
  line-height: 1;
}

.empty-hint {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 12px 0;
  color: #5a7a9a;
  font-size: 12px;
}

.empty-icon {
  font-size: 20px;
  opacity: 0.5;
}
</style>
