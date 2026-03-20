<template>
  <div class="post-card" @click="$emit('click')">
    <div class="post-body">
      <!-- 用户 + 标签 + 时间 -->
      <div class="post-meta">
        <div class="meta-left">
          <el-avatar :size="28" :src="post.avatar || undefined" class="meta-avatar">
            {{ post.avatar ? '' : (post.username || '用户')[0] }}
          </el-avatar>
          <span class="meta-author">{{ post.username || `用户 #${post.user_id}` }}</span>
          <span class="meta-sep">·</span>
          <span class="meta-time">{{ formatTime(post.created_at) }}</span>
        </div>
        <el-tag v-if="post.tag" size="small" effect="plain" round class="meta-tag">{{ post.tag }}</el-tag>
      </div>

      <!-- 标题 + 正文 -->
      <h3 class="post-title">{{ post.title }}</h3>
      <p class="post-excerpt">{{ truncate(post.body, 120) }}</p>

      <!-- 底栏统计 -->
      <div class="post-footer">
        <span class="stat-item">
          <svg class="stat-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
          </svg>
          {{ post.likes ?? 0 }}
        </span>
        <span class="stat-item">
          <svg class="stat-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
          </svg>
          {{ post.comments ?? 0 }}
        </span>
      </div>
    </div>

    <!-- 缩略图 -->
    <div v-if="firstImage" class="post-thumb">
      <el-image :src="firstImage" fit="cover" lazy class="thumb-img" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { Post } from '@/types'

const props = defineProps<{ post: Post }>()
defineEmits<{ (e: 'click'): void }>()

const firstImage = computed(() => {
  if (!props.post.image_url) return null
  if (Array.isArray(props.post.image_url)) return props.post.image_url[0] || null
  return props.post.image_url
})

function truncate(text: string, maxLen: number) {
  if (!text) return ''
  return text.length > maxLen ? text.slice(0, maxLen) + '...' : text
}

function formatTime(dateStr?: string) {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return '刚刚'
  if (mins < 60) return `${mins}分钟前`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days}天前`
  return date.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric' })
}
</script>

<style scoped>
.post-card {
  display: flex;
  gap: 16px;
  background: #fff;
  border-radius: 12px;
  border: 1px solid #f0f0f0;
  padding: 20px;
  cursor: pointer;
  transition: all 0.25s ease;
}

.post-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  border-color: #d4e5ff;
}

.post-body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

/* Meta 行 */
.post-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 10px;
}

.meta-left {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 0;
}

.meta-avatar {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: #fff;
  font-size: 12px;
  font-weight: 600;
  flex-shrink: 0;
}

.meta-author {
  font-size: 13px;
  font-weight: 600;
  color: #303133;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.meta-sep {
  color: #c0c4cc;
  font-size: 12px;
}

.meta-time {
  font-size: 12px;
  color: #a0a4ad;
  white-space: nowrap;
}

.meta-tag {
  flex-shrink: 0;
}

/* 标题 + 正文 */
.post-title {
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
  margin: 0 0 8px;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 1;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.post-excerpt {
  font-size: 13.5px;
  color: #86909c;
  margin: 0 0 14px;
  line-height: 1.65;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* 底栏统计 */
.post-footer {
  margin-top: auto;
  display: flex;
  align-items: center;
  gap: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 12.5px;
  color: #a0a4ad;
  transition: color 0.2s;
}

.stat-item:hover {
  color: #606266;
}

.stat-icon {
  width: 15px;
  height: 15px;
}

/* 缩略图 */
.post-thumb {
  width: 140px;
  height: 100px;
  flex-shrink: 0;
  border-radius: 10px;
  overflow: hidden;
}

.thumb-img {
  width: 100%;
  height: 100%;
}

@media (max-width: 768px) {
  .post-card { padding: 14px; gap: 12px; border-radius: 10px; }
  .post-thumb { width: 90px; height: 70px; border-radius: 8px; }
  .post-title { font-size: 15px; }
  .meta-avatar { display: none; }
}
</style>
