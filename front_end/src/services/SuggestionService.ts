import { httpGet, httpPost, httpDelete } from '@/network/http'
import type { FishingSuggestion, PaginatedResponse } from '@/types'

export const SuggestionService = {
  async list(params?: {
    spot_id?: number
    user_id?: number
    page?: number
    page_size?: number
  }): Promise<PaginatedResponse<FishingSuggestion>> {
    return httpGet<PaginatedResponse<FishingSuggestion>>('/api/suggestions', params)
  },

  async getLatest(limit = 5): Promise<FishingSuggestion[]> {
    return httpGet<FishingSuggestion[]>('/api/suggestions/latest', { limit })
  },

  async get(id: number): Promise<FishingSuggestion> {
    return httpGet<FishingSuggestion>(`/api/suggestions/${id}`)
  },

  async create(data: Partial<FishingSuggestion>): Promise<FishingSuggestion> {
    return httpPost<FishingSuggestion>('/api/suggestions', data)
  },

  async delete(id: number): Promise<void> {
    return httpDelete(`/api/suggestions/${id}`)
  }
}
