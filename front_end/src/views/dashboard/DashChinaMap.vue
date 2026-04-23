<template>
  <div class="panel china-map">
    <div class="panel-title">
      <span>{{ mapTitle }}</span>
      <span class="breadcrumb">
        <span v-if="drillStack.length" class="back-btn" @click="goBack">← 返回上级</span>
        <span v-if="drillStack.length" class="back-btn back-all" @click="backToChina">⌂ 全国</span>
      </span>
    </div>
    <div ref="chartRef" class="chart-container"></div>

    <!-- 点击弹出浮窗 -->
    <Transition name="popup-fade">
      <div
        v-if="popupData"
        class="spot-popup"
        :style="{ left: popupPos.x + 'px', top: popupPos.y + 'px' }"
      >
        <div class="popup-close" @click="closePopup">×</div>
        <!-- 单个水域 -->
        <template v-if="!popupData.isCluster">
          <div class="popup-header">
            <span class="popup-dot" :style="{ background: spotColor(popupData.spots[0]) }"></span>
            <span class="popup-name">{{ popupData.spots[0].name }}</span>
          </div>
          <div class="popup-body">
            <div class="popup-row">
              <span class="popup-label">状态</span>
              <span class="popup-val" :style="{ color: spotColor(popupData.spots[0]) }">{{ SPOT_STATUS_MAP[popupData.spots[0].status] }}</span>
            </div>
            <div class="popup-row">
              <span class="popup-label">类型</span>
              <span class="popup-val">{{ WATER_TYPE_MAP[popupData.spots[0].water_type] }}</span>
            </div>
            <div class="popup-row">
              <span class="popup-label">容量</span>
              <span class="popup-val">{{ popupData.spots[0].capacity }} 人</span>
            </div>
            <div class="popup-row">
              <span class="popup-label">坐标</span>
              <span class="popup-val coord">{{ popupData.spots[0].latitude.toFixed(4) }}, {{ popupData.spots[0].longitude.toFixed(4) }}</span>
            </div>
          </div>
        </template>
        <!-- 聚合点 -->
        <template v-else>
          <div class="popup-header cluster">
            <span class="popup-cluster-badge">{{ popupData.count }}</span>
            <span class="popup-name">个水域聚合</span>
          </div>
          <div class="popup-hint">双击地图可下钻查看详情</div>
          <div class="popup-list">
            <div v-for="s in popupData.spots" :key="s.id" class="popup-spot-item" @click="$emit('select', s)">
              <span class="popup-dot" :style="{ background: spotColor(s) }"></span>
              <span class="popup-spot-name">{{ s.name }}</span>
              <span class="popup-spot-status" :style="{ color: spotColor(s) }">{{ SPOT_STATUS_MAP[s.status] }}</span>
            </div>
          </div>
        </template>
      </div>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted, onUnmounted } from 'vue'
import echarts from '@/plugins/echarts'
import type { ECharts } from '@/plugins/echarts'
import type { FishingSpot, Region } from '@/types'
import { WATER_TYPE_MAP, SPOT_STATUS_MAP } from '@/types'

const props = defineProps<{
  spots: FishingSpot[]
  regions: Region[]
  selectedSpotId: number | null
}>()
const emit = defineEmits<{
  select: [spot: FishingSpot]
  'update:mapLevel': [level: number]
  'update:mapAreaName': [name: string]
}>()

const chartRef = ref<HTMLElement>()
let chart: ECharts | null = null

// ===== 弹窗状态 =====
const popupData = ref<ClusterPoint | null>(null)
const popupPos = ref({ x: 0, y: 0 })

function closePopup() {
  popupData.value = null
}

function spotColor(s: FishingSpot): string {
  return s.status === 'open' ? '#00ff88' : s.status === 'closed' ? '#ff6b6b' : '#ffd93d'
}

// ===== 下钻状态 =====
interface DrillLevel {
  name: string
  adcode: string
  type: 'province' | 'city'
}
const drillStack = ref<DrillLevel[]>([])

const mapTitle = computed(() => {
  if (drillStack.value.length === 0) return '全国垂钓水域分布'
  if (drillStack.value.length === 1) return `${drillStack.value[0].name} — 地级市分布`
  return `${drillStack.value[1].name} — 区县分布`
})

