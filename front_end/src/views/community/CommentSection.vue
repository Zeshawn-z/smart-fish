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
        <el-avatar :size="32" :src="commentAvatar(comment).src" :style="commentAvatar(comment).hasAvatar ? {} : commentAvatar(comment).style" class="comment-avatar">
          {{ commentAvatar(comment).hasAvatar ? '' : commentAvatar(comment).letter }}
        </el-avatar>
      </div>
      <div class="comment-right">
        <div class="comment-top">
          <span class="comment-user">{{ comment.username || `用户 #${comment.user_id}` }}</span>
          <span class="comment-time">{{ formatTime(comment.created_at) }}</span>
          <!-- 评论操作按钮 -->
          <div class="comment-actions">
            <button
              class="action-btn"
              :class="{ active: commentLikedMap[comment.id] }"
              @click="handleLikeComment(comment.id)"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="action-icon">
                <path d="M14 9V5a3 3 0 0 0-3-3l-4 9v11h11.28a2 2 0 0 0 2-1.7l1.38-9a2 2 0 0 0-2-2.3H14z"/>
                <path d="M7 22H4a2 2 0 0 1-2-2v-7a2 2 0 0 1 2-2h3"/>
              </svg>
              <span v-if="commentLikesMap[comment.id]" class="action-count">{{ commentLikesMap[comment.id] }}</span>
            </button>
            <button class="action-btn" @click="startReply(comment.id)">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" class="action-icon">
                <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"/>
              </svg>
            </button>
          </div>
        </div>
        <p class="comment-body">{{ comment.body }}</p>

        <!-- 子评论区（楼中楼）—— 自动加载，有则显示 -->
        <div v-if="getSubCommentsList(comment.id).length > 0" class="sub-comments">
          <div class="sub-list">
            <div
              v-for="coc in getVisibleSubComments(comment.id)"
              :key="coc.coc_id"
              class="sub-comment-item"
            >
              <span class="sub-user">{{ coc.username || `用户 #${coc.user_id}` }}</span>
              <span v-if="coc.to_username" class="reply-to"> 回复 <b>{{ coc.to_username }}</b></span>
              <span class="sub-body">：{{ coc.body }}</span>
            </div>
          </div>
          <!-- 超过2条时显示展开/收起按钮 -->
          <button
            v-if="getSubCommentsList(comment.id).length > 2"
            class="toggle-btn"
            @click="toggleExpand(comment.id)"
          >
            {{ expandedComments.has(comment.id)
              ? '收起回复 ▲'
              : `更多回复 (${getSubCommentsList(comment.id).length - 2}) ▼` }}
          </button>
        </div>

        <!-- 回复输入框（点击回复按钮后显示） -->
        <div v-if="isLoggedIn && replyingTo === comment.id" class="sub-reply-input">
          <el-input
            ref="replyInputRef"
            v-model="subReplyMap[comment.id]"
            size="small"
            placeholder="回复..."
            @keyup.enter="handleSubComment(comment.id)"
          />
          <el-button size="small" type="primary" round @click="handleSubComment(comment.id)">回复</el-button>
          <el-button size="small" round @click="replyingTo = null">取消</el-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import { CommunityService } from '@/services/CommunityService'
import { useAvatar } from '@/composables/useAvatar'
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

/** 当前正在回复的评论 ID */
const replyingTo = ref<number | null>(null)
const replyInputRef = ref<InstanceType<any> | null>(null)

/** 评论点赞数和当前用户是否已赞 */
const commentLikesMap = reactive<Record<number, number>>({})
const commentLikedMap = reactive<Record<number, boolean>>({})

/** 获取评论用户的头像信息 */
function commentAvatar(comment: CommentType) {
  return useAvatar(comment.avatar, comment.username)
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

/** 获取某评论的所有子评论 */
function getSubCommentsList(commentId: number): CommentOnComment[] {
  return subCommentsMap.get(commentId) || []
}

/** 获取可见的子评论（默认2条，展开后全部） */
function getVisibleSubComments(commentId: number): CommentOnComment[] {
  const all = getSubCommentsList(commentId)
  if (expandedComments.has(commentId)) return all
  return all.slice(0, 2)
}

/** 展开/收起更多回复 */
function toggleExpand(commentId: number) {
  if (expandedComments.has(commentId)) {
    expandedComments.delete(commentId)
  } else {
    expandedComments.add(commentId)
  }
}

/** 点击回复按钮 —— 切换回复输入框 */
function startReply(commentId: number) {
  if (!props.isLoggedIn) {
    ElMessage.warning('请先登录后再回复')
    return
  }
  if (replyingTo.value === commentId) {
    replyingTo.value = null
  } else {
    replyingTo.value = commentId
    nextTick(() => {
      replyInputRef.value?.focus?.()
    })
  }
}

/** 从嵌入数据初始化子评论和点赞（不再逐个请求） */
function initFromEmbeddedData() {
  for (const comment of props.comments) {
    // 使用后端嵌入的子评论
    if (comment.sub_comments) {
      subCommentsMap.set(comment.id, comment.sub_comments)
    }
    // 使用后端嵌入的点赞数
    if (comment.likes !== undefined) {
      commentLikesMap[comment.id] = comment.likes
    }
  }
}

// 评论列表变化时从嵌入数据初始化
watch(() => props.comments, () => {
  initFromEmbeddedData()
}, { immediate: true })

/** 发表评论 */
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

/** 评论点赞/取消 */
async function handleLikeComment(commentId: number) {
  if (!props.isLoggedIn) {
    ElMessage.warning('请先登录')
    return
  }
  try {
    if (commentLikedMap[commentId]) {
      await CommunityService.unlikeComment(commentId)
      commentLikedMap[commentId] = false
      commentLikesMap[commentId] = Math.max(0, (commentLikesMap[commentId] || 0) - 1)
    } else {
      await CommunityService.likeComment(commentId)
      commentLikedMap[commentId] = true
      commentLikesMap[commentId] = (commentLikesMap[commentId] || 0) + 1
    }
  } catch {
    ElMessage.error('操作失败')
  }
}

/** 发表子评论 */
async function handleSubComment(commentId: number) {
  const body = subReplyMap[commentId]?.trim()
  if (!body) return
  try {
    await CommunityService.createSubComment(commentId, body)
    subReplyMap[commentId] = ''
    replyingTo.value = null
    ElMessage.success('回复成功')
    // 通知父组件刷新（会重新拉取帖子详情，包含最新的评论数据）
    emit('refresh')
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

/* ===== 评论操作按钮 ===== */
.comment-actions {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-left: auto;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  background: none;
  border: none;
  cursor: pointer;
  color: #909399;
  padding: 2px 6px;
  border-radius: 4px;
  transition: all 0.2s;
  font-size: 12px;
  line-height: 1;
}

.action-btn:hover {
  color: #409eff;
  background: #ecf5ff;
}

.action-btn.active {
  color: #409eff;
}

.action-btn.active .action-icon {
  fill: #409eff;
  stroke: #409eff;
}

.action-icon {
  width: 16px;
  height: 16px;
}

.action-count {
  font-size: 12px;
  font-weight: 500;
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

.sub-list {
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

.toggle-btn {
  font-size: 12px;
  color: #909399;
  background: none;
  border: none;
  cursor: pointer;
  padding: 6px 0 2px;
  transition: color 0.2s;
}
.toggle-btn:hover {
  color: #409eff;
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
