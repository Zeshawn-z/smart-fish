// ==================== 基础类型 ====================

export interface BaseModel {
  id: number
  created_at: string
  updated_at: string
}

// ==================== 用户 ====================

export interface User {
  id: number
  username: string
  role: 'user' | 'staff' | 'admin'
  phone: string
  email: string
  avatar?: string | null
  register_time: string
}

export interface LoginInput {
  username: string
  password: string
}

export interface RegisterInput {
  username: string
  password: string
  phone?: string
  email?: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  user: User
}

// ==================== 区域 ====================

export interface Region extends BaseModel {
  name: string
  province: string
  city: string
  description: string
  fishing_spots?: FishingSpot[]
}

// ==================== 垂钓水域 ====================

export type WaterType = 'lake' | 'river' | 'reservoir' | 'pond'
export type SpotStatus = 'open' | 'closed' | 'maintenance'

export interface FishingSpot extends BaseModel {
  name: string
  region_id: number
  region?: Region
  description: string
  latitude: number
  longitude: number
  water_type: WaterType
  capacity: number
  status: SpotStatus
  bound_device_id?: number
  bound_device?: Device
}

export interface PopularSpot extends FishingSpot {
  total_fishing_count: number
}

// ==================== 设备 ====================

export type DeviceStatus = 'online' | 'offline' | 'error'
export type DeviceType = 'environment' | 'underwater' | 'fishfinder'

export interface Device extends BaseModel {
  name: string
  gateway_id?: number
  gateway?: Gateway
  status: DeviceStatus
  description: string
  device_type: DeviceType
  fishing_count: number
  water_temp: number
  air_temp: number
  humidity: number
  pressure: number
  last_active_at?: string
}

// ==================== 网关 ====================

export type GatewayStatus = 'online' | 'offline' | 'maintenance'

export interface Gateway extends BaseModel {
  name: string
  status: GatewayStatus
  mode: string
  cpu_usage: number
  memory_usage: number
  disk_usage: number
  battery_level: number
  last_active_at?: string
  devices?: Device[]
}

// ==================== 历史数据 ====================

export interface HistoricalData extends BaseModel {
  spot_id: number
  fishing_count: number
  timestamp: string
}

// ==================== 环境数据 ====================

export interface EnvironmentData extends BaseModel {
  spot_id: number
  water_temp: number
  air_temp: number
  humidity: number
  pressure: number
  ph: number
  dissolved_oxygen: number
  turbidity: number
  timestamp: string
}

// ==================== 水质数据 ====================

export interface WaterQualityData extends BaseModel {
  device_id: number
  ph: number
  dissolved_oxygen: number
  turbidity: number
  timestamp: string
}

// ==================== 提醒 ====================

export type ReminderLevel = 0 | 1 | 2 | 3

export interface Reminder extends BaseModel {
  spot_id: number
  level: ReminderLevel
  reminder_type: string
  message: string
  timestamp: string
  resolved: boolean
  publicity: boolean
}

// ==================== 通知 ====================

export interface Notice extends BaseModel {
  title: string
  content: string
  timestamp: string
  outdated: boolean
  related_spots?: FishingSpot[]
}

// ==================== 垂钓建议 ====================

export interface FishingSuggestion extends BaseModel {
  spot_id: number
  fishing_spot?: FishingSpot
  user_id?: number
  suggestion_text: string
  score: number
  timestamp: string
}

// ==================== 系统概览 ====================

export interface SystemSummary {
  total_spots: number
  open_spots: number
  total_devices: number
  online_devices: number
  total_gateways: number
  online_gateways: number
  total_users: number
  active_reminders: number
  total_fishing_count: number
  recent_notices: number
  avg_water_temp: number
  avg_air_temp: number
}

// ==================== 区域环境数据 ====================

export interface RegionEnvItem {
  region_id: number
  region_name: string
  city: string
  spot_count: number
  water_temp: number
  air_temp: number
  humidity: number
  pressure: number
  ph: number
  dissolved_oxygen: number
  turbidity: number
  timestamp: string
}

export interface RegionEnvRecord {
  water_temp: number
  air_temp: number
  humidity: number
  pressure: number
  ph: number
  dissolved_oxygen: number
  turbidity: number
  timestamp: string
}

export interface RegionEnvHistory {
  region_id: number
  region_name: string
  city: string
  records: RegionEnvRecord[]
}

// ==================== 分页 ====================

export interface PaginatedResponse<T> {
  results: T[]
  total: number
  page: number
  page_size: number
}

// ==================== 工具类型 ====================

export const WATER_TYPE_MAP: Record<WaterType, string> = {
  lake: '湖泊',
  river: '河流',
  reservoir: '水库',
  pond: '鱼塘'
}

export const SPOT_STATUS_MAP: Record<SpotStatus, string> = {
  open: '开放',
  closed: '关闭',
  maintenance: '维护中'
}

