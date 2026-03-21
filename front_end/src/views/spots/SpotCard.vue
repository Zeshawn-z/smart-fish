<template>
  <el-card shadow="hover" class="spot-card" @click="$emit('click', spot)">
    <!-- 头部：名称 + 状态 -->
    <div class="spot-card-header">
      <span class="spot-name" :title="spot.name">{{ spot.name }}</span>
      <el-tag :type="statusTagType" size="small" effect="dark" class="status-tag">
        {{ SPOT_STATUS_MAP[spot.status] }}
      </el-tag>
    </div>

    <!-- 地区 + 类型 -->
    <div class="spot-meta">
      <span><el-icon><Location /></el-icon> {{ spot.region?.province }} · {{ spot.region?.city }}</span>
      <span><el-icon><Ship /></el-icon> {{ WATER_TYPE_MAP[spot.water_type] }}</span>
    </div>

    <!-- 描述（固定两行高度，hover 显示完整内容） -->
    <el-tooltip
      :content="spot.description || '暂无描述'"
      placement="top"
      :disabled="!spot.description || spot.description.length <= 30"
      :show-after="300"
      popper-class="spot-desc-tooltip"
    >
      <p class="spot-desc">{{ spot.description || '暂无描述' }}</p>
    </el-tooltip>

    <!-- 环境数据仪表区 — 始终保持统一四宫格布局 -->
    <div class="env-panel" :class="{ 'env-panel--disabled': !deviceOnline }">
      <!-- 设备在线：显示真实数据 -->
      <template v-if="deviceOnline">
        <div class="env-item">
          <div class="env-icon water-temp-icon">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 2C10.34 2 9 3.34 9 5V14.54C7.78 15.36 7 16.74 7 18.25C7 20.87 9.24 23 12 23C14.76 23 17 20.87 17 18.25C17 16.74 16.22 15.36 15 14.54V5C15 3.34 13.66 2 12 2Z" stroke="currentColor" stroke-width="1.5" fill="none"/>
              <path d="M12 18.5C13.1 18.5 14 17.6 14 16.5C14 15.7 13.6 15 12.9 14.6L12 14V7H12V14L11.1 14.6C10.4 15 10 15.7 10 16.5C10 17.6 10.9 18.5 12 18.5Z" fill="currentColor" opacity="0.7"/>
            </svg>
          </div>
          <div class="env-data">
            <span class="env-value">{{ spot.bound_device!.water_temp }}<small>°C</small></span>
            <span class="env-label">水温</span>
          </div>
        </div>
        <div class="env-item">
          <div class="env-icon air-temp-icon">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="12" cy="12" r="5" stroke="currentColor" stroke-width="1.5" fill="none"/>
              <path d="M12 2V4M12 20V22M2 12H4M20 12H22M5.64 5.64L7.05 7.05M16.95 16.95L18.36 18.36M18.36 5.64L16.95 7.05M7.05 16.95L5.64 18.36" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            </svg>
          </div>
          <div class="env-data">
            <span class="env-value">{{ spot.bound_device!.air_temp }}<small>°C</small></span>
            <span class="env-label">气温</span>
          </div>
        </div>
        <div class="env-item">
          <div class="env-icon humidity-icon">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 2.69L17.66 8.35C19.63 10.32 20.75 13.01 20.75 15.78C20.75 20.58 16.83 24 12 24C7.17 24 3.25 20.58 3.25 15.78C3.25 13.01 4.37 10.32 6.34 8.35L12 2.69Z" stroke="currentColor" stroke-width="1.5" fill="none"/>
              <path d="M8 16C8 18.21 9.79 20 12 20" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" opacity="0.5"/>
            </svg>
          </div>
          <div class="env-data">
            <span class="env-value">{{ spot.bound_device!.humidity }}<small>%</small></span>
            <span class="env-label">湿度</span>
          </div>
        </div>
        <div class="env-item">
          <div class="env-icon fishing-icon">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 3C10.5 3 9 4.5 9 6V15C7.5 15.8 6.5 17.3 6.5 19C6.5 21.5 8.5 23.5 11 23.5" stroke="currentColor" stroke-width="1.5" fill="none" stroke-linecap="round"/>
              <circle cx="11" cy="19" r="2" fill="currentColor" opacity="0.5"/>
              <path d="M15 6L21 3V9L15 6Z" fill="currentColor" opacity="0.7"/>
              <line x1="12" y1="3" x2="12" y2="15" stroke="currentColor" stroke-width="1.5"/>
            </svg>
          </div>
          <div class="env-data">
            <span class="env-value">{{ spot.bound_device!.fishing_count }}<small>人</small></span>
            <span class="env-label">垂钓中</span>
          </div>
        </div>
      </template>

      <!-- 设备不在线：保持四宫格占位，显示状态提示 -->
      <template v-else>
        <div class="env-item env-item--placeholder">
          <div class="env-icon water-temp-icon env-icon--dim">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 2C10.34 2 9 3.34 9 5V14.54C7.78 15.36 7 16.74 7 18.25C7 20.87 9.24 23 12 23C14.76 23 17 20.87 17 18.25C17 16.74 16.22 15.36 15 14.54V5C15 3.34 13.66 2 12 2Z" stroke="currentColor" stroke-width="1.5" fill="none"/>
            </svg>
          </div>
          <div class="env-data">
            <span class="env-value env-value--dim">--<small>°C</small></span>
            <span class="env-label">水温</span>
          </div>
        </div>
        <div class="env-item env-item--placeholder">
          <div class="env-icon air-temp-icon env-icon--dim">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <circle cx="12" cy="12" r="5" stroke="currentColor" stroke-width="1.5" fill="none"/>
              <path d="M12 2V4M12 20V22M2 12H4M20 12H22" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            </svg>
          </div>
          <div class="env-data">
            <span class="env-value env-value--dim">--<small>°C</small></span>
            <span class="env-label">气温</span>
          </div>
        </div>
        <div class="env-item env-item--placeholder">
          <div class="env-icon humidity-icon env-icon--dim">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 2.69L17.66 8.35C19.63 10.32 20.75 13.01 20.75 15.78C20.75 20.58 16.83 24 12 24C7.17 24 3.25 20.58 3.25 15.78C3.25 13.01 4.37 10.32 6.34 8.35L12 2.69Z" stroke="currentColor" stroke-width="1.5" fill="none"/>
            </svg>
          </div>
          <div class="env-data">
            <span class="env-value env-value--dim">--<small>%</small></span>
            <span class="env-label">湿度</span>
          </div>
        </div>
        <div class="env-item env-item--placeholder">
          <div class="env-icon fishing-icon env-icon--dim">
            <svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 3C10.5 3 9 4.5 9 6V15C7.5 15.8 6.5 17.3 6.5 19C6.5 21.5 8.5 23.5 11 23.5" stroke="currentColor" stroke-width="1.5" fill="none" stroke-linecap="round"/>
              <line x1="12" y1="3" x2="12" y2="15" stroke="currentColor" stroke-width="1.5"/>
            </svg>
          </div>
          <div class="env-data">
            <span class="env-value env-value--dim">--<small>人</small></span>
            <span class="env-label">垂钓中</span>
          </div>
        </div>
      </template>

      <!-- 不在线时的浮层提示 -->
      <div v-if="!deviceOnline" class="env-overlay">
        <div class="env-overlay-content">
          <el-icon v-if="spot.bound_device" :size="16"><Monitor /></el-icon>
          <span class="env-overlay-device" v-if="spot.bound_device" :title="spot.bound_device.name">{{ spot.bound_device.name }}</span>
          <span class="env-overlay-text">{{ envStatusText }}</span>
        </div>
      </div>
    </div>

    <!-- 负载率进度条 -->
    <div class="spot-load-bar" v-if="spot.capacity">
      <div class="load-info">
        <span>负载: {{ loadCount }}/{{ spot.capacity }}</span>
        <span :style="{ color: loadColor }">{{ loadPercent }}%</span>
      </div>
      <el-progress :percentage="loadPercent" :color="loadColor" :stroke-width="6" :show-text="false" />
    </div>

    <!-- 底部 -->
    <div class="spot-card-footer">
      <el-button
        v-if="showFavorite"
        :type="isFavorite ? 'warning' : 'default'"
        size="small"
        @click.stop="$emit('toggle-favorite', spot.id)"
      >
        <el-icon><Star /></el-icon>
        {{ isFavorite ? '已收藏' : '收藏' }}
      </el-button>
      <el-button size="small" type="primary" @click.stop="$emit('click', spot)">
        查看详情
      </el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { FishingSpot } from '@/types'