// ===== GeoJSON 缓存 =====
const geoCache = new Map<string, any>()

async function loadGeoJSON(adcode: string): Promise<any> {
  if (geoCache.has(adcode)) return geoCache.get(adcode)
  const url = adcode === '100000'
    ? 'https://geo.datav.aliyun.com/areas_v3/bound/100000_full.json'
    : `https://geo.datav.aliyun.com/areas_v3/bound/${adcode}_full.json`
  const resp = await fetch(url, {
    method: 'GET',
    credentials: 'omit',
    referrerPolicy: 'no-referrer'
  })
  const data = await resp.json()
  geoCache.set(adcode, data)
  return data
}

// 省份名 → adcode
const PROVINCE_ADCODE: Record<string, string> = {
  '黑龙江省': '230000', '吉林省': '220000', '辽宁省': '210000',
  '北京市': '110000', '天津市': '120000', '河北省': '130000',
  '山西省': '140000', '内蒙古自治区': '150000', '上海市': '310000',
  '江苏省': '320000', '浙江省': '330000', '安徽省': '340000',
  '福建省': '350000', '江西省': '360000', '山东省': '370000',
  '河南省': '410000', '湖北省': '420000', '湖南省': '430000',
  '广东省': '440000', '广西壮族自治区': '450000', '海南省': '460000',
  '重庆市': '500000', '四川省': '510000', '贵州省': '520000',
  '云南省': '530000', '西藏自治区': '540000', '陕西省': '610000',
  '甘肃省': '620000', '青海省': '630000', '宁夏回族自治区': '640000',
  '新疆维吾尔自治区': '650000', '台湾省': '710000',
  '香港特别行政区': '810000', '澳门特别行政区': '820000',
}

// ===== 数据点聚合算法 =====
interface ClusterPoint {
  name: string
  lng: number
  lat: number
  count: number
  capacity: number
  spots: FishingSpot[]
  isCluster: boolean
}

function clusterSpots(spots: FishingSpot[], threshold: number): ClusterPoint[] {
  if (!spots.length) return []

  const clusters: ClusterPoint[] = []
  const used = new Set<number>()

  for (let i = 0; i < spots.length; i++) {
    if (used.has(i)) continue

    const s = spots[i]
    const group: FishingSpot[] = [s]
    used.add(i)

    for (let j = i + 1; j < spots.length; j++) {
      if (used.has(j)) continue
      const s2 = spots[j]
      const dist = Math.sqrt(
        Math.pow(s.longitude - s2.longitude, 2) +
        Math.pow(s.latitude - s2.latitude, 2)
      )
      if (dist < threshold) {
        group.push(s2)
        used.add(j)
      }
    }

    const avgLng = group.reduce((sum, g) => sum + g.longitude, 0) / group.length
    const avgLat = group.reduce((sum, g) => sum + g.latitude, 0) / group.length
    const totalCap = group.reduce((sum, g) => sum + g.capacity, 0)

    clusters.push({
      name: group.length === 1 ? s.name : `${group[0].name} 等 ${group.length} 个水域`,
      lng: avgLng,
      lat: avgLat,
      count: group.length,
      capacity: totalCap,
      spots: group,
      isCluster: group.length > 1
    })
  }
  return clusters
}

// ===== 辅助 =====
function getVisibleSpots(): FishingSpot[] {
  if (drillStack.value.length === 0) return props.spots

  if (drillStack.value.length >= 1) {
    const provinceName = drillStack.value[0].name
    const filtered = props.spots.filter(s => {
      const r = props.regions.find(r => r.id === s.region_id)
      if (!r) return false
      return r.province === provinceName ||
        r.province + '省' === provinceName ||
        r.province + '市' === provinceName ||
        provinceName.startsWith(r.province)
    })

    if (drillStack.value.length >= 2) {
      const cityName = drillStack.value[1].name
      return filtered.filter(s => {
        const r = props.regions.find(r => r.id === s.region_id)
        if (!r) return false
        return r.city === cityName ||
          r.city + '市' === cityName ||
          cityName.startsWith(r.city) ||
          r.city.startsWith(cityName.replace(/市$/, ''))
      })
    }
    return filtered
  }
  return props.spots
}

