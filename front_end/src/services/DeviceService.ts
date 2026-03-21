import { createResourceService } from './createResourceService'
import type { Device } from '@/types'

export const DeviceService = createResourceService({
  name: 'devices',
  model: {} as Device,
  paginated: true,
  listParams: {} as { gateway_id?: number; status?: string; device_type?: string; search?: string; page?: number; page_size?: number }
})
