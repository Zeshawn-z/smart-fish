// 工厂
export { createResourceService } from './createResourceService'
export type { ResourceServiceOptions, ResourceServiceBase, ResourceService, ServiceContext, CacheOptions } from './createResourceService'

// 业务服务
export { AuthService } from './AuthService'
export { RegionService } from './RegionService'
export { FishingSpotService } from './FishingSpotService'
export { DeviceService } from './DeviceService'
export { GatewayService } from './GatewayService'
export { ReminderService } from './ReminderService'
export { NoticeService } from './NoticeService'
export { SummaryService } from './SummaryService'
export { SuggestionService } from './SuggestionService'

// 社区 & 垂钓记录（全部 v2 接口）
export { PostService, CommentService, CommunityService } from './CommunityService'
export { FishingRecordResourceService, FishCaughtResourceService, FishingRecordService } from './FishingRecordService'
export { IoTDeviceService } from './IoTDeviceService'