function getProvinceData() {
  const map: Record<string, { count: number }> = {}
  for (const s of props.spots) {
    const region = props.regions.find(r => r.id === s.region_id)
    const province = region?.province || '未知'
    const pName = province.endsWith('省') || province.endsWith('市') || province.includes('自治')
      ? province
      : province + '省'
    if (!map[pName]) map[pName] = { count: 0 }
    map[pName].count++
  }
  return map
}

function getScatterData(thresholdDeg: number) {
  const visible = getVisibleSpots()
  const clusters = clusterSpots(visible, thresholdDeg)

  return clusters.map(c => ({
    name: c.name,
    value: [c.lng, c.lat, c.capacity],
    clusterData: c,
    spotData: c.spots.length === 1 ? c.spots[0] : null,
    itemStyle: {
      color: c.isCluster
        ? '#ffa940'
        : c.spots[0]?.id === props.selectedSpotId
          ? '#ffd93d'
          : c.spots[0]?.status === 'open' ? '#00ff88'
          : c.spots[0]?.status === 'closed' ? '#ff6b6b'
          : '#ffd93d'
    }
  }))
}

function getRippleData(thresholdDeg: number) {
  const visible = getVisibleSpots()
  const openSpots = visible.filter(s => s.status === 'open')
  const clusters = clusterSpots(openSpots, thresholdDeg)

  return clusters.map(c => ({
    name: c.name,
    value: [c.lng, c.lat, c.count > 1 ? 12 : 8],
    clusterData: c,
    spotData: c.spots.length === 1 ? c.spots[0] : null
  }))
}

// ===== Tooltip：只在 hover 时显示简短信息 =====
function tooltipFormatter(params: any) {
  const cluster: ClusterPoint | undefined = params.data?.clusterData
  if (!cluster) return ''

  if (cluster.isCluster) {
    return `<b style="color:#ffa940">${cluster.count} 个水域</b><br/><span style="color:#8ba3c7;font-size:11px">点击查看详情</span>`
  }

  const s: FishingSpot = cluster.spots[0]
  return `<b style="color:#00d4ff">${s.name}</b><br/><span style="color:#8ba3c7;font-size:11px">点击查看详情</span>`
}

function findAdcodeFromGeoJSON(geoJSON: any, areaName: string): string | null {
  const feature = geoJSON.features?.find((f: any) => {
    const name = f.properties?.name || ''
    return name === areaName || name.includes(areaName) || areaName.includes(name.replace(/[市区县]$/, ''))
  })
  return feature?.properties?.adcode?.toString() || null
}

// ===== 渲染地图 =====
async function renderMap() {
  if (!chartRef.value) return
  if (!chart) {
    chart = echarts.init(chartRef.value)
    bindEvents()
  }
  // 切换层级时关闭弹窗
  closePopup()

  const level = drillStack.value.length

  // 通知父组件当前地图层级和区域名称
  emit('update:mapLevel', level)
  const areaName = level === 0
    ? '全国'
    : level === 1
      ? drillStack.value[0].name
      : drillStack.value[1].name
  emit('update:mapAreaName', areaName)

  if (level === 0) {
    // 全国视图：中心偏北偏左，放大到 1.5 倍，裁掉南海大部分
    await renderLevel('china', '100000', 1.55, [104.5, 37.5])
  } else if (level === 1) {
    const { name, adcode } = drillStack.value[0]
    await renderLevel(name, adcode, 1.0, undefined)
  } else {
    const { name, adcode } = drillStack.value[1]
    await renderLevel(name, adcode, 1.0, undefined)
  }
}

/**
 * 根据当前地图实际缩放级别动态计算聚合阈值（度）
 * 缩放级别越大（视野越小），阈值越小，聚合越不容易发生
 */
