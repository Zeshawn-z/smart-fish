<template>
  <div v-if="fishingStore.favoriteSpots.length === 0" class="empty-state">
    <el-empty description="暂无收藏水域" />
  </div>
  <el-table v-else :data="fishingStore.favoriteSpots" stripe>
    <el-table-column prop="name" label="水域名称" min-width="150" />
    <el-table-column label="区域" min-width="120">
      <template #default="{ row }">{{ row.region?.province }} · {{ row.region?.city }}</template>
    </el-table-column>
    <el-table-column label="状态" width="90" align="center">
      <template #default="{ row }">
        <el-tag :type="row.status === 'open' ? 'success' : 'danger'" size="small">
          {{ SPOT_STATUS_MAP[row.status as SpotStatus] }}
        </el-tag>
      </template>
    </el-table-column>
    <el-table-column label="操作" width="120" align="center">
      <template #default="{ row }">
        <el-button size="small" type="danger" link @click="handleUnfavorite(row.id)">取消收藏</el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import { useFishingStore } from '@/stores/fishing'
import { SPOT_STATUS_MAP, type SpotStatus } from '@/types'
import { ElMessage } from 'element-plus'

const fishingStore = useFishingStore()

async function handleUnfavorite(spotId: number) {
  await fishingStore.toggleFavorite(spotId)
  ElMessage.success('已取消收藏')
}
</script>

<style scoped>
.empty-state { padding: 20px; }
</style>
