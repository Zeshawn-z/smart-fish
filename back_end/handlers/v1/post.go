package v1

import (
	"net/http"

	"smart-fish/back_end/database"
	"smart-fish/back_end/middleware"
	"smart-fish/back_end/models"

	"github.com/gin-gonic/gin"
)

// GetPostList GET /api/v1/post - 获取所有帖子列表
func GetPostList(c *gin.Context) {
	var posts []models.Post
	database.DB.Where("is_deleted = ?", false).Find(&posts)

	postsList := make([]gin.H, 0, len(posts))
	for _, post := range posts {
		// 获取帖子的第一张图片
		var image models.Image
		var imageURL interface{} = nil
		if err := database.DB.Where("post_id = ? AND is_deleted = ?", post.PostID, false).First(&image).Error; err == nil {
			imageURL = image.ImageURL
		}

		postsList = append(postsList, gin.H{
			"post_id":   post.PostID,
			"user_id":   post.UserID,
			"title":     post.Title,
			"body":      post.Body,
			"tag":       post.Tag,
			"image_url": imageURL,
		})
	}

	c.JSON(http.StatusOK, gin.H{"posts_list": postsList})
}

// GetPostSelf GET /api/v1/post/self - 获取当前用户的帖子
func GetPostSelf(c *gin.Context) {
	userID, ok := middleware.GetFlaskUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
		return
	}

	var posts []models.Post
	database.DB.Where("user_id = ? AND is_deleted = ?", userID, false).Find(&posts)

	postsList := make([]gin.H, 0, len(posts))
	for _, post := range posts {
		var image models.Image
		var imageURL interface{} = nil
		if err := database.DB.Where("post_id = ? AND is_deleted = ?", post.PostID, false).First(&image).Error; err == nil {
			imageURL = image.ImageURL
		}

		postsList = append(postsList, gin.H{
			"post_id":   post.PostID,
			"user_id":   post.UserID,
			"title":     post.Title,
			"body":      post.Body,
			"tag":       post.Tag,
			"image_url": imageURL,
		})
	}

	c.JSON(http.StatusOK, gin.H{"posts_list": postsList, "user_id": userID})
}

// CreatePost POST /api/v1/post - 创建帖子
func CreatePost(c *gin.Context) {
	userID, ok := middleware.GetFlaskUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"msg": "Missing or invalid JWT"})
		return
	}

	var input struct {
		Title string `json:"title"`
		Body  string `json:"body"`
		Tag   string `json:"tag"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": "Invalid request data"})
		return
	}

	post := models.Post{
		UserID: userID,
		Title:  input.Title,
		Body:   input.Body,
		Tag:    input.Tag,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": "Failed to create post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":     "New Post Created Successfully",
		"post_id": post.PostID,
	})
}

// GetPost GET /api/v1/post/:post_id - 获取单个帖子详情
func GetPost(c *gin.Context) {
	postID := c.Param("post_id")

	var post models.Post
	if err := database.DB.Where("post_id = ?", postID).First(&post).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": "Not Found"})
		return
	}

	// 获取帖子的所有图片
	var images []models.Image
	database.DB.Where("post_id = ? AND is_deleted = ?", post.PostID, false).Find(&images)
	imageURLs := make([]string, 0, len(images))
	for _, img := range images {
		imageURLs = append(imageURLs, img.ImageURL)
	}

	c.JSON(http.StatusOK, gin.H{
		"post_id":   post.PostID,
		"user_id":   post.UserID,
		"title":     post.Title,
		"body":      post.Body,
		"tag":       post.Tag,
		"image_url": imageURLs,
	})
}
