import { createResourceService } from './createResourceService'
import type { Gateway } from '@/types'

export const GatewayService = createResourceService({
  name: 'gateways',
  model: {} as Gateway,
  paginated: true,
  listParams: {} as { status?: string; search?: string; page?: number; page_size?: number }
})
