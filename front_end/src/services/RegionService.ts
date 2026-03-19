import { httpGet, httpPost, httpPut, httpDelete } from '@/network/http'
import type { Region, RegionEnvItem, RegionEnvHistory } from '@/types'

export const RegionService = {
  async list(params?: { province?: string; city?: string; search?: string }): Promise<Region[]> {
    return httpGet<Region[]>('/api/regions', params)
  },

  async get(id: number): Promise<Region> {
    return httpGet<Region>(`/api/regions/${id}`)
  },

  async getProvinces(): Promise<string[]> {
    return httpGet<string[]>('/api/regions/provinces')
  },

  /** 获取各区域最新环境数据概览 */
  async getEnvironment(): Promise<RegionEnvItem[]> {
    return httpGet<RegionEnvItem[]>('/api/regions/environment')
  },

  /** 获取某区域环境数据历史时间序列 */
  async getEnvHistory(regionId: number, hours?: number): Promise<RegionEnvHistory> {
    return httpGet<RegionEnvHistory>(`/api/regions/${regionId}/environment`, { hours })
  },

  async create(data: Partial<Region>): Promise<Region> {
    return httpPost<Region>('/api/regions', data)
  },

  async update(id: number, data: Partial<Region>): Promise<Region> {
    return httpPut<Region>(`/api/regions/${id}`, data)
  },

  async delete(id: number): Promise<void> {
    return httpDelete(`/api/regions/${id}`)
  }
}
