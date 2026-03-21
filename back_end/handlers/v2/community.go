package v2

import (
	"net/http"
	"strconv"
	"time"

	"smart-fish/back_end/database"
	"smart-fish/back_end/models"
	"smart-fish/back_end/utils"

	"github.com/gin-gonic/gin"
)

// ==================== DTO（用于兼容 BaseModel 的 id 字段输出） ====================

// PostDTO 帖子响应体（id 代替 post_id）
// 注意：底层 Post 模型没有 created_at/updated_at，DTO 中仍保留以兼容 BaseModel 要求
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

func postToDTO(post models.Post) PostDTO {
	dto := PostDTO{
		ID:        post.PostID,
		CreatedAt: post.CreatedAt.Format(time.RFC3339),
		UpdatedAt: post.UpdatedAt.Format(time.RFC3339),
		UserID:    post.UserID,
		Title:     post.Title,
		Body:      post.Body,
		Tag:       post.Tag,
	}

	// 查用户名 + 头像
	var user models.User
	if err := database.DB.First(&user, post.UserID).Error; err == nil {
		dto.Username = user.Username
		dto.Avatar = getUserAvatarURL(user.ID)
	}

	// 查第一张图片
	var image models.Image
	if err := database.DB.Where("post_id = ? AND is_deleted = ?", post.PostID, false).First(&image).Error; err == nil {
		dto.ImageURL = &image.ImageURL
	}

	// 查点赞数
	database.DB.Model(&models.LikeOnPosts{}).Where("post_id = ?", post.PostID).Count(&dto.Likes)

	// 查评论数
	database.DB.Model(&models.Comment{}).Where("post_id = ? AND is_deleted = ?", post.PostID, false).Count(&dto.Comments)

	return dto
}

func commentToDTO(comment models.Comment) CommentDTO {
	dto := CommentDTO{
		ID:        comment.CommentID,
		CreatedAt: comment.CreatedAt.Format(time.RFC3339),
		UpdatedAt: comment.UpdatedAt.Format(time.RFC3339),
		PostID:    comment.PostID,
		UserID:    comment.UserID,
		Body:      comment.Body,
	}

	var user models.User
	if err := database.DB.First(&user, comment.UserID).Error; err == nil {
		dto.Username = user.Username
		dto.Avatar = getUserAvatarURL(user.ID)
	}

	return dto
}

// getUserAvatarURL 从 Image 表查询用户头像（供 community 内部使用）
func getUserAvatarURL(userID uint) *string {
	var avatar models.Image
	if err := database.DB.Where("user_id = ? AND is_avatar = ? AND is_deleted = ?", userID, true, false).First(&avatar).Error; err == nil {
		return &avatar.ImageURL
	}
	return nil
}

// ==================== Posts Handlers ====================

