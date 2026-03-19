import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { SystemSummary } from '@/types'
import { SummaryService } from '@/services'

export const useSummaryStore = defineStore('summary', () => {
  const summary = ref<SystemSummary | null>(null)
  const isLoading = ref(false)
  const lastFetchTime = ref<number>(0)

  const CACHE_DURATION = 30 * 1000 // 30秒缓存

  async function fetchSummary(force = false) {
    const now = Date.now()
    if (!force && summary.value && now - lastFetchTime.value < CACHE_DURATION) {
      return summary.value
    }

    isLoading.value = true
    try {
      summary.value = await SummaryService.get()
      lastFetchTime.value = now
      return summary.value
    } finally {
      isLoading.value = false
    }
  }

  return {
    summary,
    isLoading,
    fetchSummary
  }
})
