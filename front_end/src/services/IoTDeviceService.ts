/**
 * IoT 设备服务 —— 对接 Go v2 接口 /api/v2/iot-devices
 *
 * 注意：IoTDevice 主键是 string (device_id)，不符合 createResourceService 的 BaseModel.id: number 要求
 * 因此采用手写服务，但结构上尽量贴近工厂模式的接口
 */

import { httpGet } from '@/network/http'
import type { IoTDeviceData } from '@/types'

export const IoTDeviceService = {
  /** 获取所有 IoT 设备列表 */
  async list(): Promise<{ data: IoTDeviceData[]; total: number }> {
    const data = await httpGet<IoTDeviceData[]>('/api/v2/iot-devices')
    return { data, total: data.length }
  },

  /** 获取单个设备数据 */
  async getDeviceData(deviceId: string): Promise<IoTDeviceData> {
    return httpGet<IoTDeviceData>(`/api/v2/iot-devices/${deviceId}`)
  }
}
