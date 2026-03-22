package services

import (
	"fmt"
	"sync"
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

// PostToDTO 将 Post 模型转换为 DTO（并发查询关联数据）
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

	// 5 个查询完全独立，并发执行
	var (
		user     *models.User
		avatar   *string
		imageURL *string
		likes    int64
		comments int64
		wg       sync.WaitGroup
	)

	wg.Add(5)
	go func() {
		defer wg.Done()
		user, _ = dao.GetUserByID(post.UserID)
	}()
	go func() {
		defer wg.Done()
		avatar = dao.GetUserAvatar(post.UserID)
	}()
	go func() {
		defer wg.Done()
		imageURL = dao.GetFirstImageByPostID(post.PostID)
	}()
	go func() {
		defer wg.Done()
		likes = dao.GetPostLikeCount(post.PostID)
	}()
	go func() {
		defer wg.Done()
		comments = dao.GetPostCommentCount(post.PostID)
	}()
	wg.Wait()

	if user != nil {
		dto.Username = user.Username
	}
	dto.Avatar = avatar
	dto.ImageURL = imageURL
	dto.Likes = likes
	dto.Comments = comments

	return dto
}

// CommentToDTO 将 Comment 模型转换为 DTO（并发查询用户和头像）
func CommentToDTO(comment models.Comment) CommentDTO {
	dto := CommentDTO{
		ID:        comment.CommentID,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt: comment.UpdatedAt.Format(time.RFC3339),
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Body:      comment.Body,
	}

	// 用户信息和头像互相独立，并发查询
	var (
		user   *models.User
		avatar *string
		wg     sync.WaitGroup
	)

	wg.Add(2)
	go func() {
		defer wg.Done()
		user, _ = dao.GetUserByID(comment.UserID)
	}()
	go func() {
		defer wg.Done()
		avatar = dao.GetUserAvatar(comment.UserID)
	}()
	wg.Wait()

	if user != nil {
		dto.Username = user.Username
	}
	dto.Avatar = avatar

	return dto
}

// GetPostDetail 获取帖子详情（含完整评论数据，带缓存 + 并发查询）
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

	// === 第一组并行：全部仅依赖 post，互相独立 ===
	var (
		user     *models.User
		avatar   *string
		imageURL *string
		images   []models.Image
		comments []models.Comment
		likes    int64
		commentCount int64
		wg1      sync.WaitGroup
	)

	wg1.Add(6)
	go func() { defer wg1.Done(); user, _ = dao.GetUserByID(post.UserID) }()
	go func() { defer wg1.Done(); imageURL = dao.GetFirstImageByPostID(post.PostID) }()
	go func() { defer wg1.Done(); images = dao.GetImagesByPostID(post.PostID) }()
	go func() { defer wg1.Done(); comments = dao.GetCommentsByPostID(post.PostID) }()
	go func() { defer wg1.Done(); likes = dao.GetPostLikeCount(post.PostID) }()
	go func() { defer wg1.Done(); commentCount = dao.GetPostCommentCount(post.PostID) }()
	wg1.Wait()

	// 用户头像（依赖 user 结果）
	username := ""
	if user != nil {
		username = user.Username
		avatar = dao.GetUserAvatar(user.ID)
	}
	_ = imageURL // PostToDTO 中用到，这里用 images 完整列表

	// 整理图片 URL
	imageURLs := make([]string, 0, len(images))
	for _, img := range images {
		imageURLs = append(imageURLs, img.ImageURL)
	}

	// === 第二组并行：依赖 comments 结果 ===
	commentIDs := make([]uint, 0, len(comments))
	userIDSet := map[uint]bool{post.UserID: true}
	for _, c := range comments {
		commentIDs = append(commentIDs, c.CommentID)
		userIDSet[c.UserID] = true
	}

	var (
		allCocs         []models.CommentOnComments
		commentLikesMap map[uint]int64
		wg2             sync.WaitGroup
	)

	wg2.Add(2)
	go func() { defer wg2.Done(); allCocs = dao.GetSubCommentsByCommentIDs(commentIDs) }()
	go func() { defer wg2.Done(); commentLikesMap = dao.BatchGetCommentLikeCounts(commentIDs) }()
	wg2.Wait()

	// 收集子评论中的用户 ID
	for _, coc := range allCocs {
		userIDSet[coc.UserID] = true
	}

	// === 第三组并行：依赖 userIDSet ===
	userIDs := make([]uint, 0, len(userIDSet))
	for uid := range userIDSet {
		userIDs = append(userIDs, uid)
	}

	var (
		userMap   map[uint]models.User
		avatarMap map[uint]string
		wg3       sync.WaitGroup
	)

	wg3.Add(2)
	go func() { defer wg3.Done(); userMap = dao.GetUsersByIDs(userIDs) }()
	go func() { defer wg3.Done(); avatarMap = dao.GetUserAvatarsBatch(userIDs) }()
	wg3.Wait()

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
		"id":            post.PostID,
		"created_at":    post.CreatedAt.Format(time.RFC3339),
		"updated_at":    post.UpdatedAt.Format(time.RFC3339),
		"user_id":       post.UserID,
		"username":      username,
		"avatar":        avatar,
		"title":         post.Title,
		"body":          post.Body,
		"tag":           post.Tag,
		"image_urls":    imageURLs,
		"likes":         likes,
		"comments":      commentCount,
		"comments_list": fullComments,
	}

	// 写入缓存
	cache.Set(cacheKey, result, cache.PostDetailTTL)

	return result, nil
}

