import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Post, FishingRecord } from '@/types'
import { CommunityService } from '@/services/CommunityService'
import { FishingRecordService } from '@/services/FishingRecordService'

/** 社区 Store —— 管理帖子列表 */
export const useCommunityStore = defineStore('community', () => {
  const posts = ref<Post[]>([])
  const myPosts = ref<Post[]>([])
  const isLoading = ref(false)

  async function fetchPosts() {
    isLoading.value = true
    try {
      const res = await CommunityService.getPosts()
      posts.value = res.posts_list || []
    } catch {
      posts.value = []
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
    fetchPosts,
    fetchMyPosts
  }
})

/** 垂钓记录 Store */
export const useFishingRecordStore = defineStore('fishingRecord', () => {
  const records = ref<FishingRecord[]>([])
  const isLoading = ref(false)

  async function fetchRecords() {
    isLoading.value = true
    try {
      const res = await FishingRecordService.getMyRecords()
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
