import { httpGet } from '@/network/http'
import type { SystemSummary } from '@/types'

export const SummaryService = {
  async get(): Promise<SystemSummary> {
    return httpGet<SystemSummary>('/api/v2/summary')
  }
}
