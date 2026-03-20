<template>
  <div class="sensor-card">
    <span class="device-title">设备实时数据</span>
    <div class="sensor-grid">
      <div class="sensor-item">
        <span class="sensor-label">温度</span>
        <span class="sensor-value">{{ data.temperature?.toFixed(1) ?? '--' }}°C</span>
      </div>
      <div class="sensor-item">
        <span class="sensor-label">湿度</span>
        <span class="sensor-value">{{ data.humidity?.toFixed(1) ?? '--' }}%</span>
      </div>
      <div class="sensor-item">
        <span class="sensor-label">拉力</span>
        <span class="sensor-value">{{ data.pulling?.toFixed(2) ?? '--' }}N</span>
      </div>
      <div class="sensor-item">
        <span class="sensor-label">气压</span>
        <span class="sensor-value">{{ data.pressure?.toFixed(0) ?? '--' }}hPa</span>
      </div>
    </div>
    <div v-if="data.gpsInfo" class="sensor-extra">
      <span class="sensor-label">GPS</span>
      <span>{{ data.gpsInfo }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { IoTDeviceData } from '@/types'
defineProps<{ data: IoTDeviceData }>()
</script>

<style scoped>
.sensor-card {
  margin-top: 12px;
  padding-top: 12px;
  border-top: 1px dashed #e5e6eb;
}
.device-title { font-size: 13px; font-weight: 600; color: #1d2129; display: block; margin-bottom: 8px; }

.sensor-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px; }
.sensor-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8px;
  background: #fff;
  border-radius: 6px;
  border: 1px solid #f0f0f0;
}
.sensor-label { font-size: 11px; color: #86909c; margin-bottom: 2px; }
.sensor-value { font-size: 14px; font-weight: 600; color: #1d2129; }
.sensor-extra { margin-top: 8px; font-size: 12.5px; color: #4e5969; display: flex; gap: 8px; }

@media (max-width: 768px) {
  .sensor-grid { grid-template-columns: repeat(2, 1fr); }
}
</style>
