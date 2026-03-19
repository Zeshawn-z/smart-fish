/**
 * Mock 数据服务 - 当后端不可用时提供本地模拟数据
 * 在后端启动后可以删除此文件并切换为真实 API
 */
import type {
  SystemSummary, FishingSpot, PopularSpot, Region, Reminder,
  Notice, HistoricalData, EnvironmentData, FishingSuggestion,
  Device, PaginatedResponse
} from '@/types'

// ===== 时间工具 =====
const hoursAgo = (h: number) => new Date(Date.now() - h * 3600000).toISOString()
const minutesAgo = (m: number) => new Date(Date.now() - m * 60000).toISOString()
const daysAgo = (d: number) => new Date(Date.now() - d * 86400000).toISOString()

// ===== 模拟区域 =====
export const MOCK_REGIONS: Region[] = [
  { id: 1, name: '哈尔滨松花江段', province: '黑龙江', city: '哈尔滨', description: '松花江哈尔滨城区段，水域开阔', created_at: daysAgo(30), updated_at: daysAgo(1) },
  { id: 2, name: '太阳岛垂钓区', province: '黑龙江', city: '哈尔滨', description: '太阳岛景区内垂钓专区', created_at: daysAgo(30), updated_at: daysAgo(1) },
  { id: 3, name: '镜泊湖风景区', province: '黑龙江', city: '牡丹江', description: '国家5A级景区，水质优良', created_at: daysAgo(30), updated_at: daysAgo(1) },
  { id: 4, name: '兴凯湖湿地', province: '黑龙江', city: '鸡西', description: '中俄界湖，淡水鱼资源丰富', created_at: daysAgo(30), updated_at: daysAgo(1) },
]

// ===== 模拟设备 =====
export const MOCK_DEVICES: Device[] = [
  { id: 1, name: 'ENV-001 多模态环境监测', status: 'online', device_type: 'environment', gateway_id: 1, description: '', fishing_count: 7, water_temp: 18.5, air_temp: 22.3, humidity: 65, pressure: 1013.2, last_active_at: minutesAgo(2), created_at: daysAgo(30), updated_at: minutesAgo(2) },
  { id: 2, name: 'ENV-002 多模态环境监测', status: 'online', device_type: 'environment', gateway_id: 1, description: '', fishing_count: 5, water_temp: 17.8, air_temp: 21.6, humidity: 68, pressure: 1012.8, last_active_at: minutesAgo(5), created_at: daysAgo(30), updated_at: minutesAgo(5) },
  { id: 3, name: 'UW-001 水下态势感知', status: 'online', device_type: 'underwater', gateway_id: 2, description: '', fishing_count: 0, water_temp: 16.2, air_temp: 0, humidity: 0, pressure: 0, last_active_at: minutesAgo(1), created_at: daysAgo(30), updated_at: minutesAgo(1) },
  { id: 4, name: 'FF-001 探鱼辅助', status: 'offline', device_type: 'fishfinder', gateway_id: 2, description: '', fishing_count: 0, water_temp: 0, air_temp: 0, humidity: 0, pressure: 0, last_active_at: hoursAgo(3), created_at: daysAgo(30), updated_at: hoursAgo(3) },
  { id: 5, name: 'ENV-003 多模态环境监测', status: 'online', device_type: 'environment', gateway_id: 3, description: '', fishing_count: 12, water_temp: 19.1, air_temp: 23.5, humidity: 62, pressure: 1014.0, last_active_at: minutesAgo(3), created_at: daysAgo(30), updated_at: minutesAgo(3) },
  { id: 6, name: 'ENV-004 多模态环境监测', status: 'online', device_type: 'environment', gateway_id: 3, description: '', fishing_count: 8, water_temp: 17.0, air_temp: 20.8, humidity: 72, pressure: 1011.5, last_active_at: minutesAgo(8), created_at: daysAgo(30), updated_at: minutesAgo(8) },
]