function getClusterThreshold(): number {
  if (!chart) return 0.8
  // 利用地图坐标系：计算当前视口宽度对应多少度
  // 取两个屏幕点（左 1/4 和右 1/4）转为地理坐标，计算度数差
  try {
    const rect = chartRef.value?.getBoundingClientRect()
    if (!rect) return 0.3
    const w = rect.width
    const ptLeft = chart.convertFromPixel('geo', [w * 0.25, rect.height * 0.5])
    const ptRight = chart.convertFromPixel('geo', [w * 0.75, rect.height * 0.5])
    if (!ptLeft || !ptRight) return 0.3
    const viewDegrees = Math.abs(ptRight[0] - ptLeft[0]) // 视口一半宽度对应的度数
    // 聚合阈值为视口宽度的 ~3%，更保守以避免意外合并
    return Math.max(viewDegrees * 0.04, 0.003)
  } catch {
    const level = drillStack.value.length
    return level === 0 ? 0.8 : level === 1 ? 0.2 : 0.04
  }
}

async function renderLevel(mapName: string, adcode: string, zoom: number, center?: [number, number]) {
  if (!chart) return

  const geoJSON = await loadGeoJSON(adcode)
  echarts.registerMap(mapName, geoJSON)

  const level = drillStack.value.length
  // 初始渲染时使用基于层级的默认阈值（保守阈值，避免过度聚合）
  const threshold = level === 0 ? 0.8 : level === 1 ? 0.2 : 0.04

  const provinceData = level === 0 ? getProvinceData() : {}
  const regionColors = level > 0 ? getAreaColors(geoJSON) : []

  chart.setOption({
    geo: {
      map: mapName,
      roam: true,
      zoom,
      center: center || undefined,
      // 全国视图使用 boundingCoords 裁掉南海
      ...(level === 0 ? {
        layoutCenter: ['50%', '50%'],
        layoutSize: '100%'
      } : {}),
      itemStyle: {
        areaColor: '#0e2341',
        borderColor: '#1e90ff',
        borderWidth: 0.8
      },
      emphasis: {
        itemStyle: {
          areaColor: '#1a3a5c',
          borderColor: '#00d4ff',
          borderWidth: 1.5
        },
        label: {
          show: true,
          color: '#00d4ff',
          fontSize: level === 0 ? 12 : 13,
          fontWeight: level > 0 ? 'bold' : 'normal'
        }
      },
      select: {
        itemStyle: { areaColor: '#1a4a6c' }
      },
      regions: level === 0
        ? Object.entries(provinceData).map(([name, data]) => ({
            name,
            itemStyle: {
              areaColor: data.count > 5
                ? 'rgba(0, 212, 255, 0.35)'
                : data.count > 0
                ? 'rgba(0, 212, 255, 0.18)'
                : '#0e2341'
            }
          }))
        : regionColors,
      label: {
        show: level > 0,
        color: '#5a7a9a',
        fontSize: 11
      }
    },
    tooltip: {
      trigger: 'item',
      backgroundColor: 'rgba(6, 30, 65, 0.95)',
      borderColor: '#1e90ff',
      textStyle: { color: '#e8f0ff', fontSize: 12 },
      formatter: tooltipFormatter,
      extraCssText: 'max-width: 280px;'
    },
    series: [
      {
        type: 'effectScatter',
        coordinateSystem: 'geo',
        data: getRippleData(threshold),
        symbolSize: (val: number[]) => Math.max(6, (val[2] || 6)),
        rippleEffect: { brushType: 'stroke', scale: 3.5, period: 4 },
        itemStyle: { color: '#00ff88', shadowBlur: 10, shadowColor: '#00ff88' },
        zlevel: 1,
        // 不显示标签
        label: { show: false }
      },
      {
        type: 'scatter',
        coordinateSystem: 'geo',
        data: getScatterData(threshold),
        symbolSize: (val: number[], params: any) => {
          const cluster = params?.data?.clusterData
          if (cluster?.isCluster) {
            return Math.max(14, Math.min(30, cluster.count * 7))
          }
          return Math.max(8, Math.min(20, (val[2] || 50) / 6))
        },
        zlevel: 2,
        // 不显示标签，只显示纯圆点
        label: {
          show: false
        }
      }
    ]
  }, true)
}

