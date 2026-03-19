<template>
  <el-table :data="regions" stripe v-loading="loading">
    <el-table-column prop="id" label="ID" width="60" />
    <el-table-column prop="name" label="名称" min-width="150" />
    <el-table-column prop="province" label="省份" width="100" />
    <el-table-column prop="city" label="城市" width="100" />
    <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
    <el-table-column label="操作" width="150" align="center">
      <template #default="{ row }">
        <el-button size="small" type="primary" link @click="$emit('edit', row)">编辑</el-button>
        <el-button size="small" type="danger" link @click="$emit('delete', 'region', row.id)">删除</el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import type { Region } from '@/types'

defineProps<{
  regions: Region[]
  loading: boolean
}>()

defineEmits<{
  (e: 'edit', row: Region): void
  (e: 'delete', type: string, id: number): void
}>()
</script>
