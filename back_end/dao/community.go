package dao

import (
	"smart-fish/back_end/database"
	"smart-fish/back_end/models"

	"gorm.io/gorm"
)

// --- Post ---

// ListPostsQuery 帖子列表查询（不执行，返回 *gorm.DB 供分页使用）
func ListPostsQuery(tag, search, userID string) *gorm.DB {
	query := database.DB.Model(&models.Post{}).Where("is_deleted = ?", false)
	if tag != "" {
		query = query.Where("tag = ?", tag)
	}
	if search != "" {
		query = query.Where("title LIKE ? OR body LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	return query
}

// GetPostByID 根据 ID 查找帖子
func GetPostByID(id int) (*models.Post, error) {
	var post models.Post
	err := database.DB.Where("post_id = ? AND is_deleted = ?", id, false).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// CreatePost 创建帖子
func CreatePost(post *models.Post) error {
	return database.DB.Create(post).Error
}

// UpdatePost 更新帖子字段
func UpdatePost(post *models.Post, updates map[string]interface{}) error {
	return database.DB.Model(post).Updates(updates).Error
}

// RefreshPost 重新加载帖子
func RefreshPost(post *models.Post, id int) error {
	return database.DB.Where("post_id = ?", id).First(post).Error
}

// SoftDeletePost 软删除帖子
func SoftDeletePost(post *models.Post) error {
	return database.DB.Model(post).Update("is_deleted", true).Error
}

// --- Comment ---

// ListCommentsQuery 评论列表查询
func ListCommentsQuery(postID string) *gorm.DB {
	query := database.DB.Model(&models.Comment{}).Where("is_deleted = ?", false)
	if postID != "" {
		query = query.Where("post_id = ?", postID)
	}
	return query
}

// GetCommentByID 根据 ID 查找评论
func GetCommentByID(id int) (*models.Comment, error) {
	var comment models.Comment
	err := database.DB.Where("comment_id = ? AND is_deleted = ?", id, false).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// CreateComment 创建评论
func CreateComment(comment *models.Comment) error {
	return database.DB.Create(comment).Error
}

// SoftDeleteComment 软删除评论
func SoftDeleteComment(comment *models.Comment) error {
	return database.DB.Model(comment).Update("is_deleted", true).Error
}

// GetCommentsByPostID 查询帖子所有评论
func GetCommentsByPostID(postID uint) []models.Comment {
	var comments []models.Comment
	database.DB.Where("post_id = ? AND is_deleted = ?", postID, false).Order("comment_id ASC").Find(&comments)
	return comments
}

// --- SubComment (CommentOnComments) ---

// GetSubCommentsByCommentIDs 批量查询子评论
func GetSubCommentsByCommentIDs(commentIDs []uint) []models.CommentOnComments {
	var cocs []models.CommentOnComments
	if len(commentIDs) > 0 {
		database.DB.Where("comment_id IN ? AND is_deleted = ?", commentIDs, false).Find(&cocs)
	}
	return cocs
}

// GetSubCommentsByCommentID 查询单个评论的子评论
func GetSubCommentsByCommentID(commentID int) []models.CommentOnComments {
	var cocs []models.CommentOnComments
	database.DB.Where("comment_id = ? AND is_deleted = ?", commentID, false).Find(&cocs)
	return cocs
}

// GetSubCommentByID 根据 ID 查找子评论
func GetSubCommentByID(cocID uint) (*models.CommentOnComments, error) {
	var coc models.CommentOnComments
	err := database.DB.Where("coc_id = ?", cocID).First(&coc).Error
	if err != nil {
		return nil, err
	}
	return &coc, nil
}

// CreateSubComment 创建子评论
func CreateSubComment(coc *models.CommentOnComments) error {
	return database.DB.Create(coc).Error
}

// --- Like ---

// GetPostLikeCount 帖子点赞数
func GetPostLikeCount(postID uint) int64 {
	var count int64
	database.DB.Model(&models.LikeOnPosts{}).Where("post_id = ?", postID).Count(&count)
	return count
}

// FindPostLike 查找帖子点赞记录
func FindPostLike(postID, userID uint) *models.LikeOnPosts {
	var like models.LikeOnPosts
	if err := database.DB.Where("post_id = ? AND user_id = ?", postID, userID).First(&like).Error; err == nil {
		return &like
	}
	return nil
}

// CreatePostLike 创建帖子点赞
func CreatePostLike(like *models.LikeOnPosts) error {
	return database.DB.Create(like).Error
}

// DeletePostLike 删除帖子点赞
func DeletePostLike(postID, userID uint) int64 {
	result := database.DB.Where("post_id = ? AND user_id = ?", postID, userID).Delete(&models.LikeOnPosts{})
	return result.RowsAffected
}

// GetCommentLikeCount 评论点赞数
func GetCommentLikeCount(commentID uint) int64 {
	var count int64
	database.DB.Model(&models.LikeOnComments{}).Where("comment_id = ?", commentID).Count(&count)
	return count
}

// BatchGetCommentLikeCounts 批量获取评论点赞数
func BatchGetCommentLikeCounts(commentIDs []uint) map[uint]int64 {
	type likeCount struct {
		CommentID uint
		Cnt       int64
	}
	result := make(map[uint]int64)
	if len(commentIDs) == 0 {
		return result
	}
	var counts []likeCount
	database.DB.Model(&models.LikeOnComments{}).
		Select("comment_id, count(*) as cnt").
		Where("comment_id IN ?", commentIDs).
		Group("comment_id").
		Scan(&counts)
	for _, c := range counts {
		result[c.CommentID] = c.Cnt
	}
	return result
}

// FindCommentLike 查找评论点赞记录
func FindCommentLike(commentID, userID uint) *models.LikeOnComments {
	var like models.LikeOnComments
	if err := database.DB.Where("comment_id = ? AND user_id = ?", commentID, userID).First(&like).Error; err == nil {
		return &like
	}
	return nil
}

// CreateCommentLike 创建评论点赞
func CreateCommentLike(like *models.LikeOnComments) error {
	return database.DB.Create(like).Error
}

// DeleteCommentLike 删除评论点赞
func DeleteCommentLike(commentID, userID uint) int64 {
	result := database.DB.Where("comment_id = ? AND user_id = ?", commentID, userID).Delete(&models.LikeOnComments{})
	return result.RowsAffected
}

// GetPostCommentCount 帖子评论数
func GetPostCommentCount(postID uint) int64 {
	var count int64
	database.DB.Model(&models.Comment{}).Where("post_id = ? AND is_deleted = ?", postID, false).Count(&count)
	return count
}

// GetUsersByIDs 批量查询用户
func GetUsersByIDs(userIDs []uint) map[uint]models.User {
	result := make(map[uint]models.User)
	if len(userIDs) == 0 {
		return result
	}
	var users []models.User
	database.DB.Where("id IN ?", userIDs).Find(&users)
	for _, u := range users {
		result[u.ID] = u
	}
	return result
}

// --- COC Likes (v1 兼容) ---

// GetCocLikeCount 子评论点赞数
func GetCocLikeCount(cocID uint) int64 {
	var count int64
	database.DB.Model(&models.LikeOnCOCS{}).Where("coc_id = ?", cocID).Count(&count)
	return count
}

// FindCocLike 查找子评论点赞记录
func FindCocLike(cocID, userID uint) *models.LikeOnCOCS {
	var like models.LikeOnCOCS
	if err := database.DB.Where("coc_id = ? AND user_id = ?", cocID, userID).First(&like).Error; err == nil {
		return &like
	}
	return nil
}

// CreateCocLike 创建子评论点赞
func CreateCocLike(like *models.LikeOnCOCS) error {
	return database.DB.Create(like).Error
}

// DeleteCocLike 删除子评论点赞
func DeleteCocLike(cocID, userID uint) int64 {
	result := database.DB.Where("coc_id = ? AND user_id = ?", cocID, userID).Delete(&models.LikeOnCOCS{})
	return result.RowsAffected
}

// --- v1 兼容查找函数 ---

// GetPostByPostIDUint 通过 post_id (uint) 查找帖子（兼容 v1 不检查 is_deleted）
func GetPostByPostIDUint(postID uint) (*models.Post, error) {
	var post models.Post
	err := database.DB.Where("post_id = ?", postID).First(&post).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

// GetCommentByCommentIDUint 通过 comment_id (uint) 查找评论
func GetCommentByCommentIDUint(commentID uint) (*models.Comment, error) {
	var comment models.Comment
	err := database.DB.Where("comment_id = ?", commentID).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// GetCocByCocID 通过 coc_id 查找子评论
func GetCocByCocID(cocID uint) (*models.CommentOnComments, error) {
	var coc models.CommentOnComments
	err := database.DB.Where("coc_id = ?", cocID).First(&coc).Error
	if err != nil {
		return nil, err
	}
	return &coc, nil
}

// GetCommentsByPostIDV1 查询帖子评论（v1 兼容，不排序）
func GetCommentsByPostIDV1(postID string) []models.Comment {
	var comments []models.Comment
	database.DB.Where("post_id = ? AND is_deleted = ?", postID, false).Find(&comments)
	return comments
}

// GetSubCommentsByCommentIDV1 查询子评论（v1 兼容）
func GetSubCommentsByCommentIDV1(commentID uint) []models.CommentOnComments {
	var cocs []models.CommentOnComments
	database.DB.Where("comment_id = ? AND is_deleted = ?", commentID, false).Find(&cocs)
	return cocs
}

// GetAllPostsV1 获取所有未删除帖子（v1 兼容）
func GetAllPostsV1() []models.Post {
	var posts []models.Post
	database.DB.Where("is_deleted = ?", false).Find(&posts)
	return posts
}

// GetPostsByUserIDV1 获取用户帖子（v1 兼容）
func GetPostsByUserIDV1(userID uint) []models.Post {
	var posts []models.Post
	database.DB.Where("user_id = ? AND is_deleted = ?", userID, false).Find(&posts)
	return posts
}