// ===== 区域着色 =====
function getAreaColors(geoJSON: any) {
  const visible = getVisibleSpots()
  const areaCount: Record<string, number> = {}
  for (const s of visible) {
    const region = props.regions.find(r => r.id === s.region_id)
    const cityOrDistrict = region?.city || '未知'
    const feature = geoJSON.features?.find((f: any) => {
      const name = f.properties?.name || ''
      return name.includes(cityOrDistrict) ||
        cityOrDistrict.includes(name.replace(/[市区县]$/, ''))
    })
    const matchName = feature?.properties?.name || cityOrDistrict
    areaCount[matchName] = (areaCount[matchName] || 0) + 1
  }

  return Object.entries(areaCount).map(([name, count]) => ({
    name,
    itemStyle: {
      areaColor: count > 3
        ? 'rgba(0, 255, 136, 0.3)'
        : count > 0
        ? 'rgba(0, 212, 255, 0.2)'
        : '#0e2341'
    }
  }))
}

// ===== 事件绑定 =====
let zoomDebounceTimer: ReturnType<typeof setTimeout> | null = null

function bindEvents() {
  if (!chart) return

  // 仅在缩放时重新计算聚合（防抖 300ms），拖动平移不触发
  chart.on('georoam', (params: any) => {
    // params.zoom 存在时表示缩放操作，否则是平移
    if (!params.zoom) return
    if (zoomDebounceTimer) clearTimeout(zoomDebounceTimer)
    zoomDebounceTimer = setTimeout(() => {
      refreshScatterOnly()
    }, 300)
  })

  // 双击区域 → 下钻
  chart.on('dblclick', (params: any) => {
    void (async () => {
      if (params.componentType === 'geo' && params.name) {
        await handleDrillDown(params.name)
      }
      if ((params.seriesType === 'scatter' || params.seriesType === 'effectScatter') && params.data?.clusterData?.isCluster) {
        const cluster: ClusterPoint = params.data.clusterData
        const firstSpot = cluster.spots[0]
        const region = props.regions.find(r => r.id === firstSpot.region_id)
        if (region) {
          if (drillStack.value.length === 0) {
            const pName = region.province.endsWith('省') || region.province.endsWith('市') || region.province.includes('自治')
              ? region.province : region.province + '省'
            await handleDrillDown(pName)
          } else if (drillStack.value.length === 1) {
            await handleDrillDown(region.city)
          }
        }
      }
    })()
  })

  // 单击散点 → 弹出浮窗
  chart.on('click', (params: any) => {
    if (params.seriesType === 'scatter' || params.seriesType === 'effectScatter') {
      const cluster: ClusterPoint | undefined = params.data?.clusterData
      if (!cluster) return

      // 计算弹窗位置（基于鼠标事件坐标，限制在容器内）
      const containerRect = chartRef.value?.getBoundingClientRect()
      if (!containerRect) return

      const mouseX = (params.event?.offsetX ?? params.event?.event?.offsetX ?? 0)
      const mouseY = (params.event?.offsetY ?? params.event?.event?.offsetY ?? 0)

      // 弹窗宽约 220px，高约 180px，做边界保护
      const maxX = containerRect.width - 240
      const maxY = containerRect.height - 200

      popupPos.value = {
        x: Math.min(Math.max(mouseX + 12, 0), maxX),
        y: Math.min(Math.max(mouseY - 60, 0), maxY)
      }
      popupData.value = cluster

      // 同时 emit select 事件
      if (cluster.spots.length === 1) {
        emit('select', cluster.spots[0])
      }
    } else {
      // 点击地图其他区域关闭弹窗
      closePopup()
    }
  })
}

// ===== 下钻 =====
async function handleDrillDown(areaName: string) {
  const level = drillStack.value.length

  if (level === 0) {
    const adcode = PROVINCE_ADCODE[areaName]
    if (adcode) {
      drillStack.value.push({ name: areaName, adcode, type: 'province' })
      await renderMap()
    }
  } else if (level === 1) {
    const provinceAdcode = drillStack.value[0].adcode
    const provinceGeo = await loadGeoJSON(provinceAdcode)
    const cityAdcode = findAdcodeFromGeoJSON(provinceGeo, areaName)
    if (cityAdcode) {
      try {
        await loadGeoJSON(cityAdcode)
        drillStack.value.push({ name: areaName, adcode: cityAdcode, type: 'city' })
        await renderMap()
      } catch {
        console.warn(`无法加载 ${areaName}(${cityAdcode}) 的区县地图`)
      }
    }
  }
}

function goBack() {
  drillStack.value.pop()
  renderMap()
}

function backToChina() {
  drillStack.value = []
  renderMap()
}

