<template>
  <el-dialog :model-value="visible" :title="spot?.name" width="720px" top="5vh" destroy-on-close
    @update:model-value="$emit('update:visible', $event)">
    <div v-if="spot" class="spot-detail">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="区域">{{ spot.region?.province }} · {{ spot.region?.city }}</el-descriptions-item>
        <el-descriptions-item label="水域类型">
          <el-tag size="small">{{ WATER_TYPE_MAP[spot.water_type] }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="spot.status === 'open' ? 'success' : 'danger'" size="small" effect="dark">
            {{ SPOT_STATUS_MAP[spot.status] }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="最大容纳">{{ spot.capacity }} 人</el-descriptions-item>
        <el-descriptions-item label="经度">{{ spot.longitude }}</el-descriptions-item>
        <el-descriptions-item label="纬度">{{ spot.latitude }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ spot.description || '暂无描述' }}</el-descriptions-item>
      </el-descriptions>

      <!-- 设备信息 -->
      <div v-if="spot.bound_device" class="device-info">
        <h4><el-icon><Monitor /></el-icon> 绑定设备信息</h4>
        <el-descriptions :column="3" size="small" border>
          <el-descriptions-item label="设备名">{{ spot.bound_device.name }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="spot.bound_device.status === 'online' ? 'success' : 'info'" size="small" effect="dark">
              {{ spot.bound_device.status === 'online' ? '在线' : '离线' }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="当前垂钓人数">
            <span style="font-weight: 600; color: #409eff">{{ spot.bound_device.fishing_count }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="水温">{{ spot.bound_device.water_temp }}°C</el-descriptions-item>
          <el-descriptions-item label="气温">{{ spot.bound_device.air_temp }}°C</el-descriptions-item>
          <el-descriptions-item label="湿度">{{ spot.bound_device.humidity }}%</el-descriptions-item>
        </el-descriptions>
      </div>

      <!-- 环境趋势图 -->
      <div class="detail-chart" v-if="chartOption">
        <h4><el-icon><TrendCharts /></el-icon> 环境数据趋势（24h）</h4>
        <BaseChart :option="chartOption" height="300px" />
      </div>
      <div v-else-if="chartLoading" class="chart-loading">
        <el-skeleton :rows="4" animated />
      </div>
    </div>
  </el-dialog>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import BaseChart from '@/components/chart/BaseChart.vue'
import type { FishingSpot, EnvironmentData } from '@/types'
import { WATER_TYPE_MAP, SPOT_STATUS_MAP } from '@/types'
import type { EChartsOption } from 'echarts'
import { Monitor, TrendCharts } from '@element-plus/icons-vue'

const props = defineProps<{
  visible: boolean
  spot: FishingSpot | null
  environmentData: EnvironmentData[]
  chartLoading: boolean
}>()

defineEmits<{
  (e: 'update:visible', value: boolean): void
}>()

const chartOption = computed<EChartsOption | null>(() => {
  if (!props.spot || !props.environmentData.length) return null

  const sorted = [...props.environmentData].sort(
    (a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
  )
  const times = sorted.map(d =>
    new Date(d.timestamp).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  )

  return {
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(15, 23, 42, 0.9)',
      textStyle: { color: '#fff', fontSize: 12 }
    },
    legend: { data: ['水温', '气温', '湿度'], textStyle: { fontSize: 12 } },
    xAxis: { type: 'category', data: times, axisLabel: { rotate: 45, fontSize: 10 } },
    yAxis: [
      { type: 'value', name: '温度(°C)', axisLabel: { fontSize: 10 } },
      { type: 'value', name: '湿度(%)', position: 'right', axisLabel: { fontSize: 10 } }
    ],
    series: [
      {
        name: '水温', type: 'line', data: sorted.map(d => d.water_temp), smooth: true,
        lineStyle: { width: 2.5 }, itemStyle: { color: '#409eff' },
        areaStyle: {
          color: {
            type: 'linear', x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(64,158,255,0.2)' },
              { offset: 1, color: 'rgba(64,158,255,0.02)' }
            ]
          } as any
        }
      },
      {
        name: '气温', type: 'line', data: sorted.map(d => d.air_temp), smooth: true,
        lineStyle: { width: 2.5 }, itemStyle: { color: '#e6a23c' }
      },
      {
        name: '湿度', type: 'line', yAxisIndex: 1, data: sorted.map(d => d.humidity), smooth: true,
        lineStyle: { width: 2, type: 'dashed' as const }, itemStyle: { color: '#67c23a' }
      }
    ],
    grid: { left: '3%', right: '4%', bottom: '15%', top: '15%', containLabel: true }
  }
})
</script>

<style scoped>
.spot-detail {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.device-info h4, .detail-chart h4 {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 10px;
  font-size: 15px;
  color: #303133;
}

.chart-loading { padding: 20px; }
</style>
