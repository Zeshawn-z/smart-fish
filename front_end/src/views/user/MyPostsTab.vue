<template>
  <div class="my-posts-tab">
    <el-empty v-if="communityStore.myPosts.length === 0" description="还没有发过帖子" />

    <div v-for="post in communityStore.myPosts" :key="post.id" class="my-post-item" @click="goToPost(post.id)">
      <div class="post-info">
        <el-tag v-if="post.tag" size="small" type="info" effect="plain">{{ post.tag }}</el-tag>
        <h4 class="post-title">{{ post.title }}</h4>
        <p class="post-excerpt">{{ truncate(post.body, 100) }}</p>
      </div>
      <el-icon class="arrow-icon"><ArrowRight /></el-icon>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ArrowRight } from '@element-plus/icons-vue'
import { useCommunityStore } from '@/stores/community'

const router = useRouter()
const communityStore = useCommunityStore()

function truncate(text: string, max: number) {
  if (!text) return ''
  return text.length > max ? text.slice(0, max) + '...' : text
}

function goToPost(id: number) {
  router.push(`/community/${id}`)
}
</script>

<style scoped>
.my-posts-tab {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.my-post-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  background: #fafbfc;
  border-radius: 8px;
  border: 1px solid #f0f0f0;
  cursor: pointer;
  transition: all 0.2s;
}

.my-post-item:hover {
  border-color: #d9ecff;
  background: #f5f8ff;
}

.post-info {
  min-width: 0;
  flex: 1;
}

.post-title {
  font-size: 15px;
  font-weight: 600;
  color: #1d2129;
  margin: 4px 0 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.post-excerpt {
  font-size: 13px;
  color: #86909c;
  margin: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.arrow-icon {
  color: #c0c4cc;
  font-size: 16px;
  flex-shrink: 0;
  margin-left: 12px;
}
</style>
