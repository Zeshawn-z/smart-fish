<template>
  <div class="fishing-spots-page">
    <el-row :gutter="16">
      <!-- 左侧：省份/区域树 -->
      <el-col :xs="24" :sm="7" :lg="6">
        <RegionTree
          ref="regionTreeRef"
          :regions-by-province="regionsByProvince"
          :selected-region-id="selectedRegionId"
          :loading="treeLoading"
          @select-region="selectRegion"
        />
      </el-col>

      <!-- 右侧：水域列表和详情 -->
      <el-col :xs="24" :sm="17" :lg="18">
        <!-- 筛选栏 -->
        <SpotFilterBar
          v-model:status="filterStatus"
          v-model:water-type="filterWaterType"
          v-model:keyword="searchKeyword"
          v-model:view-mode="viewMode"
          :has-region-filter="!!selectedRegionId"
          @search="loadSpots()"
          @clear-region="clearRegionFilter"
        />

        <!-- 水域卡片视图 -->
        <div v-if="viewMode === 'card'" class="spot-grid" v-loading="spotsLoading">
          <SpotCard
            v-for="spot in spotsList"
            :key="spot.id"
            :spot="spot"
            :is-favorite="isFavorite(spot.id)"
            :show-favorite="authStore.isLoggedIn"
            @click="openSpotDetail"
            @toggle-favorite="handleFavorite"
          />
          <el-empty v-if="spotsList.length === 0 && !spotsLoading" description="暂无水域数据" />
        </div>

        <!-- 水域表格视图 -->
        <el-card v-else shadow="hover" class="table-card">
          <el-table :data="spotsList" stripe v-loading="spotsLoading" empty-text="暂无水域数据">
            <el-table-column prop="name" label="名称" min-width="150" />
            <el-table-column label="区域" min-width="130">
              <template #default="{ row }">{{ row.region?.province }} · {{ row.region?.city }}</template>
            </el-table-column>
            <el-table-column label="类型" width="80">
              <template #default="{ row }"><el-tag size="small">{{ WATER_TYPE_MAP[row.water_type as WaterType] }}</el-tag></template>
            </el-table-column>
            <el-table-column label="状态" width="90" align="center">
              <template #default="{ row }">
                <el-tag :type="row.status === 'open' ? 'success' : row.status === 'closed' ? 'danger' : 'warning'" size="small" effect="dark">
                  {{ SPOT_STATUS_MAP[row.status as SpotStatus] }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="负载" width="120" align="center">
              <template #default="{ row }">
                <span v-if="row.capacity">{{ row.bound_device?.fishing_count ?? 0 }}/{{ row.capacity }}</span>
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column prop="capacity" label="容纳" width="80" align="center" />
            <el-table-column label="操作" width="140" align="center">
              <template #default="{ row }">
                <el-button size="small" type="primary" link @click="openSpotDetail(row)">详情</el-button>
                <el-button v-if="authStore.isLoggedIn" size="small" :type="isFavorite(row.id) ? 'warning' : 'default'" link @click="handleFavorite(row.id)">
                  {{ isFavorite(row.id) ? '取消收藏' : '收藏' }}
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-card>

        <!-- 分页 -->
        <div class="pagination-wrap" v-if="spotsTotal > 20">
          <el-pagination
            v-model:current-page="currentPage"
            :page-size="20"
            :total="spotsTotal"
            layout="total, prev, pager, next"
            @current-change="handlePageChange"
          />
        </div>
      </el-col>
    </el-row>

    <!-- 水域详情弹窗 -->
    <SpotDetailDialog
      v-model:visible="detailVisible"
      :spot="currentSpot"
      :environment-data="detailEnvData"
      :chart-loading="detailChartLoading"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useFishingStore } from '@/stores/fishing'
import { useEnvironmentStore } from '@/stores/environment'
import { MockAPI } from '@/services/MockDataService'
import {
  WATER_TYPE_MAP, SPOT_STATUS_MAP,
  type FishingSpot, type Region, type WaterType, type SpotStatus
} from '@/types'

// 子组件
import RegionTree from './spots/RegionTree.vue'
import SpotFilterBar from './spots/SpotFilterBar.vue'
import SpotCard from './spots/SpotCard.vue'
import SpotDetailDialog from './spots/SpotDetailDialog.vue'

const authStore = useAuthStore()
const fishingStore = useFishingStore()
const environmentStore = useEnvironmentStore()

const regionTreeRef = ref<InstanceType<typeof RegionTree>>()

const searchKeyword = ref('')
const filterStatus = ref('')
const filterWaterType = ref('')
const viewMode = ref<'card' | 'table'>('card')
const currentPage = ref(1)
const selectedRegionId = ref<number | null>(null)