// ===== 模拟垂钓水域 =====
export const MOCK_SPOTS: FishingSpot[] = [
  { id: 1, name: '松花江北岸休闲钓场', region_id: 1, region: MOCK_REGIONS[0], description: '松花江北岸优质钓场，交通便利，水质清澈，鱼种丰富', latitude: 45.75, longitude: 126.65, water_type: 'river', capacity: 80, status: 'open', bound_device_id: 1, bound_device: MOCK_DEVICES[0], created_at: daysAgo(30), updated_at: hoursAgo(1) },
  { id: 2, name: '太阳岛湖心垂钓', region_id: 2, region: MOCK_REGIONS[1], description: '太阳岛核心区域，环境优美，适合家庭垂钓', latitude: 45.77, longitude: 126.62, water_type: 'lake', capacity: 50, status: 'open', bound_device_id: 2, bound_device: MOCK_DEVICES[1], created_at: daysAgo(30), updated_at: hoursAgo(2) },
  { id: 3, name: '镜泊湖深水区', region_id: 3, region: MOCK_REGIONS[2], description: '镜泊湖深水垂钓区，大鱼多，需一定技术', latitude: 43.88, longitude: 128.98, water_type: 'lake', capacity: 120, status: 'open', bound_device_id: 5, bound_device: MOCK_DEVICES[4], created_at: daysAgo(30), updated_at: hoursAgo(3) },
  { id: 4, name: '兴凯湖湿地钓场', region_id: 4, region: MOCK_REGIONS[3], description: '兴凯湖边缘湿地区域，生态丰富', latitude: 45.32, longitude: 132.40, water_type: 'lake', capacity: 60, status: 'open', bound_device_id: 6, bound_device: MOCK_DEVICES[5], created_at: daysAgo(30), updated_at: hoursAgo(4) },
  { id: 5, name: '呼兰河水库', region_id: 1, region: MOCK_REGIONS[0], description: '呼兰河上游水库，水域平静，适合初学者', latitude: 45.95, longitude: 126.58, water_type: 'reservoir', capacity: 40, status: 'maintenance', bound_device_id: 3, bound_device: MOCK_DEVICES[2], created_at: daysAgo(30), updated_at: daysAgo(2) },
  { id: 6, name: '松北鱼塘乐园', region_id: 1, region: MOCK_REGIONS[0], description: '专业鱼塘，提供钓具租赁和指导服务', latitude: 45.80, longitude: 126.55, water_type: 'pond', capacity: 30, status: 'closed', created_at: daysAgo(30), updated_at: daysAgo(5) },
]

// ===== 模拟热门水域 =====
export const MOCK_POPULAR: PopularSpot[] = MOCK_SPOTS.filter(s => s.status === 'open').map((s, i) => ({
  ...s,
  total_fishing_count: [12, 9, 7, 5][i] || 3
}))

// ===== 模拟提醒 =====
export const MOCK_REMINDERS: Reminder[] = [
  { id: 1, spot_id: 1, level: 3, reminder_type: 'weather', message: '松花江北岸今日午后有雷阵雨，请注意安全撤离', timestamp: minutesAgo(30), resolved: false, publicity: true, created_at: minutesAgo(30), updated_at: minutesAgo(30) },
  { id: 2, spot_id: 3, level: 2, reminder_type: 'capacity', message: '镜泊湖深水区当前垂钓人数接近上限（31/40），请合理安排', timestamp: hoursAgo(1), resolved: false, publicity: true, created_at: hoursAgo(1), updated_at: hoursAgo(1) },
  { id: 3, spot_id: 2, level: 1, reminder_type: 'environment', message: '太阳岛湖心水温偏低（16.2°C），建议调整钓饵策略', timestamp: hoursAgo(2), resolved: false, publicity: true, created_at: hoursAgo(2), updated_at: hoursAgo(2) },
  { id: 4, spot_id: 4, level: 0, reminder_type: 'info', message: '兴凯湖湿地钓场设备巡检完成，一切正常', timestamp: hoursAgo(4), resolved: true, publicity: true, created_at: hoursAgo(4), updated_at: hoursAgo(3) },
]

