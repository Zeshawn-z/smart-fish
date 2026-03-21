<template>
  <div class="dashboard">
    <!-- 顶部标题栏 -->
    <DashHeader :summary="summary" />

    <!-- 三栏主体 -->
    <div class="dashboard-body">
      <!-- 左侧栏 -->
      <div class="dash-left">
        <DashStatCards :summary="summary" />
        <DashPopularRank :spots="popularSpots" @select="onSpotSelect" />
        <div class="dash-row">
          <DashDeviceStatus :devices="devices" />
          <DashGatewayStatus :gateways="gateways" />
        </div>
      </div>

      <!-- 中央区域 -->
      <div class="dash-center">
        <DashChinaMap
          :spots="allSpots"
          :regions="regions"
          :selected-spot-id="selectedSpotId"
          @select="onSpotSelect"
          @update:map-level="onMapLevelChange"
          @update:map-area-name="onMapAreaNameChange"
        />
        <div class="dash-center-bottom">
          <DashEnvTrend :spot-id="selectedSpotId" :spot-name="selectedSpotName" />
          <DashEnvRadar :env-items="filteredEnvItems" :map-level="mapLevel" :map-area-name="mapAreaName" />
        </div>
      </div>

      <!-- 右侧栏 -->
      <div class="dash-right">
        <DashSpotDetail :spot="selectedSpot" :regions="regions" />
        <DashAlertList :reminders="reminders" />
        <DashSuggestion :suggestions="suggestions" />
        <DashNoticeList :notices="notices" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import type {
  SystemSummary, FishingSpot, PopularSpot, Device, Gateway,
  Region, Reminder, Notice, FishingSuggestion, RegionEnvItem
} from '@/types'
import {
  SummaryService, FishingSpotService, DeviceService, GatewayService,
  RegionService, ReminderService, NoticeService, SuggestionService
} from '@/services'

import DashHeader from './dashboard/DashHeader.vue'
import DashStatCards from './dashboard/DashStatCards.vue'
import DashPopularRank from './dashboard/DashPopularRank.vue'
import DashDeviceStatus from './dashboard/DashDeviceStatus.vue'
import DashGatewayStatus from './dashboard/DashGatewayStatus.vue'
import DashChinaMap from './dashboard/DashChinaMap.vue'
import DashEnvTrend from './dashboard/DashEnvTrend.vue'
import DashEnvRadar from './dashboard/DashEnvRadar.vue'
import DashAlertList from './dashboard/DashAlertList.vue'
import DashSpotDetail from './dashboard/DashSpotDetail.vue'
import DashSuggestion from './dashboard/DashSuggestion.vue'
import DashNoticeList from './dashboard/DashNoticeList.vue'

// ===== 数据 =====
const summary = ref<SystemSummary | null>(null)
const popularSpots = ref<PopularSpot[]>([])
const allSpots = ref<FishingSpot[]>([])
const devices = ref<Device[]>([])
const gateways = ref<Gateway[]>([])
const regions = ref<Region[]>([])
const regionEnvItems = ref<RegionEnvItem[]>([])
const reminders = ref<Reminder[]>([])
const notices = ref<Notice[]>([])
const suggestions = ref<FishingSuggestion[]>([])

// ===== 选中水域 =====
const selectedSpotId = ref<number | null>(null)

// ===== 地图层级感知 =====
const mapLevel = ref(0)       // 0=全国, 1=省级, 2=市级
const mapAreaName = ref('全国')

function onMapLevelChange(level: number) {
  mapLevel.value = level
}
function onMapAreaNameChange(name: string) {
  mapAreaName.value = name
}

/**
 * 按省份汇总辅助函数：将多个 RegionEnvItem 聚合为省级平均
 */
