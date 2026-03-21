<template>
  <div class="panel env-trend">
    <div class="panel-title">
      {{ spotName ? `${spotName} — 48h 环境趋势` : '环境趋势（请选择水域）' }}
    </div>
    <div ref="chartRef" class="chart-container"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import echarts from '@/plugins/echarts'
import type { EnvironmentData } from '@/types'
import { FishingSpotService } from '@/services'

const props = defineProps<{
  spotId: number | null
  spotName: string
}>()

const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null
const envData = ref<EnvironmentData[]>([])

async function loadData() {
  if (!props.spotId) return
  try {
    envData.value = await FishingSpotService.getEnvironment(props.spotId, 48)
  } catch {
    envData.value = []
  }
}

function renderChart() {
  if (!chartRef.value) return
  if (!chart) chart = echarts.init(chartRef.value)

  if (!envData.value.length) {
    chart.setOption({
      graphic: {
        type: 'text',
        left: 'center', top: 'middle',
        style: { text: '暂无数据', fill: '#5a7a9a', fontSize: 14 }
      },
      xAxis: { show: false }, yAxis: { show: false }, series: []
    }, true)
    return
  }

  const times = envData.value.map(d =>
    new Date(d.timestamp).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
  )

  chart.setOption({
    graphic: [],
    tooltip: {
      trigger: 'axis',
      backgroundColor: 'rgba(6, 30, 65, 0.95)',
      borderColor: '#1e90ff',
      textStyle: { color: '#e8f0ff', fontSize: 11 }
    },
    legend: {
      top: 2,
      right: 10,
      textStyle: { color: '#8ba3c7', fontSize: 10 },
      itemWidth: 12,
      itemHeight: 8
    },
    grid: { left: 40, right: 40, top: 30, bottom: 20, containLabel: false },
    xAxis: {
      type: 'category',
      data: times,
      axisLine: { lineStyle: { color: '#2a4a6b' } },
      axisLabel: { color: '#5a7a9a', fontSize: 9, interval: Math.floor(envData.value.length / 8) },
      axisTick: { show: false }
    },
    yAxis: [
      {
        type: 'value',
        name: '温度(℃)',
        nameTextStyle: { color: '#5a7a9a', fontSize: 10 },
        axisLine: { show: false },
        axisLabel: { color: '#5a7a9a', fontSize: 10 },
        splitLine: { lineStyle: { color: '#1a2a3a' } }
      },
      {
        type: 'value',
        name: '湿度(%)',
        nameTextStyle: { color: '#5a7a9a', fontSize: 10 },
        axisLine: { show: false },
        axisLabel: { color: '#5a7a9a', fontSize: 10 },
        splitLine: { show: false }
      }
    ],
    series: [
      {
        name: '水温',
        type: 'line',
        data: envData.value.map(d => d.water_temp),
        smooth: true,
        symbol: 'none',
        lineStyle: { color: '#00d4ff', width: 2 },
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(0, 212, 255, 0.25)' },
            { offset: 1, color: 'rgba(0, 212, 255, 0)' }
          ])
        }
      },
      {
        name: '气温',
        type: 'line',
        data: envData.value.map(d => d.air_temp),
        smooth: true,
        symbol: 'none',
        lineStyle: { color: '#ffd93d', width: 1.5 }
      },
      {
        name: '湿度',
        type: 'line',
        yAxisIndex: 1,
        data: envData.value.map(d => d.humidity),
        smooth: true,
        symbol: 'none',
        lineStyle: { color: '#a78bfa', width: 1.5, type: 'dashed' }
      }
    ]
  }, true)
}

watch(() => props.spotId, async () => {
  await loadData()
  renderChart()
})

watch(envData, renderChart, { deep: true })

onMounted(async () => {
  await loadData()
  renderChart()
  window.addEventListener('resize', () => chart?.resize())
})

onUnmounted(() => chart?.dispose())
</script>

<style scoped>
.env-trend {
  flex: 1.5;
  padding: 10px;
  display: flex;
  flex-direction: column;
}
.chart-container {
  flex: 1;
  width: 100%;
  min-height: 0;
}
</style>