// ListPosts GET /api/posts - 帖子列表（分页）
func ListPosts(c *gin.Context) {
	query := database.DB.Model(&models.Post{}).Where("is_deleted = ?", false)

	if tag := c.Query("tag"); tag != "" {
		query = query.Where("tag = ?", tag)
	}
	if search := c.Query("search"); search != "" {
		query = query.Where("title LIKE ? OR body LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if userID := c.Query("user_id"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	utils.PaginateMap[models.Post, PostDTO](c, query, "post_id DESC", postToDTO)
}

// FullCommentDTO 帖子详情中嵌入的完整评论（含子评论 + 点赞数）
type FullCommentDTO struct {
	CommentDTO
	Likes       int64    `json:"likes"`
	SubComments []CocDTO `json:"sub_comments"`
}

// GetPost GET /api/posts/:id - 帖子详情（含完整评论数据）
func GetPostByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var post models.Post
	if err := database.DB.Where("post_id = ? AND is_deleted = ?", id, false).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}

	dto := postToDTO(post)

	// 详情页加载所有图片
	var images []models.Image
	database.DB.Where("post_id = ? AND is_deleted = ?", post.PostID, false).Find(&images)
	imageURLs := make([]string, 0, len(images))
	for _, img := range images {
		imageURLs = append(imageURLs, img.ImageURL)
	}

	// ===== 嵌入完整评论数据（含子评论 + 点赞数），消灭前端 N+1 请求 =====

	// 1. 查询所有评论
	var comments []models.Comment
	database.DB.Where("post_id = ? AND is_deleted = ?", post.PostID, false).Order("comment_id ASC").Find(&comments)

	// 2. 批量收集所有需要查询的用户 ID
	userIDSet := map[uint]bool{}
	for _, c := range comments {
		userIDSet[c.UserID] = true
	}

	// 3. 查询所有评论的子评论
	commentIDs := make([]uint, 0, len(comments))
	for _, c := range comments {
		commentIDs = append(commentIDs, c.CommentID)
	}

	var allCocs []models.CommentOnComments
	if len(commentIDs) > 0 {
		database.DB.Where("comment_id IN ? AND is_deleted = ?", commentIDs, false).Find(&allCocs)
	}

	// 收集子评论中的用户 ID
	for _, coc := range allCocs {
		userIDSet[coc.UserID] = true
	}

	// 4. 批量查询所有相关用户（一次 DB 查询）
	userIDs := make([]uint, 0, len(userIDSet))
	for uid := range userIDSet {
		userIDs = append(userIDs, uid)
	}
	var users []models.User
	userMap := map[uint]models.User{}
	if len(userIDs) > 0 {
		database.DB.Where("id IN ?", userIDs).Find(&users)
		for _, u := range users {
			userMap[u.ID] = u
		}
	}

	// 5. 批量查询评论点赞数（一次聚合查询）
	type likeCount struct {
		CommentID uint
		Cnt       int64
	}
	var commentLikes []likeCount
	commentLikesMap := map[uint]int64{}
	if len(commentIDs) > 0 {
		database.DB.Model(&models.LikeOnComments{}).
			Select("comment_id, count(*) as cnt").
			Where("comment_id IN ?", commentIDs).
			Group("comment_id").
			Scan(&commentLikes)
		for _, cl := range commentLikes {
			commentLikesMap[cl.CommentID] = cl.Cnt
		}
	}

	// 6. 批量查询用户头像（一次查询）
	avatarMap := map[uint]string{}
	if len(userIDs) > 0 {
		var avatars []models.Image
		database.DB.Where("user_id IN ? AND is_avatar = ? AND is_deleted = ?", userIDs, true, false).Find(&avatars)
		for _, a := range avatars {
			avatarMap[a.UserID] = a.ImageURL
		}
	}

	// 7. 按 comment_id 分组子评论
	cocsByComment := map[uint][]models.CommentOnComments{}
	for _, coc := range allCocs {
		cocsByComment[coc.CommentID] = append(cocsByComment[coc.CommentID], coc)
	}

	// 辅助：根据 userID 获取头像指针
	getAvatar := func(uid uint) *string {
		if url, ok := avatarMap[uid]; ok {
			return &url
		}
		return nil
	}

	// 8. 组装完整评论 DTO
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

		// 子评论
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
					// 找到被回复的子评论的用户
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

	c.JSON(http.StatusOK, gin.H{
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
	})
}

