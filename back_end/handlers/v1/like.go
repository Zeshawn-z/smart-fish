package v1

import (
	"net/http"

	"smart-fish/back_end/database"
	"smart-fish/back_end/middleware"
	"smart-fish/back_end/models"

	"github.com/gin-gonic/gin"
)

// ===== Post Likes =====

// GetPostLikes GET /api/v1/post/:post_id/like
func GetPostLikes(c *gin.Context) {
	postID := c.Param("post_id")

	// 验证帖子存在
	var post models.Post
	if err := database.DB.Where("post_id = ?", postID).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	var count int64
	database.DB.Model(&models.LikeOnPosts{}).Where("post_id = ?", postID).Count(&count)
	c.JSON(http.StatusOK, gin.H{"likes": count})
}

// CreatePostLike POST /api/v1/post/:post_id/like
func CreatePostLike(c *gin.Context) {
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

	// 检查是否已点赞
	var existing models.LikeOnPosts
	if err := database.DB.Where("post_id = ? AND user_id = ?", post.PostID, userID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Already Liked"})
		return
	}

	like := models.LikeOnPosts{
		PostID: post.PostID,
		UserID: userID,
	}
	database.DB.Create(&like)
	c.JSON(http.StatusCreated, gin.H{"msg": "success"})
}

// DeletePostLike DELETE /api/v1/post/:post_id/like
func DeletePostLike(c *gin.Context) {
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

	result := database.DB.Where("post_id = ? AND user_id = ?", post.PostID, userID).Delete(&models.LikeOnPosts{})
	if result.RowsAffected > 0 {
		c.JSON(http.StatusNoContent, gin.H{"msg": "success"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"msg": "not found"})
	}
}

// ===== Comment Likes =====

// GetCommentLikes GET /api/v1/comment/:comment_id/like
func GetCommentLikes(c *gin.Context) {
	commentID := c.Param("comment_id")

	var comment models.Comment
	if err := database.DB.Where("comment_id = ?", commentID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	var count int64
	database.DB.Model(&models.LikeOnComments{}).Where("comment_id = ?", commentID).Count(&count)
	c.JSON(http.StatusOK, gin.H{"likes": count})
}

// CreateCommentLike POST /api/v1/comment/:comment_id/like
func CreateCommentLike(c *gin.Context) {
	commentID := c.Param("comment_id")
	userID, ok := middleware.GetFlaskUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
		return
	}

	var comment models.Comment
	if err := database.DB.Where("comment_id = ?", commentID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	var existing models.LikeOnComments
	if err := database.DB.Where("comment_id = ? AND user_id = ?", comment.CommentID, userID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Already Liked"})
		return
	}

	like := models.LikeOnComments{
		CommentID: comment.CommentID,
		UserID:    userID,
	}
	database.DB.Create(&like)
	c.JSON(http.StatusCreated, gin.H{"msg": "success"})
}

// DeleteCommentLike DELETE /api/v1/comment/:comment_id/like
func DeleteCommentLike(c *gin.Context) {
	commentID := c.Param("comment_id")
	userID, ok := middleware.GetFlaskUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
		return
	}

	var comment models.Comment
	if err := database.DB.Where("comment_id = ?", commentID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	result := database.DB.Where("comment_id = ? AND user_id = ?", comment.CommentID, userID).Delete(&models.LikeOnComments{})
	if result.RowsAffected > 0 {
		c.JSON(http.StatusNoContent, gin.H{"msg": "success"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"msg": "not found"})
	}
}

// ===== COC Likes =====

// GetCocLikes GET /api/v1/comment_on_comments/:coc_id/like
func GetCocLikes(c *gin.Context) {
	cocID := c.Param("coc_id")

	var coc models.CommentOnComments
	if err := database.DB.Where("coc_id = ?", cocID).First(&coc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	var count int64
	database.DB.Model(&models.LikeOnCOCS{}).Where("coc_id = ?", cocID).Count(&count)
	c.JSON(http.StatusOK, gin.H{"likes": count})
}

// CreateCocLike POST /api/v1/comment_on_comments/:coc_id/like
func CreateCocLike(c *gin.Context) {
	cocID := c.Param("coc_id")
	userID, ok := middleware.GetFlaskUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
		return
	}

	var coc models.CommentOnComments
	if err := database.DB.Where("coc_id = ?", cocID).First(&coc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	var existing models.LikeOnCOCS
	if err := database.DB.Where("coc_id = ? AND user_id = ?", coc.CocID, userID).First(&existing).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Already Liked"})
		return
	}

	like := models.LikeOnCOCS{
		CocID:  coc.CocID,
		UserID: userID,
	}
	database.DB.Create(&like)
	c.JSON(http.StatusCreated, gin.H{"msg": "success"})
}

// DeleteCocLike DELETE /api/v1/comment_on_comments/:coc_id/like
func DeleteCocLike(c *gin.Context) {
	cocID := c.Param("coc_id")
	userID, ok := middleware.GetFlaskUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
		return
	}

	// Flask 原始代码这里有 bug：用 Comment.query.get_or_404(coc_id) 查错了模型
	// 我们用正确的模型 CommentOnComments，但保持相同的接口行为
	var coc models.CommentOnComments
	if err := database.DB.Where("coc_id = ?", cocID).First(&coc).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	result := database.DB.Where("coc_id = ? AND user_id = ?", coc.CocID, userID).Delete(&models.LikeOnCOCS{})
	if result.RowsAffected > 0 {
		c.JSON(http.StatusNoContent, gin.H{"msg": "success"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"msg": "not found"})
	}
}
