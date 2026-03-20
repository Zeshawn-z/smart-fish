import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { FishingSuggestion } from '@/types'
import { SuggestionService } from '@/services'

export const useSuggestionStore = defineStore('suggestion', () => {
  const suggestions = ref<FishingSuggestion[]>([])
  const total = ref(0)
  const isLoading = ref(false)

  async function fetchLatest(limit = 5): Promise<FishingSuggestion[]> {
    isLoading.value = true
    try {
      const data = await SuggestionService.getLatest(limit)
      suggestions.value = data
      return data
    } finally {
      isLoading.value = false
    }
  }

  async function fetchSuggestions(params?: {
    spot_id?: number
    user_id?: number
    page?: number
    page_size?: number
  }) {
    isLoading.value = true
    try {
      const res = await SuggestionService.list(params)
      suggestions.value = res.data
      total.value = res.total
    } finally {
      isLoading.value = false
    }
  }

  return {
    suggestions,
    total,
    isLoading,
    fetchLatest,
    fetchSuggestions
  }
})
