package v2

import (
	"net/http"
	"strconv"

	"smart-fish/back_end/dao"
	"smart-fish/back_end/models"
	"smart-fish/back_end/services"
	"smart-fish/back_end/utils"

	"github.com/gin-gonic/gin"
)

// ==================== Posts Handlers ====================

// ListPosts GET /api/posts - 帖子列表（分页）
func ListPosts(c *gin.Context) {
	query := dao.ListPostsQuery(c.Query("tag"), c.Query("search"), c.Query("user_id"))
	utils.PaginateMapConcurrent[models.Post, services.PostDTO](c, query, "post_id DESC", services.PostToDTO)
}

// GetPostByID GET /api/posts/:id - 帖子详情（含完整评论数据）
func GetPostByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	result, err := services.GetPostDetail(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}

	c.JSON(http.StatusOK, result)
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

	if err := dao.CreatePost(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	go services.InvalidatePostCache(post.PostID)
	c.JSON(http.StatusCreated, services.PostToDTO(post))
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

	post, err := dao.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}

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
		dao.UpdatePost(post, updates)
	}

	dao.RefreshPost(post, id)
	go services.InvalidatePostCache(post.PostID)
	c.JSON(http.StatusOK, services.PostToDTO(*post))
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

	post, err := dao.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}

	if post.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除他人帖子"})
		return
	}

	dao.SoftDeletePost(post)
	go services.InvalidatePostCache(post.PostID)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ==================== Comments Handlers ====================

// ListComments GET /api/comments - 评论列表（支持分页）
func ListComments(c *gin.Context) {
	query := dao.ListCommentsQuery(c.Query("post_id"))
	utils.PaginateMapConcurrent[models.Comment, services.CommentDTO](c, query, "comment_id ASC", services.CommentToDTO)
}

// GetComment GET /api/comments/:id
func GetComment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	comment, err := dao.GetCommentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	c.JSON(http.StatusOK, services.CommentToDTO(*comment))
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

	if _, err := dao.GetPostByID(int(input.PostID)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "帖子不存在"})
		return
	}

	comment := models.Comment{
		PostID: input.PostID,
		UserID: userID.(uint),
		Body:   input.Body,
	}

	if err := dao.CreateComment(&comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, services.CommentToDTO(comment))
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

	comment, err := dao.GetCommentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	if comment.UserID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除他人评论"})
		return
	}

	dao.SoftDeleteComment(comment)
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

	post, err := dao.GetPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
		return
	}

	if existing := dao.FindPostLike(post.PostID, userID.(uint)); existing != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已点赞"})
		return
	}

	dao.CreatePostLike(&models.LikeOnPosts{PostID: post.PostID, UserID: userID.(uint)})
	go services.InvalidatePostLikesCache(post.PostID)
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

	affected := dao.DeletePostLike(uint(id), userID.(uint))
	if affected > 0 {
		go services.InvalidatePostLikesCache(uint(id))
		c.JSON(http.StatusOK, gin.H{"message": "取消点赞成功"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到点赞记录"})
	}
}

// GetPostLikesV2 GET /api/v2/posts/:id/like - 获取帖子点赞数
func GetPostLikesV2(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	userID, exists := c.Get("userID")
	result := services.GetPostLikes(id, userID, exists)
	c.JSON(http.StatusOK, result)
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

	comment, err := dao.GetCommentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	if existing := dao.FindCommentLike(comment.CommentID, userID.(uint)); existing != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已点赞"})
		return
	}

	dao.CreateCommentLike(&models.LikeOnComments{CommentID: comment.CommentID, UserID: userID.(uint)})
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

	affected := dao.DeleteCommentLike(uint(id), userID.(uint))
	if affected > 0 {
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

	count := dao.GetCommentLikeCount(uint(id))
	c.JSON(http.StatusOK, gin.H{"likes": count})
}

// ==================== Sub-Comments (楼中楼) Handlers ====================

// ListSubComments GET /api/v2/comments/:id/replies - 获取评论的子评论
func ListSubComments(c *gin.Context) {
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 ID"})
		return
	}

	result := services.GetSubComments(commentID)
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
		targetCoc, err := dao.GetSubCommentByID(*input.ToCocID)
		if err != nil {
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

	if err := dao.CreateSubComment(&coc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建失败"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "回复成功",
		"coc_id":  coc.CocID,
	})
}
