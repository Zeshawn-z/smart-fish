<template>
  <el-card class="dashboard-card">
    <template #header>
      <div class="card-header">
        <el-icon class="card-icon suggest-icon"><Promotion /></el-icon>
        <span class="card-title">智能垂钓建议</span>
      </div>
    </template>
    <div v-loading="loading">
      <div v-if="suggestions.length > 0" class="suggestion-list">
        <div v-for="s in suggestions" :key="s.id" class="suggestion-item">
          <div class="suggestion-header">
            <el-tag type="success" effect="dark" size="small" round>
              评分 {{ s.score }}
            </el-tag>
            <span class="suggestion-spot">{{ s.fishing_spot?.name }}</span>
          </div>
          <p class="suggestion-text">{{ s.suggestion_text }}</p>
        </div>
      </div>
      <el-empty v-else description="暂无垂钓建议" :image-size="60" />
    </div>
  </el-card>
</template>

<script setup lang="ts">
import type { FishingSuggestion } from '@/types'
import { Promotion } from '@element-plus/icons-vue'

defineProps<{
  suggestions: FishingSuggestion[]
  loading: boolean
}>()
</script>

<style scoped>
.dashboard-card {
  margin-bottom: 25px;
  border-radius: 12px !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08) !important;
  border: none;
  transition: all 0.3s;
}

.dashboard-card:hover {
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1) !important;
}

.dashboard-card :deep(.el-card__header) {
  padding: 18px 24px;
  background: linear-gradient(45deg, #fafafa, #f6f9ff) !important;
  border-bottom: 1px solid #e4e7ed;
}

.card-header {
  display: flex;
  align-items: center;
}

.card-icon {
  margin-right: 8px;
  font-size: 20px;
  padding: 3px;
  border-radius: 8px;
  color: white;
}

.suggest-icon { background-color: #67c23a; }

.card-title {
  font-size: 17px !important;
  font-weight: 600;
  color: #333;
  letter-spacing: 0.5px;
}

.suggestion-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.suggestion-item {
  padding: 12px 16px;
  border-radius: 8px;
  background: #f9fbfd;
  border: 1px solid #ebeef5;
  transition: all 0.2s;
}

.suggestion-item:hover {
  border-color: #c6e2ff;
  box-shadow: 0 2px 8px rgba(64, 158, 255, 0.1);
}

.suggestion-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
}

.suggestion-spot {
  font-weight: 600;
  font-size: 14px;
  color: #303133;
}

.suggestion-text {
  font-size: 13px;
  color: #606266;
  line-height: 1.6;
  margin: 0;
}
</style>
