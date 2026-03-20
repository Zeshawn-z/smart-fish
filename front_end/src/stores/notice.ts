import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Notice } from '@/types'
import { NoticeService } from '@/services'

export const useNoticeStore = defineStore('notice', () => {
  const notices = ref<Notice[]>([])
  const total = ref(0)
  const page = ref(1)
  const isLoading = ref(false)

  async function fetchNotices(params?: {
    outdated?: string
    search?: string
    page?: number
    page_size?: number
  }) {
    isLoading.value = true
    try {
      const res = await NoticeService.list(params)
      notices.value = res.data
      total.value = res.total
      page.value = params?.page ?? 1
    } finally {
      isLoading.value = false
    }
  }

  async function fetchNotice(id: number) {
    return NoticeService.get(id)
  }

  return {
    notices,
    total,
    page,
    isLoading,
    fetchNotices,
    fetchNotice
  }
})
