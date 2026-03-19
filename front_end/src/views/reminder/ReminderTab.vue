<template>
  <div class="filter-bar">
    <div class="filter-left">
      <el-select v-model="level" placeholder="级别" clearable size="default" style="width: 110px" @change="reload">
        <el-option label="信息" :value="0" />
        <el-option label="提示" :value="1" />
        <el-option label="重要" :value="2" />
        <el-option label="紧急" :value="3" />
      </el-select>
      <el-select v-model="resolved" placeholder="状态" clearable size="default" style="width: 110px" @change="reload">
        <el-option label="未处理" value="false" />
        <el-option label="已处理" value="true" />
      </el-select>
    </div>
    <div class="filter-right">
      <el-button-group>
        <el-button :type="viewMode === 'table' ? 'primary' : 'default'" size="default" @click="viewMode = 'table'">
          <el-icon><List /></el-icon>
        </el-button>
        <el-button :type="viewMode === 'card' ? 'primary' : 'default'" size="default" @click="viewMode = 'card'">
          <el-icon><Grid /></el-icon>
        </el-button>
      </el-button-group>
    </div>
  </div>

  <!-- 表格视图 -->
  <div v-if="viewMode === 'table'" class="data-section">
    <el-table :data="reminderStore.reminders" v-loading="reminderStore.isLoading" stripe empty-text="暂无提醒" class="data-table">
      <el-table-column label="级别" width="80" align="center">
        <template #default="{ row }">
          <span class="level-dot" :class="'level-' + row.level"></span>
          <span class="level-text">{{ REMINDER_LEVEL_MAP[row.level as ReminderLevel] }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="reminder_type" label="类型" width="100" />
      <el-table-column prop="message" label="内容" min-width="260" show-overflow-tooltip />
      <el-table-column label="时间" width="170">
        <template #default="{ row }">
          <span class="time-text">{{ formatTime(row.timestamp) }}</span>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="90" align="center">
        <template #default="{ row }">
          <el-tag :type="row.resolved ? 'success' : 'warning'" size="small" effect="plain">
            {{ row.resolved ? '已处理' : '待处理' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100" align="center">
        <template #default="{ row }">
          <el-button v-if="!row.resolved && authStore.isStaff" size="small" type="success" plain @click="handleResolve(row.id)">
            处理
          </el-button>
        </template>
      </el-table-column>
    </el-table>
  </div>

  <!-- 卡片视图 -->
  <div v-else class="data-section">
    <div v-if="reminderStore.isLoading" v-loading="true" style="height: 200px"></div>
    <div v-else-if="reminderStore.reminders.length === 0" class="empty-box">
      <el-empty description="暂无提醒" :image-size="80" />
    </div>
    <div v-else class="card-grid">
      <div v-for="r in reminderStore.reminders" :key="r.id" class="info-card" :class="{ 'is-resolved': r.resolved }">
        <div class="card-top">
          <el-tag :type="levelTagType(r.level)" size="small" effect="plain" round>
            {{ REMINDER_LEVEL_MAP[r.level as ReminderLevel] }}
          </el-tag>
          <el-tag size="small" effect="plain" type="info" round>{{ r.reminder_type }}</el-tag>
          <span class="card-status" :class="r.resolved ? 'resolved' : 'pending'">
            {{ r.resolved ? '已处理' : '待处理' }}
          </span>
        </div>
        <p class="card-body">{{ r.message }}</p>
        <div class="card-bottom">
          <span class="card-time">{{ formatTime(r.timestamp) }}</span>
          <el-button v-if="!r.resolved && authStore.isStaff" size="small" type="success" plain @click="handleResolve(r.id)">
            标记处理
          </el-button>
        </div>
      </div>
    </div>
  </div>

  <!-- 分页 -->
  <div class="pagination-bar" v-if="reminderStore.total > 20">
    <el-pagination
      v-model:current-page="currentPage"
      :page-size="20"
      :total="reminderStore.total"
      layout="total, prev, pager, next"
      background
      @current-change="reload"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useReminderStore } from '@/stores/reminder'
import { REMINDER_LEVEL_MAP, type ReminderLevel } from '@/types'
import { List, Grid } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { formatTime } from './utils.ts'

const authStore = useAuthStore()
const reminderStore = useReminderStore()

const viewMode = ref<'table' | 'card'>('card')
const currentPage = ref(1)
const level = ref<number | ''>('')
const resolved = ref('')

function reload() {
  const params: Record<string, unknown> = { page: currentPage.value, page_size: 20 }
  if (level.value !== '') params.level = level.value
  if (resolved.value) params.resolved = resolved.value
  reminderStore.fetchReminders(params as any)
}

function levelTagType(level: number) {
  if (level >= 2) return 'danger'
  if (level === 1) return 'warning'
  return 'info'
}

async function handleResolve(id: number) {
  await reminderStore.resolveReminder(id)
  ElMessage.success('已标记为已处理')
}

defineExpose({ reload })
</script>

<style scoped>
.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 0 12px;
  gap: 12px;
  flex-wrap: wrap;
}
.filter-left {
  display: flex;
  gap: 8px;
  align-items: center;
}
.filter-right {
  display: flex;
  align-items: center;
}
.data-section {
  min-height: 200px;
}
.data-table {
  border-radius: 6px;
  overflow: hidden;
}
.data-table :deep(.el-table__header th) {
  background-color: #fafafa;
  color: #606266;
  font-weight: 600;
  font-size: 13px;
}
.level-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-right: 6px;
  vertical-align: middle;
}
.level-0 { background-color: #909399; }
.level-1 { background-color: #e6a23c; }
.level-2 { background-color: #f56c6c; }
.level-3 { background-color: #f56c6c; box-shadow: 0 0 0 3px rgba(245, 108, 108, 0.2); }
.level-text {
  font-size: 13px;
  vertical-align: middle;
}
.time-text {
  font-size: 13px;
  color: #909399;
}

/* 卡片网格 */
.card-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 12px;
}
.info-card {
  background: #fff;
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 14px 16px;
}
.info-card.is-resolved {
  opacity: 0.5;
}
.card-top {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 10px;
  flex-wrap: wrap;
}
.card-status {
  margin-left: auto;
  font-size: 12px;
  font-weight: 500;
}
.card-status.pending { color: #e6a23c; }
.card-status.resolved { color: #67c23a; }
.card-body {
  font-size: 14px;
  color: #4a4a4a;
  line-height: 1.6;
  margin: 0 0 10px;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.card-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 10px;
  border-top: 1px solid #f5f5f5;
}
.card-time {
  font-size: 12px;
  color: #b0b3ba;
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
  .filter-bar {
    flex-direction: column;
    align-items: stretch;
  }
  .filter-left,
  .filter-right {
    width: 100%;
    justify-content: space-between;
  }
  .card-grid {
    grid-template-columns: 1fr;
  }
}
</style>
