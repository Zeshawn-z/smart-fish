package services

import (
	"fmt"
	"time"

	"smart-fish/back_end/cache"
	"smart-fish/back_end/dao"
	"smart-fish/back_end/models"
)

// ==================== DTO ====================

// PostDTO 帖子响应体
type PostDTO struct {
	ID        uint    `json:"id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	UserID    uint    `json:"user_id"`
	Username  string  `json:"username"`
	Avatar    *string `json:"avatar"`
	Title     string  `json:"title"`
	Body      string  `json:"body"`
	Tag       string  `json:"tag"`
	ImageURL  *string `json:"image_url"`
	Likes     int64   `json:"likes"`
	Comments  int64   `json:"comments"`
}

// CommentDTO 评论响应体
type CommentDTO struct {
	ID        uint    `json:"id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
	PostID    uint    `json:"post_id"`
	UserID    uint    `json:"user_id"`
	Username  string  `json:"username"`
	Avatar    *string `json:"avatar"`
	Body      string  `json:"body"`
}

// CocDTO 子评论响应体
type CocDTO struct {
	CocID      uint    `json:"coc_id"`
	CommentID  uint    `json:"comment_id"`
	UserID     uint    `json:"user_id"`
	Username   string  `json:"username"`
	Avatar     *string `json:"avatar"`
	Body       string  `json:"body"`
	ToCocID    *uint   `json:"to_coc_id"`
	ToUserID   *uint   `json:"to_user_id"`
	ToUsername *string `json:"to_username"`
}

// FullCommentDTO 帖子详情中嵌入的完整评论
type FullCommentDTO struct {
	CommentDTO
	Likes       int64    `json:"likes"`
	SubComments []CocDTO `json:"sub_comments"`
}

// ==================== Post Service ====================

// PostToDTO 将 Post 模型转换为 DTO
func PostToDTO(post models.Post) PostDTO {
	dto := PostDTO{
		ID:        post.PostID,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
		UserID:    post.UserID,
		Title:     post.Title,
		Body:      post.Body,
		Tag:       post.Tag,
	}

	user, err := dao.GetUserByID(post.UserID)
	if err == nil {
		dto.Username = user.Username
		dto.Avatar = dao.GetUserAvatar(user.ID)
	}

	dto.ImageURL = dao.GetFirstImageByPostID(post.PostID)
	dto.Likes = dao.GetPostLikeCount(post.PostID)
	dto.Comments = dao.GetPostCommentCount(post.PostID)

	return dto
}

// CommentToDTO 将 Comment 模型转换为 DTO
func CommentToDTO(comment models.Comment) CommentDTO {
	dto := CommentDTO{
		ID:        comment.CommentID,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt: comment.UpdatedAt.Format(time.RFC3339),
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Body:      comment.Body,
	}

	user, err := dao.GetUserByID(comment.UserID)
	if err == nil {
		dto.Username = user.Username
		dto.Avatar = dao.GetUserAvatar(user.ID)
	}

	return dto
}

