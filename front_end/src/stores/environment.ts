import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { HistoricalData, EnvironmentData } from '@/types'
import { FishingSpotService } from '@/services'

export const useEnvironmentStore = defineStore('environment', () => {
  // 按 spotId 缓存环境数据
  const historicalDataMap = ref<Record<number, HistoricalData[]>>({})
  const environmentDataMap = ref<Record<number, EnvironmentData[]>>({})
  const isLoading = ref(false)

  async function fetchHistorical(spotId: number, limit?: number) {
    isLoading.value = true
    try {
      const data = await FishingSpotService.getHistorical(spotId, limit)
      historicalDataMap.value[spotId] = data
      return data
    } finally {
      isLoading.value = false
    }
  }

  async function fetchEnvironment(spotId: number, limit?: number) {
    isLoading.value = true
    try {
      const data = await FishingSpotService.getEnvironment(spotId, limit)
      environmentDataMap.value[spotId] = data
      return data
    } finally {
      isLoading.value = false
    }
  }

  function getHistorical(spotId: number): HistoricalData[] {
    return historicalDataMap.value[spotId] || []
  }

  function getEnvironment(spotId: number): EnvironmentData[] {
    return environmentDataMap.value[spotId] || []
  }

  return {
    historicalDataMap,
    environmentDataMap,
    isLoading,
    fetchHistorical,
    fetchEnvironment,
    getHistorical,
    getEnvironment
  }
})
