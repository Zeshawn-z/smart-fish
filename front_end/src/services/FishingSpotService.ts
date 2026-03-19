import { httpGet, httpPost, httpPut, httpDelete } from '@/network/http'
import type { FishingSpot, PopularSpot, HistoricalData, EnvironmentData, PaginatedResponse } from '@/types'

export const FishingSpotService = {
  async list(params?: {
    region_id?: number
    status?: string
    water_type?: string
    search?: string
    page?: number
    page_size?: number
  }): Promise<PaginatedResponse<FishingSpot>> {
    return httpGet<PaginatedResponse<FishingSpot>>('/api/spots', params)
  },

  async get(id: number): Promise<FishingSpot> {
    return httpGet<FishingSpot>(`/api/spots/${id}`)
  },

  async getPopular(limit?: number): Promise<PopularSpot[]> {
    return httpGet<PopularSpot[]>('/api/spots/popular', { limit })
  },

  async getHistorical(spotId: number, limit?: number): Promise<HistoricalData[]> {
    return httpGet<HistoricalData[]>(`/api/spots/${spotId}/historical`, { limit })
  },

  async getEnvironment(spotId: number, limit?: number): Promise<EnvironmentData[]> {
    return httpGet<EnvironmentData[]>(`/api/spots/${spotId}/environment`, { limit })
  },

  async toggleFavorite(spotId: number): Promise<{ message: string; favorited: boolean }> {
    return httpPost(`/api/spots/${spotId}/favor`)
  },

  async getFavorites(): Promise<FishingSpot[]> {
    return httpGet<FishingSpot[]>('/api/spots/favorites')
  },

  async create(data: Partial<FishingSpot>): Promise<FishingSpot> {
    return httpPost<FishingSpot>('/api/spots', data)
  },

  async update(id: number, data: Partial<FishingSpot>): Promise<FishingSpot> {
    return httpPut<FishingSpot>(`/api/spots/${id}`, data)
  },

  async delete(id: number): Promise<void> {
    return httpDelete(`/api/spots/${id}`)
  }
}