/**
 * 仅更新散点图数据（缩放后重新聚合），不重新渲染底图
 */
function refreshScatterOnly() {
  if (!chart) return
  const threshold = getClusterThreshold()
  chart.setOption({
    series: [
      {
        type: 'effectScatter',
        data: getRippleData(threshold)
      },
      {
        type: 'scatter',
        data: getScatterData(threshold)
      }
    ]
  })
}

// ===== 监听 =====
// 数据源变化（水域列表、区域列表）时需要重新渲染整个地图
watch([() => props.spots, () => props.regions], () => {
  renderMap()
}, { deep: true })

// 选中水域变化时，只更新散点颜色，不重置地图缩放和位移
watch(() => props.selectedSpotId, () => {
  refreshScatterOnly()
})

onMounted(() => {
  renderMap()
  window.addEventListener('resize', () => chart?.resize())
})

onUnmounted(() => {
  chart?.dispose()
})
</script>

<style scoped>
.china-map {
  flex: 1;
  padding: 10px;
  display: flex;
  flex-direction: column;
  position: relative;
}

.panel-title {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.breadcrumb {
  display: flex;
  gap: 6px;
}

.back-btn {
  font-size: 12px;
  color: #00d4ff;
  cursor: pointer;
  padding: 2px 8px;
  border: 1px solid rgba(0, 212, 255, 0.3);
  border-radius: 4px;
  transition: all 0.3s;
}

.back-btn:hover {
  background: rgba(0, 212, 255, 0.15);
  border-color: #00d4ff;
}

.back-all {
  color: #8ba3c7;
  border-color: rgba(139, 163, 199, 0.3);
}

.back-all:hover {
  color: #00d4ff;
  border-color: #00d4ff;
}

.chart-container {
  flex: 1;
  width: 100%;
  min-height: 300px;
}

/* ===== 弹出浮窗 ===== */
.spot-popup {
  position: absolute;
  z-index: 100;
  width: 220px;
  background: rgba(6, 24, 54, 0.96);
  border: 1px solid rgba(0, 212, 255, 0.35);
  border-radius: 10px;
  padding: 14px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.5), 0 0 20px rgba(0, 212, 255, 0.08);
  backdrop-filter: blur(12px);
}

.popup-close {
  position: absolute;
  top: 6px;
  right: 10px;
  font-size: 16px;
  color: #5a7a9a;
  cursor: pointer;
  line-height: 1;
  transition: color 0.2s;
}

.popup-close:hover {
  color: #ff6b6b;
}

.popup-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.popup-header.cluster {
  gap: 4px;
}

.popup-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  flex-shrink: 0;
}

.popup-name {
  font-size: 14px;
  font-weight: 600;
  color: #e8f0ff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.popup-cluster-badge {
  font-size: 18px;
  font-weight: 700;
  color: #ffa940;
  font-family: 'DIN Alternate', 'Courier New', monospace;
}

.popup-body {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.popup-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.popup-label {
  font-size: 11px;
  color: #5a7a9a;
}

.popup-val {
  font-size: 12px;
  color: #b8cce0;
  font-weight: 500;
}

.popup-val.coord {
  font-size: 10px;
  font-family: 'Courier New', monospace;
  color: #5a7a9a;
}

.popup-hint {
  font-size: 10px;
  color: #5a7a9a;
  margin-bottom: 8px;
}

.popup-list {
  display: flex;
  flex-direction: column;
  gap: 5px;
  max-height: 150px;
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: #2a4a6b transparent;
}

.popup-spot-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 6px;
  border-radius: 4px;
  background: rgba(255, 255, 255, 0.03);
  cursor: pointer;
  transition: background 0.2s;
}

.popup-spot-item:hover {
  background: rgba(0, 212, 255, 0.1);
}

.popup-spot-name {
  flex: 1;
  font-size: 11px;
  color: #b8cce0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.popup-spot-status {
  font-size: 10px;
  flex-shrink: 0;
}

/* 浮窗动画 */
.popup-fade-enter-active {
  transition: all 0.25s ease-out;
}
.popup-fade-leave-active {
  transition: all 0.15s ease-in;
}
.popup-fade-enter-from {
  opacity: 0;
  transform: scale(0.9) translateY(8px);
}
.popup-fade-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
