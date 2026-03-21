<template>
  <div class="stats-panel">
    <div class="tab-header">
      <h3>垂钓数据统计</h3>
    </div>

    <div v-loading="isLoading" class="stats-content">
      <el-empty v-if="!isLoading && !stats" description="暂无垂钓统计数据" />

      <template v-if="stats">
        <!-- 核心指标卡片 -->
        <div class="stat-cards">
          <div class="stat-card">
            <span class="stat-icon">🎣</span>
            <div class="stat-data">
              <span class="stat-value">{{ stats.total_trips }}</span>
              <span class="stat-name">出钓次数</span>
            </div>
          </div>
          <div class="stat-card">
            <span class="stat-icon">🐟</span>
            <div class="stat-data">
              <span class="stat-value">{{ stats.total_fish }}</span>
              <span class="stat-name">总渔获数</span>
            </div>
          </div>
          <div class="stat-card">
            <span class="stat-icon">⚖️</span>
            <div class="stat-data">
              <span class="stat-value">{{ stats.total_kg.toFixed(1) }}<small>kg</small></span>
              <span class="stat-name">总重量</span>
            </div>
          </div>
          <div class="stat-card">
            <span class="stat-icon">🏆</span>
            <div class="stat-data">
              <span class="stat-value">{{ stats.max_kg.toFixed(1) }}<small>kg</small></span>
              <span class="stat-name">最大单条</span>
            </div>
          </div>
          <div class="stat-card">
            <span class="stat-icon">⏱️</span>
            <div class="stat-data">
              <span class="stat-value">{{ stats.total_hours.toFixed(1) }}<small>h</small></span>
              <span class="stat-name">总时长</span>
            </div>
          </div>
        </div>

        <!-- 鱼种分布 -->
        <div v-if="stats.fish_types && stats.fish_types.length > 0" class="fish-types-section">
          <h4 class="section-title">鱼种分布 Top 5</h4>
          <div class="fish-type-bars">
            <div v-for="item in stats.fish_types" :key="item.fish_type" class="fish-type-row">
              <span class="fish-type-name">{{ item.fish_type || '未知' }}</span>
              <div class="fish-type-bar-bg">
                <div
                  class="fish-type-bar"
                  :style="{ width: getBarWidth(item.count) + '%' }"
                />
              </div>
              <span class="fish-type-count">{{ item.count }} 条</span>
            </div>
          </div>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { FishingRecordService } from '@/services/FishingRecordService'

interface FishingStats {
  total_trips: number
  total_fish: number
  total_kg: number
  max_kg: number
  total_hours: number
  fish_types: { fish_type: string; count: number }[]
}

const isLoading = ref(false)
const stats = ref<FishingStats | null>(null)

const maxCount = computed(() => {
  if (!stats.value?.fish_types?.length) return 1
  return Math.max(...stats.value.fish_types.map(t => t.count), 1)
})

function getBarWidth(count: number) {
  return Math.round((count / maxCount.value) * 100)
}

onMounted(async () => {
  isLoading.value = true
  try {
    stats.value = await FishingRecordService.getMyStats()
  } catch {
    stats.value = null
  } finally {
    isLoading.value = false
  }
})
</script>

<style scoped>
.stats-panel {
  min-height: 200px;
}

.tab-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.tab-header h3 {
  margin: 0;
  font-size: 16px;
  font-weight: 600;
  color: #1d2129;
}

/* ===== Stat Cards ===== */
.stat-cards {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 12px;
  margin-bottom: 24px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px;
  background: #fafbfc;
  border-radius: 10px;
  border: 1px solid #f0f0f0;
  transition: box-shadow 0.2s;
}

.stat-card:hover {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

.stat-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.stat-data {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 20px;
  font-weight: 700;
  color: #1d2129;
  line-height: 1.2;
}

.stat-value small {
  font-size: 12px;
  font-weight: 400;
  color: #86909c;
  margin-left: 2px;
}

.stat-name {
  font-size: 12px;
  color: #a0a4ad;
  margin-top: 2px;
}

/* ===== Fish Types ===== */
.fish-types-section {
  padding-top: 16px;
  border-top: 1px solid #f0f0f0;
}

.section-title {
  font-size: 14px;
  font-weight: 600;
  color: #1d2129;
  margin: 0 0 12px;
}

.fish-type-bars {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.fish-type-row {
  display: flex;
  align-items: center;
  gap: 10px;
}

.fish-type-name {
  width: 64px;
  font-size: 13px;
  color: #4e5969;
  text-align: right;
  flex-shrink: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.fish-type-bar-bg {
  flex: 1;
  height: 18px;
  background: #f2f3f5;
  border-radius: 9px;
  overflow: hidden;
}

.fish-type-bar {
  height: 100%;
  border-radius: 9px;
  background: linear-gradient(90deg, #5b8ff9, #61c0bf);
  transition: width 0.5s ease;
  min-width: 4px;
}

.fish-type-count {
  width: 48px;
  font-size: 13px;
  color: #86909c;
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .stat-cards {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
