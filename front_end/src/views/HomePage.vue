<script lang="ts" setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useFishingStore } from '@/stores/fishing'
import { useSummaryStore } from '@/stores/summary'
import { useReminderStore } from '@/stores/reminder'
import { useNoticeStore } from '@/stores/notice'
import { useEnvironmentStore } from '@/stores/environment'
import { useSuggestionStore } from '@/stores/suggestion'
import SpotList from '@/components/data/SpotList.vue'
import { TrophyBase, Star } from '@element-plus/icons-vue'
import { MockAPI } from '@/services/MockDataService'
import type { SystemSummary, Reminder, Notice, FishingSuggestion, PopularSpot } from '@/types'

// 子组件
import HomeBanner from './home/HomeBanner.vue'
import HomeStatsPanel from './home/HomeStatsPanel.vue'
import HomeSuggestionList from './home/HomeSuggestionList.vue'
import HomeReminderPanel from './home/HomeReminderPanel.vue'
import HomeNoticePanel from './home/HomeNoticePanel.vue'
import HomeTrendChart from './home/HomeTrendChart.vue'
import HomeEnvGauge from './home/HomeEnvGauge.vue'

const router = useRouter()
const authStore = useAuthStore()
const fishingStore = useFishingStore()
const summaryStore = useSummaryStore()
const reminderStore = useReminderStore()
const noticeStore = useNoticeStore()
const environmentStore = useEnvironmentStore()
const suggestionStore = useSuggestionStore()

const isMobile = ref(false)
const isFirstLoad = ref(true)
const loadingSummary = ref(false)
const loadingHot = ref(false)
const loadingSuggestions = ref(false)
const loadingReminders = ref(false)
const loadingNotices = ref(false)

// 直接管理的数据（用 mock 作为后备）
const summary = ref<SystemSummary | null>(null)
const hotSpots = ref<PopularSpot[]>([])
const suggestions = ref<FishingSuggestion[]>([])
const reminders = ref<Reminder[]>([])
const notices = ref<Notice[]>([])
const favoriteSpots = ref<PopularSpot[]>([])

const trendSpotId = ref<number | null>(null)

// ===== 数据加载（先走真实 API，失败用 Mock） =====
const fetchSummary = async () => {
  loadingSummary.value = isFirstLoad.value
  try {
    const data = await summaryStore.fetchSummary()
    summary.value = data ?? null
  } catch {
    summary.value = await MockAPI.getSummary()
  } finally {
    loadingSummary.value = false
  }
}

const fetchHotSpots = async () => {
  loadingHot.value = isFirstLoad.value
  try {
    await fishingStore.fetchPopularSpots(8)
    hotSpots.value = fishingStore.popularSpots as PopularSpot[]
  } catch {
    hotSpots.value = await MockAPI.getPopularSpots(8)
  } finally {
    loadingHot.value = false
  }
}

const fetchSuggestions = async () => {
  loadingSuggestions.value = isFirstLoad.value
  try {
    await suggestionStore.fetchLatest(5)
    suggestions.value = suggestionStore.suggestions
  } catch {
    suggestions.value = await MockAPI.getSuggestions()
  } finally {
    loadingSuggestions.value = false
  }
}

const fetchReminders = async () => {
  loadingReminders.value = isFirstLoad.value
  try {
    await reminderStore.fetchReminders({ page_size: 5 })
    reminders.value = reminderStore.reminders
  } catch {
    reminders.value = await MockAPI.getReminders()
  } finally {
    loadingReminders.value = false
  }
}

const fetchNotices = async () => {
  loadingNotices.value = isFirstLoad.value
  try {
    await noticeStore.fetchNotices({ page_size: 5 })
    notices.value = noticeStore.notices
  } catch {
    notices.value = await MockAPI.getNotices()
  } finally {
    loadingNotices.value = false
  }
}

const fetchFavorites = async () => {
  if (!authStore.isLoggedIn) return
  try {
    await fishingStore.fetchFavorites()
    favoriteSpots.value = fishingStore.favoriteSpots as PopularSpot[]
  } catch {
    favoriteSpots.value = (await MockAPI.getFavorites()) as PopularSpot[]
  }
}

// ===== 趋势数据 =====
const trendHistorical = ref<import('@/types').HistoricalData[]>([])

const loadTrendData = async (spotId: number) => {
  try {
    await environmentStore.fetchHistorical(spotId, 48)
    trendHistorical.value = environmentStore.getHistorical(spotId)
  } catch {
    const mockData = await MockAPI.getHistorical(spotId)
    environmentStore.historicalDataMap[spotId] = mockData
    trendHistorical.value = mockData
  }
}

// ===== 环境监测数据（与趋势图共用 spotId） =====
const envSpotId = ref<number | null>(null)
const envData = ref<import('@/types').EnvironmentData[]>([])

const loadEnvData = async (spotId: number) => {
  try {
    await environmentStore.fetchEnvironment(spotId, 48)
    envData.value = environmentStore.getEnvironment(spotId)
  } catch {
    const mockData = await MockAPI.getEnvironment(spotId)
    environmentStore.environmentDataMap[spotId] = mockData
    envData.value = mockData
  }
}

