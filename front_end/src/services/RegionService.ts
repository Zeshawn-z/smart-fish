import { createResourceService } from './createResourceService'
import type { Region, RegionEnvItem, RegionEnvHistory } from '@/types'

export const RegionService = createResourceService({
  name: 'regions',
  model: {} as Region,
  paginated: true,
  listParams: {} as { province?: string; city?: string; search?: string; page?: number; page_size?: number },
  extend: (ctx) => ({
    /** 获取省份列表 */
    async getProvinces(): Promise<string[]> {
      return ctx.http.get<string[]>(`${ctx.baseURL}/provinces`)
    },

    /** 获取各区域最新环境数据概览 */
    async getEnvironment(): Promise<RegionEnvItem[]> {
      return ctx.http.get<RegionEnvItem[]>(`${ctx.baseURL}/environment`)
    },

    /** 获取某区域环境数据历史时间序列 */
    async getEnvHistory(regionId: number, hours?: number): Promise<RegionEnvHistory> {
      return ctx.http.get<RegionEnvHistory>(`${ctx.baseURL}/${regionId}/environment`, { hours })
    }
  })
})
