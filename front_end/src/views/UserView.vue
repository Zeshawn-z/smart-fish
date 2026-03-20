<template>
  <div class="user-center">
    <!-- 个人信息卡片 -->
    <div class="profile-card">
      <div class="profile-main">
        <div class="avatar-box" @click="handleAvatarClick">
          <img v-if="authStore.user?.avatar" :src="authStore.user.avatar" class="avatar-img" />
          <span v-else class="avatar-letter">{{ avatarLetter }}</span>
          <div class="avatar-overlay">
            <span>更换</span>
          </div>
          <input ref="avatarInputRef" type="file" accept="image/*" style="display: none" @change="handleAvatarUpload" />
        </div>
        <div class="profile-info">
          <h2 class="profile-name">{{ authStore.user?.username || '未登录' }}</h2>
          <p class="profile-meta">
            <span class="role-label">{{ roleText }}</span>
            <span class="meta-sep">·</span>
            <span>注册于 {{ registerDate }}</span>
          </p>
        </div>
      </div>
      <div class="profile-stats">
        <div class="stat-item">
          <span class="stat-num">{{ fishingStore.favoriteSpots.length }}</span>
          <span class="stat-label">收藏水域</span>
        </div>
        <div class="stat-item">
          <span class="stat-num">{{ recordStore.records.length }}</span>
          <span class="stat-label">垂钓记录</span>
        </div>
        <div class="stat-item">
          <span class="stat-num">{{ communityStore.myPosts.length }}</span>
          <span class="stat-label">发布帖子</span>
        </div>
      </div>
    </div>

    <!-- 标签页面板 -->
    <div class="main-panel">
      <el-tabs v-model="activeTab" class="user-tabs">
        <el-tab-pane label="个人信息" name="info">
          <ProfileForm />
        </el-tab-pane>

        <el-tab-pane label="修改密码" name="password">
          <PasswordForm />
        </el-tab-pane>

        <el-tab-pane name="favorites">
          <template #label>
            <span class="tab-label">
              我的收藏
              <em v-if="fishingStore.favoriteSpots.length > 0" class="tab-count">{{ fishingStore.favoriteSpots.length }}</em>
            </span>
          </template>
          <FavoriteList />
        </el-tab-pane>

        <el-tab-pane name="records">
          <template #label>
            <span class="tab-label">
              垂钓记录
              <em v-if="recordStore.records.length > 0" class="tab-count">{{ recordStore.records.length }}</em>
            </span>
          </template>
          <FishingRecordTab />
        </el-tab-pane>

        <el-tab-pane name="myposts">
          <template #label>
            <span class="tab-label">
              我的帖子
              <em v-if="communityStore.myPosts.length > 0" class="tab-count">{{ communityStore.myPosts.length }}</em>
            </span>
          </template>
          <MyPostsTab />
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '@/stores/auth'
import { useFishingStore } from '@/stores/fishing'
import { useCommunityStore, useFishingRecordStore } from '@/stores/community'
import { CommunityService } from '@/services/CommunityService'

import ProfileForm from './user/ProfileForm.vue'
import PasswordForm from './user/PasswordForm.vue'
import FavoriteList from './user/FavoriteList.vue'
import FishingRecordTab from './user/FishingRecordTab.vue'
import MyPostsTab from './user/MyPostsTab.vue'

const authStore = useAuthStore()
const fishingStore = useFishingStore()
const communityStore = useCommunityStore()
const recordStore = useFishingRecordStore()
const activeTab = ref('info')
const avatarInputRef = ref<HTMLInputElement | null>(null)

const avatarLetter = computed(() => {
  const name = authStore.user?.username
  if (!name) return '?'
  return name.charAt(0).toUpperCase()
})

const roleText = computed(() => {
  const map: Record<string, string> = { admin: '管理员', staff: '工作人员', user: '普通用户' }
  return map[authStore.user?.role || 'user'] || '普通用户'
})

const registerDate = computed(() => {
  const t = authStore.user?.register_time
  if (!t) return '--'
  return t.substring(0, 10)
})