// GetPostDetail 获取帖子详情（含完整评论数据，带缓存）
func GetPostDetail(postID int) (map[string]interface{}, error) {
	// 尝试从缓存获取
	cacheKey := fmt.Sprintf(cache.KeyPostDetail, postID)
	var cached map[string]interface{}
	if err := cache.Get(cacheKey, &cached); err == nil {
		return cached, nil
	}

	post, err := dao.GetPostByID(postID)
	if err != nil {
		return nil, err
	}

	dto := PostToDTO(*post)

	// 获取所有图片
	images := dao.GetImagesByPostID(post.PostID)
	imageURLs := make([]string, 0, len(images))
	for _, img := range images {
		imageURLs = append(imageURLs, img.ImageURL)
	}

	// 获取评论
	comments := dao.GetCommentsByPostID(post.PostID)

	// 批量收集用户 ID
	userIDSet := map[uint]bool{}
	for _, c := range comments {
		userIDSet[c.UserID] = true
	}

	// 批量获取子评论
	commentIDs := make([]uint, 0, len(comments))
	for _, c := range comments {
		commentIDs = append(commentIDs, c.CommentID)
	}

	allCocs := dao.GetSubCommentsByCommentIDs(commentIDs)
	for _, coc := range allCocs {
		userIDSet[coc.UserID] = true
	}

	// 批量查询用户
	userIDs := make([]uint, 0, len(userIDSet))
	for uid := range userIDSet {
		userIDs = append(userIDs, uid)
	}
	userMap := dao.GetUsersByIDs(userIDs)

	// 批量查询评论点赞数
	commentLikesMap := dao.BatchGetCommentLikeCounts(commentIDs)

	// 批量查询头像
	avatarMap := dao.GetUserAvatarsBatch(userIDs)

	// 按 comment_id 分组子评论
	cocsByComment := map[uint][]models.CommentOnComments{}
	for _, coc := range allCocs {
		cocsByComment[coc.CommentID] = append(cocsByComment[coc.CommentID], coc)
	}

	getAvatar := func(uid uint) *string {
		if url, ok := avatarMap[uid]; ok {
			return &url
		}
		return nil
	}

	// 组装完整评论
	fullComments := make([]FullCommentDTO, 0, len(comments))
	for _, comment := range comments {
		cdto := CommentDTO{
			ID:        comment.CommentID,
			CreatedAt: comment.CreatedAt.Format(time.RFC3339),
			UpdatedAt: comment.UpdatedAt.Format(time.RFC3339),
			PostID:    comment.PostID,
			UserID:    comment.UserID,
			Body:      comment.Body,
			Avatar:    getAvatar(comment.UserID),
		}
		if u, ok := userMap[comment.UserID]; ok {
			cdto.Username = u.Username
		}

		subDtos := make([]CocDTO, 0)
		if cocs, ok := cocsByComment[comment.CommentID]; ok {
			for _, coc := range cocs {
				sd := CocDTO{
					CocID:     coc.CocID,
					CommentID: coc.CommentID,
					UserID:    coc.UserID,
					Body:      coc.Body,
					ToCocID:   coc.ToCocID,
					Avatar:    getAvatar(coc.UserID),
				}
				if u, ok := userMap[coc.UserID]; ok {
					sd.Username = u.Username
				}
				if coc.ToCocID != nil {
					for _, tc := range allCocs {
						if tc.CocID == *coc.ToCocID {
							sd.ToUserID = &tc.UserID
							if tu, ok := userMap[tc.UserID]; ok {
								sd.ToUsername = &tu.Username
							}
							break
						}
					}
				}
				subDtos = append(subDtos, sd)
			}
		}

		fullComments = append(fullComments, FullCommentDTO{
			CommentDTO:  cdto,
			Likes:       commentLikesMap[comment.CommentID],
			SubComments: subDtos,
		})
	}

	result := map[string]interface{}{
		"id":            dto.ID,
		"created_at":    dto.CreatedAt,
		"updated_at":    dto.UpdatedAt,
		"user_id":       dto.UserID,
		"username":      dto.Username,
		"avatar":        dto.Avatar,
		"title":         dto.Title,
		"body":          dto.Body,
		"tag":           dto.Tag,
		"image_urls":    imageURLs,
		"likes":         dto.Likes,
		"comments":      dto.Comments,
		"comments_list": fullComments,
	}

	// 写入缓存
	cache.Set(cacheKey, result, cache.PostDetailTTL)

	return result, nil
}

// GetSubComments 获取评论的子评论列表
func GetSubComments(commentID int) []CocDTO {
	cocs := dao.GetSubCommentsByCommentID(commentID)

	result := make([]CocDTO, 0, len(cocs))
	for _, coc := range cocs {
		dto := CocDTO{
			CocID:     coc.CocID,
			CommentID: coc.CommentID,
			UserID:    coc.UserID,
			Body:      coc.Body,
			ToCocID:   coc.ToCocID,
		}

		user, err := dao.GetUserByID(coc.UserID)
		if err == nil {
			dto.Username = user.Username
			dto.Avatar = dao.GetUserAvatar(user.ID)
		}

		if coc.ToCocID != nil {
			toCoc, err := dao.GetSubCommentByID(*coc.ToCocID)
			if err == nil {
				dto.ToUserID = &toCoc.UserID
				toUser, err := dao.GetUserByID(toCoc.UserID)
				if err == nil {
					dto.ToUsername = &toUser.Username
				}
			}
		}

		result = append(result, dto)
	}

	return result
}

// GetPostLikes 获取帖子点赞信息
func GetPostLikes(postID int, userID interface{}, authenticated bool) map[string]interface{} {
	count := dao.GetPostLikeCount(uint(postID))
	liked := false
	if authenticated {
		if existing := dao.FindPostLike(uint(postID), userID.(uint)); existing != nil {
			liked = true
		}
	}
	return map[string]interface{}{"likes": count, "liked": liked}
}

// ==================== 缓存失效 ====================

// InvalidatePostCache 清除帖子相关缓存
func InvalidatePostCache(postID uint) {
	cache.Del(fmt.Sprintf(cache.KeyPostDetail, postID))
	cache.Del(fmt.Sprintf(cache.KeyPostLikes, postID))
	// 帖子列表缓存用模式清除
	cache.DelByPattern("posts:list:*")
}

// InvalidatePostLikesCache 清除帖子点赞缓存
func InvalidatePostLikesCache(postID uint) {
	cache.Del(fmt.Sprintf(cache.KeyPostLikes, postID))
	cache.Del(fmt.Sprintf(cache.KeyPostDetail, postID))
}

// InvalidateSummaryCache 清除系统概览缓存
func InvalidateSummaryCache() {
	cache.Del(cache.KeySummary)
}

// InvalidateRegionEnvCache 清除区域环境缓存
func InvalidateRegionEnvCache() {
	cache.Del(cache.KeyRegionEnv)
	cache.DelByPattern("region:env:hist:*")
}
