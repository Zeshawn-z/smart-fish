<template>
  <div class="data-section">
    <div v-if="suggestionStore.isLoading" v-loading="true" style="height: 200px"></div>
    <div v-else-if="suggestionStore.suggestions.length === 0" class="empty-box">
      <el-empty description="暂无垂钓建议" :image-size="80" />
    </div>
    <div v-else class="suggestion-list">
      <div v-for="s in suggestionStore.suggestions" :key="s.id" class="suggestion-item">
        <div class="suggestion-header">
          <div class="suggestion-score" :class="getScoreClass(s.score)">{{ s.score }}分</div>
          <span v-if="s.fishing_spot" class="suggestion-spot">{{ s.fishing_spot.name }}</span>
          <span class="suggestion-time">{{ formatTime(s.timestamp) }}</span>
        </div>
        <p class="suggestion-text">{{ s.suggestion_text }}</p>
      </div>
    </div>
  </div>

  <div class="pagination-bar" v-if="suggestionStore.total > 20">
    <el-pagination
      v-model:current-page="currentPage"
      :page-size="20"
      :total="suggestionStore.total"
      layout="total, prev, pager, next"
      background
      @current-change="reload"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useSuggestionStore } from '@/stores/suggestion'
import { formatTime, getScoreClass } from './utils.ts'

const suggestionStore = useSuggestionStore()
const currentPage = ref(1)

function reload() {
  suggestionStore.fetchSuggestions({ page: currentPage.value, page_size: 20 })
}

defineExpose({ reload })
</script>

<style scoped>
.data-section {
  min-height: 200px;
  padding-top: 12px;
}
.suggestion-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
.suggestion-item {
  background: #fff;
  border: 1px solid #eee;
  border-radius: 8px;
  padding: 14px 18px;
}
.suggestion-header {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 8px;
  flex-wrap: wrap;
}
.suggestion-score {
  font-size: 12px;
  font-weight: 600;
  padding: 2px 10px;
  border-radius: 10px;
}
.suggestion-score.high {
  background-color: #e8f5e9;
  color: #43a047;
}
.suggestion-score.mid {
  background-color: #fff8e1;
  color: #f57c00;
}
.suggestion-score.low {
  background-color: #f5f5f5;
  color: #909399;
}
.suggestion-spot {
  font-size: 13px;
  color: #606266;
  font-weight: 500;
}
.suggestion-time {
  font-size: 12px;
  color: #b0b3ba;
  margin-left: auto;
}
.suggestion-text {
  font-size: 14px;
  color: #4a4a4a;
  line-height: 1.7;
  margin: 0;
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
  .suggestion-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
  .suggestion-time {
    margin-left: 0;
  }
}
</style>