function aggregateByProvince(items: RegionEnvItem[]): RegionEnvItem[] {
  const map = new Map<string, { items: RegionEnvItem[]; province: string }>()
  for (const item of items) {
    const region = regions.value.find(r => r.id === item.region_id)
    const province = region?.province || '未知'
    if (!map.has(province)) map.set(province, { items: [], province })
    map.get(province)!.items.push(item)
  }
  return Array.from(map.values()).map(({ items: group, province }) => {
    const n = group.length
    return {
      region_id: group[0].region_id,
      region_name: province,
      city: '',
      spot_count: group.reduce((s, g) => s + g.spot_count, 0),
      water_temp: +(group.reduce((s, g) => s + g.water_temp, 0) / n).toFixed(1),
      air_temp: +(group.reduce((s, g) => s + g.air_temp, 0) / n).toFixed(1),
      humidity: +(group.reduce((s, g) => s + g.humidity, 0) / n).toFixed(1),
      pressure: +(group.reduce((s, g) => s + g.pressure, 0) / n).toFixed(1),
      ph: +(group.reduce((s, g) => s + g.ph, 0) / n).toFixed(2),
      dissolved_oxygen: +(group.reduce((s, g) => s + g.dissolved_oxygen, 0) / n).toFixed(1),
      turbidity: +(group.reduce((s, g) => s + g.turbidity, 0) / n).toFixed(1),
      timestamp: group[0].timestamp
    } as RegionEnvItem
  })
}

/**
 * 按城市汇总辅助函数：将多个 RegionEnvItem 聚合为市级平均
 */
function aggregateByCity(items: RegionEnvItem[]): RegionEnvItem[] {
  const map = new Map<string, { items: RegionEnvItem[]; city: string }>()
  for (const item of items) {
    const city = item.city || '未知'
    if (!map.has(city)) map.set(city, { items: [], city })
    map.get(city)!.items.push(item)
  }
  return Array.from(map.values()).map(({ items: group, city }) => {
    const n = group.length
    return {
      region_id: group[0].region_id,
      region_name: city,
      city,
      spot_count: group.reduce((s, g) => s + g.spot_count, 0),
      water_temp: +(group.reduce((s, g) => s + g.water_temp, 0) / n).toFixed(1),
      air_temp: +(group.reduce((s, g) => s + g.air_temp, 0) / n).toFixed(1),
      humidity: +(group.reduce((s, g) => s + g.humidity, 0) / n).toFixed(1),
      pressure: +(group.reduce((s, g) => s + g.pressure, 0) / n).toFixed(1),
      ph: +(group.reduce((s, g) => s + g.ph, 0) / n).toFixed(2),
      dissolved_oxygen: +(group.reduce((s, g) => s + g.dissolved_oxygen, 0) / n).toFixed(1),
      turbidity: +(group.reduce((s, g) => s + g.turbidity, 0) / n).toFixed(1),
      timestamp: group[0].timestamp
    } as RegionEnvItem
  })
}

/**
 * 根据地图层级聚合环境数据：
 * - 全国（level 0）：按省份聚合为省平均
 * - 省级（level 1）：筛选该省内区域，按城市聚合为市平均
 * - 市级（level 2）：筛选该市内区域，显示具体钓场（区域）数据
 */
const filteredEnvItems = computed(() => {
  if (mapLevel.value === 0) {
    // 全国 → 省级平均
    return aggregateByProvince(regionEnvItems.value)
  }
  if (mapLevel.value === 1) {
    // 省级 → 先筛选该省区域，再按城市聚合
    const provinceName = mapAreaName.value
    const provinceItems = regionEnvItems.value.filter(item => {
      const region = regions.value.find(r => r.id === item.region_id)
      if (!region) return false
      return region.province === provinceName ||
        region.province + '省' === provinceName ||
        region.province + '市' === provinceName ||
        provinceName.startsWith(region.province)
    })
    return aggregateByCity(provinceItems)
  }
  // 市级 → 直接显示该市下的各个区域（钓场级别）
  const cityName = mapAreaName.value
  return regionEnvItems.value.filter(item => {
    const region = regions.value.find(r => r.id === item.region_id)
    if (!region) return false
    return region.city === cityName ||
      region.city + '市' === cityName ||
      cityName.startsWith(region.city) ||
      region.city.startsWith(cityName.replace(/市$/, ''))
  })
})
const selectedSpotName = computed(() => {
  if (!selectedSpotId.value) return ''
  return allSpots.value.find(s => s.id === selectedSpotId.value)?.name ?? ''
})

