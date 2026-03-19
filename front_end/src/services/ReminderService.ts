import { httpGet, httpPost, httpPatch, httpDelete } from '@/network/http'
import type { Reminder, PaginatedResponse } from '@/types'

export const ReminderService = {
  async list(params?: {
    spot_id?: number
    level?: number
    resolved?: string
    page?: number
    page_size?: number
  }): Promise<PaginatedResponse<Reminder>> {
    return httpGet<PaginatedResponse<Reminder>>('/api/reminders', params)
  },

  async get(id: number): Promise<Reminder> {
    return httpGet<Reminder>(`/api/reminders/${id}`)
  },

  async create(data: Partial<Reminder>): Promise<Reminder> {
    return httpPost<Reminder>('/api/reminders', data)
  },

  async resolve(id: number): Promise<{ message: string }> {
    return httpPatch(`/api/reminders/${id}/resolve`)
  },

  async delete(id: number): Promise<void> {
    return httpDelete(`/api/reminders/${id}`)
  }
}
