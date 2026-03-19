<template>
  <div class="filter-bar">
    <div class="filter-left">
      <el-input v-model="search" placeholder="搜索通知..." clearable style="width: 220px" @input="debounceReload">
        <template #prefix><el-icon><Search /></el-icon></template>
      </el-input>
    </div>
  </div>

  <div class="data-section">
    <div v-if="noticeStore.isLoading" v-loading="true" style="height: 200px"></div>
    <div v-else-if="noticeStore.notices.length === 0" class="empty-box">
      <el-empty description="暂无通知" :image-size="80" />
    </div>
    <div v-else class="notice-list">
      <div
        v-for="notice in noticeStore.notices"
        :key="notice.id"
        class="notice-item"
        :class="{ 'is-outdated': notice.outdated }"
        @click="expandedId = expandedId === notice.id ? null : notice.id"
      >
        <div class="notice-header">
          <h4 class="notice-title">{{ notice.title }}</h4>
          <span class="notice-time">{{ formatTime(notice.timestamp) }}</span>
          <el-tag v-if="notice.outdated" size="small" type="info" effect="plain">已过期</el-tag>
        </div>
        <div v-if="expandedId === notice.id" class="notice-content">
          <p>{{ notice.content }}</p>
          <div v-if="notice.related_spots && notice.related_spots.length > 0" class="notice-spots">
            <span class="spots-label">关联水域：</span>
            <el-tag v-for="spot in notice.related_spots" :key="spot.id" size="small" effect="plain" type="info" class="spot-tag">
              {{ spot.name }}
            </el-tag>
          </div>
        </div>
        <p v-else class="notice-preview">{{ notice.content }}</p>
      </div>
    </div>
  </div>

  <div class="pagination-bar" v-if="noticeStore.total > 20">
    <el-pagination
      v-model:current-page="currentPage"
      :page-size="20"
      :total="noticeStore.total"
      layout="total, prev, pager, next"
      background
      @current-change="reload"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useNoticeStore } from '@/stores/notice'
import { Search } from '@element-plus/icons-vue'
import { formatTime } from './utils.ts'

const noticeStore = useNoticeStore()

const currentPage = ref(1)
const search = ref('')
const expandedId = ref<number | null>(null)
let debounceTimer: ReturnType<typeof setTimeout> | null = null

function reload() {
  const params: Record<string, unknown> = { page: currentPage.value, page_size: 20 }
  if (search.value) params.search = search.value
  noticeStore.fetchNotices(params as any)
}

function debounceReload() {
  if (debounceTimer) clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    currentPage.value = 1
    reload()
  }, 300)
}

defineExpose({ reload })
</script>

<style scoped>
.filter-bar {
  display: flex;
  align-items: center;
  padding: 16px 0 12px;
}
.filter-left {
  display: flex;
  gap: 8px;
  align-items: center;
}
.data-section {
  min-height: 200px;
}
.notice-list {
  display: flex;
  flex-direction: column;
  gap: 1px;
  background: #f0f0f0;
  border-radius: 8px;
  overflow: hidden;
}
.notice-item {
  background: #fff;
  padding: 14px 18px;
  cursor: pointer;
  transition: background-color 0.15s;
}
.notice-item:active {
  background-color: #fafafa;
}
.notice-item.is-outdated {
  opacity: 0.55;
}
.notice-header {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
.notice-title {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
  margin: 0;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.notice-time {
  font-size: 12px;
  color: #b0b3ba;
  flex-shrink: 0;
}
.notice-preview {
  font-size: 13px;
  color: #909399;
  margin: 6px 0 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
.notice-content {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px dashed #ebeef5;
}
.notice-content p {
  font-size: 14px;
  color: #4a4a4a;
  line-height: 1.7;
  margin: 0 0 8px;
  white-space: pre-wrap;
}
.notice-spots {
  display: flex;
  align-items: center;
  gap: 6px;
  flex-wrap: wrap;
  margin-top: 6px;
}
.spots-label {
  font-size: 12px;
  color: #909399;
}
.spot-tag {
  font-size: 11px;
}
.pagination-bar {
  display: flex;
  justify-content: center;
  padding: 18px 0 4px;
}
.empty-box {
  padding: 40px 0;
}

@media (max-width: 768px) {
  .notice-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
}
</style>
