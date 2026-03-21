<template>
  <div class="panel env-radar">
    <div class="panel-title">{{ panelTitle }}</div>
    <div ref="chartRef" class="chart-container"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import echarts from '@/plugins/echarts'
import type { RegionEnvItem } from '@/types'

const props = withDefaults(defineProps<{
  envItems: RegionEnvItem[]
  /** 地图层级: 0=全国, 1=省级, 2=市级 */
  mapLevel?: number
  /** 当前地图区域名称 */
  mapAreaName?: string
}>(), {
  mapLevel: 0,
  mapAreaName: ''
})

const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

// 指标元信息（不含 max，max 由数据动态计算）
const INDICATOR_META = [
  { key: 'water_temp', name: '水温', unit: '℃', minRange: 5 },
  { key: 'air_temp', name: '气温', unit: '℃', minRange: 5 },
  { key: 'humidity', name: '湿度', unit: '%', minRange: 10 },
  { key: 'ph', name: 'pH值', unit: '', minRange: 1 },
  { key: 'dissolved_oxygen', name: '溶氧', unit: 'mg/L', minRange: 2 },
  { key: 'turbidity', name: '浊度', unit: 'NTU', minRange: 5 }
] as const

const COLORS = [
  '#00d4ff', '#00ff88', '#ffd93d', '#a78bfa', '#ff6b6b', '#ff9f43'
]

/** 面板标题根据地图层级动态变化 */
const panelTitle = computed(() => {
  if (props.mapLevel === 0) return '全国环境概览'
  if (props.mapLevel === 1) return `${props.mapAreaName || '省级'}环境概览`
  return `${props.mapAreaName || '市级'}环境详情`
})

/**
 * 根据实际数据动态计算各维度最大值
 * 取数据中最大值 + 30% padding，并 ceil 到整数，确保雷达图饱满
 */
function computeDynamicMax(items: RegionEnvItem[]) {
  return INDICATOR_META.map(meta => {
    const key = meta.key as keyof RegionEnvItem
    const values = items.map(it => Number(it[key]) || 0)
    const dataMax = Math.max(...values, 0)
    // 最大值加 30% padding，但不低于 minRange（避免数据为零时 max=0）
    const padded = dataMax * 1.3
    const result = Math.max(padded, meta.minRange)
    // ceil 到好看的整数
    return Math.ceil(result)
  })
}

function renderChart() {
  if (!chartRef.value) return
  if (!chart) chart = echarts.init(chartRef.value)

  if (!props.envItems.length) {
    chart.setOption({
      graphic: {
        type: 'text',
        left: 'center', top: 'middle',
        style: { text: '暂无数据', fill: '#5a7a9a', fontSize: 14 }
      },
      radar: { show: false }, series: []
    }, true)
    return
  }

  // 最多取前5个区域，避免雷达图过于拥挤
  const items = props.envItems.slice(0, 5)
  const dynamicMax = computeDynamicMax(items)

  const seriesData = items.map((item, i) => ({
    name: item.region_name,
    value: [
      item.water_temp,
      item.air_temp,
      item.humidity,
      item.ph,
      item.dissolved_oxygen,
      item.turbidity
    ],
    lineStyle: { color: COLORS[i % COLORS.length], width: 1.5 },
    areaStyle: { color: COLORS[i % COLORS.length], opacity: 0.08 },
    itemStyle: { color: COLORS[i % COLORS.length] },
    symbol: 'circle',
    symbolSize: 4
  }))

  chart.setOption({
    graphic: [],
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(6, 30, 65, 0.95)',
      borderColor: '#1e90ff',
      textStyle: { color: '#e8f0ff', fontSize: 11 },
      formatter: (params: any) => {
        if (!params.value) return ''
        const name = params.name
        const lines = INDICATOR_META.map((ind, i) =>
          `${ind.name}：<b>${params.value[i]?.toFixed(1) ?? '-'}</b> ${ind.unit}`
        )
        return `<b style="color:#00d4ff">${name}</b><br/>` + lines.join('<br/>')
      }
    },
    legend: {
      bottom: 2,
      left: 'center',
      textStyle: { color: '#8ba3c7', fontSize: 10 },
      itemWidth: 10,
      itemHeight: 6,
      itemGap: 10
    },
    radar: {
      center: ['50%', '46%'],
      radius: '60%',
      indicator: INDICATOR_META.map((ind, i) => ({
        name: ind.name,
        max: dynamicMax[i]
      })),
      axisName: {
        color: '#5a7a9a',
        fontSize: 10
      },
      splitArea: {
        areaStyle: {
          color: [
            'rgba(0, 212, 255, 0.02)',
            'rgba(0, 212, 255, 0.04)',
            'rgba(0, 212, 255, 0.06)',
            'rgba(0, 212, 255, 0.08)',
            'rgba(0, 212, 255, 0.10)'
          ]
        }
      },
      splitLine: {
        lineStyle: { color: '#1a2a3a' }
      },
      axisLine: {
        lineStyle: { color: '#1a2a3a' }
      }
    },
    series: [{
      type: 'radar',
      data: seriesData
    }]
  }, true)
}

watch(() => props.envItems, renderChart, { deep: true })

onMounted(() => {
  renderChart()
  window.addEventListener('resize', () => chart?.resize())
})

onUnmounted(() => chart?.dispose())
</script>

<style scoped>
.env-radar {
  flex: 1;
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