onMounted(() => {
  fishingStore.fetchFavorites()
  if (authStore.user?.id) {
    communityStore.fetchMyPosts(authStore.user.id)
  }
  recordStore.fetchRecords()
})

function handleAvatarClick() {
  avatarInputRef.value?.click()
}

async function handleAvatarUpload(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  if (!file) return
  try {
    await CommunityService.uploadAvatarImage(file)
    // 重新获取用户信息以刷新头像
    await authStore.fetchUser()
    ElMessage.success('头像更新成功')
  } catch {
    ElMessage.error('头像上传失败')
  }
  // 重置 input 以允许再次选择同一文件
  target.value = ''
}
</script>

<style scoped>
.user-center {
  max-width: 860px;
  margin: 0 auto;
  padding: 24px;
}

/* ===== 个人信息卡片 ===== */
.profile-card {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff;
  border-radius: 10px;
  border: 1px solid #f0f0f0;
  padding: 22px 26px;
  margin-bottom: 20px;
}

.profile-main {
  display: flex;
  align-items: center;
  gap: 16px;
}

.avatar-box {
  width: 52px;
  height: 52px;
  border-radius: 12px;
  background: #e8f0fe;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  position: relative;
  cursor: pointer;
  overflow: hidden;
}

.avatar-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 11px;
  color: #fff;
  font-weight: 600;
  opacity: 0;
  transition: opacity 0.2s;
  border-radius: 12px;
}

.avatar-box:hover .avatar-overlay {
  opacity: 1;
}

.avatar-letter {
  font-size: 22px;
  font-weight: 700;
  color: #3b7dd8;
  line-height: 1;
}

.profile-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.profile-name {
  font-size: 18px;
  font-weight: 700;
  color: #1d2129;
  margin: 0;
  line-height: 1.3;
}

.profile-meta {
  font-size: 12.5px;
  color: #a0a4ad;
  margin: 0;
  display: flex;
  align-items: center;
}

.role-label {
  color: #5b8ff9;
  font-weight: 500;
}

.meta-sep {
  margin: 0 6px;
  color: #d9dbe0;
}

.profile-stats {
  display: flex;
  gap: 24px;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.stat-num {
  font-size: 20px;
  font-weight: 700;
  color: #1d2129;
  line-height: 1.2;
}

.stat-label {
  font-size: 12px;
  color: #a0a4ad;
}

/* ===== 主面板 ===== */
.main-panel {
  background: #fff;
  border-radius: 10px;
  border: 1px solid #f0f0f0;
  padding: 20px 24px;
}

/* ===== Tabs ===== */
.user-tabs :deep(.el-tabs__header) {
  margin-bottom: 0;
}

.user-tabs :deep(.el-tabs__nav-wrap::after) {
  height: 1px;
  background-color: #ebeef5;
}

.user-tabs :deep(.el-tabs__item) {
  height: 44px;
  line-height: 44px;
  font-size: 14px;
  color: #86909c;
}

.user-tabs :deep(.el-tabs__item.is-active) {
  color: #1d2129;
  font-weight: 600;
}

.user-tabs :deep(.el-tabs__content) {
  padding-top: 20px;
}

.tab-label {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.tab-count {
  font-style: normal;
  font-size: 11px;
  min-width: 18px;
  height: 18px;
  line-height: 18px;
  text-align: center;
  border-radius: 9px;
  background-color: #5b8ff9;
  color: #fff;
  padding: 0 5px;
}

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .user-center {
    padding: 14px;
  }

  .profile-card {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
    padding: 18px 16px;
  }

  .avatar-box {
    width: 44px;
    height: 44px;
    border-radius: 10px;
  }

  .avatar-letter {
    font-size: 18px;
  }

  .profile-name {
    font-size: 16px;
  }

  .profile-meta {
    font-size: 11.5px;
  }

  .profile-stats {
    padding-left: 60px;
  }

  .main-panel {
    padding: 14px 12px;
  }
}
</style>
