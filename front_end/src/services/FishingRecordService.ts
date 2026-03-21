/**
 * 垂钓记录服务 —— 全部使用 v2 接口
 *
 * 标准 CRUDL 由 createResourceService 工厂自动生成
 */

import { createResourceService } from './createResourceService'
import { httpGet, httpPost } from '@/network/http'
import type { FishingRecord, FishCaught } from '@/types'

/** 垂钓记录服务（工厂模式） */
export const FishingRecordResourceService = createResourceService({
  name: 'fishing-records',
  model: {} as FishingRecord,
  paginated: true,
  listParams: {} as { user_id?: number; page?: number; page_size?: number }
})

/** 渔获服务（工厂模式） */
export const FishCaughtResourceService = createResourceService({
  name: 'fish-caught',
  model: {} as FishCaught,
  paginated: true,
  listParams: {} as { record_id?: number; page?: number; page_size?: number }
})

/**
 * 垂钓记录服务 —— 聚合接口
 */
export const FishingRecordService = {
  async getMyRecords(userId: number) {
    try {
      const res = await FishingRecordResourceService.list({ user_id: userId } as any)
      return { records: res.data }
    } catch {
      return { records: [] as FishingRecord[] }
    }
  },

  async getMyStats(): Promise<{
    total_trips: number
    total_fish: number
    total_kg: number
    max_kg: number
    total_hours: number
    fish_types: { fish_type: string; count: number }[]
  }> {
    return httpGet('/api/v2/fishing-records/stats')
  },

  async getRecord(recordId: number): Promise<FishingRecord> {
    return FishingRecordResourceService.get(recordId)
  },

  async createRecord(data: {
    start_time: string
    end_time: string
    latitude: number
    longitude: number
    device_id?: string
  }): Promise<FishingRecord> {
    return FishingRecordResourceService.create(data as Partial<FishingRecord>)
  },

  async addFishCaught(data: {
    record_id: number
    caught_time: string
    fish_type: string
    weight: number
    bait_type?: string
    bait_weight?: number
    fishing_depth?: number
  }): Promise<FishCaught> {
    return FishCaughtResourceService.create(data as Partial<FishCaught>)
  },

  /** 上传渔获图片（v2 接口） */
  async uploadFishImage(fishId: number, file: File): Promise<{ message: string; image_id: number; image_url: string }> {
    const form = new FormData()
    form.append('entity_type', 'fish')
    form.append('entity_id', String(fishId))
    form.append('file', file)
    return httpPost('/api/v2/upload/image', form, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  }
}
