import { createResourceService } from './createResourceService'
import type { Reminder } from '@/types'

export const ReminderService = createResourceService({
  name: 'reminders',
  model: {} as Reminder,
  paginated: true,
  cache: { ttl: 30_000 },
  listParams: {} as {
    spot_id?: number
    level?: number
    resolved?: string
    page?: number
    page_size?: number
  },
  extend: (ctx) => ({
    /** 标记提醒为已解决 */
    async resolve(id: number): Promise<{ message: string }> {
      const result = await ctx.http.patch<{ message: string }>(`${ctx.baseURL}/${id}/resolve`)
      ctx.cache.removeItem(id)
      ctx.cache.invalidateList()
      return result
    }
  })
})
