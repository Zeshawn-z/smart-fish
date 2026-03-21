package v1

import (
	"net/http"

	"smart-fish/back_end/database"
	"smart-fish/back_end/middleware"
	"smart-fish/back_end/models"

	"github.com/gin-gonic/gin"
)

// GetCommentList GET /api/v1/post/:post_id/comment - 获取帖子的评论列表
func GetCommentList(c *gin.Context) {
	postID := c.Param("post_id")

	var comments []models.Comment
	database.DB.Where("post_id = ? AND is_deleted = ?", postID, false).Find(&comments)

	commentsList := make([]gin.H, 0, len(comments))
	for _, comment := range comments {
		// 查询用户名（Go User 表主键是 id）
		var user models.User
		username := ""
		if err := database.DB.First(&user, comment.UserID).Error; err == nil {
			username = user.Username
		}

		commentsList = append(commentsList, gin.H{
			"comment_id": comment.CommentID,
			"post_id":    comment.CommentID, // Flask 原始代码的 bug: 返回的是 comment.comment_id 而非 comment.post_id
			"user_id":    comment.UserID,
			"username":   username,
			"body":       comment.Body,
		})
	}

	c.JSON(http.StatusOK, gin.H{"comments_list": commentsList})
}

// CreateComment POST /api/v1/post/:post_id/comment - 创建评论
func CreateComment(c *gin.Context) {
	postID := c.Param("post_id")
	userID, ok := middleware.GetFlaskUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
		return
	}

	// 验证帖子存在
	var post models.Post
	if err := database.DB.Where("post_id = ?", postID).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	var input struct {
		Body string `json:"body"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request data"})
		return
	}

	comment := models.Comment{
		PostID: post.PostID,
		UserID: userID,
		Body:   input.Body,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":        "New Comment Created Successfully",
		"comment_id": comment.CommentID,
	})
}

// GetCommentOnComments GET /api/v1/comment/:comment_id - 获取评论的子评论列表
func GetCommentOnComments(c *gin.Context) {
	commentIDStr := c.Param("comment_id")

	// 解析 comment_id 为 uint，与 Flask 返回的 int 类型一致
	commentID, err := parseUintParam(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid comment_id"})
		return
	}

	var cocs []models.CommentOnComments
	database.DB.Where("comment_id = ? AND is_deleted = ?", commentID, false).Find(&cocs)

	result := make([]gin.H, 0, len(cocs))
	for _, coc := range cocs {
		// 查询评论者用户名（Go User 表主键是 id）
		var user models.User
		username := ""
		if err := database.DB.First(&user, coc.UserID).Error; err == nil {
			username = user.Username
		}

		item := gin.H{
			"coc_id":      coc.CocID,
			"to_coc_id":   nil,
			"comment_id":  coc.CommentID,
			"user_id":     coc.UserID,
			"to_user_id":  nil,
			"username":    username,
			"to_username": nil,
			"body":        coc.Body,
		}

		// 处理回复评论
		if coc.ToCocID != nil {
			item["to_coc_id"] = *coc.ToCocID

			var toCoc models.CommentOnComments
			if err := database.DB.Where("coc_id = ?", *coc.ToCocID).First(&toCoc).Error; err == nil {
				item["to_user_id"] = toCoc.UserID

				var toUser models.User
				if err := database.DB.First(&toUser, toCoc.UserID).Error; err == nil {
					item["to_username"] = toUser.Username
				}
			}
		}

		result = append(result, item)
	}

	// Flask 返回 comment_id 为 int 类型
	c.JSON(http.StatusOK, gin.H{
		"comment_id": commentID,
		"comments":   result,
	})
}

// CreateCommentOnComments POST /api/v1/comment/:comment_id - 在评论下创建子评论
func CreateCommentOnComments(c *gin.Context) {
	commentID := c.Param("comment_id")
	userID, ok := middleware.GetFlaskUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
		return
	}

	var input struct {
		Body string `json:"body"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request data"})
		return
	}

	// 解析 comment_id 字符串为 uint
	var cid uint
	if _, err := parseUintParam(commentID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid comment_id"})
		return
	}
	cid, _ = parseUintParam(commentID)

	coc := models.CommentOnComments{
		CommentID: cid,
		UserID:    userID,
		Body:      input.Body,
	}

	if err := database.DB.Create(&coc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create COC"})
		return
	}

	// Flask 原始代码返回 comment_id 而非 coc_id（这是一个 bug，但我们保持兼容）
	c.JSON(http.StatusOK, gin.H{
		"msg":                    "New COC Created Successfully",
		"comment_on_comments_id": coc.CommentID,
	})
}

// CreateCommentOnCocs POST /api/v1/coc/:coc_id - 回复子评论
func CreateCommentOnCocs(c *gin.Context) {
	cocIDStr := c.Param("coc_id")
	userID, ok := middleware.GetFlaskUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
		return
	}

	var input struct {
		Body string `json:"body"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request data"})
		return
	}

	cocID, err := parseUintParam(cocIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid coc_id"})
		return
	}

	// 查找原始 COC 获取其 comment_id
	var originalCoc models.CommentOnComments
	if err := database.DB.Where("coc_id = ?", cocID).First(&originalCoc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "COC not found"})
		return
	}

	newCoc := models.CommentOnComments{
		CommentID: originalCoc.CommentID,
		UserID:    userID,
		Body:      input.Body,
		ToCocID:   &cocID,
	}

	if err := database.DB.Create(&newCoc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create COC"})
		return
	}

	// Flask 原始代码返回 comment_id 而非 coc_id（保持兼容）
	c.JSON(http.StatusOK, gin.H{
		"msg":                    "New COC Created Successfully",
		"comment_on_comments_id": newCoc.CommentID,
	})
}

// parseUintParam 从路径参数字符串解析 uint
func parseUintParam(s string) (uint, error) {
	var n uint
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0, gin.Error{Err: nil}
		}
		n = n*10 + uint(ch-'0')
	}
	return n, nil
}
