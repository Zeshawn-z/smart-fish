<template>
  <div class="info-center">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="title-row">
        <span class="title-accent"></span>
        <h2 class="page-title">信息中心</h2>
        <el-button text @click="refreshData" class="refresh-btn">
          <el-icon :class="{ spinning: isRefreshing }"><Refresh /></el-icon>
        </el-button>
      </div>
      <div class="header-stats">
        <span v-if="reminderStore.unresolvedCount > 0" class="stat-item warn">
          {{ reminderStore.unresolvedCount }} 条待处理提醒
        </span>
        <span v-else class="stat-item">暂无待处理提醒</span>
      </div>
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
  margin-bottom: 20px;
}
.title-row {
  display: flex;
  align-items: center;
  gap: 10px;
}
.title-accent {
  display: block;
  width: 4px;
  height: 22px;
  border-radius: 2px;
  background: #2b6e3f;
  flex-shrink: 0;
}
.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #1d2129;
  margin: 0;
}
.refresh-btn {
  margin-left: auto;
  color: #909399;
  font-size: 17px;
  padding: 6px;
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
.header-stats {
  margin-top: 6px;
  padding-left: 14px;
}
.stat-item {
  font-size: 13px;
  color: #a0a4ad;
}
.stat-item.warn {
  color: #cf7a30;
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
  .main-panel {
    padding: 14px 12px;
  }
  .page-title {
    font-size: 18px;
  }
  .title-accent {
    height: 18px;
  }
}
</style>