const treeLoading = ref(false)
const spotsLoading = ref(false)
const detailChartLoading = ref(false)

const detailVisible = ref(false)
const currentSpot = ref<FishingSpot | null>(null)

// Mock 回退的本地数据
const localRegions = ref<Region[]>([])
const localSpots = ref<FishingSpot[]>([])
const localSpotsTotal = ref(0)
const localFavoriteIds = ref<Set<number>>(new Set())

// 计算属性
const regionsByProvince = computed(() => {
  const regions = localRegions.value.length > 0 ? localRegions.value : fishingStore.regions
  const map: Record<string, Region[]> = {}
  for (const region of regions) {
    if (!map[region.province]) map[region.province] = []
    map[region.province].push(region)
  }
  return map
})

const spotsList = computed(() =>
  localSpots.value.length > 0 ? localSpots.value : fishingStore.spots
)

const spotsTotal = computed(() =>
  localSpots.value.length > 0 ? localSpotsTotal.value : fishingStore.spotsTotal
)

const detailEnvData = computed(() => {
  if (!currentSpot.value) return []
  return environmentStore.getEnvironment(currentSpot.value.id)
})

// 区域选择
function selectRegion(region: Region) {
  selectedRegionId.value = region.id
  loadSpots()
}

function clearRegionFilter() {
  selectedRegionId.value = null
  loadSpots()
}

// 收藏
function isFavorite(spotId: number): boolean {
  return localFavoriteIds.value.has(spotId) || fishingStore.isFavorite(spotId)
}

async function handleFavorite(spotId: number) {
  try {
    const result = await fishingStore.toggleFavorite(spotId)
    if (result.favorited) {
      localFavoriteIds.value.add(spotId)
    } else {
      localFavoriteIds.value.delete(spotId)
    }
  } catch {
    if (localFavoriteIds.value.has(spotId)) {
      localFavoriteIds.value.delete(spotId)
    } else {
      localFavoriteIds.value.add(spotId)
    }
  }
}

// 加载区域
async function loadRegions() {
  treeLoading.value = true
  try {
    await fishingStore.fetchRegions()
    localRegions.value = []
  } catch {
    localRegions.value = await MockAPI.getRegions()
  } finally {
    treeLoading.value = false
  }
}

// 加载水域
async function loadSpots() {
  spotsLoading.value = true
  try {
    const params: any = { page: currentPage.value }
    if (selectedRegionId.value) params.region_id = selectedRegionId.value
    if (filterStatus.value) params.status = filterStatus.value
    if (filterWaterType.value) params.water_type = filterWaterType.value
    if (searchKeyword.value) params.search = searchKeyword.value

    await fishingStore.fetchSpots(params)
    localSpots.value = []
    localSpotsTotal.value = 0
  } catch {
    const data = await MockAPI.getSpots({
      region_id: selectedRegionId.value,
      status: filterStatus.value || undefined,
      water_type: filterWaterType.value || undefined,
      search: searchKeyword.value || undefined
    })
    localSpots.value = data.results
    localSpotsTotal.value = data.total
  } finally {
    spotsLoading.value = false
  }
}

function handlePageChange(page: number) {
  currentPage.value = page
  loadSpots()
}

async function openSpotDetail(spot: FishingSpot) {
  detailChartLoading.value = true
  try {
    currentSpot.value = await fishingStore.fetchSpot(spot.id)
  } catch {
    const mockSpot = await MockAPI.getSpot(spot.id)
    currentSpot.value = mockSpot || spot
  }
  detailVisible.value = true

  try {
    await environmentStore.fetchEnvironment(spot.id, 48)
  } catch {
    const mockEnv = await MockAPI.getEnvironment(spot.id)
    environmentStore.environmentDataMap[spot.id] = mockEnv
  }
  detailChartLoading.value = false
}

// 筛选变化时重新加载
watch([filterStatus, filterWaterType], () => {
  currentPage.value = 1
  loadSpots()
})

onMounted(async () => {
  await loadRegions()
  await loadSpots()

  if (authStore.isLoggedIn) {
    try {
      await fishingStore.fetchFavorites()
    } catch {
      localFavoriteIds.value = new Set([1, 3])
    }
  }

  // 默认展开第一个省份
  regionTreeRef.value?.expandFirst()
})
</script>

<style scoped>
.fishing-spots-page {
  min-height: calc(100vh - 160px);
}

.spot-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 18px;
}

.spot-grid :deep(.spot-card) {
  height: 100%;
}

.table-card {
  margin-bottom: 16px;
  border-radius: 12px !important;
}

.pagination-wrap {
  display: flex;
  justify-content: center;
  margin-top: 16px;
}

@media (max-width: 768px) {
  .spot-grid { grid-template-columns: 1fr; }
}
</style>
