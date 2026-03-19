<template>
  <div class="admin-page">
    <el-container>
      <AdminSidebar v-model:active="activeMenu" :show-users="authStore.isAdmin" />

      <el-main class="admin-main">
        <el-card shadow="hover">
          <template #header>
            <div class="admin-header">
              <span class="card-title">{{ menuTitles[activeMenu] }}</span>
              <el-button type="primary" @click="handleCreate" v-if="activeMenu !== 'users'">
                <el-icon><Plus /></el-icon> 新增
              </el-button>
            </div>
          </template>

          <AdminRegionTable
            v-if="activeMenu === 'regions'"
            :regions="regionList"
            :loading="loading"
            @edit="handleEdit"
            @delete="handleDelete"
          />

          <AdminSpotTable
            v-else-if="activeMenu === 'spots'"
            :spots="spotList"
            :loading="loading"
            @edit="handleEdit"
            @delete="handleDelete"
          />

          <AdminDeviceTable
            v-else-if="activeMenu === 'devices'"
            :devices="deviceList"
            :loading="loading"
            @edit="handleEdit"
            @delete="handleDelete"
          />

          <AdminGatewayTable
            v-else-if="activeMenu === 'gateways'"
            :gateways="gatewayList"
            :loading="loading"
            @edit="handleEdit"
            @delete="handleDelete"
          />

          <template v-else>
            <el-empty :description="`${menuTitles[activeMenu]}模块开发中...`" />
          </template>
        </el-card>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { RegionService, DeviceService, FishingSpotService } from '@/services'
import { httpGet, httpDelete } from '@/network/http'
import type { Region, FishingSpot, Device, Gateway } from '@/types'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

// 子组件
import AdminSidebar from './admin/AdminSidebar.vue'
import AdminRegionTable from './admin/AdminRegionTable.vue'
import AdminSpotTable from './admin/AdminSpotTable.vue'
import AdminDeviceTable from './admin/AdminDeviceTable.vue'
import AdminGatewayTable from './admin/AdminGatewayTable.vue'

const authStore = useAuthStore()
const activeMenu = ref('regions')
const loading = ref(false)

const regionList = ref<Region[]>([])
const spotList = ref<FishingSpot[]>([])
const deviceList = ref<Device[]>([])
const gatewayList = ref<Gateway[]>([])

const menuTitles: Record<string, string> = {
  regions: '区域管理',
  spots: '水域管理',
  devices: '设备管理',
  gateways: '网关管理',
  reminders: '提醒管理',
  notices: '通知管理',
  users: '用户管理'
}

async function loadData() {
  loading.value = true
  try {
    switch (activeMenu.value) {
      case 'regions':
        regionList.value = await RegionService.list()
        break
      case 'spots': {
        const data = await FishingSpotService.list({ page_size: 100 })
        spotList.value = data.results
        break
      }
      case 'devices':
        deviceList.value = await DeviceService.list()
        break
      case 'gateways':
        gatewayList.value = await httpGet<Gateway[]>('/api/gateways')
        break
    }
  } finally {
    loading.value = false
  }
}

function handleCreate() {
  ElMessage.info('新增功能开发中')
}

function handleEdit(_row: any) {
  ElMessage.info('编辑功能开发中')
}

async function handleDelete(type: string, id: number) {
  await ElMessageBox.confirm('确认删除？此操作不可恢复。', '警告', {
    confirmButtonText: '确认',
    cancelButtonText: '取消',
    type: 'warning'
  })

  const urlMap: Record<string, string> = {
    region: `/api/regions/${id}`,
    spot: `/api/spots/${id}`,
    device: `/api/devices/${id}`,
    gateway: `/api/gateways/${id}`
  }

  await httpDelete(urlMap[type])
  ElMessage.success('删除成功')
  loadData()
}

watch(activeMenu, () => loadData())
onMounted(() => loadData())
</script>

<style scoped>
.admin-page {
  min-height: calc(100vh - 160px);
}

.admin-main { padding: 0; }

.admin-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-weight: 600;
  font-size: 16px;
}
</style>
