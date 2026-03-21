<template>
  <div class="dash-header">
    <div class="header-left">
      <img class="logo-icon" src="" alt="" style="display:none" />
      <h1 class="title">
        <span class="title-icon">🐟</span>
        智钓蓝海
        <span class="title-sub">数据监控平台</span>
      </h1>
    </div>
    <div class="header-center">
      <div class="summary-items" v-if="summary">
        <div class="s-item">
          <span class="s-label">水域总数</span>
          <span class="s-value">{{ summary.total_spots }}</span>
        </div>
        <span class="s-divider">|</span>
        <div class="s-item">
          <span class="s-label">在线设备</span>
          <span class="s-value highlight">{{ summary.online_devices }}<small>/{{ summary.total_devices }}</small></span>
        </div>
        <span class="s-divider">|</span>
        <div class="s-item">
          <span class="s-label">在线网关</span>
          <span class="s-value highlight">{{ summary.online_gateways }}<small>/{{ summary.total_gateways }}</small></span>
        </div>
        <span class="s-divider">|</span>
        <div class="s-item">
          <span class="s-label">当前垂钓</span>
          <span class="s-value accent">{{ summary.total_fishing_count }}<small>人</small></span>
        </div>
      </div>
    </div>
    <div class="header-right">
      <div class="time-display">
        <div class="date">{{ dateStr }}</div>
        <div class="time">{{ timeStr }}</div>
      </div>
      <button class="back-to-main" @click="router.push('/')">
        <span class="back-icon">←</span> 返回主站
      </button>
    </div>
    <!-- 底部装饰线 -->
    <div class="header-line"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import type { SystemSummary } from '@/types'

const router = useRouter()
defineProps<{ summary: SystemSummary | null }>()

const dateStr = ref('')
const timeStr = ref('')

function updateTime() {
  const now = new Date()
  dateStr.value = now.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit', weekday: 'short' })
  timeStr.value = now.toLocaleTimeString('zh-CN', { hour12: false })
}

let timer: ReturnType<typeof setInterval>
onMounted(() => {
  updateTime()
  timer = setInterval(updateTime, 1000)
})
onUnmounted(() => clearInterval(timer))
</script>

<style scoped>
.dash-header {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 24px 8px;
  background: linear-gradient(180deg, rgba(6, 30, 65, 0.95) 0%, rgba(6, 30, 65, 0.6) 100%);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 8px;
}

.title {
  font-size: 22px;
  font-weight: 700;
  letter-spacing: 4px;
  background: linear-gradient(90deg, #00d4ff, #1e90ff, #00d4ff);
  background-size: 200% auto;
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  animation: titleShimmer 3s linear infinite;
}

.title-icon {
  font-size: 24px;
  -webkit-text-fill-color: initial;
  margin-right: 4px;
}

.title-sub {
  font-size: 14px;
  font-weight: 400;
  letter-spacing: 2px;
  margin-left: 8px;
  opacity: 0.7;
}

@keyframes titleShimmer {
  to { background-position: 200% center; }
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
}

.summary-items {
  display: flex;
  align-items: center;
  gap: 12px;
}

.s-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
}

.s-label {
  font-size: 11px;
  color: #8ba3c7;
  letter-spacing: 1px;
}

.s-value {
  font-size: 18px;
  font-weight: 700;
  color: #e8f0ff;
  font-family: 'DIN Alternate', 'Courier New', monospace;
}

.s-value small {
  font-size: 12px;
  font-weight: 400;
  color: #8ba3c7;
}

.s-value.highlight { color: #00d4ff; }
.s-value.accent { color: #00ff88; }

.s-divider {
  color: #2a4a6b;
  font-size: 20px;
  margin: 0 4px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header-right .time-display {
  text-align: right;
}

.back-to-main {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 14px;
  border: 1px solid rgba(0, 212, 255, 0.3);
  border-radius: 6px;
  background: rgba(0, 212, 255, 0.08);
  color: #00d4ff;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.3s;
  white-space: nowrap;
}

.back-to-main:hover {
  background: rgba(0, 212, 255, 0.2);
  border-color: #00d4ff;
}

.back-icon {
  font-size: 14px;
}

.date {
  font-size: 12px;
  color: #8ba3c7;
}

.time {
  font-size: 20px;
  font-weight: 700;
  font-family: 'DIN Alternate', 'Courier New', monospace;
  color: #00d4ff;
  letter-spacing: 2px;
}

.header-line {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, #1e90ff, #00d4ff, #1e90ff, transparent);
}
</style>
