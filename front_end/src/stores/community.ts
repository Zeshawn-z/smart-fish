import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Post, FishingRecord } from '@/types'
import { CommunityService } from '@/services/CommunityService'
import { FishingRecordService } from '@/services/FishingRecordService'

/** 社区 Store —— 管理帖子列表（支持分页+搜索+标签过滤） */
export const useCommunityStore = defineStore('community', () => {
  const posts = ref<Post[]>([])
  const myPosts = ref<Post[]>([])
  const isLoading = ref(false)
  const total = ref(0)
  const currentPage = ref(1)
  const pageSize = ref(20)
  const hasMore = ref(true)
  const searchKeyword = ref('')
  const activeTag = ref('')

  /** 获取帖子（首次加载或条件变更时调用） */
  async function fetchPosts() {
    isLoading.value = true
    currentPage.value = 1
    try {
      const params: Record<string, any> = { page: 1, page_size: pageSize.value }
      if (activeTag.value) params.tag = activeTag.value
      if (searchKeyword.value) params.search = searchKeyword.value
      const res = await CommunityService.getPosts(params)
      posts.value = res.posts_list || []
      total.value = res.total || 0
      hasMore.value = posts.value.length < total.value
    } catch {
      posts.value = []
      total.value = 0
      hasMore.value = false
    } finally {
      isLoading.value = false
    }
  }

  /** 加载更多帖子（追加） */
  async function loadMore() {
    if (!hasMore.value || isLoading.value) return
    isLoading.value = true
    currentPage.value++
    try {
      const params: Record<string, any> = { page: currentPage.value, page_size: pageSize.value }
      if (activeTag.value) params.tag = activeTag.value
      if (searchKeyword.value) params.search = searchKeyword.value
      const res = await CommunityService.getPosts(params)
      const newPosts = res.posts_list || []
      posts.value = [...posts.value, ...newPosts]
      hasMore.value = posts.value.length < (res.total || 0)
    } catch {
      // 失败时回退页码
      currentPage.value--
    } finally {
      isLoading.value = false
    }
  }

  async function fetchMyPosts(userId: number) {
    try {
      const res = await CommunityService.getMyPosts(userId)
      myPosts.value = res.posts_list || []
    } catch {
      myPosts.value = []
    }
  }

  return {
    posts,
    myPosts,
    isLoading,
    total,
    currentPage,
    hasMore,
    searchKeyword,
    activeTag,
    fetchPosts,
    loadMore,
    fetchMyPosts
  }
})

/** 垂钓记录 Store */
export const useFishingRecordStore = defineStore('fishingRecord', () => {
  const records = ref<FishingRecord[]>([])
  const isLoading = ref(false)

  async function fetchRecords(userId: number) {
    isLoading.value = true
    try {
      const res = await FishingRecordService.getMyRecords(userId)
      records.value = res.records || []
    } catch {
      records.value = []
    } finally {
      isLoading.value = false
    }
  }

  return {
    records,
    isLoading,
    fetchRecords
  }
})
