import { createResourceService } from './createResourceService'
import type { Device } from '@/types'

export const DeviceService = createResourceService({
  name: 'devices',
  model: {} as Device,
  listParams: {} as { gateway_id?: number; status?: string; device_type?: string; search?: string }
})