watch(trendSpotId, (id) => {
  if (id) loadTrendData(id)
})

watch(envSpotId, (id) => {
  if (id) loadEnvData(id)
})

// ===== 定时器 =====
let refreshTimer: ReturnType<typeof setInterval> | null = null

const checkScreenSize = () => {
  isMobile.value = window.innerWidth < 992
}

onMounted(async () => {
  checkScreenSize()
  window.addEventListener('resize', checkScreenSize)

  await Promise.all([
    fetchSummary(),
    fetchHotSpots(),
    fetchSuggestions(),
    fetchReminders(),
    fetchNotices(),
    fetchFavorites()
  ])
  isFirstLoad.value = false

  // 默认选第一个热门水域作为趋势图和环境监测图
  if (hotSpots.value.length > 0) {
    if (!trendSpotId.value) {
      trendSpotId.value = hotSpots.value[0].id
    }
    if (!envSpotId.value) {
      envSpotId.value = hotSpots.value[0].id
    }
  }

  // 30 秒轮询
  refreshTimer = setInterval(() => {
    fetchSummary()
    fetchHotSpots()
    fetchReminders()
    fetchNotices()
  }, 30000)
})

onBeforeUnmount(() => {
  if (refreshTimer) clearInterval(refreshTimer)
  window.removeEventListener('resize', checkScreenSize)
})
</script>

<template>
  <div class="home-container">
    <!-- 顶部 Banner -->
    <HomeBanner />

    <!-- 统计卡片 -->
    <HomeStatsPanel :summary="summary" :loading="loadingSummary" class="mt-20" />

    <el-row :gutter="20" class="mt-20">
      <el-col :span="isMobile ? 24 : 16" :xs="24" :sm="24" :md="16" :lg="16">
        <!-- 热门水域排行 -->
        <el-card class="dashboard-card">
          <template #header>
            <div class="card-header">
              <el-icon class="card-icon hot-icon"><TrophyBase /></el-icon>
              <span class="card-title">热门水域实时排行</span>
            </div>
          </template>
          <SpotList
            :spots="hotSpots"
            :loading="loadingHot"
            empty-text="暂无热门水域数据"
            :max-height="hotSpots.length > 8 ? '200px' : 'auto'"
          />
        </el-card>

        <!-- 垂钓建议 -->
        <HomeSuggestionList :suggestions="suggestions" :loading="loadingSuggestions" />

        <!-- 我的收藏（已登录） -->
        <el-card v-if="authStore.isLoggedIn" class="dashboard-card">
          <template #header>
            <div class="card-header">
              <el-icon class="card-icon favorite-icon"><Star /></el-icon>
              <span class="card-title">我的收藏水域</span>
            </div>
          </template>
          <SpotList
            :spots="favoriteSpots"
            :loading="false"
            :max-height="favoriteSpots.length > 6 ? '193px' : 'auto'"
            empty-text="暂无收藏水域"
          />
        </el-card>

        <!-- 垂钓数量趋势 -->
        <HomeTrendChart
          v-model:spotId="trendSpotId"
          :spots="hotSpots"
          :historical-data="trendHistorical"
        />

        <!-- 环境数据监测 -->
        <HomeEnvGauge
          v-model:spotId="envSpotId"
          :spots="hotSpots"
          :environment-data="envData"
        />
      </el-col>

      <el-col :span="isMobile ? 24 : 8" :xs="24" :sm="24" :md="8" :lg="8">
        <!-- 公开提醒 -->
        <HomeReminderPanel :reminders="reminders" :loading="loadingReminders" />

        <!-- 近期通知 -->
        <HomeNoticePanel :notices="notices" :loading="loadingNotices" />
      </el-col>
    </el-row>
  </div>
</template>

<style scoped>
.home-container {
  max-width: 1300px;
  margin: 0 auto;
  padding: 30px;
  overflow-x: hidden;
}

.mt-20 { margin-top: 20px; }

/* Dashboard 卡片通用 */
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

.dashboard-card :deep(.el-card__body) {
  flex: 1;
  display: flex;
  flex-direction: column;
  padding: 15px 20px;
  overflow: hidden;
}

.card-header { display: flex; align-items: center; }

.card-icon {
  margin-right: 8px;
  font-size: 20px;
  padding: 3px;
  border-radius: 8px;
  color: white;
}

.hot-icon { background-color: #f56c6c; }
.favorite-icon { background-color: #e6a23c; }

.card-title {
  font-size: 17px !important;
  font-weight: 600;
  color: #333;
  letter-spacing: 0.5px;
}

@media (max-width: 768px) {
  .home-container { padding: 15px; max-width: 100%; }
  .dashboard-card { margin-bottom: 15px; }
  .dashboard-card :deep(.el-card__header) { padding: 12px 15px; }
  .card-title { font-size: 16px !important; }
  .mt-20 { margin-top: 15px; }
}
</style>