import { WATER_TYPE_MAP, SPOT_STATUS_MAP } from '@/types'
import { Location, Ship, Monitor, Star } from '@element-plus/icons-vue'

const props = defineProps<{
  spot: FishingSpot
  isFavorite: boolean
  showFavorite: boolean
}>()

defineEmits<{
  (e: 'click', spot: FishingSpot): void
  (e: 'toggle-favorite', spotId: number): void
}>()

const statusTagType = computed(() => {
  return props.spot.status === 'open' ? 'success' : props.spot.status === 'closed' ? 'danger' : 'warning'
})

const deviceOnline = computed(() => {
  return !!(props.spot.bound_device && props.spot.bound_device.status === 'online')
})

const envStatusText = computed(() => {
  if (!props.spot.bound_device) {
    if (props.spot.status === 'maintenance') return '维护中'
    if (props.spot.status === 'closed') return '已关闭'
    return '未绑定设备'
  }
  return '设备离线'
})

const loadCount = computed(() => props.spot.bound_device?.fishing_count ?? 0)

const loadPercent = computed(() => {
  if (!props.spot.capacity) return 0
  return Math.min(100, Math.round((loadCount.value / props.spot.capacity) * 100))
})

const loadColor = computed(() => {
  const pct = loadPercent.value
  if (pct >= 90) return '#F56C6C'
  if (pct >= 70) return '#E6A23C'
  if (pct >= 30) return '#67C23A'
  return '#409EFF'
})
</script>

