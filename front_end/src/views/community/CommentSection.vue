<template>
  <div class="comment-section">
    <h3 class="section-title">
      <svg class="title-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
      </svg>
      评论 ({{ comments.length }})
    </h3>

    <!-- 发评论 -->
    <div v-if="isLoggedIn" class="comment-input">
      <el-input
        v-model="newComment"
        type="textarea"
        :rows="2"
        placeholder="写下你的评论..."
        maxlength="1000"
        resize="none"
      />
      <div class="input-actions">
        <span class="char-count">{{ newComment.length }}/1000</span>
        <el-button type="primary" size="small" :loading="isCommenting" round @click="handleComment">
          发表评论
        </el-button>
      </div>
    </div>
    <div v-else class="login-hint">
      <router-link to="/auth" class="login-link">登录</router-link>
      后才能发表评论
    </div>

    <!-- 评论列表 -->
    <el-empty v-if="comments.length === 0" description="暂无评论，快来抢沙发！" :image-size="80" />

    <div v-for="comment in comments" :key="comment.id" class="comment-item">
      <div class="comment-left">
        <el-avatar :size="32" :src="comment.avatar || undefined" class="comment-avatar">
          {{ comment.avatar ? '' : (comment.username || '用户')[0] }}
        </el-avatar>
      </div>
      <div class="comment-right">
        <div class="comment-top">
          <span class="comment-user">{{ comment.username || `用户 #${comment.user_id}` }}</span>
          <span class="comment-time">{{ formatTime(comment.created_at) }}</span>
        </div>
        <p class="comment-body">{{ comment.body }}</p>

        <!-- 子评论区（楼中楼） -->
        <div class="sub-comments">
          <button class="toggle-btn" @click="toggleSubComments(comment.id)">
            {{ expandedComments.has(comment.id) ? '收起回复 ▲' : '查看回复 ▼' }}
          </button>

          <template v-if="expandedComments.has(comment.id)">
            <div class="sub-list">
              <div v-for="coc in subCommentsMap.get(comment.id) || []" :key="coc.coc_id" class="sub-comment-item">
                <span class="sub-user">{{ coc.username || `用户 #${coc.user_id}` }}</span>
                <span v-if="coc.to_username" class="reply-to"> 回复 <b>{{ coc.to_username }}</b></span>
                <span class="sub-body">：{{ coc.body }}</span>
              </div>

              <div v-if="(subCommentsMap.get(comment.id) || []).length === 0" class="sub-empty">
                暂无回复
              </div>
            </div>

            <div v-if="isLoggedIn" class="sub-reply-input">
              <el-input
                v-model="subReplyMap[comment.id]"
                size="small"
                placeholder="回复..."
                @keyup.enter="handleSubComment(comment.id)"
              />
              <el-button size="small" type="primary" round @click="handleSubComment(comment.id)">回复</el-button>
            </div>
          </template>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { CommunityService } from '@/services/CommunityService'
import type { Comment as CommentType, CommentOnComment } from '@/types'

const props = defineProps<{
  postId: number
  comments: CommentType[]
  isLoggedIn: boolean
}>()

const emit = defineEmits<{ (e: 'refresh'): void }>()

const newComment = ref('')
const isCommenting = ref(false)

const expandedComments = reactive(new Set<number>())
const subCommentsMap = reactive(new Map<number, CommentOnComment[]>())
const subReplyMap = reactive<Record<number, string>>({})

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

async function handleComment() {
  if (!newComment.value.trim()) return
  isCommenting.value = true
  try {
    await CommunityService.createComment(props.postId, newComment.value)
    newComment.value = ''
    ElMessage.success('评论成功')
    emit('refresh')
  } catch {
    ElMessage.error('评论失败')
  } finally {
    isCommenting.value = false
  }
}

