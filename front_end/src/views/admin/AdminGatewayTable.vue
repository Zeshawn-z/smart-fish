<template>
  <el-table :data="gateways" stripe v-loading="loading">
    <el-table-column prop="id" label="ID" width="60" />
    <el-table-column prop="name" label="名称" min-width="150" />
    <el-table-column label="状态" width="90">
      <template #default="{ row }">
        <el-tag :type="row.status === 'online' ? 'success' : 'info'" size="small">{{ row.status }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="cpu_usage" label="CPU%" width="80" align="center" />
    <el-table-column prop="memory_usage" label="内存%" width="80" align="center" />
    <el-table-column prop="battery_level" label="电量%" width="80" align="center" />
    <el-table-column label="设备数" width="80" align="center">
      <template #default="{ row }">{{ row.devices?.length || 0 }}</template>
    </el-table-column>
    <el-table-column label="操作" width="150" align="center">
      <template #default="{ row }">
        <el-button size="small" type="primary" link @click="$emit('edit', row)">编辑</el-button>
        <el-button size="small" type="danger" link @click="$emit('delete', 'gateway', row.id)">删除</el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import type { Gateway } from '@/types'

defineProps<{
  gateways: Gateway[]
  loading: boolean
}>()

defineEmits<{
  (e: 'edit', row: Gateway): void
  (e: 'delete', type: string, id: number): void
}>()
</script>
