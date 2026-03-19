<template>
  <el-table :data="spots" stripe v-loading="loading">
    <el-table-column prop="id" label="ID" width="60" />
    <el-table-column prop="name" label="名称" min-width="150" />
    <el-table-column label="区域" min-width="130">
      <template #default="{ row }">{{ row.region?.province }} · {{ row.region?.city }}</template>
    </el-table-column>
    <el-table-column label="类型" width="80">
      <template #default="{ row }"><el-tag size="small">{{ WATER_TYPE_MAP[row.water_type as WaterType] }}</el-tag></template>
    </el-table-column>
    <el-table-column label="状态" width="90">
      <template #default="{ row }">
        <el-tag :type="row.status === 'open' ? 'success' : 'danger'" size="small">{{ SPOT_STATUS_MAP[row.status as SpotStatus] }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column label="操作" width="150" align="center">
      <template #default="{ row }">
        <el-button size="small" type="primary" link @click="$emit('edit', row)">编辑</el-button>
        <el-button size="small" type="danger" link @click="$emit('delete', 'spot', row.id)">删除</el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import { WATER_TYPE_MAP, SPOT_STATUS_MAP, type FishingSpot, type WaterType, type SpotStatus } from '@/types'

defineProps<{
  spots: FishingSpot[]
  loading: boolean
}>()

defineEmits<{
  (e: 'edit', row: FishingSpot): void
  (e: 'delete', type: string, id: number): void
}>()
</script>
