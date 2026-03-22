package v1

import (
	"net/http"

	"smart-fish/back_end/dao"
	"smart-fish/back_end/middleware"
	"smart-fish/back_end/models"

	"github.com/gin-gonic/gin"
)

// GetCommentList GET /api/v1/post/:post_id/comment - 获取帖子的评论列表
func GetCommentList(c *gin.Context) {
	postID := c.Param("post_id")

	comments := dao.GetCommentsByPostIDV1(postID)

	commentsList := make([]gin.H, 0, len(comments))
	for _, comment := range comments {
		// 查询用户名
		username := ""
		if user, err := dao.GetUserByID(comment.UserID); err == nil {
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
	post, err := dao.GetPostByPostID(parseUintSafe(postID))
	if err != nil {
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

	if err := dao.CreateComment(&comment); err != nil {
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

	commentID, err := parseUintParam(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid comment_id"})
		return
	}

	cocs := dao.GetSubCommentsByCommentIDV1(commentID)

	result := make([]gin.H, 0, len(cocs))
	for _, coc := range cocs {
		// 查询评论者用户名
		username := ""
		if user, err := dao.GetUserByID(coc.UserID); err == nil {
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

			if toCoc, err := dao.GetCocByCocID(*coc.ToCocID); err == nil {
				item["to_user_id"] = toCoc.UserID

				if toUser, err := dao.GetUserByID(toCoc.UserID); err == nil {
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

	cid, err := parseUintParam(commentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid comment_id"})
		return
	}

	coc := models.CommentOnComments{
		CommentID: cid,
		UserID:    userID,
		Body:      input.Body,
	}

	if err := dao.CreateSubComment(&coc); err != nil {
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
	originalCoc, err := dao.GetCocByCocID(cocID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "COC not found"})
		return
	}

	newCoc := models.CommentOnComments{
		CommentID: originalCoc.CommentID,
		UserID:    userID,
		Body:      input.Body,
		ToCocID:   &cocID,
	}

	if err := dao.CreateSubComment(&newCoc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create COC"})
		return
	}

	// Flask 原始代码返回 comment_id 而非 coc_id（保持兼容）
	c.JSON(http.StatusOK, gin.H{
		"msg":                    "New COC Created Successfully",
		"comment_on_comments_id": newCoc.CommentID,
	})
}

// parseUintSafe 安全地将 string 转为 uint，出错返回 0
func parseUintSafe(s string) uint {
	n, _ := parseUintParam(s)
	return n
}
