import { httpGet, httpPost, httpPut, httpDelete } from '@/network/http'
import type { Device } from '@/types'

export const DeviceService = {
  async list(params?: {
    gateway_id?: number
    status?: string
    device_type?: string
    search?: string
  }): Promise<Device[]> {
    return httpGet<Device[]>('/api/devices', params)
  },

  async get(id: number): Promise<Device> {
    return httpGet<Device>(`/api/devices/${id}`)
  },

  async create(data: Partial<Device>): Promise<Device> {
    return httpPost<Device>('/api/devices', data)
  },

  async update(id: number, data: Partial<Device>): Promise<Device> {
    return httpPut<Device>(`/api/devices/${id}`, data)
  },

  async delete(id: number): Promise<void> {
    return httpDelete(`/api/devices/${id}`)
  }
}
