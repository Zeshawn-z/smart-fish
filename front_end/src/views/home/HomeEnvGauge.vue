<template>
  <el-card class="dashboard-card">
    <template #header>
      <div class="chart-header">
        <div class="card-header">
          <el-icon class="card-icon env-icon"><Odometer /></el-icon>
          <span class="card-title">环境数据监测</span>
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

    <!-- 概览指标行 -->
    <div class="env-overview" v-if="latestEnv">
      <div class="env-stat">
        <span class="env-stat-value water-temp">{{ latestEnv.water_temp }}°C</span>
        <span class="env-stat-label">水温</span>
      </div>
      <div class="env-stat">
        <span class="env-stat-value air-temp">{{ latestEnv.air_temp }}°C</span>
        <span class="env-stat-label">气温</span>
      </div>
      <div class="env-stat">
        <span class="env-stat-value humidity">{{ latestEnv.humidity }}%</span>
        <span class="env-stat-label">湿度</span>
      </div>
      <div class="env-stat">
        <span class="env-stat-value ph">{{ latestEnv.ph }}</span>
        <span class="env-stat-label">pH</span>
      </div>
      <div class="env-stat">
        <span class="env-stat-value do-val">{{ latestEnv.dissolved_oxygen }}</span>
        <span class="env-stat-label">溶氧 mg/L</span>
      </div>
      <div class="env-stat">
        <span class="env-stat-value turbidity">{{ latestEnv.turbidity }}</span>
        <span class="env-stat-label">浑浊度 NTU</span>
      </div>
    </div>

    <!-- 历史折线图 -->
    <BaseChart v-if="chartOption" :option="chartOption" height="280px" />
    <BaseChart
      v-else
      :option="{ xAxis: { type: 'category', data: [] }, yAxis: { type: 'value' }, series: [] }"
      height="280px"
    />
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BaseChart from '@/components/chart/BaseChart.vue'
import type { EnvironmentData, PopularSpot } from '@/types'
import type { EChartsOption } from 'echarts'
import { Odometer } from '@element-plus/icons-vue'

const props = defineProps<{
  spotId: number | null
  spots: PopularSpot[]
  environmentData: EnvironmentData[]
}>()

defineEmits<{
  (e: 'update:spotId', value: number | null): void
}>()

// 最新一条数据用于概览指标
const latestEnv = computed(() => {
  if (!props.environmentData.length) return null
  const sorted = [...props.environmentData].sort(
    (a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
  )
  return sorted[0]
})

const chartOption = computed<EChartsOption | null>(() => {
  if (!props.spotId || !props.environmentData.length) return null

  const sorted = [...props.environmentData].sort(
    (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
  )

  const timestamps = sorted.map(r => {
    const d = new Date(r.timestamp)
    return `${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')} ${String(d.getHours()).padStart(2, '0')}:00`
  })

  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(255,255,255,0.96)',
      borderColor: '#e4e7ed',
      textStyle: { fontSize: 12, color: '#333' }
    },
    legend: {
      data: ['水温', '气温', 'pH', '溶氧'],
      top: 0,
      left: 'center',
      textStyle: { fontSize: 12 },
      itemWidth: 20,
      itemHeight: 10,
      itemGap: 20
    },
    grid: { left: 50, right: 50, top: 35, bottom: 40 },
    xAxis: {
      type: 'category',
      data: timestamps,
      axisLabel: { fontSize: 10, color: '#999', rotate: 30 },
      axisLine: { lineStyle: { color: '#e4e7ed' } },
      axisTick: { show: false }
    },
    yAxis: [
      {
        type: 'value',
        name: '°C',
        nameTextStyle: { fontSize: 11, color: '#999' },
        axisLabel: { fontSize: 11, color: '#999' },
        splitLine: { lineStyle: { type: 'dashed', color: '#f0f0f0' } }
      },
      {
        type: 'value',
        name: '',
        axisLabel: { fontSize: 11, color: '#999' },
        splitLine: { show: false }
      }
    ],
    series: [
      {
        name: '水温',
        type: 'line',
        smooth: true,
        symbol: 'none',
        lineStyle: { width: 2 },
        itemStyle: { color: '#409eff' },
        areaStyle: { color: 'rgba(64,158,255,0.08)' },
        data: sorted.map(r => r.water_temp)
      },
      {
        name: '气温',
        type: 'line',
        smooth: true,
        symbol: 'none',
        lineStyle: { width: 2 },
        itemStyle: { color: '#e6a23c' },
        areaStyle: { color: 'rgba(230,162,60,0.08)' },
        data: sorted.map(r => r.air_temp)
      },
      {
        name: 'pH',
        type: 'line',
        smooth: true,
        symbol: 'none',
        yAxisIndex: 1,
        lineStyle: { width: 2 },
        itemStyle: { color: '#67c23a' },
        data: sorted.map(r => r.ph)
      },
      {
        name: '溶氧',
        type: 'line',
        smooth: true,
        symbol: 'none',
        yAxisIndex: 1,
        lineStyle: { width: 2, type: 'dashed' },
        itemStyle: { color: '#17a2b8' },
        data: sorted.map(r => r.dissolved_oxygen)
      }
    ]
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

.env-icon { background-color: #17a2b8; }

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

/* 概览指标行 */
.env-overview {
  display: flex;
  justify-content: space-around;
  padding: 10px 0 14px;
  border-bottom: 1px solid #f5f5f5;
  margin-bottom: 10px;
  flex-wrap: wrap;
  gap: 4px 0;
}

.env-stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 70px;
}

.env-stat-value {
  font-size: 18px;
  font-weight: 700;
  line-height: 1.3;
}

.env-stat-label {
  font-size: 11px;
  color: #909399;
  margin-top: 2px;
}

.water-temp { color: #409eff; }
.air-temp { color: #e6a23c; }
.humidity { color: #909399; }
.ph { color: #67c23a; }
.do-val { color: #17a2b8; }
.turbidity { color: #f56c6c; }

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

  .env-overview { gap: 8px; }
  .env-stat { min-width: 60px; }
  .env-stat-value { font-size: 15px; }
}
</style>