const selectedSpot = computed(() => {
  if (!selectedSpotId.value) return null
  return allSpots.value.find(s => s.id === selectedSpotId.value) ?? null
})

function onSpotSelect(spot: FishingSpot | PopularSpot) {
  selectedSpotId.value = spot.id
}

// ===== 数据加载 =====
async function loadAll() {
  const results = await Promise.allSettled([
    SummaryService.get(),
    FishingSpotService.getPopular(8),
    FishingSpotService.list({ page: 1, page_size: 100 }),
    DeviceService.list(),
    GatewayService.list(),
    RegionService.list(),
    RegionService.getEnvironment(),
    ReminderService.list({ resolved: 'false', page: 1, page_size: 20 }),
    NoticeService.list({ outdated: 'false', page: 1, page_size: 10 }),
    SuggestionService.getLatest(10)
  ])

  if (results[0].status === 'fulfilled') summary.value = results[0].value
  if (results[1].status === 'fulfilled') popularSpots.value = results[1].value
  if (results[2].status === 'fulfilled') {
    allSpots.value = results[2].value.data
    // 默认选中第一个开放水域
    if (!selectedSpotId.value && allSpots.value.length) {
      const openSpot = allSpots.value.find(s => s.status === 'open')
      if (openSpot) selectedSpotId.value = openSpot.id
    }
  }
  if (results[3].status === 'fulfilled') devices.value = results[3].value.data
  if (results[4].status === 'fulfilled') gateways.value = results[4].value.data
  if (results[5].status === 'fulfilled') regions.value = results[5].value.data
  if (results[6].status === 'fulfilled') regionEnvItems.value = results[6].value
  if (results[7].status === 'fulfilled') reminders.value = results[7].value.data
  if (results[8].status === 'fulfilled') notices.value = results[8].value.data
  if (results[9].status === 'fulfilled') suggestions.value = results[9].value
}

// ===== 定时刷新 =====
let timer: ReturnType<typeof setInterval>

onMounted(() => {
  loadAll()
  timer = setInterval(loadAll, 30_000)
})

onUnmounted(() => {
  clearInterval(timer)
})
</script>

<style scoped>
.dashboard {
  width: 100vw;
  height: 100vh;
  background: linear-gradient(135deg, #0d1b2a 0%, #1b2838 50%, #0d1b2a 100%);
  color: #e8f0ff;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;
}

.dashboard-body {
  flex: 1;
  display: flex;
  gap: 12px;
  padding: 0 16px 12px;
  min-height: 0;
}

.dash-left,
.dash-right {
  width: 350px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-height: 0;
}

.dash-center {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-width: 0;
}

.dash-center-bottom {
  display: flex;
  gap: 10px;
  height: 220px;
  flex-shrink: 0;
}

.dash-row {
  display: flex;
  gap: 10px;
}
</style>

<!-- 面板全局样式（非 scoped，供所有子组件使用） -->
<style>
.dashboard .panel {
  background: rgba(6, 30, 65, 0.75);
  border: 1px solid rgba(30, 144, 255, 0.2);
  border-radius: 8px;
  backdrop-filter: blur(8px);
  transition: border-color 0.3s;
}

.dashboard .panel:hover {
  border-color: rgba(0, 212, 255, 0.35);
}

.dashboard .panel-title {
  font-size: 14px;
  font-weight: 600;
  color: #00d4ff;
  margin-bottom: 8px;
  padding-left: 8px;
  border-left: 3px solid #00d4ff;
  letter-spacing: 1px;
  line-height: 1.4;
}
</style>
