<template>
  <el-card shadow="hover" class="filter-bar">
    <el-row :gutter="12" align="middle">
      <el-col :xs="12" :sm="6" :md="6">
        <el-select :model-value="status" placeholder="水域状态" clearable size="default"
          @update:model-value="$emit('update:status', $event)">
          <el-option label="开放" value="open" />
          <el-option label="关闭" value="closed" />
          <el-option label="维护中" value="maintenance" />
        </el-select>
      </el-col>
      <el-col :xs="12" :sm="6" :md="6">
        <el-select :model-value="waterType" placeholder="水域类型" clearable size="default"
          @update:model-value="$emit('update:waterType', $event)">
          <el-option label="湖泊" value="lake" />
          <el-option label="河流" value="river" />
          <el-option label="水库" value="reservoir" />
          <el-option label="鱼塘" value="pond" />
        </el-select>
      </el-col>
      <el-col :xs="16" :sm="6" :md="6">
        <el-input :model-value="keyword" placeholder="搜索水域名称..." clearable size="default"
          @update:model-value="$emit('update:keyword', $event)"
          @keyup.enter="$emit('search')">
          <template #prefix><el-icon><Search /></el-icon></template>
        </el-input>
      </el-col>
      <el-col :xs="8" :sm="6" :md="6" class="text-right">
        <el-button v-if="hasRegionFilter" size="small" @click="$emit('clear-region')">
          清除区域
        </el-button>
        <el-button-group>
          <el-button :type="viewMode === 'card' ? 'primary' : 'default'" @click="$emit('update:viewMode', 'card')">
            <el-icon><Grid /></el-icon>
          </el-button>
          <el-button :type="viewMode === 'table' ? 'primary' : 'default'" @click="$emit('update:viewMode', 'table')">
            <el-icon><List /></el-icon>
          </el-button>
        </el-button-group>
      </el-col>
    </el-row>
  </el-card>
</template>

<script setup lang="ts">
import { Search, Grid, List } from '@element-plus/icons-vue'

defineProps<{
  status: string
  waterType: string
  keyword: string
  viewMode: 'card' | 'table'
  hasRegionFilter: boolean
}>()

defineEmits<{
  (e: 'update:status', value: string): void
  (e: 'update:waterType', value: string): void
  (e: 'update:keyword', value: string): void
  (e: 'update:viewMode', value: 'card' | 'table'): void
  (e: 'search'): void
  (e: 'clear-region'): void
}>()
</script>

<style scoped>
.filter-bar {
  margin-bottom: 16px;
  border-radius: 12px !important;
}

.text-right {
  text-align: right;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 8px;
}

@media (max-width: 768px) {
  .filter-bar :deep(.el-col) { margin-bottom: 8px; }
  .text-right { flex-wrap: wrap; }
}
</style>
