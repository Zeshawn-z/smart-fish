import { createResourceService } from './createResourceService'
import type { FishingSpot, PopularSpot, HistoricalData, EnvironmentData } from '@/types'

export const FishingSpotService = createResourceService({
  name: 'spots',
  model: {} as FishingSpot,
  paginated: true,
  listParams: {} as {
    region_id?: number
    status?: string
    water_type?: string
    search?: string
    page?: number
    page_size?: number
  },
  extend: (ctx) => ({
    /** 热门水域 */
    async getPopular(limit?: number): Promise<PopularSpot[]> {
      return ctx.http.get<PopularSpot[]>(`${ctx.baseURL}/popular`, { limit })
    },

    /** 历史垂钓数据 */
    async getHistorical(spotId: number, limit?: number): Promise<HistoricalData[]> {
      return ctx.http.get<HistoricalData[]>(`${ctx.baseURL}/${spotId}/historical`, { limit })
    },

    /** 环境数据 */
    async getEnvironment(spotId: number, limit?: number): Promise<EnvironmentData[]> {
      return ctx.http.get<EnvironmentData[]>(`${ctx.baseURL}/${spotId}/environment`, { limit })
    },

    /** 收藏/取消收藏 */
    async toggleFavorite(spotId: number): Promise<{ message: string; favorited: boolean }> {
      return ctx.http.post<{ message: string; favorited: boolean }>(`${ctx.baseURL}/${spotId}/favor`)
    },

    /** 收藏列表 */
    async getFavorites(): Promise<FishingSpot[]> {
      return ctx.http.get<FishingSpot[]>(`${ctx.baseURL}/favorites`)
    }
  })
})
