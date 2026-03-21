<template>
  <div v-loading="isLoading" class="post-detail-view">
    <div class="back-bar">
      <el-button text :icon="ArrowLeft" @click="router.push('/community')">返回社区</el-button>
    </div>

    <template v-if="post">
      <!-- 帖子主体 -->
      <div class="post-main">
        <!-- 作者信息 -->
        <div class="author-bar">
          <el-avatar :size="40" :src="authorAvatar.src" :style="authorAvatar.hasAvatar ? {} : authorAvatar.style" class="author-avatar">
            {{ authorAvatar.hasAvatar ? '' : authorAvatar.letter }}
          </el-avatar>
          <div class="author-info">
            <span class="author-name">{{ post.username || `用户 #${post.user_id}` }}</span>
            <span class="author-time">{{ formatTime(post.created_at) }}</span>
          </div>
          <el-tag v-if="post.tag" size="small" effect="plain" round>{{ post.tag }}</el-tag>
          <!-- 作者操作按钮 -->
          <div v-if="isAuthor" class="author-actions">
            <el-button text type="primary" :icon="Edit" size="small" @click="showEditDialog = true">编辑</el-button>
            <el-button text type="danger" :icon="Delete" size="small" @click="handleDelete">删除</el-button>
          </div>
        </div>

        <!-- 标题 -->
        <h1 class="post-title">{{ post.title }}</h1>

        <!-- 正文 -->
        <div class="post-content">
          <p v-for="(line, i) in contentLines" :key="i">{{ line }}</p>
        </div>

        <!-- 图片 -->
        <div v-if="imageUrls.length > 0" class="post-images">
          <el-image
            v-for="(url, idx) in imageUrls"
            :key="idx"
            :src="url"
            :preview-src-list="imageUrls"
            :initial-index="idx"
            fit="cover"
            class="post-image"
            lazy
          />
        </div>

        <!-- 互动栏 -->
        <div class="interact-bar">
          <button class="interact-btn" :class="{ liked: hasLiked }" @click="toggleLike">
            <svg class="interact-icon" viewBox="0 0 24 24" :fill="hasLiked ? 'currentColor' : 'none'" stroke="currentColor" stroke-width="2">
              <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
            </svg>
            <span>{{ likeCount }}</span>
          </button>
          <span class="interact-btn" style="cursor: default;">
            <svg class="interact-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
            </svg>
            <span>{{ comments.length }} 评论</span>
          </span>
        </div>
      </div>

      <!-- 评论区 -->
      <CommentSection :post-id="postId" :comments="comments" :is-logged-in="authStore.isLoggedIn" @refresh="refreshComments" />

      <!-- 编辑帖子对话框 -->
      <EditPostDialog v-model="showEditDialog" :post="post" @updated="loadPost" />
    </template>

    <div v-else-if="!isLoading" class="not-found">
      <el-empty description="帖子不存在或已被删除" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ArrowLeft, Edit, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { CommunityService, PostService } from '@/services/CommunityService'
import { useAvatar } from '@/composables/useAvatar'
import type { Post, Comment as CommentType } from '@/types'

import CommentSection from './community/CommentSection.vue'
import EditPostDialog from './community/EditPostDialog.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const post = ref<Post | null>(null)
const comments = ref<CommentType[]>([])
const isLoading = ref(false)
const likeCount = ref(0)
const hasLiked = ref(false)
const showEditDialog = ref(false)

const postId = computed(() => Number(route.params.id))
const contentLines = computed(() => post.value?.body?.split('\n') || [])
const imageUrls = computed(() => {
  if (!post.value) return []
  const urls = (post.value as any).image_urls || post.value.image_url
  if (!urls) return []
  return Array.isArray(urls) ? urls : [urls]
})
const isAuthor = computed(() => {
  return authStore.isLoggedIn && post.value && authStore.user?.id === post.value.user_id
})

const authorAvatar = computed(() => useAvatar((post.value as any)?.avatar, post.value?.username))

function formatTime(dateStr?: string) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
    + ' ' + d.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

