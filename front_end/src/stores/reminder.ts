import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { Reminder, PaginatedResponse } from '@/types'
import { ReminderService } from '@/services'

export const useReminderStore = defineStore('reminder', () => {
  const reminders = ref<Reminder[]>([])
  const total = ref(0)
  const page = ref(1)
  const isLoading = ref(false)

  const unresolvedCount = computed(() =>
    reminders.value.filter(r => !r.resolved).length
  )

  const urgentReminders = computed(() =>
    reminders.value.filter(r => !r.resolved && r.level >= 2)
  )

  async function fetchReminders(params?: {
    spot_id?: number
    level?: number
    resolved?: string
    page?: number
    page_size?: number
  }) {
    isLoading.value = true
    try {
      const data: PaginatedResponse<Reminder> = await ReminderService.list(params)
      reminders.value = data.results
      total.value = data.total
      page.value = data.page
    } finally {
      isLoading.value = false
    }
  }

  async function resolveReminder(id: number) {
    await ReminderService.resolve(id)
    const reminder = reminders.value.find(r => r.id === id)
    if (reminder) {
      reminder.resolved = true
    }
  }

  return {
    reminders,
    total,
    page,
    isLoading,
    unresolvedCount,
    urgentReminders,
    fetchReminders,
    resolveReminder
  }
})
