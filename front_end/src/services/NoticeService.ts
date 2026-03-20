import { createResourceService } from './createResourceService'
import type { Notice } from '@/types'

export const NoticeService = createResourceService({
  name: 'notices',
  model: {} as Notice,
  paginated: true,
  listParams: {} as {
    outdated?: string
    search?: string
    page?: number
    page_size?: number
  }
})
