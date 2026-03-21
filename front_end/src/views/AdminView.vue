<template>
  <div class="admin-layout">
    <!-- 左侧侧边栏 -->
    <AdminSidebar
      v-model:active="activeMenu"
      :show-users="authStore.isAdmin"
      :username="authStore.username"
      :role-name="authStore.isAdmin ? '管理员' : '工作人员'"
      @logout="handleLogout"
    />

    <!-- 右侧内容区 -->
    <div class="admin-content">
      <ResourceManager
        v-if="currentConfig"
        :key="activeMenu"
        :config="currentConfig"
      />
      <div v-else class="admin-empty">
        <el-empty :description="`${activeMenu} 模块开发中...`" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  RegionService, FishingSpotService, DeviceService,
  GatewayService, ReminderService, NoticeService,
  PostService, CommentService, IoTDeviceService,
  FishingRecordResourceService, FishCaughtResourceService,
  SuggestionService
} from '@/services'
import { httpGet, httpPatch, httpDelete } from '@/network/http'
import {
  WATER_TYPE_MAP, SPOT_STATUS_MAP, DEVICE_TYPE_MAP,
  REMINDER_LEVEL_MAP,
  type WaterType, type SpotStatus, type DeviceType,
  type ReminderLevel, type ResourceConfig
} from '@/types'

import AdminSidebar from './admin/AdminSidebar.vue'
import ResourceManager from './admin/ResourceManager.vue'

const router = useRouter()
const authStore = useAuthStore()
const activeMenu = ref('regions')

function handleLogout() {
  authStore.logout()
  router.push('/')
}

