<template>
  <div class="panel popular-rank">
    <div class="panel-title">热门水域 TOP {{ spots.length }}</div>
    <div ref="chartRef" class="chart-container"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import echarts from '@/plugins/echarts'
import type { PopularSpot } from '@/types'
import { WATER_TYPE_MAP } from '@/types'

const props = defineProps<{ spots: PopularSpot[] }>()
const emit = defineEmits<{ select: [spot: PopularSpot] }>()

const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

function renderChart() {
  if (!chartRef.value || !props.spots.length) return
  if (!chart) {
    chart = echarts.init(chartRef.value)
    chart.on('click', (params: any) => {
      if (params.dataIndex !== undefined) {
        emit('select', props.spots[props.spots.length - 1 - params.dataIndex])
      }
    })
  }

  const sorted = [...props.spots].sort((a, b) => a.total_fishing_count - b.total_fishing_count)
  const maxVal = Math.max(...sorted.map(s => s.total_fishing_count))

  chart.setOption({
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' },
      formatter: (params: any) => {
        const s = sorted[params[0].dataIndex]
        return `<div style="font-size:13px"><b>${s.name}</b><br/>
          类型：${WATER_TYPE_MAP[s.water_type] || s.water_type}<br/>
          垂钓总数：<b style="color:#00ff88">${s.total_fishing_count}</b> 人次<br/>
          容量：${s.capacity} 人</div>`
      }
    },
    grid: { left: 10, right: 50, top: 4, bottom: 4, containLabel: true },
    xAxis: {
      type: 'value',
      max: maxVal * 1.2,
      show: false
    },
    yAxis: {
      type: 'category',
      data: sorted.map(s => s.name.length > 8 ? s.name.slice(0, 8) + '…' : s.name),
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: '#8ba3c7', fontSize: 11 }
    },
    series: [{
      type: 'bar',
      data: sorted.map((s, i) => ({
        value: s.total_fishing_count,
        itemStyle: {
          borderRadius: [0, 4, 4, 0],
          color: new echarts.graphic.LinearGradient(0, 0, 1, 0, [
            { offset: 0, color: 'rgba(0, 212, 255, 0.3)' },
            { offset: 1, color: i === sorted.length - 1 ? '#00ff88' : '#00d4ff' }
          ])
        }
      })),
      barWidth: 16,
      label: {
        show: true,
        position: 'right',
        color: '#e8f0ff',
        fontSize: 12,
        fontWeight: 'bold',
        fontFamily: 'DIN Alternate, Courier New, monospace'
      }
    }]
  }, true)
}

watch(() => props.spots, renderChart, { deep: true })

onMounted(() => {
  renderChart()
  window.addEventListener('resize', () => chart?.resize())
})
onUnmounted(() => {
  chart?.dispose()
  window.removeEventListener('resize', () => chart?.resize())
})
</script>

<style scoped>
.popular-rank {
  padding: 10px;
  flex: 1;
}

.chart-container {
  width: 100%;
  height: calc(100% - 28px);
  min-height: 140px;
}
</style>
