/**
 * 资源服务工厂 —— 统一生成 CRUDL + ID 级别缓存 + 自定义扩展
 *
 * 设计思路：
 *   1. 标准 CRUDL 由工厂自动生成，无需每个 Service 手写
 *   2. 单级缓存：所有数据以 id 为 key 存入一个 Map，列表查询的结果也拆进去
 *      - list() 缓存的是「这组查询参数 → 对应的 id 列表 + total」
 *      - 实际数据只在 id → item 这一层存
 *   3. 写操作（create/update/delete）自动失效相关缓存
 *   4. 通过 extend 回调注入自定义方法，可复用工厂内的 http 和缓存设施
 *   5. URL 自动拼接：默认 /api/v2/{name}，也可通过 baseURL 覆盖
 */

import { httpGet, httpPost, httpPut, httpPatch, httpDelete } from '@/network/http'
import type { BaseModel, PaginatedResponse } from '@/types'

// ==================== 类型定义 ====================

/** 缓存配置 */
export interface CacheOptions {
  /** 缓存 TTL（毫秒），0 = 不缓存。默认 60_000 (60s) */
  ttl?: number
}

/** 资源服务工厂配置 */
export interface ResourceServiceOptions<
  T extends BaseModel,
  TListParams extends Record<string, any> = Record<string, any>
> {
  /** 资源名称（用于自动拼接 URL 和缓存键，如 'regions'） */
  name: string
  /** API 前缀覆盖。不指定时自动拼为 /api/v2/{name} */
  baseURL?: string
  /** 列表响应是否分页包裹（默认 false → 裸数组） */
  paginated?: boolean
  /** 缓存配置 */
  cache?: CacheOptions
  /**
   * 类型标记：用 `{} as MyModel` 标注资源类型 T。
   * 运行时不使用，仅供 TypeScript 推断。
   */
  model: T
  /**
   * 类型标记：用 `{} as MyListParams` 标注列表查询参数类型。
   * 运行时不使用，仅供 TypeScript 推断。可不填（默认 Record<string, any>）。
   */
  listParams?: TListParams
}

/** 传给 extend 回调的上下文工具集 */
export interface ServiceContext<T extends BaseModel> {
  /** 该资源的 API 基础路径 */
  baseURL: string
  http: {
    get: typeof httpGet
    post: typeof httpPost
    put: typeof httpPut
    patch: typeof httpPatch
    delete: typeof httpDelete
  }
  cache: {
    /** 获取缓存中的单条数据 */
    getItem: (id: number) => T | null
    /** 写入单条缓存 */
    setItem: (id: number, data: T) => void
    /** 移除单条缓存 */
    removeItem: (id: number) => void
    /** 清除所有列表查询索引（下次 list 会重新请求） */
    invalidateList: () => void
    /** 清除全部缓存（数据 + 列表索引） */
    invalidateAll: () => void
  }
}

/** 标准 CRUDL 方法签名 */
export interface ResourceServiceBase<
  T extends BaseModel,
  TListParams extends Record<string, any> = Record<string, any>
> {
  /** 列表查询（带缓存） */
  list(params?: TListParams): Promise<{ data: T[]; total: number }>
  /** 获取单条（优先从缓存取） */
  get(id: number): Promise<T>
  /** 创建，自动失效列表缓存 */
  create(data: Partial<T>): Promise<T>
  /** 更新，自动失效列表 + 对应 ID 缓存 */
  update(id: number, data: Partial<T>): Promise<T>
  /** 删除，自动失效列表 + 对应 ID 缓存 */
  delete(id: number): Promise<void>
  /** 手动失效所有缓存 */
  invalidateCache(): void
}

/** 最终服务类型 = 标准 CRUDL + 自定义扩展 */
export type ResourceService<
  T extends BaseModel,
  TListParams extends Record<string, any> = Record<string, any>,
  TExt extends Record<string, (...args: any[]) => any> = Record<string, never>
> = ResourceServiceBase<T, TListParams> & TExt

// ==================== 缓存管理器（单级 ID 缓存） ====================

/** 列表查询索引：记录某组查询参数对应哪些 id + total */
interface ListIndex {
  ids: number[]
  total: number
  timestamp: number
}

class ResourceCache<T extends BaseModel> {
  /** ID → 数据 */
  private items = new Map<number, { data: T; timestamp: number }>()
  /** 查询参数序列化 key → 对应的 id 列表 */
  private listIndex = new Map<string, ListIndex>()
  private ttl: number

  constructor(ttl: number) {
    this.ttl = ttl
  }

  // ---- 工具 ----

  private isExpired(timestamp: number): boolean {
    return this.ttl > 0 && Date.now() - timestamp > this.ttl
  }

  private serializeParams(params?: Record<string, any>): string {
    if (!params) return '__all__'
    const sorted = Object.keys(params).sort()
      .filter(k => params[k] !== undefined && params[k] !== null && params[k] !== '')
      .map(k => `${k}=${params[k]}`)
    return sorted.join('&') || '__all__'
  }

  // ---- 单条操作 ----

  getItem(id: number): T | null {
    if (this.ttl <= 0) return null
    const entry = this.items.get(id)
    if (!entry) return null
    if (this.isExpired(entry.timestamp)) {
      this.items.delete(id)
      return null
    }
    return entry.data
  }