// ===== 各模型的配置表 =====
const resourceConfigs: Record<string, ResourceConfig> = {
  regions: {
    title: '区域管理',
    resource: 'regions',
    searchable: true,
    creatable: true,
    columns: [
      { prop: 'id', label: 'ID', width: 60 },
      { prop: 'name', label: '名称', minWidth: 150 },
      { prop: 'province', label: '省份', width: 100 },
      { prop: 'city', label: '城市', width: 100 },
      { prop: 'description', label: '描述', minWidth: 200 }
    ],
    formFields: [
      { prop: 'name', label: '名称', type: 'input', required: true },
      { prop: 'province', label: '省份', type: 'input', required: true },
      { prop: 'city', label: '城市', type: 'input', required: true },
      { prop: 'description', label: '描述', type: 'textarea' }
    ],
    loadFn: (params) => RegionService.list(params),
    createFn: (data) => RegionService.create(data),
    updateFn: (id, data) => RegionService.update(id, data),
    deleteFn: (id) => RegionService.delete(id)
  },

  spots: {
    title: '水域管理',
    resource: 'spots',
    searchable: true,
    creatable: true,
    paginated: true,
    columns: [
      { prop: 'id', label: 'ID', width: 60 },
      { prop: 'name', label: '名称', minWidth: 150 },
      {
        prop: 'region',
        label: '区域',
        minWidth: 130,
        formatter: (row) => `${row.region?.province || ''} · ${row.region?.city || ''}`
      },
      {
        prop: 'water_type',
        label: '类型',
        width: 90,
        tag: { label: (row) => WATER_TYPE_MAP[row.water_type as WaterType] || row.water_type }
      },
      {
        prop: 'status',
        label: '状态',
        width: 90,
        tag: {
          label: (row) => SPOT_STATUS_MAP[row.status as SpotStatus] || row.status,
          type: (row) => row.status === 'open' ? 'success' : row.status === 'maintenance' ? 'warning' : 'danger'
        }
      }
    ],
    formFields: [
      { prop: 'name', label: '名称', type: 'input', required: true },
      { prop: 'description', label: '描述', type: 'textarea' },
      { prop: 'latitude', label: '纬度', type: 'number' },
      { prop: 'longitude', label: '经度', type: 'number' },
      {
        prop: 'water_type', label: '水域类型', type: 'select',
        options: Object.entries(WATER_TYPE_MAP).map(([k, v]) => ({ label: v, value: k }))
      },
      { prop: 'capacity', label: '容量', type: 'number' },
      {
        prop: 'status', label: '状态', type: 'select',
        options: Object.entries(SPOT_STATUS_MAP).map(([k, v]) => ({ label: v, value: k }))
      }
    ],
    loadFn: (params) => FishingSpotService.list(params),
    createFn: (data) => FishingSpotService.create(data),
    updateFn: (id, data) => FishingSpotService.update(id, data),
    deleteFn: (id) => FishingSpotService.delete(id)
  },

  devices: {
    title: '设备管理',
    resource: 'devices',
    searchable: true,
    creatable: true,
    columns: [
      { prop: 'id', label: 'ID', width: 60 },
      { prop: 'name', label: '名称', minWidth: 180 },
      {
        prop: 'device_type', label: '类型', width: 100,
        tag: { label: (row) => DEVICE_TYPE_MAP[row.device_type as DeviceType] || row.device_type }
      },
      {
        prop: 'status', label: '状态', width: 90,
        tag: {
          label: (row) => row.status,
          type: (row) => row.status === 'online' ? 'success' : 'info'
        }
      },
      { prop: 'fishing_count', label: '垂钓人数', width: 100, align: 'center' }
    ],
    formFields: [
      { prop: 'name', label: '名称', type: 'input', required: true },
      { prop: 'description', label: '描述', type: 'textarea' },
      {
        prop: 'device_type', label: '类型', type: 'select',
        options: Object.entries(DEVICE_TYPE_MAP).map(([k, v]) => ({ label: v, value: k }))
      },
      {
        prop: 'status', label: '状态', type: 'select',
        options: [
          { label: '在线', value: 'online' },
          { label: '离线', value: 'offline' },
          { label: '错误', value: 'error' }
        ]
      }
    ],
    loadFn: (params) => DeviceService.list(params),
    createFn: (data) => DeviceService.create(data),
    updateFn: (id, data) => DeviceService.update(id, data),
    deleteFn: (id) => DeviceService.delete(id)
  },

  gateways: {
    title: '网关管理',
    resource: 'gateways',
    searchable: false,
    creatable: true,
    columns: [
      { prop: 'id', label: 'ID', width: 60 },
      { prop: 'name', label: '名称', minWidth: 150 },
      {
        prop: 'status', label: '状态', width: 90,
        tag: {
          label: (row) => row.status,
          type: (row) => row.status === 'online' ? 'success' : 'info'
        }
      },
      { prop: 'cpu_usage', label: 'CPU%', width: 80, align: 'center' },
      { prop: 'memory_usage', label: '内存%', width: 80, align: 'center' },
      { prop: 'battery_level', label: '电量%', width: 80, align: 'center' },
      {
        prop: 'devices', label: '设备数', width: 80, align: 'center',
        formatter: (row) => String(row.devices?.length || 0)
      }
    ],
    formFields: [
      { prop: 'name', label: '名称', type: 'input', required: true },
      {
        prop: 'status', label: '状态', type: 'select',
        options: [
          { label: '在线', value: 'online' },
          { label: '离线', value: 'offline' },
          { label: '维护中', value: 'maintenance' }
        ]
      },
      {
        prop: 'mode', label: '模式', type: 'select',
        options: [
          { label: '在线', value: 'online' },
          { label: '维护', value: 'maintenance' }
        ]
      }
    ],
    loadFn: (params) => GatewayService.list(params),
    createFn: (data) => GatewayService.create(data),
    updateFn: (id, data) => GatewayService.update(id, data),
    deleteFn: (id) => GatewayService.delete(id)
  },

  reminders: {
    title: '提醒管理',
    resource: 'reminders',
    searchable: false,
    creatable: true,
    paginated: true,
    columns: [
      { prop: 'id', label: 'ID', width: 60 },
      {
        prop: 'level', label: '等级', width: 80,
        tag: {
          label: (row) => REMINDER_LEVEL_MAP[row.level as ReminderLevel] || String(row.level),
          type: (row) => row.level >= 3 ? 'danger' : row.level >= 2 ? 'warning' : row.level >= 1 ? '' : 'info'
        }
      },
      { prop: 'reminder_type', label: '类型', width: 90 },
      { prop: 'message', label: '内容', minWidth: 250 },
      {
        prop: 'resolved', label: '状态', width: 80,
        tag: {
          label: (row) => row.resolved ? '已解决' : '待处理',
          type: (row) => row.resolved ? 'success' : 'danger'
        }
      },
      { prop: 'timestamp', label: '时间', width: 170, formatter: (row) => row.timestamp?.replace('T', ' ').slice(0, 19) }
    ],
    formFields: [
      {
        prop: 'level', label: '等级', type: 'select', required: true,
        options: Object.entries(REMINDER_LEVEL_MAP).map(([k, v]) => ({ label: v, value: Number(k) }))
      },
      {
        prop: 'reminder_type', label: '类型', type: 'select', required: true,
        options: [
          { label: '天气', value: 'weather' },
          { label: '垂钓', value: 'fishing' },
          { label: '安全', value: 'safety' },
          { label: '环境', value: 'environment' }
        ]
      },
      { prop: 'message', label: '内容', type: 'textarea', required: true },
      { prop: 'spot_id', label: '关联水域ID', type: 'number', required: true },
      { prop: 'resolved', label: '已解决', type: 'switch' }
    ],
    loadFn: (params) => ReminderService.list(params),
    createFn: (data) => ReminderService.create(data),
    deleteFn: (id) => ReminderService.delete(id)
  },

  notices: {
    title: '通知管理',
    resource: 'notices',
    searchable: true,
    creatable: true,
    paginated: true,
    columns: [
      { prop: 'id', label: 'ID', width: 60 },
      { prop: 'title', label: '标题', minWidth: 200 },
      { prop: 'content', label: '内容', minWidth: 300 },
      {
        prop: 'outdated', label: '状态', width: 80,
        tag: {
          label: (row) => row.outdated ? '已过期' : '有效',
          type: (row) => row.outdated ? 'info' : 'success'
        }
      },
      { prop: 'timestamp', label: '时间', width: 170, formatter: (row) => row.timestamp?.replace('T', ' ').slice(0, 19) }
    ],
    formFields: [
      { prop: 'title', label: '标题', type: 'input', required: true },
      { prop: 'content', label: '内容', type: 'textarea', required: true },
      { prop: 'outdated', label: '已过期', type: 'switch' }
    ],
    loadFn: (params) => NoticeService.list(params),
    createFn: (data) => NoticeService.create(data),
    updateFn: (id, data) => NoticeService.update(id, data),
    deleteFn: (id) => NoticeService.delete(id)
  },

  // ===== 社区/SFR 管理 =====
  users: {
    title: '用户管理',
    resource: 'users',
    searchable: false,
    creatable: false,
    columns: [
      { prop: 'id', label: 'ID', width: 60 },
      { prop: 'username', label: '用户名', minWidth: 120 },
      { prop: 'email', label: '邮箱', minWidth: 180 },
      { prop: 'phone', label: '手机', width: 130 },
      {
        prop: 'role', label: '角色', width: 100,
        tag: {
          label: (row) => ({ admin: '管理员', staff: '工作人员', user: '普通用户' }[row.role as string] || row.role),
          type: (row) => row.role === 'admin' ? 'danger' : row.role === 'staff' ? 'warning' : 'info'
        }
      },
      {
        prop: 'register_time', label: '注册时间', width: 170,
        formatter: (row) => row.register_time?.replace('T', ' ').slice(0, 19) || '--'
      }
    ],
    formFields: [
      { prop: 'username', label: '用户名', type: 'input', readonly: true },
      {
        prop: 'role', label: '角色', type: 'select', required: true,
        options: [
          { label: '管理员', value: 'admin' },
          { label: '工作人员', value: 'staff' },
          { label: '普通用户', value: 'user' }
        ]
      }
    ],
    loadFn: async (params) => {
      const res = await httpGet<{ results: any[]; total: number }>('/api/v2/users', params)
      return { data: res.results, total: res.total }
    },
    updateFn: async (id, data) => {
      await httpPatch(`/api/v2/users/${id}/role`, { role: data.role })
    },
    deleteFn: async (id) => {
      await httpDelete(`/api/v2/users/${id}`)
    }
  },

  posts: {
    title: '帖子管理',
    resource: 'posts',
    searchable: false,
    creatable: false,
    paginated: true,
    columns: [
      { prop: 'id', label: 'ID', width: 60 },
      { prop: 'title', label: '标题', minWidth: 200 },
      { prop: 'username', label: '作者', width: 100 },
      { prop: 'user_id', label: '用户ID', width: 80 },
      {
        prop: 'tag', label: '标签', width: 100,
        tag: { label: (row) => row.tag || '无' }
      },
      { prop: 'likes', label: '点赞', width: 70, align: 'center' },
      { prop: 'comments', label: '评论', width: 70, align: 'center' },
      {
        prop: 'body', label: '内容摘要', minWidth: 200,
        formatter: (row) => {
          const body = row.body || ''
          return body.length > 50 ? body.slice(0, 50) + '...' : body
        }
      }
    ],
    loadFn: async (params) => {
      return PostService.list(params)
    },
    deleteFn: async (id) => {
      await PostService.delete(id)
    }
  },

  comments: {
    title: '评论管理',
    resource: 'comments',
    searchable: false,
    creatable: false,
    columns: [
      { prop: 'id', label: 'ID', width: 60 },
      { prop: 'post_id', label: '帖子ID', width: 80 },
      { prop: 'username', label: '用户', width: 100 },
      { prop: 'user_id', label: '用户ID', width: 80 },
      { prop: 'body', label: '评论内容', minWidth: 300 }
    ],
    loadFn: async (params) => {
      return CommentService.list(params)
    },
    deleteFn: async (id) => {
      await CommentService.delete(id)
    }
  },

  suggestions: {
    title: '垂钓建议管理',
    resource: 'suggestions',
    searchable: false,
    creatable: true,
    columns: [
      { prop: 'id', label: 'ID', width: 60 },
      { prop: 'spot_id', label: '水域ID', width: 80 },
      { prop: 'suggestion_text', label: '建议内容', minWidth: 300 },
      { prop: 'score', label: '评分', width: 80, align: 'center' },
      { prop: 'timestamp', label: '时间', width: 170, formatter: (row) => row.timestamp?.replace('T', ' ').slice(0, 19) }
    ],
    formFields: [
      { prop: 'spot_id', label: '水域ID', type: 'number', required: true },
      { prop: 'suggestion_text', label: '建议内容', type: 'textarea', required: true },
      { prop: 'score', label: '评分', type: 'number' }
    ],
    loadFn: async (params) => {
      return SuggestionService.list(params)
    },
    createFn: async (data) => {
      return SuggestionService.create(data)
    },
    deleteFn: async (id) => {
      await SuggestionService.delete(id)
    }
  },

  fish_caught: {
    title: '渔获管理',
    resource: 'fish_caught',
    searchable: false,
    creatable: false,
    columns: [
      { prop: 'id', label: 'ID', width: 60 },
      { prop: 'record_id', label: '记录ID', width: 80 },
      { prop: 'fish_type', label: '鱼种', width: 100 },
      { prop: 'weight', label: '重量(kg)', width: 100, align: 'center' },
      { prop: 'bait_type', label: '饵料', width: 100 },
      { prop: 'fishing_depth', label: '深度(m)', width: 100, align: 'center' },
      { prop: 'caught_time', label: '时间', width: 170, formatter: (row) => row.caught_time?.replace('T', ' ').slice(0, 19) || '--' }
    ],
    loadFn: async (params) => {
      return FishCaughtResourceService.list(params)
    },
    deleteFn: async (_id) => {
      throw new Error('渔获记录不支持删除')
    }
  },

  iot_devices: {
    title: 'IoT设备管理',
    resource: 'iot_devices',
    searchable: false,
    creatable: false,
    columns: [
      { prop: 'device_id', label: '设备ID', minWidth: 200 },
      { prop: 'temperature', label: '温度(°C)', width: 100, align: 'center' },
      { prop: 'humidity', label: '湿度(%)', width: 100, align: 'center' },
      { prop: 'pulling', label: '拉力(N)', width: 100, align: 'center' },
      { prop: 'pressure', label: '气压(hPa)', width: 100, align: 'center' },
      {
        prop: 'last_update', label: '最后更新', width: 170,
        formatter: (row) => row.last_update?.replace('T', ' ').slice(0, 19) || '--'
      }
    ],
    loadFn: async (params) => {
      return IoTDeviceService.list(params)
    },
    deleteFn: async (_id) => {
      throw new Error('IoT设备不支持删除')
    }
  },

  fishing_records: {
    title: '垂钓记录管理',
    resource: 'fishing_records',
    searchable: false,
    creatable: false,
    columns: [
      { prop: 'id', label: 'ID', width: 60 },
      { prop: 'user_id', label: '用户ID', width: 80 },
      { prop: 'device_id', label: '设备ID', minWidth: 140, formatter: (row) => row.device_id || '--' },
      {
        prop: 'start_time', label: '开始时间', width: 160,
        formatter: (row) => row.start_time?.replace('T', ' ').slice(0, 16) || '--'
      },
      {
        prop: 'end_time', label: '结束时间', width: 160,
        formatter: (row) => row.end_time?.replace('T', ' ').slice(0, 16) || '--'
      },
      { prop: 'latitude', label: '纬度', width: 100 },
      { prop: 'longitude', label: '经度', width: 100 },
      {
        prop: 'fish_caught', label: '渔获数', width: 80, align: 'center',
        formatter: (row) => String(row.fish_caught?.length || 0)
      }
    ],
    loadFn: async (params) => {
      return FishingRecordResourceService.list(params)
    },
    deleteFn: async (id) => {
      await FishingRecordResourceService.delete(id)
    }
  }
}

const currentConfig = computed(() => resourceConfigs[activeMenu.value] || null)
</script>

<style scoped>
.admin-layout {
  display: flex;
  height: calc(100vh - 60px);
  overflow: hidden;
  background: #f0f2f5;
}

.admin-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
  overflow: hidden;
}

.admin-empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #fff;
}
</style>
