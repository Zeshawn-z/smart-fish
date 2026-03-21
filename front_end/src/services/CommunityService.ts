/**
 * 社区服务 —— 全部使用 v2 接口
 *
 * 标准 CRUDL 由 createResourceService 工厂自动生成
 * 点赞、子评论、图片上传通过 extend 或独立方法实现
 */

import { createResourceService } from './createResourceService'
import { httpGet, httpPost, httpDelete } from '@/network/http'
import type { Post, Comment as CommentType, CocListResponse } from '@/types'

/** 帖子服务（工厂模式） */
export const PostService = createResourceService({
  name: 'posts',
  paginated: true,
  model: {} as Post,
  listParams: {} as { tag?: string; search?: string; user_id?: number; page?: number; page_size?: number },
  extend: (ctx) => ({
    /** 获取帖子详情（含完整图片列表） */
    async getDetail(id: number) {
      return ctx.http.get<Post & { image_urls?: string[] }>(`${ctx.baseURL}/${id}`)
    }
  })
})

/** 评论服务（工厂模式） */
export const CommentService = createResourceService({
  name: 'comments',
  model: {} as CommentType,
  paginated: true,
  listParams: {} as { post_id?: number; page?: number; page_size?: number }
})

/**
 * 社区服务 —— 聚合帖子 + 评论 + 点赞 + 子评论 + 图片上传
 * 全部使用 v2 接口
 */
export const CommunityService = {
  // ===== 帖子 CRUD =====

  async getPosts(params?: { tag?: string; search?: string; page?: number; page_size?: number }) {
    const res = await PostService.list(params)
    return { posts_list: res.data, total: res.total }
  },

  async getMyPosts(userId: number) {
    const res = await PostService.list({ user_id: userId })
    return { posts_list: res.data, total: res.total }
  },

  async getPost(postId: number): Promise<Post & { image_urls?: string[] }> {
    return PostService.getDetail(postId)
  },

  async createPost(data: { title: string; body: string; tag?: string }): Promise<Post> {
    return PostService.create(data as Partial<Post>)
  },

  // ===== 评论 =====

  async getComments(postId: number) {
    // 绕过缓存直接请求，确保评论数据实时准确
    CommentService.invalidateCache()
    const res = await CommentService.list({ post_id: postId })
    return { comments_list: res.data }
  },

  async createComment(postId: number, body: string): Promise<CommentType> {
    return CommentService.create({ post_id: postId, body } as Partial<CommentType>)
  },

  // ===== 子评论（楼中楼）=====

  async getSubComments(commentId: number): Promise<CocListResponse> {
    return httpGet<CocListResponse>(`/api/v2/comments/${commentId}/replies`)
  },

  async createSubComment(commentId: number, body: string, toCocId?: number): Promise<{ message: string; coc_id: number }> {
    return httpPost(`/api/v2/comments/${commentId}/replies`, { body, to_coc_id: toCocId || undefined })
  },

  // ===== 点赞 =====

  async getPostLikes(postId: number): Promise<{ likes: number; liked: boolean }> {
    return httpGet(`/api/v2/posts/${postId}/like`)
  },

  async likePost(postId: number): Promise<void> {
    await httpPost(`/api/v2/posts/${postId}/like`)
  },

  async unlikePost(postId: number): Promise<void> {
    await httpDelete(`/api/v2/posts/${postId}/like`)
  },

  async getCommentLikes(commentId: number): Promise<{ likes: number }> {
    return httpGet(`/api/v2/comments/${commentId}/like`)
  },

  async likeComment(commentId: number): Promise<void> {
    await httpPost(`/api/v2/comments/${commentId}/like`)
  },

  async unlikeComment(commentId: number): Promise<void> {
    await httpDelete(`/api/v2/comments/${commentId}/like`)
  },

  // ===== 图片上传 =====

  async uploadPostImage(postId: number, file: File): Promise<{ message: string; image_id: number; image_url: string }> {
    const form = new FormData()
    form.append('entity_type', 'post')
    form.append('entity_id', String(postId))
    form.append('file', file)
    return httpPost('/api/v2/upload/image', form, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  },

  async uploadAvatarImage(file: File): Promise<{ message: string; image_id: number; image_url: string }> {
    const form = new FormData()
    form.append('entity_type', 'avatar')
    form.append('file', file)
    return httpPost('/api/v2/upload/image', form, {
      headers: { 'Content-Type': 'multipart/form-data' }
    })
  }
}
