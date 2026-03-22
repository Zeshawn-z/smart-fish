package v1

import (
	"net/http"

	"smart-fish/back_end/dao"
	"smart-fish/back_end/middleware"
	"smart-fish/back_end/models"

	"github.com/gin-gonic/gin"
)

// ===== Post Likes =====

// GetPostLikes GET /api/v1/post/:post_id/like
func GetPostLikes(c *gin.Context) {
	postID := c.Param("post_id")

	// 验证帖子存在
	post, err := dao.GetPostByPostIDUint(parseUintSafe(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	count := dao.GetPostLikeCount(post.PostID)
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
	post, err := dao.GetPostByPostIDUint(parseUintSafe(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	// 检查是否已点赞
	if existing := dao.FindPostLike(post.PostID, userID); existing != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Already Liked"})
		return
	}

	like := models.LikeOnPosts{
		PostID: post.PostID,
		UserID: userID,
	}
	dao.CreatePostLike(&like)
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
	post, err := dao.GetPostByPostIDUint(parseUintSafe(postID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	affected := dao.DeletePostLike(post.PostID, userID)
	if affected > 0 {
		c.JSON(http.StatusNoContent, gin.H{"msg": "success"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"msg": "not found"})
	}
}

// ===== Comment Likes =====

// GetCommentLikes GET /api/v1/comment/:comment_id/like
func GetCommentLikes(c *gin.Context) {
	commentID := c.Param("comment_id")

	comment, err := dao.GetCommentByCommentIDUint(parseUintSafe(commentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	count := dao.GetCommentLikeCount(comment.CommentID)
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

	comment, err := dao.GetCommentByCommentIDUint(parseUintSafe(commentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	if existing := dao.FindCommentLike(comment.CommentID, userID); existing != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Already Liked"})
		return
	}

	like := models.LikeOnComments{
		CommentID: comment.CommentID,
		UserID:    userID,
	}
	dao.CreateCommentLike(&like)
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

	comment, err := dao.GetCommentByCommentIDUint(parseUintSafe(commentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	affected := dao.DeleteCommentLike(comment.CommentID, userID)
	if affected > 0 {
		c.JSON(http.StatusNoContent, gin.H{"msg": "success"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"msg": "not found"})
	}
}

// ===== COC Likes =====

// GetCocLikes GET /api/v1/comment_on_comments/:coc_id/like
func GetCocLikes(c *gin.Context) {
	cocID := c.Param("coc_id")

	coc, err := dao.GetCocByCocID(parseUintSafe(cocID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	count := dao.GetCocLikeCount(coc.CocID)
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

	coc, err := dao.GetCocByCocID(parseUintSafe(cocID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	if existing := dao.FindCocLike(coc.CocID, userID); existing != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Already Liked"})
		return
	}

	like := models.LikeOnCOCS{
		CocID:  coc.CocID,
		UserID: userID,
	}
	dao.CreateCocLike(&like)
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
	coc, err := dao.GetCocByCocID(parseUintSafe(cocID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	affected := dao.DeleteCocLike(coc.CocID, userID)
	if affected > 0 {
		c.JSON(http.StatusNoContent, gin.H{"msg": "success"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"msg": "not found"})
	}
}
