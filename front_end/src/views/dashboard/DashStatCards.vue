<template>
  <div class="panel stat-cards">
    <div class="panel-title">核心指标</div>
    <div class="cards-grid">
      <div v-for="item in cards" :key="item.label" class="stat-card" :style="{ '--accent': item.color }">
        <div class="card-icon">{{ item.icon }}</div>
        <div class="card-info">
          <div class="card-value">{{ item.value }}</div>
          <div class="card-label">{{ item.label }}</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import type { SystemSummary } from '@/types'

const props = defineProps<{ summary: SystemSummary | null }>()

const cards = computed(() => {
  const s = props.summary
  if (!s) return []
  return [
    { icon: '🎣', label: '开放水域', value: s.open_spots, color: '#00ff88' },
    { icon: '📡', label: '在线设备', value: s.online_devices, color: '#00d4ff' },
    { icon: '🌐', label: '在线网关', value: s.online_gateways, color: '#1e90ff' },
    { icon: '👥', label: '注册用户', value: s.total_users, color: '#a78bfa' },
    { icon: '🐠', label: '垂钓人数', value: s.total_fishing_count, color: '#ffd93d' },
    { icon: '🔔', label: '活跃提醒', value: s.active_reminders, color: '#ff6b6b' },
  ]
})
</script>

<style scoped>
.stat-cards { padding: 10px; }

.cards-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 8px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  background: linear-gradient(135deg, rgba(0, 212, 255, 0.08), rgba(30, 144, 255, 0.04));
  border: 1px solid rgba(0, 212, 255, 0.15);
  transition: all 0.3s;
}

.stat-card:hover {
  border-color: var(--accent);
  box-shadow: 0 0 12px color-mix(in srgb, var(--accent) 30%, transparent);
}

.card-icon {
  font-size: 24px;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  background: rgba(0, 212, 255, 0.1);
}

.card-value {
  font-size: 22px;
  font-weight: 700;
  color: var(--accent);
  font-family: 'DIN Alternate', 'Courier New', monospace;
  line-height: 1.2;
}

.card-label {
  font-size: 11px;
  color: #8ba3c7;
  letter-spacing: 1px;
}
</style>
