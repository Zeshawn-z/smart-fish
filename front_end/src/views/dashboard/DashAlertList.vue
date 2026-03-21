<template>
  <div class="panel alert-list">
    <div class="panel-title">
      <span>实时告警</span>
      <span class="badge" v-if="reminders.length">{{ reminders.length }}</span>
      <span class="marquee-status" v-if="reminders.length > 2">
        <span class="status-dot" :class="{ paused: isPaused }"></span>
        {{ isPaused ? '已暂停' : '轮播中' }}
      </span>
    </div>
    <div
      class="alert-scroll"
      ref="scrollRef"
      @mouseenter="pause"
      @mouseleave="resume"
    >
      <!-- 真实列表 + 克隆副本实现无缝滚动 -->
      <div class="scroll-inner" ref="innerRef">
        <div
          v-for="(r, idx) in displayList"
          :key="r.id + '-' + idx"
          class="alert-item"
          :class="['level-' + r.level]"
        >
          <div class="alert-top">
            <span class="alert-type-tag" :class="['type-' + r.reminder_type]">
              {{ typeIcon(r.reminder_type) }} {{ typeLabel(r.reminder_type) }}
            </span>
            <span class="alert-level-tag" :class="['lvl-' + r.level]">
              {{ REMINDER_LEVEL_MAP[r.level as ReminderLevel] }}
            </span>
          </div>
          <div class="alert-msg">{{ r.message }}</div>
          <div class="alert-bottom">
            <span class="alert-time">{{ relativeTime(r.timestamp) }}</span>
          </div>
        </div>
      </div>
      <div v-if="!reminders.length" class="empty-hint">暂无活跃告警 ✅</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import type { Reminder, ReminderLevel } from '@/types'
import { REMINDER_LEVEL_MAP } from '@/types'

const props = defineProps<{ reminders: Reminder[] }>()

const scrollRef = ref<HTMLElement>()
const innerRef = ref<HTMLElement>()
const isPaused = ref(false)

// 无缝滚动：原始列表 + 克隆一份
const displayList = computed(() => {
  if (props.reminders.length <= 2) return props.reminders
  return [...props.reminders, ...props.reminders]
})

// ===== 自动滚动逻辑 =====
let rafId: number | null = null
let scrollSpeed = 0.5 // px/frame

function scrollStep() {
  const el = scrollRef.value
  const inner = innerRef.value
  if (!el || !inner || isPaused.value || props.reminders.length <= 2) {
    rafId = requestAnimationFrame(scrollStep)
    return
  }

  // 单份内容高度
  const singleHeight = inner.scrollHeight / 2

  el.scrollTop += scrollSpeed

  // 当滚到克隆区域时，瞬间跳回
  if (el.scrollTop >= singleHeight) {
    el.scrollTop -= singleHeight
  }

  rafId = requestAnimationFrame(scrollStep)
}

function pause() { isPaused.value = true }
function resume() { isPaused.value = false }

// 数据变化时重置滚动位置
watch(() => props.reminders, () => {
  if (scrollRef.value) {
    scrollRef.value.scrollTop = 0
  }
}, { deep: true })

onMounted(() => {
  rafId = requestAnimationFrame(scrollStep)
})

onUnmounted(() => {
  if (rafId !== null) cancelAnimationFrame(rafId)
})

// ===== 辅助函数 =====
function typeIcon(type: string): string {
  const icons: Record<string, string> = {
    weather: '🌤️', fishing: '🎣', safety: '⚠️', environment: '🌊',
    capacity: '👥', info: '📋'
  }
  return icons[type] || '📢'
}

function typeLabel(type: string): string {
  const labels: Record<string, string> = {
    weather: '天气', fishing: '垂钓', safety: '安全',
    environment: '环境', capacity: '容量', info: '通知'
  }
  return labels[type] || '系统'
}

function relativeTime(ts: string): string {
  const diff = Date.now() - new Date(ts).getTime()
  const mins = Math.floor(diff / 60000)
  if (mins < 1) return '刚刚'
  if (mins < 60) return `${mins}分钟前`
  const hours = Math.floor(mins / 60)
  if (hours < 24) return `${hours}小时前`
  return `${Math.floor(hours / 24)}天前`
}
</script>