export const REMINDER_LEVEL_MAP: Record<ReminderLevel, string> = {
  0: '信息',
  1: '提示',
  2: '重要',
  3: '紧急'
}

export const REMINDER_LEVEL_COLOR: Record<ReminderLevel, string> = {
  0: '#409eff',
  1: '#e6a23c',
  2: '#f56c6c',
  3: '#f56c6c'
}

export const DEVICE_TYPE_MAP: Record<DeviceType, string> = {
  environment: '环境监测',
  underwater: '水下感知',
  fishfinder: '探鱼辅助'
}

// ==================== 社区帖子 ====================

export interface Post extends BaseModel {
  user_id: number
  username?: string
  avatar?: string | null
  title: string
  body: string
  tag: string
  image_url?: string | string[] | null
  likes?: number
  comments?: number
  is_deleted?: boolean
  /** @deprecated 旧 v1 字段，新接口使用 id */
  post_id?: number
}

export interface PostListResponse {
  posts_list: Post[]
}

export interface PostDetail extends Post {
  image_urls?: string[]
}

// ==================== 评论 ====================

export interface Comment extends BaseModel {
  post_id: number
  user_id: number
  username?: string
  avatar?: string | null
  body: string
  is_deleted?: boolean
  /** 帖子详情嵌入时返回的评论点赞数 */
  likes?: number
  /** 帖子详情嵌入时返回的子评论列表 */
  sub_comments?: CommentOnComment[]
  /** @deprecated 旧 v1 字段，新接口使用 id */
  comment_id?: number
}

export interface CommentListResponse {
  comments_list: Comment[]
}

export interface CommentOnComment {
  coc_id: number
  comment_id: number
  user_id: number
  username: string
  avatar?: string | null
  body: string
  to_coc_id: number | null
  to_user_id: number | null
  to_username: string | null
}

export interface CocListResponse {
  comment_id: string
  comments: CommentOnComment[]
}

// ==================== 垂钓记录 ====================

export interface FishCaught extends BaseModel {
  record_id: number
  caught_time: string
  fish_type: string
  weight: number
  bait_type: string
  bait_weight: number
  fishing_depth: number
  image_url?: string | null
  /** @deprecated 旧 v1 字段，新接口使用 id */
  fish_id?: number
}

export interface IoTDeviceData {
  device_id: string
  temperature: number
  humidity: number
  pulling: number
  pressure: number
  gps_info?: string
  gpsInfo?: string
  imu_data: string
  last_update: string
}

export interface FishingRecord extends BaseModel {
  user_id: number
  device_id: string
  start_time: string
  end_time: string
  latitude: number
  longitude: number
  fish_caught?: FishCaught[]
  device_data?: IoTDeviceData
  /** @deprecated 旧 v1 字段 */
  caught?: FishCaught[]
  /** @deprecated 旧 v1 字段，新接口使用 id */
  record_id?: number
}

export interface FishingRecordListResponse {
  records: FishingRecord[]
}

// ==================== 后台管理配置类型 ====================

/** 表格列配置 */
export interface ColumnConfig {
  prop: string
  label: string
  width?: number | string
  minWidth?: number | string
  align?: 'left' | 'center' | 'right'
  showOverflow?: boolean
  formatter?: (row: any) => string
  render?: (row: any) => string
  tag?: {
    label?: (row: any) => string
    type?: (row: any) => '' | 'success' | 'warning' | 'danger' | 'info'
  }
}

/** 表单字段配置（用于新增/编辑） */
export interface FormFieldConfig {
  /** 字段标识（对应数据属性名） */
  prop: string
  /** 显示标签 */
  label: string
  /** 表单控件类型 */
  type: 'input' | 'textarea' | 'number' | 'select' | 'switch' | 'datetime'
  /** select 选项 */
  options?: { label: string; value: string | number }[]
  /** 是否必填 */
  required?: boolean
  /** 占位提示 */
  placeholder?: string
  /** 是否只读 */
  readonly?: boolean
}

/** 资源管理配置 */
export interface ResourceConfig {
  /** 页面标题 */
  title: string
  /** 资源标识 (用于缓存键和路由，如 "regions") */
  resource: string
  /** 表格列配置 */
  columns: ColumnConfig[]
  /** 加载数据 */
  loadFn: (params: Record<string, any>) => Promise<{ data: any[]; total: number }>
  /** 删除单条记录 */
  deleteFn: (id: number) => Promise<void>
  /** 更新单条记录（如果提供则启用编辑功能） */
  updateFn?: (id: number, data: Record<string, any>) => Promise<any>
  /** 创建记录的函数 */
  createFn?: (data: Record<string, any>) => Promise<any>
  /** 表单字段配置（用于新增/编辑对话框） */
  formFields?: FormFieldConfig[]
  /** 是否支持搜索 */
  searchable?: boolean
  /** 是否可新增 */
  creatable?: boolean
  /** 是否分页 */
  paginated?: boolean
  /** 操作列宽度 */
  actionWidth?: number
}
