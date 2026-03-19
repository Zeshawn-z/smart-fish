<template>
  <div ref="chartRef" :style="{ width: width, height: height }"></div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted, onUnmounted, watch } from 'vue'
import * as echarts from 'echarts'

const props = withDefaults(defineProps<{
  option: echarts.EChartsOption
  width?: string
  height?: string
}>(), {
  width: '100%',
  height: '300px'
})

const chartRef = ref<HTMLDivElement>()
let chart: echarts.ECharts | null = null
let resizeObserver: ResizeObserver | null = null

function initChart() {
  if (!chartRef.value) return
  chart = echarts.init(chartRef.value)
  chart.setOption(props.option)
}

function handleResize() {
  chart?.resize()
}

onMounted(() => {
  initChart()
  window.addEventListener('resize', handleResize)

  // 监听容器尺寸变化（解决 dialog 动画导致初始宽度不正确的问题）
  if (chartRef.value) {
    resizeObserver = new ResizeObserver(() => {
      chart?.resize()
    })
    resizeObserver.observe(chartRef.value)
  }

  // dialog 打开动画结束后再 resize 一次，确保宽度正确
  nextTick(() => {
    setTimeout(() => chart?.resize(), 350)
  })
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  resizeObserver?.disconnect()
  resizeObserver = null
  chart?.dispose()
})

watch(() => props.option, (newOption) => {
  if (chart) {
    chart.setOption(newOption, true)
  }
  // option 变化时（数据加载完成）也 resize 一次
  nextTick(() => chart?.resize())
}, { deep: true })
</script>
