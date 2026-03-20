import { createResourceService } from './createResourceService'
import type { FishingSuggestion } from '@/types'

export const SuggestionService = createResourceService({
  name: 'suggestions',
  model: {} as FishingSuggestion,
  paginated: true,
  listParams: {} as {
    spot_id?: number
    user_id?: number
    page?: number
    page_size?: number
  },
  extend: (ctx) => ({
    /** 获取最新建议 */
    async getLatest(limit = 5): Promise<FishingSuggestion[]> {
      return ctx.http.get<FishingSuggestion[]>(`${ctx.baseURL}/latest`, { limit })
    }
  })
})