// ===== 模拟通知 =====
export const MOCK_NOTICES: Notice[] = [
  { id: 1, title: '智钓蓝海平台春季功能更新通知', content: '新增水质实时监测仪表盘，优化垂钓建议算法，提升设备数据传输稳定性。', timestamp: daysAgo(1), outdated: false, created_at: daysAgo(1), updated_at: daysAgo(1) },
  { id: 2, title: '黑龙江省休闲垂钓安全管理新规', content: '根据省渔政部门最新通知，所有开放水域需配备安全警示标识和应急救援设备。', timestamp: daysAgo(3), outdated: false, created_at: daysAgo(3), updated_at: daysAgo(3) },
  { id: 3, title: '镜泊湖垂钓区开放时间调整', content: '自即日起，镜泊湖深水区开放时间调整为 06:00-18:00，请垂钓爱好者注意。', timestamp: daysAgo(5), outdated: false, created_at: daysAgo(5), updated_at: daysAgo(5) },
]

// ===== 模拟垂钓建议 =====
export const MOCK_SUGGESTIONS: FishingSuggestion[] = [
  { id: 1, spot_id: 1, fishing_spot: MOCK_SPOTS[0], suggestion_text: '今日松花江水温适宜（18.5°C），建议使用蚯蚓或红虫钓底层鲤鱼，最佳时段 6:00-9:00', score: 92, timestamp: hoursAgo(1), created_at: hoursAgo(1), updated_at: hoursAgo(1) },
  { id: 2, spot_id: 3, fishing_spot: MOCK_SPOTS[2], suggestion_text: '镜泊湖深水区近期白鲢活跃，推荐浮钓配合雾化饵料，水深建议 3-5 米', score: 88, timestamp: hoursAgo(2), created_at: hoursAgo(2), updated_at: hoursAgo(2) },
]

// ===== 生成历史数据（24h 每30分钟一条） =====
function generateHistorical(spotId: number): HistoricalData[] {
  const data: HistoricalData[] = []
  const now = Date.now()
  const baseCounts: Record<number, number> = { 1: 20, 2: 12, 3: 25, 4: 10, 5: 8, 6: 0 }
  const base = baseCounts[spotId] ?? 15

  for (let i = 48; i >= 0; i--) {
    const ts = now - i * 30 * 60000
    const hour = new Date(ts).getHours()
    let count = base

    // 早晚高峰
    if ((hour >= 6 && hour <= 9) || (hour >= 16 && hour <= 19)) {
      count += Math.floor(Math.random() * 15) + 10
    }
    // 午间
    if (hour >= 11 && hour <= 14) {
      count += Math.floor(Math.random() * 8) + 5
    }
    // 深夜
    if (hour >= 22 || hour <= 4) {
      count = Math.max(0, Math.floor(count * 0.2))
    }

    count = Math.max(0, count + Math.floor(Math.random() * 6) - 3)

    data.push({
      id: spotId * 100 + (48 - i),
      spot_id: spotId,
      fishing_count: count,
      timestamp: new Date(ts).toISOString(),
      created_at: new Date(ts).toISOString(),
      updated_at: new Date(ts).toISOString()
    })
  }
  return data
}

