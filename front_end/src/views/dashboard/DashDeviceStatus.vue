<template>
  <div class="panel mini-chart">
    <div class="panel-title">设备状态</div>
    <div ref="chartRef" class="chart-container"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import echarts from '@/plugins/echarts'
import type { ECharts } from '@/plugins/echarts'
import type { Device } from '@/types'

const props = defineProps<{ devices: Device[] }>()
const chartRef = ref<HTMLElement>()
let chart: ECharts | null = null

const STATUS_COLORS: Record<string, string> = {
  online: '#00ff88',
  offline: '#8ba3c7',
  error: '#ff6b6b'
}

const STATUS_LABELS: Record<string, string> = {
  online: '在线',
  offline: '离线',
  error: '异常'
}

function renderChart() {
  if (!chartRef.value || !props.devices.length) return
  if (!chart) chart = echarts.init(chartRef.value)

  const counts: Record<string, number> = {}
  props.devices.forEach(d => { counts[d.status] = (counts[d.status] || 0) + 1 })

  const data = Object.entries(counts).map(([k, v]) => ({
    name: STATUS_LABELS[k] || k,
    value: v,
    itemStyle: { color: STATUS_COLORS[k] || '#8ba3c7' }
  }))

  chart.setOption({
    tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
    series: [{
      type: 'pie',
      radius: ['50%', '72%'],
      center: ['50%', '52%'],
      label: { show: true, color: '#8ba3c7', fontSize: 10, formatter: '{b}\n{c}台' },
      labelLine: { length: 8, length2: 6, lineStyle: { color: '#2a4a6b' } },
      itemStyle: { borderColor: '#0d1b2a', borderWidth: 2 },
      data
    }]
  }, true)
}

watch(() => props.devices, renderChart, { deep: true })
onMounted(() => {
  renderChart()
  window.addEventListener('resize', () => chart?.resize())
})
onUnmounted(() => chart?.dispose())
</script>

<style scoped>
.mini-chart {
  padding: 8px 10px;
  flex: 1;
}
.chart-container {
  width: 100%;
  height: calc(100% - 24px);
  min-height: 100px;
}
</style>
