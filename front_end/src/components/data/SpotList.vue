<script setup lang="ts">
import type { PopularSpot } from '@/types'
import { WATER_TYPE_MAP } from '@/types'
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  spots: PopularSpot[]
  loading?: boolean
  emptyText?: string
  maxHeight?: string
}>(), {
  loading: false,
  emptyText: '暂无水域数据',
  maxHeight: '400px'
})

// 负载率相关
const calculateRatio = (count: number | undefined, capacity: number | undefined) => {
  if (!capacity) return -1
  return (count || 0) / capacity
}

const getProgressColor = (rate: number) => {
  if (rate >= 0.9) return '#F56C6C'
  if (rate >= 0.7) return '#E6A23C'
  if (rate >= 0.3) return '#67C23A'
  return '#409EFF'
}

const getTagType = (rate: number) => {
  if (rate === -1) return 'info'
  if (rate > 0.9) return 'danger'
  if (rate > 0.7) return 'warning'
  if (rate > 0.3) return 'success'
  return 'info'
}

const getStatusText = (rate: number) => {
  if (rate === -1) return '未知'
  if (rate >= 0.9) return '拥挤'
  if (rate >= 0.7) return '较忙'
  if (rate >= 0.3) return '适中'
  return '空闲'
}

const processedSpots = computed(() => {
  return props.spots.map(spot => {
    const fishingCount = spot.bound_device?.fishing_count ?? spot.total_fishing_count ?? 0
    const ratio = calculateRatio(fishingCount, spot.capacity)
    return {
      ...spot,
      fishingCount,
      ratio,
      progressColor: getProgressColor(ratio),
      tagType: getTagType(ratio),
      statusText: getStatusText(ratio),
      percentDisplay: ratio === -1 ? '未知' : `${Math.floor(ratio * 100)}%`
    }
  })
})
</script>

<template>
  <el-skeleton :rows="4" animated :loading="loading">
    <template #default>
      <div v-if="spots && spots.length > 0" class="spot-list-container">
        <el-scrollbar :height="maxHeight === 'auto' ? undefined : maxHeight" class="no-horizontal-scroll">
          <div class="row-container">
            <el-row :gutter="14">
              <el-col
                v-for="spot in processedSpots"
                :key="spot.id"
                :xs="24" :sm="12" :md="12" :lg="8"
                class="spot-col"
              >
                <div class="spot-card">
                  <div class="card-line">
                    <h4 class="spot-name">{{ spot.name }}</h4>

                    <div class="spot-data-section">
                      <div class="data-chip">
                        <span class="count-value" :style="{ color: spot.progressColor }">
                          {{ spot.fishingCount }}
                        </span>
                        <span v-if="spot.capacity" class="capacity-indicator">/ {{ spot.capacity }}</span>
                        <span class="stat-badge" :class="spot.tagType">
                          {{ spot.percentDisplay }}
                        </span>
                      </div>
                    </div>
                  </div>

                  <div class="card-meta">
                    <div class="meta-left">
                      <el-tag v-if="spot.region?.city" size="small" type="info" effect="plain" class="meta-tag">
                        {{ spot.region.city }}
                      </el-tag>
                      <el-tag size="small" effect="plain" class="meta-tag">{{ WATER_TYPE_MAP[spot.water_type] || spot.water_type }}</el-tag>
                    </div>
                    <span class="status-text" :style="{ color: spot.progressColor }">{{ spot.statusText }}</span>
                  </div>

                  <div class="progress-border-container">
                    <div
                      class="progress-border"
                      :style="{
                        width: spot.ratio === -1 ? '0%' : Math.floor(spot.ratio * 100) + '%',
                        backgroundColor: spot.progressColor
                      }"
                    ></div>
                  </div>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-scrollbar>
      </div>

      <div v-else class="no-data-message">
        <el-empty :description="emptyText" />
      </div>
    </template>
  </el-skeleton>
</template>

<style scoped>
.spot-list-container {
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.row-container {
  width: 100%;
  padding: 0 7px;
  box-sizing: border-box;
  overflow: hidden;
}

.no-horizontal-scroll {
  overflow-x: hidden !important;
}

.no-horizontal-scroll :deep(.el-scrollbar__wrap) {
  overflow-x: hidden !important;
}

.spot-col {
  margin-bottom: 14px;
}

.spot-card {
  background-color: #ffffff;
  border-radius: 8px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.05);
  transition: all 0.2s;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  border: 1.5px solid #ebeef5;
  position: relative;
  height: 100%;
  cursor: pointer;
}

.spot-card:hover {
  box-shadow: 0 4px 12px rgba(64, 158, 255, 0.12);
  border-color: #bde3ff;
}

.card-line {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 14px 4px;
  min-height: 36px;
  gap: 4px;
  overflow: hidden;
}

.spot-name {
  font-size: 13px;
  font-weight: 600;
  color: #303133;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  min-width: 0;
  margin-right: 4px;
}

.card-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 2px 14px 10px;
  font-size: 12px;
}

.meta-left {
  display: flex;
  align-items: center;
  gap: 6px;
}

.meta-tag {
  font-size: 11px;
  padding: 0 6px;
  height: 20px;
  line-height: 20px;
}

.status-text {
  font-weight: 600;
  font-size: 12px;
  flex-shrink: 0;
}

.progress-border-container {
  height: 3px;
  width: 100%;
  background-color: #f0f0f0;
  position: absolute;
  bottom: 0;
  left: 0;
}

.progress-border {
  height: 100%;
  transition: all 0.3s ease;
  border-radius: 0 2px 2px 0;
}

.spot-data-section {
  display: flex;
  align-items: center;
  white-space: nowrap;
}

.data-chip {
  display: flex;
  align-items: center;
}

.count-value {
  font-size: 16px;
  font-weight: bold;
}

.capacity-indicator {
  font-size: 12px;
  opacity: 0.7;
  margin-left: 2px;
  margin-right: 8px;
}

.stat-badge {
  font-size: 10px;
  padding: 1px 6px;
  border-radius: 10px;
  color: #fff;
  margin-left: 4px;
}

.stat-badge.success { background-color: #67C23A; }
.stat-badge.info { background-color: #409EFF; }
.stat-badge.warning { background-color: #E6A23C; }
.stat-badge.danger { background-color: #F56C6C; }

.no-data-message {
  padding: 30px 0;
  text-align: center;
}

@media (max-width: 768px) {
  .card-line {
    padding: 10px 10px 3px;
  }
  .spot-name {
    font-size: 12px;
  }
  .count-value {
    font-size: 14px;
  }
  .capacity-indicator {
    font-size: 11px;
  }
  .card-meta {
    padding: 2px 10px 8px;
  }
  .meta-tag {
    font-size: 10px;
    padding: 0 4px;
    height: 18px;
    line-height: 18px;
  }
}
</style>
