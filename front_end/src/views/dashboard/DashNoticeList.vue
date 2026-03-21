<template>
  <div class="panel notice-list">
    <div class="panel-title">通知公告</div>
    <div class="notice-carousel">
      <Transition name="slide" mode="out-in">
        <div v-if="current" :key="current.id" class="notice-card">
          <div class="notice-title">{{ current.title }}</div>
          <div class="notice-content">{{ current.content }}</div>
          <div class="notice-footer">
            <span class="notice-time">{{ formatDate(current.timestamp) }}</span>
            <span class="indicator">{{ activeIndex + 1 }} / {{ notices.length }}</span>
          </div>
        </div>
      </Transition>
      <div v-if="!notices.length" class="empty-hint">暂无公告</div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import type { Notice } from '@/types'

const props = defineProps<{ notices: Notice[] }>()

const activeIndex = ref(0)
const current = computed(() => props.notices[activeIndex.value] || null)

function formatDate(ts: string): string {
  const d = new Date(ts)
  return d.toLocaleDateString('zh-CN', { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' })
}

let timer: ReturnType<typeof setInterval>

function startCarousel() {
  clearInterval(timer)
  if (props.notices.length > 1) {
    timer = setInterval(() => {
      activeIndex.value = (activeIndex.value + 1) % props.notices.length
    }, 5000)
  }
}

watch(() => props.notices.length, () => {
  activeIndex.value = 0
  startCarousel()
})

onMounted(startCarousel)
onUnmounted(() => clearInterval(timer))
</script>

<style scoped>
.notice-list {
  padding: 10px;
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.notice-carousel {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
  margin-top: 6px;
}

.notice-card {
  padding: 10px;
  border-radius: 6px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(30, 144, 255, 0.08);
}

.notice-title {
  font-size: 13px;
  font-weight: 600;
  color: #e8f0ff;
  margin-bottom: 6px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.notice-content {
  font-size: 11px;
  color: #8ba3c7;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notice-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
}

.notice-time {
  font-size: 10px;
  color: #5a7a9a;
}

.indicator {
  font-size: 10px;
  color: #5a7a9a;
}

.empty-hint {
  text-align: center;
  color: #5a7a9a;
  font-size: 12px;
  padding: 20px 0;
}

.slide-enter-active, .slide-leave-active { transition: all 0.4s ease; }
.slide-enter-from { opacity: 0; transform: translateY(12px); }
.slide-leave-to { opacity: 0; transform: translateY(-12px); }
</style>
