<template>
  <el-card class="dashboard-card">
    <template #header>
      <div class="chart-header">
        <div class="card-header">
          <el-icon class="card-icon trend-icon"><Histogram /></el-icon>
          <span class="card-title">垂钓数量趋势</span>
        </div>
        <div class="chart-controls">
          <el-select
            :model-value="spotId"
            placeholder="选择水域"
            size="small"
            style="width: 180px"
            clearable
            @update:model-value="$emit('update:spotId', $event)"
          >
            <el-option
              v-for="spot in spots"
              :key="spot.id"
              :label="spot.name"
              :value="spot.id"
            />
          </el-select>
        </div>
      </div>
    </template>
    <BaseChart
      v-if="chartOption"
      :option="chartOption"
      height="320px"
    />
    <BaseChart
      v-else
      :option="{ xAxis: { type: 'category', data: [] }, yAxis: { type: 'value' }, series: [] }"
      height="320px"
    />
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BaseChart from '@/components/chart/BaseChart.vue'
import type { HistoricalData, PopularSpot } from '@/types'
import type { EChartsOption } from '@/plugins/echarts'
import { Histogram } from '@element-plus/icons-vue'

const props = defineProps<{
  spotId: number | null
  spots: PopularSpot[]
  historicalData: HistoricalData[]
}>()

defineEmits<{
  (e: 'update:spotId', value: number | null): void
}>()

const chartOption = computed<EChartsOption | null>(() => {
  if (!props.spotId || !props.historicalData.length) return null

  const sorted = [...props.historicalData].sort(
    (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
  )
  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(15, 23, 42, 0.9)',
      textStyle: { color: '#fff', fontSize: 12 }
    },
    xAxis: {
      type: 'category',
      data: sorted.map(d => new Date(d.timestamp).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })),
      axisLabel: { color: '#666', fontSize: 11 },
      axisLine: { lineStyle: { color: '#e0e0e0' } }
    },
    yAxis: {
      type: 'value',
      name: '垂钓人数',
      axisLabel: { color: '#666' },
      splitLine: { lineStyle: { color: '#f0f0f0', type: 'dashed' } }
    },
    series: [{
      name: '垂钓人数',
      type: 'line',
      data: sorted.map(d => d.fishing_count),
      smooth: true,
      lineStyle: { width: 3, color: '#409eff' },
      areaStyle: {
        color: {
          type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
          colorStops: [
            { offset: 0, color: 'rgba(64,158,255,0.4)' },
            { offset: 1, color: 'rgba(64,158,255,0.05)' }
          ]
        } as any
      },
      itemStyle: { color: '#409eff' },
      symbol: 'circle',
      symbolSize: 6
    }],
    grid: { left: '3%', right: '4%', bottom: '10%', top: '15%', containLabel: true }
  }
})
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

.trend-icon { background-color: #409eff; }

.card-title {
  font-size: 17px !important;
  font-weight: 600;
  color: #333;
  letter-spacing: 0.5px;
}

.chart-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.chart-controls {
  display: flex;
  align-items: center;
}

@media (max-width: 768px) {
  .chart-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .chart-controls {
    width: 100%;
    justify-content: flex-end;
  }
}
</style>