  setItem(id: number, data: T): void {
    if (this.ttl <= 0) return
    this.items.set(id, { data, timestamp: Date.now() })
  }

  removeItem(id: number): void {
    this.items.delete(id)
  }

  // ---- 列表操作 ----

  /** 尝试从缓存还原列表结果。列表索引有效 + 所有 item 都还在 → 命中 */
  getList(params?: Record<string, any>): { data: T[]; total: number } | null {
    if (this.ttl <= 0) return null
    const key = this.serializeParams(params)
    const index = this.listIndex.get(key)
    if (!index) return null
    if (this.isExpired(index.timestamp)) {
      this.listIndex.delete(key)
      return null
    }
    // 尝试从 item 缓存中组装完整列表
    const data: T[] = []
    for (const id of index.ids) {
      const item = this.getItem(id)
      if (!item) {
        // 有一条缺失就放弃，重新请求
        this.listIndex.delete(key)
        return null
      }
      data.push(item)
    }
    return { data, total: index.total }
  }

  /** 缓存列表结果：拆解每个 item 到 ID Map，同时记录索引 */
  setList(params: Record<string, any> | undefined, data: T[], total: number): void {
    if (this.ttl <= 0) return
    const now = Date.now()
    const ids: number[] = []
    for (const item of data) {
      ids.push(item.id)
      this.items.set(item.id, { data: item, timestamp: now })
    }
    const key = this.serializeParams(params)
    this.listIndex.set(key, { ids, total, timestamp: now })
  }

  /** 清除所有列表索引（数据保留，下次 list 重新请求但 get 仍可命中） */
  invalidateList(): void {
    this.listIndex.clear()
  }

  /** 全清 */
  invalidateAll(): void {
    this.items.clear()
    this.listIndex.clear()
  }
}

// ==================== 工厂函数 ====================

/** 带扩展的重载：TExt 从 extend 返回值推断 */
export function createResourceService<
  T extends BaseModel,
  TListParams extends Record<string, any>,
  TExt extends Record<string, (...args: any[]) => any>
>(
  options: ResourceServiceOptions<T, TListParams> & {
    extend: (ctx: ServiceContext<T>) => TExt
  }
): ResourceServiceBase<T, TListParams> & TExt

/** 不带扩展的重载 */
export function createResourceService<
  T extends BaseModel,
  TListParams extends Record<string, any> = Record<string, any>
>(
  options: ResourceServiceOptions<T, TListParams>
): ResourceServiceBase<T, TListParams>

/** 实现 */
export function createResourceService<
  T extends BaseModel,
  TListParams extends Record<string, any> = Record<string, any>
>(
  options: ResourceServiceOptions<T, TListParams> & {
    extend?: (ctx: ServiceContext<T>) => Record<string, (...args: any[]) => any>
  }
): any {

  const {
    name,
    baseURL = `/api/v2/${name}`,
    paginated = false,
    cache: cacheOpts = {}
  } = options

  const ttl = cacheOpts.ttl ?? 60_000
  const cache = new ResourceCache<T>(ttl)

  // ---- 构建上下文工具集 ----
  const ctx: ServiceContext<T> = {
    baseURL,
    http: {
      get: httpGet,
      post: httpPost,
      put: httpPut,
      patch: httpPatch,
      delete: httpDelete
    },
    cache: {
      getItem: (id) => cache.getItem(id),
      setItem: (id, data) => cache.setItem(id, data),
      removeItem: (id) => cache.removeItem(id),
      invalidateList: () => cache.invalidateList(),
      invalidateAll: () => cache.invalidateAll()
    }
  }

  // ---- 标准 CRUDL ----
  const base: ResourceServiceBase<T, TListParams> = {
    async list(params?: TListParams) {
      const cached = cache.getList(params as Record<string, any>)
      if (cached) return cached

      let data: T[]
      let total: number
      if (paginated) {
        const res = await httpGet<PaginatedResponse<T>>(baseURL, params)
        data = res.results
        total = res.total
      } else {
        data = await httpGet<T[]>(baseURL, params)
        total = data.length
      }

      cache.setList(params as Record<string, any>, data, total)
      return { data, total }
    },

    async get(id: number) {
      // 优先从缓存取（可能是之前 list 拆进来的）
      const cached = cache.getItem(id)
      if (cached) return cached

      const data = await httpGet<T>(`${baseURL}/${id}`)
      cache.setItem(id, data)
      return data
    },

    async create(data: Partial<T>) {
      const result = await httpPost<T>(baseURL, data)
      cache.invalidateList()
      return result
    },

    async update(id: number, data: Partial<T>) {
      const result = await httpPut<T>(`${baseURL}/${id}`, data)
      cache.removeItem(id)
      cache.invalidateList()
      return result
    },

    async delete(id: number) {
      await httpDelete(`${baseURL}/${id}`)
      cache.removeItem(id)
      cache.invalidateList()
    },

    invalidateCache() {
      cache.invalidateAll()
    }
  }

  // ---- 合并自定义扩展 ----
  const extensions = options.extend ? options.extend(ctx) : ({} as TExt)

  return Object.assign(base, extensions) as any
}