<style scoped>
.alert-list {
  padding: 10px;
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
}

.panel-title {
  display: flex;
  align-items: center;
  gap: 6px;
}

.badge {
  font-size: 11px;
  background: #ff6b6b;
  color: #fff;
  padding: 0 6px;
  border-radius: 10px;
  font-weight: 700;
  line-height: 18px;
}

.marquee-status {
  margin-left: auto;
  font-size: 10px;
  color: #5a7a9a;
  display: flex;
  align-items: center;
  gap: 4px;
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: #00ff88;
  animation: pulse-dot 1.5s ease-in-out infinite;
}

.status-dot.paused {
  background: #ffd93d;
  animation: none;
}

@keyframes pulse-dot {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.3; }
}

.alert-scroll {
  flex: 1;
  overflow: hidden;
  margin-top: 6px;
  position: relative;
  /* 渐隐遮罩 */
  mask-image: linear-gradient(
    to bottom,
    transparent 0%,
    black 4%,
    black 92%,
    transparent 100%
  );
  -webkit-mask-image: linear-gradient(
    to bottom,
    transparent 0%,
    black 4%,
    black 92%,
    transparent 100%
  );
}

.scroll-inner {
  display: flex;
  flex-direction: column;
}

.alert-item {
  padding: 10px 12px;
  margin-bottom: 8px;
  border-radius: 8px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.04);
  transition: background 0.3s, border-color 0.3s;
  flex-shrink: 0;
}

.alert-item:hover {
  background: rgba(255, 255, 255, 0.06);
}

/* 级别底色 — 微妙的渐变 */
.alert-item.level-3 {
  background: linear-gradient(135deg, rgba(255, 68, 68, 0.12), rgba(255, 68, 68, 0.04));
  border-color: rgba(255, 68, 68, 0.15);
}
.alert-item.level-2 {
  background: linear-gradient(135deg, rgba(245, 108, 108, 0.08), rgba(245, 108, 108, 0.02));
  border-color: rgba(245, 108, 108, 0.1);
}
.alert-item.level-1 {
  background: linear-gradient(135deg, rgba(230, 162, 60, 0.08), rgba(230, 162, 60, 0.02));
  border-color: rgba(230, 162, 60, 0.1);
}
.alert-item.level-0 {
  border-color: rgba(64, 158, 255, 0.08);
}

.alert-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 6px;
}

.alert-type-tag {
  font-size: 10px;
  padding: 2px 8px;
  border-radius: 10px;
  background: rgba(0, 212, 255, 0.1);
  color: #8ba3c7;
}
.alert-type-tag.type-weather { background: rgba(0, 212, 255, 0.1); color: #00d4ff; }
.alert-type-tag.type-environment { background: rgba(0, 255, 136, 0.08); color: #00ff88; }
.alert-type-tag.type-safety { background: rgba(255, 217, 61, 0.1); color: #ffd93d; }
.alert-type-tag.type-capacity { background: rgba(167, 139, 250, 0.1); color: #a78bfa; }
.alert-type-tag.type-fishing { background: rgba(0, 212, 255, 0.1); color: #00d4ff; }
.alert-type-tag.type-info { background: rgba(139, 163, 199, 0.1); color: #8ba3c7; }

.alert-level-tag {
  font-size: 10px;
  padding: 1px 8px;
  border-radius: 10px;
  font-weight: 600;
}
.alert-level-tag.lvl-0 { background: rgba(64, 158, 255, 0.12); color: #409eff; }
.alert-level-tag.lvl-1 { background: rgba(230, 162, 60, 0.12); color: #e6a23c; }
.alert-level-tag.lvl-2 { background: rgba(245, 108, 108, 0.15); color: #f56c6c; }
.alert-level-tag.lvl-3 { background: rgba(255, 68, 68, 0.18); color: #ff4444; }

.alert-msg {
  font-size: 12px;
  color: #d0dbe8;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.alert-bottom {
  display: flex;
  justify-content: flex-end;
  margin-top: 6px;
}

.alert-time {
  font-size: 10px;
  color: #4a6a8a;
}

.empty-hint {
  text-align: center;
  color: #5a7a9a;
  font-size: 13px;
  padding: 20px 0;
}
</style>