<style scoped>
.spot-card {
  cursor: pointer;
  transition: all 0.3s;
  border-radius: 12px !important;
  overflow: hidden;
  position: relative;
  height: 100%;
  display: flex;
  flex-direction: column;
}

.spot-card:hover {
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.1) !important;
}

.spot-card :deep(.el-card__body) {
  padding: 16px 18px;
  display: flex;
  flex-direction: column;
  flex: 1;
}

/* 头部 */
.spot-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
  gap: 8px;
}

.spot-name {
  font-weight: 600;
  font-size: 15px;
  color: #303133;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  min-width: 0;
}

.status-tag {
  flex-shrink: 0;
}

/* 元信息 */
.spot-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  color: #909399;
  font-size: 13px;
  margin-bottom: 8px;
}

.spot-meta span {
  display: flex;
  align-items: center;
  gap: 3px;
}

/* 描述 — 固定两行高度 */
.spot-desc {
  font-size: 13px;
  color: #606266;
  line-height: 1.5;
  height: calc(13px * 1.5 * 2);  /* 固定两行高度 */
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin: 0 0 12px;
}

/* ===== 环境数据仪表面板 ===== */
.env-panel {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 6px;
  padding: 10px 8px;
  background: linear-gradient(135deg, #f6f9ff 0%, #f0f7ff 100%);
  border-radius: 10px;
  margin-bottom: 12px;
  border: 1px solid #e8f0fe;
  position: relative;
  min-height: 82px;
}

.env-panel--disabled {
  background: linear-gradient(135deg, #f9f9fb 0%, #f5f5f7 100%);
  border-color: #ebeef5;
}

.env-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.env-item--placeholder {
  opacity: 0.25;
}

.env-icon {
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 5px;
}

.env-icon svg {
  width: 100%;
  height: 100%;
}

.env-icon--dim {
  filter: grayscale(1);
}

.water-temp-icon {
  color: #409eff;
  background: rgba(64, 158, 255, 0.1);
}

.air-temp-icon {
  color: #e6a23c;
  background: rgba(230, 162, 60, 0.1);
}

.humidity-icon {
  color: #67c23a;
  background: rgba(103, 194, 58, 0.1);
}

.fishing-icon {
  color: #f56c6c;
  background: rgba(245, 108, 108, 0.1);
}

.env-data {
  display: flex;
  flex-direction: column;
  align-items: center;
  line-height: 1.2;
}

.env-value {
  font-size: 15px;
  font-weight: 700;
  color: #303133;
}

.env-value--dim {
  color: #c0c4cc;
}

.env-value small {
  font-size: 11px;
  font-weight: 400;
  color: #909399;
  margin-left: 1px;
}

.env-label {
  font-size: 11px;
  color: #909399;
}

/* 不在线时的浮层提示 */
.env-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.6);
  backdrop-filter: blur(1px);
  border-radius: 10px;
  z-index: 1;
}

.env-overlay-content {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  max-width: 90%;
}

.env-overlay-device {
  font-size: 12px;
  color: #606266;
  max-width: 120px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.env-overlay-text {
  font-size: 12px;
  color: #909399;
  white-space: nowrap;
}

/* 负载条 */
.spot-load-bar {
  margin-bottom: 4px;
}

.load-info {
  display: flex;
  justify-content: space-between;
  font-size: 12px;
  color: #909399;
  margin-bottom: 4px;
}

/* 底部 */
.spot-card-footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  margin-top: auto;
  padding-top: 12px;
}

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .spot-card :deep(.el-card__body) {
    padding: 12px 14px;
  }

  .spot-name { font-size: 14px; }
  .spot-meta { gap: 8px; font-size: 12px; }
  .spot-desc { font-size: 12px; }

  .env-panel {
    grid-template-columns: repeat(2, 1fr);
    gap: 8px;
    padding: 8px;
    min-height: 140px;
  }

  .env-icon { width: 28px; height: 28px; }
  .env-value { font-size: 14px; }

  .env-overlay-device { max-width: 80px; }
}
</style>

<!-- tooltip popper 挂在 body 上，需非 scoped 样式 -->
<style>
.spot-desc-tooltip {
  max-width: 320px !important;
  line-height: 1.6 !important;
}
</style>
