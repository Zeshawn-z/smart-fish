<template>
  <el-card shadow="hover" class="tree-card">
    <template #header>
      <div class="tree-header">
        <el-icon class="tree-header-icon"><MapLocation /></el-icon>
        <span class="card-title">垂钓水域</span>
      </div>
    </template>
    <el-input v-model="localSearch" placeholder="搜索水域..." clearable class="search-input">
      <template #prefix><el-icon><Search /></el-icon></template>
    </el-input>
    <div class="province-tree" v-loading="loading">
      <div v-for="province in filteredProvinces" :key="province" class="province-group">
        <div class="province-header" @click="toggleProvince(province)">
          <el-icon class="arrow-icon" :class="{ rotated: expandedProvinces.has(province) }"><ArrowRight /></el-icon>
          <span>{{ province }}</span>
          <el-badge :value="getProvinceCount(province)" type="info" class="province-badge" />
        </div>
        <transition name="collapse">
          <div v-show="expandedProvinces.has(province)" class="region-list">
            <div
              v-for="region in getRegions(province)"
              :key="region.id"
              class="region-item"
              :class="{ active: selectedRegionId === region.id }"
              @click="$emit('select-region', region)"
            >
              <el-icon><Location /></el-icon>
              <span>{{ region.city }} - {{ region.name }}</span>
            </div>
          </div>
        </transition>
      </div>
      <el-empty v-if="filteredProvinces.length === 0 && !loading" description="暂无区域数据" :image-size="60" />
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import type { Region } from '@/types'
import { MapLocation, Search, ArrowRight, Location } from '@element-plus/icons-vue'

const props = defineProps<{
  regionsByProvince: Record<string, Region[]>
  selectedRegionId: number | null
  loading: boolean
}>()

defineEmits<{
  (e: 'select-region', region: Region): void
}>()

const localSearch = ref('')
const expandedProvinces = ref<Set<string>>(new Set())

const filteredProvinces = computed(() => {
  const provinces = Object.keys(props.regionsByProvince)
  if (!localSearch.value) return provinces
  return provinces.filter(p => {
    const regions = props.regionsByProvince[p]
    return p.includes(localSearch.value) ||
      regions.some(r => r.name.includes(localSearch.value) || r.city.includes(localSearch.value))
  })
})

function getRegions(province: string): Region[] {
  return props.regionsByProvince[province] || []
}

function getProvinceCount(province: string): number {
  return props.regionsByProvince[province]?.length || 0
}

function toggleProvince(province: string) {
  if (expandedProvinces.value.has(province)) {
    expandedProvinces.value.delete(province)
  } else {
    expandedProvinces.value.add(province)
  }
}

// 默认展开第一个省份
function expandFirst() {
  const provinces = Object.keys(props.regionsByProvince)
  if (provinces.length > 0) {
    expandedProvinces.value.add(provinces[0])
  }
}

defineExpose({ expandFirst })
</script>

<style scoped>
.tree-card {
  position: sticky;
  top: 80px;
  border-radius: 12px !important;
}

.tree-header {
  display: flex;
  align-items: center;
  gap: 8px;
}

.tree-header-icon {
  font-size: 20px;
  color: #409eff;
}

.card-title {
  font-weight: 600;
  font-size: 16px;
}

.search-input { margin-bottom: 12px; }

.province-tree {
  max-height: calc(100vh - 280px);
  overflow-y: auto;
}

.province-group { margin-bottom: 4px; }

.province-header {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  cursor: pointer;
  border-radius: 6px;
  font-weight: 600;
  font-size: 14px;
  transition: background 0.2s;
}

.province-header:hover { background: #f5f7fa; }
.province-badge { margin-left: auto; }

.arrow-icon { transition: transform 0.2s; }
.arrow-icon.rotated { transform: rotate(90deg); }

.region-list { padding-left: 24px; }

.region-item {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 12px;
  cursor: pointer;
  border-radius: 6px;
  font-size: 13px;
  color: #606266;
  transition: all 0.2s;
}

.region-item:hover { background: #ecf5ff; color: #409eff; }
.region-item.active { background: #409eff; color: #fff; }

.collapse-enter-active, .collapse-leave-active {
  transition: all 0.2s;
  overflow: hidden;
}

.collapse-enter-from, .collapse-leave-to {
  max-height: 0;
  opacity: 0;
}
</style>
