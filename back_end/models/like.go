package models

// LikeOnPosts 帖子点赞（兼容 Flask SFR 数据结构）
type LikeOnPosts struct {
	LikeID uint `json:"like_id" gorm:"primaryKey;autoIncrement;column:like_id"`
	PostID uint `json:"post_id" gorm:"index;column:post_id"`
	UserID uint `json:"user_id" gorm:"index;column:user_id"`
}

func (LikeOnPosts) TableName() string {
	return "like_on_posts"
}

// LikeOnComments 评论点赞（兼容 Flask SFR 数据结构）
type LikeOnComments struct {
	LikeID    uint `json:"like_id" gorm:"primaryKey;autoIncrement;column:like_id"`
	CommentID uint `json:"comment_id" gorm:"index;column:comment_id"`
	UserID    uint `json:"user_id" gorm:"index;column:user_id"`
}

func (LikeOnComments) TableName() string {
	return "like_on_comments"
}

// LikeOnCOCS 评论的评论点赞（兼容 Flask SFR 数据结构）
type LikeOnCOCS struct {
	LikeID uint `json:"like_id" gorm:"primaryKey;autoIncrement;column:like_id"`
	CocID  uint `json:"coc_id" gorm:"index;column:coc_id"`
	UserID uint `json:"user_id" gorm:"index;column:user_id"`
}

func (LikeOnCOCS) TableName() string {
	return "like_on_cocs"
}
