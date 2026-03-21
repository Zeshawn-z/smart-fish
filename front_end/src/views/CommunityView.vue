<template>
  <div class="community-view">
    <!-- 头部 Banner -->
    <div class="community-header">
      <div class="header-content">
        <div class="header-text">
          <h1>🎣 钓友社区</h1>
          <p>分享你的垂钓故事，交流经验心得</p>
        </div>
        <el-button
          v-if="authStore.isLoggedIn"
          type="primary"
          :icon="EditPen"
          size="large"
          round
          @click="showCreateDialog = true"
        >
          发布帖子
        </el-button>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="search-bar">
      <el-input
        v-model="searchInput"
        placeholder="搜索帖子标题或内容..."
        :prefix-icon="Search"
        clearable
        size="large"
        @keyup.enter="handleSearch"
        @clear="handleSearch"
      />
    </div>

    <!-- 标签过滤 -->
    <div class="tag-bar">
      <div class="tag-list">
        <button
          class="tag-btn"
          :class="{ active: communityStore.activeTag === '' }"
          @click="setTag('')"
        >
          全部
        </button>
        <button
          v-for="t in TAGS"
          :key="t"
          class="tag-btn"
          :class="{ active: communityStore.activeTag === t }"
          @click="setTag(t)"
        >
          {{ TAG_ICONS[t] || '' }} {{ t }}
        </button>
      </div>
    </div>

    <!-- 帖子列表 -->
    <div v-loading="communityStore.isLoading && communityStore.posts.length === 0" class="post-list">
      <el-empty
        v-if="communityStore.posts.length === 0 && !communityStore.isLoading"
        description="暂无帖子，快来发第一篇吧！"
        :image-size="120"
      />
      <PostCard
        v-for="post in communityStore.posts"
        :key="post.id"
        :post="post"
        @click="router.push(`/community/${post.id}`)"
      />
    </div>

    <!-- 加载更多 -->
    <div v-if="communityStore.posts.length > 0" class="load-more-wrap">
      <el-button
        v-if="communityStore.hasMore"
        :loading="communityStore.isLoading"
        round
        @click="communityStore.loadMore()"
      >
        {{ communityStore.isLoading ? '加载中...' : '加载更多' }}
      </el-button>
      <span v-else class="no-more">没有更多了</span>
    </div>

    <!-- 发帖对话框 -->
    <CreatePostDialog v-model="showCreateDialog" @created="communityStore.fetchPosts()" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { EditPen, Search } from '@element-plus/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { useCommunityStore } from '@/stores/community'

import PostCard from './community/PostCard.vue'
import CreatePostDialog from './community/CreatePostDialog.vue'

const TAGS = ['钓鱼日记', '经验分享', '装备测评', '钓点推荐', '问答求助']
const TAG_ICONS: Record<string, string> = {
  '钓鱼日记': '📖',
  '经验分享': '💡',
  '装备测评': '🔧',
  '钓点推荐': '📍',
  '问答求助': '❓'
}

const router = useRouter()
const authStore = useAuthStore()
const communityStore = useCommunityStore()

const showCreateDialog = ref(false)
const searchInput = ref('')

function setTag(tag: string) {
  communityStore.activeTag = tag
  communityStore.fetchPosts()
}

function handleSearch() {
  communityStore.searchKeyword = searchInput.value.trim()
  communityStore.fetchPosts()
}

onMounted(() => {
  communityStore.fetchPosts()
})
</script>

<style scoped>
.community-view {
  max-width: 900px;
  margin: 0 auto;
  padding: 0 16px;
}

/* ===== Header Banner ===== */
.community-header {
  background: linear-gradient(135deg, #1a73e8 0%, #4fc3f7 100%);
  border-radius: 16px;
  padding: 32px;
  margin-bottom: 24px;
  position: relative;
  overflow: hidden;
}

.community-header::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -20%;
  width: 300px;
  height: 300px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
}

.community-header::after {
  content: '';
  position: absolute;
  bottom: -30%;
  left: 10%;
  width: 200px;
  height: 200px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.06);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: relative;
  z-index: 1;
}

.header-text h1 {
  font-size: 26px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 6px;
  letter-spacing: 0.5px;
}

.header-text p {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.85);
  margin: 0;
}

/* ===== Search Bar ===== */
.search-bar {
  margin-bottom: 16px;
}

.search-bar :deep(.el-input__wrapper) {
  border-radius: 12px;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.06);
}

/* ===== Tag Bar ===== */
.tag-bar {
  margin-bottom: 20px;
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.tag-btn {
  padding: 6px 16px;
  border-radius: 20px;
  border: 1px solid #e4e7ed;
  background: #fff;
  font-size: 13px;
  color: #606266;
  cursor: pointer;
  transition: all 0.25s;
  white-space: nowrap;
}

.tag-btn:hover {
  color: #1a73e8;
  border-color: #b3d8ff;
  background: #ecf5ff;
}

.tag-btn.active {
  color: #fff;
  background: #1a73e8;
  border-color: #1a73e8;
  box-shadow: 0 2px 8px rgba(26, 115, 232, 0.3);
}

/* ===== Post List ===== */
.post-list {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

/* ===== Load More ===== */
.load-more-wrap {
  display: flex;
  justify-content: center;
  padding: 24px 0 40px;
}

.no-more {
  font-size: 13px;
  color: #a0a4ad;
}

@media (max-width: 768px) {
  .community-view { padding: 0 12px; }
  .community-header { padding: 24px 20px; border-radius: 12px; }
  .header-content { flex-direction: column; align-items: flex-start; gap: 16px; }
  .header-text h1 { font-size: 22px; }
  .tag-list { gap: 6px; }
  .tag-btn { padding: 5px 12px; font-size: 12px; }
}
</style>
