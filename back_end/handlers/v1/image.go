package v1

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"smart-fish/back_end/database"
	"smart-fish/back_end/middleware"
	"smart-fish/back_end/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadPostImage POST /api/v1/image/post - 上传帖子图片
func UploadPostImage(c *gin.Context) {
	postIDStr := c.PostForm("post_id")
	if postIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing or invalid post_id"})
		return
	}
	postID, err := parseUintParam(postIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing or invalid post_id"})
		return
	}
	handleImageUpload(c, "post", postID)
}

// UploadCommentImage POST /api/v1/image/comment - 上传评论图片
func UploadCommentImage(c *gin.Context) {
	commentIDStr := c.PostForm("comment_id")
	if commentIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing or invalid comment_id"})
		return
	}
	commentID, err := parseUintParam(commentIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing or invalid comment_id"})
		return
	}
	handleImageUpload(c, "comment", commentID)
}

// UploadFishImage POST /api/v1/image/fish - 上传鱼的图片
func UploadFishImage(c *gin.Context) {
	fishIDStr := c.PostForm("fish_id")
	if fishIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing or invalid fish_id"})
		return
	}
	fishID, err := parseUintParam(fishIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing or invalid fish_id"})
		return
	}
	handleImageUpload(c, "fish", fishID)
}

// UploadAvatarImage POST /api/v1/image/avatar - 上传用户头像
func UploadAvatarImage(c *gin.Context) {
	userID, ok := middleware.GetFlaskUserID(c)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing or invalid fish_id"}) // Flask 原始错误消息
		return
	}
	handleImageUpload(c, "avatar", userID)
}

// handleImageUpload 通用图片上传处理
func handleImageUpload(c *gin.Context, entityType string, entityID uint) {
	userID, ok := middleware.GetFlaskUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
		return
	}

	// 验证实体存在和权限
	if !validateEntity(c, entityType, entityID, userID) {
		return
	}

	file, err := c.FormFile("picbed")
	if err != nil || file == nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Error, No File"})
		return
	}

	// 保存到本地
	uploadDir := "static/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create upload directory"})
		return
	}

	filename := fmt.Sprintf("%s_%s", uuid.New().String(), filepath.Base(file.Filename))
	savePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to save file"})
		return
	}

	// 构造图片 URL（使用相对路径，由服务端静态文件服务提供）
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	imageURL := fmt.Sprintf("%s://%s/static/uploads/%s", scheme, c.Request.Host, filename)

	// 处理头像：软删除旧头像
	if entityType == "avatar" {
		database.DB.Model(&models.Image{}).
			Where("user_id = ? AND is_avatar = ? AND is_deleted = ?", userID, true, false).
			Update("is_deleted", true)
	}

	// 创建图片记录
	image := models.Image{
		UserID:   userID,
		ImageURL: imageURL,
	}

	switch entityType {
	case "post":
		image.PostID = &entityID
	case "comment":
		image.CommentID = &entityID
	case "fish":
		image.FishID = &entityID
	case "avatar":
		image.IsAvatar = true
	}

	if err := database.DB.Create(&image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to save image record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":       "success",
		"image_id":  image.ImageID,
		"image_url": imageURL,
	})
}

// validateEntity 验证实体存在和权限
func validateEntity(c *gin.Context, entityType string, entityID, userID uint) bool {
	switch entityType {
	case "post":
		var post models.Post
		if err := database.DB.Where("post_id = ?", entityID).First(&post).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing or invalid post_id"})
			return false
		}
		if post.UserID != userID {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Permission denied"})
			return false
		}
	case "comment":
		var comment models.Comment
		if err := database.DB.Where("comment_id = ?", entityID).First(&comment).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing or invalid comment_id"})
			return false
		}
		if comment.UserID != userID {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Permission denied"})
			return false
		}
	case "fish":
		var fish models.FishCaught
		if err := database.DB.Where("fish_id = ?", entityID).First(&fish).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing or invalid fish_id"})
			return false
		}
		// 验证鱼属于当前用户的记录
		var record models.FishingRecord
		if err := database.DB.Where("record_id = ? AND is_deleted = ?", fish.RecordID, false).First(&record).Error; err != nil || record.UserID != userID {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Permission denied"})
			return false
		}
	case "avatar":
		var user models.User
		if err := database.DB.First(&user, entityID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Missing or invalid user"})
			return false
		}
		// Flask 中 avatar 的 entityID 就是 userID
		if user.ID != userID {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "Permission denied"})
			return false
		}
	}
	return true
}
