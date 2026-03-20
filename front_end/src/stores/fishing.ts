import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Region, FishingSpot, PopularSpot } from '@/types'
import { RegionService, FishingSpotService } from '@/services'

export const useFishingStore = defineStore('fishing', () => {
  // ===== State =====
  const regions = ref<Region[]>([])
  const provinces = ref<string[]>([])
  const currentRegion = ref<Region | null>(null)

  const spots = ref<FishingSpot[]>([])
  const spotsTotal = ref(0)
  const spotsPage = ref(1)
  const currentSpot = ref<FishingSpot | null>(null)

  const popularSpots = ref<PopularSpot[]>([])
  const favoriteSpots = ref<FishingSpot[]>([])

  const isLoading = ref(false)

  // ===== Getters =====
  const regionsByProvince = computed(() => {
    const map: Record<string, Region[]> = {}
    for (const region of regions.value) {
      if (!map[region.province]) {
        map[region.province] = []
      }
      map[region.province].push(region)
    }
    return map
  })

  const openSpots = computed(() => spots.value.filter(s => s.status === 'open'))

  // ===== Actions =====
  async function fetchRegions(params?: { province?: string; search?: string }) {
    isLoading.value = true
    try {
      const res = await RegionService.list(params)
      regions.value = res.data
    } finally {
      isLoading.value = false
    }
  }

  async function fetchProvinces() {
    provinces.value = await RegionService.getProvinces()
  }

  async function fetchRegion(id: number) {
    currentRegion.value = await RegionService.get(id)
    return currentRegion.value
  }

  async function fetchSpots(params?: {
    region_id?: number
    status?: string
    water_type?: string
    search?: string
    page?: number
    page_size?: number
  }) {
    isLoading.value = true
    try {
      const res = await FishingSpotService.list(params)
      spots.value = res.data
      spotsTotal.value = res.total
      spotsPage.value = params?.page ?? 1
    } finally {
      isLoading.value = false
    }
  }

  async function fetchSpot(id: number) {
    currentSpot.value = await FishingSpotService.get(id)
    return currentSpot.value
  }

  async function fetchPopularSpots(limit?: number) {
    popularSpots.value = await FishingSpotService.getPopular(limit)
  }

  async function fetchFavorites() {
    try {
      favoriteSpots.value = await FishingSpotService.getFavorites()
    } catch {
      favoriteSpots.value = []
    }
  }

  async function toggleFavorite(spotId: number) {
    const result = await FishingSpotService.toggleFavorite(spotId)
    // 更新本地收藏列表
    if (result.favorited) {
      const spot = spots.value.find(s => s.id === spotId)
      if (spot && !favoriteSpots.value.find(s => s.id === spotId)) {
        favoriteSpots.value.push(spot)
      }
    } else {
      favoriteSpots.value = favoriteSpots.value.filter(s => s.id !== spotId)
    }
    return result
  }

  function isFavorite(spotId: number): boolean {
    return favoriteSpots.value.some(s => s.id === spotId)
  }

  return {
    // State
    regions,
    provinces,
    currentRegion,
    spots,
    spotsTotal,
    spotsPage,
    currentSpot,
    popularSpots,
    favoriteSpots,
    isLoading,
    // Getters
    regionsByProvince,
    openSpots,
    // Actions
    fetchRegions,
    fetchProvinces,
    fetchRegion,
    fetchSpots,
    fetchSpot,
    fetchPopularSpots,
    fetchFavorites,
    toggleFavorite,
    isFavorite
  }
})