// CreatePostV2 POST /api/posts - 创建帖子
func CreatePostV2(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var input struct {
		Title string `json:"title" binding:"required"`
		Body  string `json:"body" binding:"required"`
		Tag   string `json:"tag"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	post := models.Post{
		UserID: userID.(uint),
		Title:  input.Title,
		Body:   input.Body,
		Tag:    input.Tag,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, postToDTO(post))
}

// UpdatePostV2 PUT /api/posts/:id - 更新帖子
func UpdatePostV2(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var post models.Post
	if err := database.DB.Where("post_id = ? AND is_deleted = ?", id, false).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}

	// 权限校验：只有帖子作者才能修改
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改他人帖子"})
		return
	}

	var input struct {
		Title string `json:"title"`
		Body  string `json:"body"`
		Tag   string `json:"tag"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败"})
		return
	}

	updates := map[string]interface{}{}
	if input.Title != "" {
		updates["title"] = input.Title
	}
	if input.Body != "" {
		updates["body"] = input.Body
	}
	if input.Tag != "" {
		updates["tag"] = input.Tag
	}

	if len(updates) > 0 {
		database.DB.Model(&post).Updates(updates)
	}

	database.DB.Where("post_id = ?", id).First(&post)
	c.JSON(http.StatusOK, postToDTO(post))
}

// DeletePostV2 DELETE /api/posts/:id - 软删除帖子
func DeletePostV2(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var post models.Post
	if err := database.DB.Where("post_id = ? AND is_deleted = ?", id, false).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}

	// 权限校验：只有帖子作者才能删除
	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除他人帖子"})
		return
	}

	database.DB.Model(&post).Update("is_deleted", true)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ==================== Comments Handlers ====================

// ListComments GET /api/comments - 评论列表（支持分页）
func ListComments(c *gin.Context) {
	query := database.DB.Model(&models.Comment{}).Where("is_deleted = ?", false)

	if postID := c.Query("post_id"); postID != "" {
		query = query.Where("post_id = ?", postID)
	}

	utils.PaginateMap[models.Comment, CommentDTO](c, query, "comment_id ASC", commentToDTO)
}

// GetComment GET /api/comments/:id
func GetComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var comment models.Comment
	if err := database.DB.Where("comment_id = ? AND is_deleted = ?", id, false).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	c.JSON(http.StatusOK, commentToDTO(comment))
}

// CreateCommentV2 POST /api/comments
func CreateCommentV2(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var input struct {
		PostID uint   `json:"post_id" binding:"required"`
		Body   string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	// 验证帖子存在
	var post models.Post
	if err := database.DB.Where("post_id = ? AND is_deleted = ?", input.PostID, false).First(&post).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "帖子不存在"})
		return
	}

	comment := models.Comment{
		PostID: input.PostID,
		UserID: userID.(uint),
		Body:   input.Body,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, commentToDTO(comment))
}

// DeleteCommentV2 DELETE /api/comments/:id - 软删除评论
func DeleteCommentV2(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var comment models.Comment
	if err := database.DB.Where("comment_id = ? AND is_deleted = ?", id, false).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 权限校验：只有评论作者才能删除
	if comment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除他人评论"})
		return
	}

	database.DB.Model(&comment).Update("is_deleted", true)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ==================== Likes Handlers ====================

// LikePost POST /api/v2/posts/:id/like - 点赞帖子
func LikePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var post models.Post
	if err := database.DB.Where("post_id = ? AND is_deleted = ?", id, false).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}

	var existing models.LikeOnPosts
	if err := database.DB.Where("post_id = ? AND user_id = ?", post.PostID, userID.(uint)).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已点赞"})
		return
	}

	like := models.LikeOnPosts{PostID: post.PostID, UserID: userID.(uint)}
	database.DB.Create(&like)
	c.JSON(http.StatusCreated, gin.H{"message": "点赞成功"})
}

// UnlikePost DELETE /api/v2/posts/:id/like - 取消帖子点赞
func UnlikePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	result := database.DB.Where("post_id = ? AND user_id = ?", id, userID.(uint)).Delete(&models.LikeOnPosts{})
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "取消点赞成功"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到点赞记录"})
	}
}

