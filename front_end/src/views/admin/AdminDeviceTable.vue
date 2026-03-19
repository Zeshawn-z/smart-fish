<template>
  <el-table :data="devices" stripe v-loading="loading">
    <el-table-column prop="id" label="ID" width="60" />
    <el-table-column prop="name" label="名称" min-width="150" />
    <el-table-column label="类型" width="100">
      <template #default="{ row }"><el-tag size="small">{{ DEVICE_TYPE_MAP[row.device_type as DeviceType] || row.device_type }}</el-tag></template>
    </el-table-column>
    <el-table-column label="状态" width="90">
      <template #default="{ row }">
        <el-tag :type="row.status === 'online' ? 'success' : 'info'" size="small">{{ row.status }}</el-tag>
      </template>
    </el-table-column>
    <el-table-column prop="fishing_count" label="垂钓人数" width="100" align="center" />
    <el-table-column label="操作" width="150" align="center">
      <template #default="{ row }">
        <el-button size="small" type="primary" link @click="$emit('edit', row)">编辑</el-button>
        <el-button size="small" type="danger" link @click="$emit('delete', 'device', row.id)">删除</el-button>
      </template>
    </el-table-column>
  </el-table>
</template>

<script setup lang="ts">
import { DEVICE_TYPE_MAP, type Device, type DeviceType } from '@/types'

defineProps<{
  devices: Device[]
  loading: boolean
}>()

defineEmits<{
  (e: 'edit', row: Device): void
  (e: 'delete', type: string, id: number): void
}>()
</script>