// GetSubComments 获取评论的子评论列表（批量优化）
func GetSubComments(commentID int) []CocDTO {
	cocs := dao.GetSubCommentsByCommentID(commentID)
	if len(cocs) == 0 {
		return []CocDTO{}
	}

	// 收集所有涉及的 userID 和 toCocID
	userIDSet := map[uint]bool{}
	toCocIDs := map[uint]bool{}
	for _, coc := range cocs {
		userIDSet[coc.UserID] = true
		if coc.ToCocID != nil {
			toCocIDs[*coc.ToCocID] = true
		}
	}

	// 构建 coc 查找 map（用于解析 ToCocID → UserID）
	cocMap := map[uint]models.CommentOnComments{}
	for _, coc := range cocs {
		cocMap[coc.CocID] = coc
	}
	// 不在当前列表中的 toCocID，需要额外查询
	for tocID := range toCocIDs {
		if _, ok := cocMap[tocID]; !ok {
			if toCoc, err := dao.GetSubCommentByID(tocID); err == nil {
				cocMap[tocID] = *toCoc
				userIDSet[toCoc.UserID] = true
			}
		} else {
			userIDSet[cocMap[tocID].UserID] = true
		}
	}

	// 批量查询用户和头像
	userIDs := make([]uint, 0, len(userIDSet))
	for uid := range userIDSet {
		userIDs = append(userIDs, uid)
	}

	var (
		userMap   map[uint]models.User
		avatarMap map[uint]string
		wg        sync.WaitGroup
	)
	wg.Add(2)
	go func() { defer wg.Done(); userMap = dao.GetUsersByIDs(userIDs) }()
	go func() { defer wg.Done(); avatarMap = dao.GetUserAvatarsBatch(userIDs) }()
	wg.Wait()

	getAvatar := func(uid uint) *string {
		if url, ok := avatarMap[uid]; ok {
			return &url
		}
		return nil
	}

	// 组装结果
	result := make([]CocDTO, 0, len(cocs))
	for _, coc := range cocs {
		dto := CocDTO{
			CocID:     coc.CocID,
			CommentID: coc.CommentID,
			UserID:    coc.UserID,
			Body:      coc.Body,
			ToCocID:   coc.ToCocID,
			Avatar:    getAvatar(coc.UserID),
		}

		if u, ok := userMap[coc.UserID]; ok {
			dto.Username = u.Username
		}

		if coc.ToCocID != nil {
			if tc, ok := cocMap[*coc.ToCocID]; ok {
				dto.ToUserID = &tc.UserID
				if tu, ok := userMap[tc.UserID]; ok {
					dto.ToUsername = &tu.Username
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