// ===== 生成环境数据（24h 每30分钟一条） =====
function generateEnvironment(spotId: number): EnvironmentData[] {
  const data: EnvironmentData[] = []
  const now = Date.now()
  const baseTemp: Record<number, number> = { 1: 18.5, 2: 17.8, 3: 19.1, 4: 17.0, 5: 16.2, 6: 15.0 }
  const baseWaterTemp = baseTemp[spotId] ?? 17

  for (let i = 48; i >= 0; i--) {
    const ts = now - i * 30 * 60000
    const hour = new Date(ts).getHours()

    // 温度日变化
    const tempOffset = Math.sin((hour - 6) / 24 * Math.PI * 2) * 3
    const waterTemp = Math.round((baseWaterTemp + tempOffset * 0.5 + (Math.random() - 0.5)) * 10) / 10
    const airTemp = Math.round((baseWaterTemp + 3 + tempOffset + (Math.random() - 0.5) * 2) * 10) / 10
    const humidity = Math.round(65 + Math.sin((hour - 14) / 24 * Math.PI * 2) * 10 + (Math.random() - 0.5) * 5)

    data.push({
      id: spotId * 100 + (48 - i),
      spot_id: spotId,
      water_temp: waterTemp,
      air_temp: airTemp,
      humidity: Math.min(95, Math.max(30, humidity)),
      pressure: Math.round((1013 + (Math.random() - 0.5) * 4) * 10) / 10,
      ph: Math.round((7.2 + (Math.random() - 0.5) * 0.6) * 10) / 10,
      dissolved_oxygen: Math.round((8.5 + (Math.random() - 0.5) * 2) * 10) / 10,
      turbidity: Math.round((15 + (Math.random() - 0.5) * 10) * 10) / 10,
      timestamp: new Date(ts).toISOString(),
      created_at: new Date(ts).toISOString(),
      updated_at: new Date(ts).toISOString()
    })
  }
  return data
}

// ===== 模拟系统概览 =====
export const MOCK_SUMMARY: SystemSummary = {
  total_spots: 6,
  open_spots: 4,
  total_devices: 6,
  online_devices: 5,
  total_gateways: 3,
  online_gateways: 3,
  total_users: 128,
  active_reminders: 3,
  total_fishing_count: 55,
  recent_notices: 3,
  avg_water_temp: 17.7,
  avg_air_temp: 21.6
}

// ===== 对外 API 模拟 =====
const delay = (ms = 300) => new Promise(r => setTimeout(r, ms + Math.random() * 200))

export const MockAPI = {
  async getSummary(): Promise<SystemSummary> {
    await delay(200)
    return { ...MOCK_SUMMARY }
  },

  async getRegions(): Promise<Region[]> {
    await delay()
    return [...MOCK_REGIONS]
  },

  async getProvinces(): Promise<string[]> {
    await delay(100)
    return [...new Set(MOCK_REGIONS.map(r => r.province))]
  },

  async getSpots(params?: any): Promise<PaginatedResponse<FishingSpot>> {
    await delay()
    let list = [...MOCK_SPOTS]
    if (params?.region_id) list = list.filter(s => s.region_id === params.region_id)
    if (params?.status) list = list.filter(s => s.status === params.status)
    if (params?.water_type) list = list.filter(s => s.water_type === params.water_type)
    if (params?.search) {
      const q = params.search.toLowerCase()
      list = list.filter(s => s.name.toLowerCase().includes(q) || s.description.toLowerCase().includes(q))
    }
    return { results: list, total: list.length, page: 1, page_size: 20 }
  },

  async getSpot(id: number): Promise<FishingSpot | null> {
    await delay(200)
    return MOCK_SPOTS.find(s => s.id === id) || null
  },

  async getPopularSpots(limit = 5): Promise<PopularSpot[]> {
    await delay()
    return MOCK_POPULAR.slice(0, limit)
  },

  async getHistorical(spotId: number): Promise<HistoricalData[]> {
    await delay()
    return generateHistorical(spotId)
  },

  async getEnvironment(spotId: number): Promise<EnvironmentData[]> {
    await delay()
    return generateEnvironment(spotId)
  },

  async getReminders(): Promise<Reminder[]> {
    await delay()
    return [...MOCK_REMINDERS]
  },

  async getNotices(): Promise<Notice[]> {
    await delay()
    return [...MOCK_NOTICES]
  },

  async getSuggestions(): Promise<FishingSuggestion[]> {
    await delay()
    return [...MOCK_SUGGESTIONS]
  },

  async getFavorites(): Promise<FishingSpot[]> {
    await delay(200)
    return [MOCK_SPOTS[0], MOCK_SPOTS[2]]
  },

  async toggleFavorite(spotId: number): Promise<{ message: string; favorited: boolean }> {
    await delay(200)
    return { message: 'ok', favorited: true }
  }
}
