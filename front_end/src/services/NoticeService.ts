import { httpGet, httpPost, httpPut, httpDelete } from '@/network/http'
import type { Notice, PaginatedResponse } from '@/types'

export const NoticeService = {
  async list(params?: {
    outdated?: string
    search?: string
    page?: number
    page_size?: number
  }): Promise<PaginatedResponse<Notice>> {
    return httpGet<PaginatedResponse<Notice>>('/api/notices', params)
  },

  async get(id: number): Promise<Notice> {
    return httpGet<Notice>(`/api/notices/${id}`)
  },

  async create(data: { title: string; content: string; spot_ids?: number[] }): Promise<Notice> {
    return httpPost<Notice>('/api/notices', data)
  },

  async update(id: number, data: Partial<Notice & { spot_ids?: number[] }>): Promise<Notice> {
    return httpPut<Notice>(`/api/notices/${id}`, data)
  },

  async delete(id: number): Promise<void> {
    return httpDelete(`/api/notices/${id}`)
  }
}
