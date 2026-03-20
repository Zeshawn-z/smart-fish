import { createResourceService } from './createResourceService'
import type { Gateway } from '@/types'

export const GatewayService = createResourceService({
  name: 'gateways',
  model: {} as Gateway,
  listParams: {} as { status?: string; search?: string }
})
