<template>
  <div class="info-center">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-left">
        <div class="title-icon">
          <svg viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M14 3C8 3 3 7 3 12c0 3 1.5 5.5 4 7.5V24l4-2.5c1 .3 2 .5 3 .5 6 0 11-4 11-9S20 3 14 3z" stroke="currentColor" stroke-width="1.8" fill="none"/>
            <circle cx="9" cy="12" r="1.2" fill="currentColor"/>
            <circle cx="14" cy="12" r="1.2" fill="currentColor"/>
            <circle cx="19" cy="12" r="1.2" fill="currentColor"/>
          </svg>
        </div>
        <div class="title-text">
          <h2 class="page-title">信息中心</h2>
          <p class="page-desc">
            <span v-if="reminderStore.unresolvedCount > 0" class="stat-warn">
              {{ reminderStore.unresolvedCount }} 条待处理
            </span>
            <span v-else>全部已处理</span>
            <span class="stat-sep">·</span>
            <span>提醒、通知与垂钓建议</span>
          </p>
        </div>
      </div>
      <el-button text @click="refreshData" class="refresh-btn">
        <el-icon :class="{ spinning: isRefreshing }"><Refresh /></el-icon>
      </el-button>
    </div>

    <!-- 主体内容 -->
    <div class="main-panel">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange" class="info-tabs">
        <el-tab-pane name="reminders">
          <template #label>
            <span class="tab-label">
              提醒
              <em v-if="reminderStore.unresolvedCount > 0" class="tab-count">{{ reminderStore.unresolvedCount }}</em>
            </span>
          </template>
          <ReminderTab ref="reminderTabRef" />
        </el-tab-pane>

        <el-tab-pane name="notices">
          <template #label>
            <span class="tab-label">通知</span>
          </template>
          <NoticeTab ref="noticeTabRef" />
        </el-tab-pane>

        <el-tab-pane name="suggestions">
          <template #label>
            <span class="tab-label">垂钓建议</span>
          </template>
          <SuggestionTab ref="suggestionTabRef" />
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useReminderStore } from '@/stores/reminder'
import { Refresh } from '@element-plus/icons-vue'
import ReminderTab from './reminder/ReminderTab.vue'
import NoticeTab from './reminder/NoticeTab.vue'
import SuggestionTab from './reminder/SuggestionTab.vue'

const route = useRoute()
const router = useRouter()
const reminderStore = useReminderStore()

const activeTab = ref('reminders')
const reminderTabRef = ref<InstanceType<typeof ReminderTab> | null>(null)
const noticeTabRef = ref<InstanceType<typeof NoticeTab> | null>(null)
const suggestionTabRef = ref<InstanceType<typeof SuggestionTab> | null>(null)

const isRefreshing = ref(false)
let refreshTimer: ReturnType<typeof setInterval> | null = null

function refreshData() {
  isRefreshing.value = true
  setTimeout(() => { isRefreshing.value = false }, 600)
  if (activeTab.value === 'reminders') reminderTabRef.value?.reload()
  else if (activeTab.value === 'notices') noticeTabRef.value?.reload()
  else suggestionTabRef.value?.reload()
}

function handleTabChange(tab: string | number) {
  router.replace({ query: { ...route.query, tab: String(tab) } })
  if (tab === 'notices') noticeTabRef.value?.reload()
  if (tab === 'suggestions') suggestionTabRef.value?.reload()
}

onMounted(() => {
  const tabParam = route.query.tab as string
  if (tabParam && ['reminders', 'notices', 'suggestions'].includes(tabParam)) {
    activeTab.value = tabParam
  }

  // 首次加载三个 tab 的数据
  reminderTabRef.value?.reload()
  noticeTabRef.value?.reload()
  suggestionTabRef.value?.reload()

  refreshTimer = setInterval(refreshData, 30000)
})

onBeforeUnmount(() => {
  if (refreshTimer) clearInterval(refreshTimer)
})
</script>

<style scoped>
.info-center {
  max-width: 1100px;
  margin: 0 auto;
  padding: 24px;
}

/* ===== 标题区 ===== */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  background: #fff;
  border-radius: 10px;
  border: 1px solid #f0f0f0;
  padding: 18px 22px;
}
.header-left {
  display: flex;
  align-items: center;
  gap: 14px;
}
.title-icon {
  width: 42px;
  height: 42px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #eef6fb;
  border-radius: 10px;
  color: #2196f3;
  flex-shrink: 0;
}
.title-icon svg {
  width: 22px;
  height: 22px;
}
.title-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.page-title {
  font-size: 18px;
  font-weight: 700;
  color: #1d2129;
  margin: 0;
  line-height: 1.3;
}
.page-desc {
  font-size: 12.5px;
  color: #a0a4ad;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 0;
}
.stat-warn {
  color: #cf7a30;
  font-weight: 500;
}
.stat-sep {
  margin: 0 6px;
  color: #d9dbe0;
}
.refresh-btn {
  color: #b0b3ba;
  font-size: 17px;
  padding: 8px;
}
.refresh-btn:hover {
  color: #606266;
}
.spinning {
  animation: spin 0.6s ease;
}
@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* ===== 主面板 ===== */
.main-panel {
  background: #fff;
  border-radius: 10px;
  border: 1px solid #f0f0f0;
  padding: 20px 24px;
}

/* ===== Tabs ===== */
.info-tabs :deep(.el-tabs__header) {
  margin-bottom: 0;
}
.info-tabs :deep(.el-tabs__nav-wrap::after) {
  height: 1px;
  background-color: #ebeef5;
}
.info-tabs :deep(.el-tabs__item) {
  height: 44px;
  line-height: 44px;
  font-size: 14px;
  color: #86909c;
}
.info-tabs :deep(.el-tabs__item.is-active) {
  color: #1d2129;
  font-weight: 600;
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
  background-color: #f56c6c;
  color: #fff;
  padding: 0 5px;
}

@media (max-width: 768px) {
  .info-center {
    padding: 14px;
  }
  .page-header {
    padding: 14px 16px;
  }
  .title-icon {
    width: 36px;
    height: 36px;
    border-radius: 8px;
  }
  .title-icon svg {
    width: 18px;
    height: 18px;
  }
  .page-title {
    font-size: 16px;
  }
  .page-desc {
    font-size: 11.5px;
  }
  .main-panel {
    padding: 14px 12px;
  }
}
</style>
