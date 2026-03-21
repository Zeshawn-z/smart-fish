<template>
  <div class="panel suggestion-panel">
    <div class="panel-title">垂钓建议</div>
    <div class="suggestion-carousel">
      <Transition name="fade" mode="out-in">
        <div v-if="current" :key="current.id" class="suggestion-card">
          <div class="suggestion-header">
            <span class="spot-name">{{ current.fishing_spot?.name || `水域 #${current.spot_id}` }}</span>
            <span class="score" :style="{ color: scoreColor(current.score) }">{{ current.score.toFixed(0) }}分</span>
          </div>
          <div
            class="suggestion-text"
            :title="current.suggestion_text"
            @mouseenter="showTooltip($event, current.suggestion_text)"
            @mouseleave="hideTooltip"
          >{{ current.suggestion_text }}</div>
          <div class="suggestion-footer">
            <span class="indicator">{{ activeIndex + 1 }} / {{ suggestions.length }}</span>
          </div>
        </div>
      </Transition>
      <div v-if="!suggestions.length" class="empty-hint">暂无建议</div>
    </div>
    <!-- 浮窗 tooltip -->
    <Teleport to="body">
      <div
        v-if="tooltipVisible"
        class="suggestion-tooltip"
        :style="tooltipStyle"
      >{{ tooltipText }}</div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted, reactive } from 'vue'
import type { FishingSuggestion } from '@/types'

const props = defineProps<{ suggestions: FishingSuggestion[] }>()

const activeIndex = ref(0)
const current = computed(() => props.suggestions[activeIndex.value] || null)

// tooltip 状态
const tooltipVisible = ref(false)
const tooltipText = ref('')
const tooltipStyle = reactive<Record<string, string>>({})

function showTooltip(e: MouseEvent, text: string) {
  const el = e.currentTarget as HTMLElement
  // 只有内容被截断时才显示 tooltip
  if (el.scrollHeight <= el.clientHeight) return
  tooltipText.value = text
  const rect = el.getBoundingClientRect()
  tooltipStyle.top = `${rect.top - 6}px`
  tooltipStyle.left = `${rect.left}px`
  tooltipStyle.maxWidth = `${Math.min(rect.width + 40, 380)}px`
  tooltipStyle.transform = 'translateY(-100%)'
  tooltipVisible.value = true
}

function hideTooltip() {
  tooltipVisible.value = false
}

function scoreColor(score: number): string {
  if (score >= 90) return '#00ff88'
  if (score >= 75) return '#00d4ff'
  if (score >= 60) return '#ffd93d'
  return '#ff6b6b'
}

let timer: ReturnType<typeof setInterval>

function startCarousel() {
  clearInterval(timer)
  if (props.suggestions.length > 1) {
    timer = setInterval(() => {
      activeIndex.value = (activeIndex.value + 1) % props.suggestions.length
    }, 6000)
  }
}

watch(() => props.suggestions.length, () => {
  activeIndex.value = 0
  startCarousel()
})

onMounted(startCarousel)
onUnmounted(() => clearInterval(timer))
</script>

<style scoped>
.suggestion-panel {
  padding: 10px;
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.suggestion-carousel {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  margin-top: 6px;
  min-height: 0;
}

.suggestion-card {
  padding: 10px;
  border-radius: 6px;
  background: linear-gradient(135deg, rgba(0, 212, 255, 0.06), rgba(0, 255, 136, 0.04));
  border: 1px solid rgba(0, 212, 255, 0.12);
  display: flex;
  flex-direction: column;
  max-height: 100%;
}

.suggestion-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  flex-shrink: 0;
}

.spot-name {
  font-size: 13px;
  font-weight: 600;
  color: #00d4ff;
}

.score {
  font-size: 18px;
  font-weight: 700;
  font-family: 'DIN Alternate', 'Courier New', monospace;
}

.suggestion-text {
  font-size: 12px;
  color: #b8cce0;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  cursor: default;
}

.suggestion-footer {
  text-align: center;
  flex-shrink: 0;
}

.indicator {
  font-size: 10px;
  color: #5a7a9a;
}

.empty-hint {
  text-align: center;
  color: #5a7a9a;
  font-size: 12px;
}

.fade-enter-active, .fade-leave-active { transition: opacity 0.4s ease; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>

<style>
/* tooltip 不 scoped，因为 teleport 到 body */
.suggestion-tooltip {
  position: fixed;
  z-index: 9999;
  padding: 10px 12px;
  font-size: 12px;
  line-height: 1.6;
  color: #e8f0ff;
  background: rgba(10, 25, 50, 0.95);
  border: 1px solid rgba(0, 212, 255, 0.25);
  border-radius: 6px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.5);
  pointer-events: none;
  word-break: break-all;
}</style>