async function toggleSubComments(commentId: number) {
  if (expandedComments.has(commentId)) {
    expandedComments.delete(commentId)
  } else {
    expandedComments.add(commentId)
    if (!subCommentsMap.has(commentId)) {
      try {
        const res = await CommunityService.getSubComments(commentId)
        subCommentsMap.set(commentId, res.comments || [])
      } catch {
        subCommentsMap.set(commentId, [])
      }
    }
  }
}

async function handleSubComment(commentId: number) {
  const body = subReplyMap[commentId]?.trim()
  if (!body) return
  try {
    await CommunityService.createSubComment(commentId, body)
    subReplyMap[commentId] = ''
    const res = await CommunityService.getSubComments(commentId)
    subCommentsMap.set(commentId, res.comments || [])
    ElMessage.success('回复成功')
  } catch {
    ElMessage.error('回复失败')
  }
}
</script>

<style scoped>
.comment-section {
  background: #fff;
  border-radius: 14px;
  border: 1px solid #f0f0f0;
  padding: 24px 28px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
  margin: 0 0 20px;
}

.title-icon {
  width: 20px;
  height: 20px;
  color: #86909c;
}

/* ===== 评论输入 ===== */
.comment-input {
  margin-bottom: 24px;
  background: #f7f8fa;
  border-radius: 10px;
  padding: 14px;
}

.comment-input :deep(.el-textarea__inner) {
  background: transparent;
  border: none;
  box-shadow: none;
  padding: 0;
}

.input-actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-top: 10px;
}

.char-count {
  font-size: 12px;
  color: #c0c4cc;
}

.login-hint {
  font-size: 13px;
  color: #86909c;
  margin-bottom: 20px;
  padding: 12px 16px;
  background: #f7f8fa;
  border-radius: 8px;
}

.login-link {
  color: #409eff;
  text-decoration: none;
  font-weight: 600;
}
.login-link:hover {
  text-decoration: underline;
}

/* ===== 评论项 ===== */
.comment-item {
  display: flex;
  gap: 12px;
  padding: 16px 0;
  border-bottom: 1px solid #f2f3f5;
}
.comment-item:last-child { border-bottom: none; }

.comment-left { flex-shrink: 0; }

.comment-avatar {
  background: linear-gradient(135deg, #36d1dc, #5b86e5);
  color: #fff;
  font-size: 13px;
  font-weight: 600;
}

.comment-right {
  flex: 1;
  min-width: 0;
}

.comment-top {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 6px;
}

.comment-user {
  font-size: 14px;
  font-weight: 600;
  color: #1d2129;
}

.comment-time {
  font-size: 12px;
  color: #c0c4cc;
}

.comment-body {
  font-size: 14px;
  color: #4e5969;
  line-height: 1.7;
  margin: 0 0 8px;
  word-break: break-word;
}

/* ===== 子评论 ===== */
.sub-comments {
  margin-top: 4px;
}

.toggle-btn {
  font-size: 12px;
  color: #909399;
  background: none;
  border: none;
  cursor: pointer;
  padding: 4px 0;
  transition: color 0.2s;
}
.toggle-btn:hover {
  color: #409eff;
}

.sub-list {
  margin-top: 8px;
  padding: 10px 14px;
  background: #f7f8fa;
  border-radius: 8px;
}

.sub-comment-item {
  font-size: 13px;
  color: #4e5969;
  padding: 6px 0;
  line-height: 1.6;
}
.sub-comment-item + .sub-comment-item {
  border-top: 1px dashed #e8e8e8;
}

.sub-user {
  font-weight: 600;
  color: #303133;
}

.reply-to {
  color: #86909c;
}
.reply-to b {
  color: #409eff;
  font-weight: 600;
}

.sub-body {
  color: #4e5969;
}

.sub-empty {
  font-size: 13px;
  color: #c0c4cc;
  text-align: center;
  padding: 8px 0;
}

.sub-reply-input {
  display: flex;
  gap: 8px;
  margin-top: 10px;
}

@media (max-width: 768px) {
  .comment-section { padding: 18px 16px; border-radius: 10px; }
  .comment-avatar { display: none; }
  .comment-item { gap: 0; }
}
</style>