async function loadPost() {
  isLoading.value = true
  try {
    // 帖子详情现在直接包含完整评论数据（含子评论 + 点赞数）
    const detail = await CommunityService.getPost(postId.value)
    post.value = detail as Post

    // 评论从帖子响应中直接获取（后端已嵌入）
    comments.value = (detail as any).comments_list || []

    try {
      const likesRes = await CommunityService.getPostLikes(postId.value)
      likeCount.value = likesRes.likes || 0
      hasLiked.value = likesRes.liked || false
    } catch { /* 默认0/false */ }
  } catch {
    post.value = null
  } finally {
    isLoading.value = false
  }
}

async function refreshComments() {
  // 刷新时重新拉取帖子详情（评论已嵌入其中）
  const detail = await CommunityService.getPost(postId.value)
  comments.value = (detail as any).comments_list || []
}

async function toggleLike() {
  if (!authStore.isLoggedIn) { ElMessage.warning('请先登录'); return }
  try {
    if (hasLiked.value) {
      await CommunityService.unlikePost(postId.value)
      likeCount.value--
      hasLiked.value = false
    } else {
      await CommunityService.likePost(postId.value)
      likeCount.value++
      hasLiked.value = true
    }
  } catch { /* Already Liked / not found */ }
}

async function handleDelete() {
  try {
    await ElMessageBox.confirm('确定要删除这篇帖子吗？删除后无法恢复。', '删除帖子', {
      confirmButtonText: '确认删除',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await PostService.delete(postId.value)
    ElMessage.success('帖子已删除')
    router.push('/community')
  } catch (e) {
    if (e !== 'cancel') {
      ElMessage.error('删除失败，请稍后再试')
    }
  }
}

// 监听 postId 变化自动加载（含首次），支持帖子间直接导航
watch(postId, () => loadPost(), { immediate: true })
</script>

<style scoped>
.post-detail-view {
  max-width: 800px;
  margin: 0 auto;
  padding: 0 16px;
}

.back-bar { margin-bottom: 12px; }

/* ===== Author Bar ===== */
.author-bar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 20px;
}

.author-actions {
  margin-left: auto;
  display: flex;
  gap: 4px;
}

.author-avatar {
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  flex-shrink: 0;
}

.author-info {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-width: 0;
}

.author-name {
  font-size: 15px;
  font-weight: 600;
  color: #1d2129;
}

.author-time {
  font-size: 12px;
  color: #a0a4ad;
  margin-top: 2px;
}

/* ===== Post Main ===== */
.post-main {
  background: #fff;
  border-radius: 14px;
  border: 1px solid #f0f0f0;
  padding: 28px;
  margin-bottom: 16px;
}

.post-title {
  font-size: 24px;
  font-weight: 700;
  color: #1d2129;
  margin: 0 0 20px;
  line-height: 1.45;
}

.post-content {
  font-size: 15px;
  line-height: 1.85;
  color: #4e5969;
}
.post-content p { margin: 0 0 10px; }

/* ===== Images ===== */
.post-images {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 10px;
  margin-top: 20px;
}
.post-image {
  width: 100%;
  height: 160px;
  border-radius: 10px;
  overflow: hidden;
  cursor: zoom-in;
}

/* ===== Interact Bar ===== */
.interact-bar {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-top: 24px;
  padding-top: 16px;
  border-top: 1px solid #f2f3f5;
}

.interact-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #86909c;
  background: none;
  border: none;
  cursor: pointer;
  padding: 6px 12px;
  border-radius: 8px;
  transition: all 0.2s;
}

.interact-btn:hover {
  background: #f7f8fa;
  color: #4e5969;
}

.interact-btn.liked {
  color: #f56c6c;
}
.interact-btn.liked:hover {
  background: #fef0f0;
}

.interact-icon {
  width: 18px;
  height: 18px;
}

.not-found { padding: 80px 0; }

@media (max-width: 768px) {
  .post-detail-view { padding: 0 12px; }
  .post-main { padding: 20px; border-radius: 10px; }
  .post-title { font-size: 20px; }
  .post-image { height: 200px; }
}
</style>
