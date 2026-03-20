<template>
  <div class="favorite-list">
    <div v-if="fishingStore.favoriteSpots.length === 0" class="empty-state">
      <el-empty description="暂无收藏水域">
        <template #description>
          <p class="empty-text">还没有收藏任何水域，去水域列表看看吧</p>
        </template>
        <el-button type="primary" size="small" @click="router.push('/spots')">浏览水域</el-button>
      </el-empty>
    </div>

    <div v-else class="spot-grid">
      <SpotCard
        v-for="spot in fishingStore.favoriteSpots"
        :key="spot.id"
        :spot="spot"
        :is-favorite="true"
        :show-favorite="true"
        @click="handleSpotClick"
        @toggle-favorite="handleToggleFavorite"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { useFishingStore } from '@/stores/fishing'
import { ElMessage } from 'element-plus'
import type { FishingSpot } from '@/types'
import SpotCard from '@/views/spots/SpotCard.vue'

const router = useRouter()
const fishingStore = useFishingStore()

function handleSpotClick(spot: FishingSpot) {
  router.push(`/spots?id=${spot.id}`)
}

async function handleToggleFavorite(spotId: number) {
  await fishingStore.toggleFavorite(spotId)
  ElMessage.success('已取消收藏')
}
</script>

<style scoped>
.favorite-list {
  min-height: 200px;
}

.spot-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.empty-state {
  padding: 40px 0;
}

.empty-text {
  color: #a0a4ad;
  font-size: 13px;
  margin: 0;
}

@media (max-width: 768px) {
  .spot-grid {
    grid-template-columns: 1fr;
    gap: 12px;
  }
}
</style>