// GetPostLikesV2 GET /api/v2/posts/:id/like - 获取帖子点赞数（支持可选认证）
func GetPostLikesV2(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var count int64
	database.DB.Model(&models.LikeOnPosts{}).Where("post_id = ?", id).Count(&count)

	liked := false
	if userID, exists := c.Get("userID"); exists {
		var existing models.LikeOnPosts
		if err := database.DB.Where("post_id = ? AND user_id = ?", id, userID.(uint)).First(&existing).Error; err == nil {
			liked = true
		}
	}

	c.JSON(http.StatusOK, gin.H{"likes": count, "liked": liked})
}

// LikeComment POST /api/v2/comments/:id/like - 点赞评论
func LikeComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var comment models.Comment
	if err := database.DB.Where("comment_id = ? AND is_deleted = ?", id, false).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	var existing models.LikeOnComments
	if err := database.DB.Where("comment_id = ? AND user_id = ?", comment.CommentID, userID.(uint)).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已点赞"})
		return
	}

	like := models.LikeOnComments{CommentID: comment.CommentID, UserID: userID.(uint)}
	database.DB.Create(&like)
	c.JSON(http.StatusCreated, gin.H{"message": "点赞成功"})
}

// UnlikeComment DELETE /api/v2/comments/:id/like - 取消评论点赞
func UnlikeComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	result := database.DB.Where("comment_id = ? AND user_id = ?", id, userID.(uint)).Delete(&models.LikeOnComments{})
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "取消点赞成功"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到点赞记录"})
	}
}

// GetCommentLikesV2 GET /api/v2/comments/:id/like - 获取评论点赞数
func GetCommentLikesV2(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var count int64
	database.DB.Model(&models.LikeOnComments{}).Where("comment_id = ?", id).Count(&count)
	c.JSON(http.StatusOK, gin.H{"likes": count})
}

// ==================== Sub-Comments (楼中楼) Handlers ====================

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

// ListSubComments GET /api/v2/comments/:id/replies - 获取评论的子评论
func ListSubComments(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	var cocs []models.CommentOnComments
	database.DB.Where("comment_id = ? AND is_deleted = ?", commentID, false).Find(&cocs)

	result := make([]CocDTO, 0, len(cocs))
	for _, coc := range cocs {
		dto := CocDTO{
			CocID:     coc.CocID,
			CommentID: coc.CommentID,
			UserID:    coc.UserID,
			Body:      coc.Body,
			ToCocID:   coc.ToCocID,
		}

		var user models.User
		if err := database.DB.First(&user, coc.UserID).Error; err == nil {
			dto.Username = user.Username
			dto.Avatar = getUserAvatarURL(user.ID)
		}

		if coc.ToCocID != nil {
			var toCoc models.CommentOnComments
			if err := database.DB.Where("coc_id = ?", *coc.ToCocID).First(&toCoc).Error; err == nil {
				dto.ToUserID = &toCoc.UserID
				var toUser models.User
				if err := database.DB.First(&toUser, toCoc.UserID).Error; err == nil {
					dto.ToUsername = &toUser.Username
				}
			}
		}

		result = append(result, dto)
	}

	c.JSON(http.StatusOK, gin.H{
		"comment_id": commentID,
		"comments":   result,
	})
}

// CreateSubComment POST /api/v2/comments/:id/replies - 创建子评论
func CreateSubComment(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证"})
		return
	}

	var input struct {
		Body    string `json:"body" binding:"required"`
		ToCocID *uint  `json:"to_coc_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数验证失败: " + err.Error()})
		return
	}

	if input.ToCocID != nil {
		var targetCoc models.CommentOnComments
		if err := database.DB.Where("coc_id = ?", *input.ToCocID).First(&targetCoc).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "目标子评论不存在"})
			return
		}
		if targetCoc.CommentID != uint(commentID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "目标子评论不属于该评论"})
			return
		}
	}

	coc := models.CommentOnComments{
		CommentID: uint(commentID),
		UserID:    userID.(uint),
		Body:      input.Body,
		ToCocID:   input.ToCocID,
	}

	if err := database.DB.Create(&coc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "回复成功",
		"coc_id":  coc.CocID,
	})
}
