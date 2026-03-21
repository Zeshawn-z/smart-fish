package models

import "time"

// Comment 评论（兼容 Flask SFR 数据结构）
type Comment struct {
	CommentID uint      `json:"comment_id" gorm:"primaryKey;autoIncrement;column:comment_id"`
	PostID    uint      `json:"post_id" gorm:"not null;index;column:post_id"`
	UserID    uint      `json:"user_id" gorm:"not null;index;column:user_id"`
	Body      string    `json:"body" gorm:"type:text;column:body"`
	IsDeleted bool      `json:"is_deleted" gorm:"default:false;column:is_deleted"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Comment) TableName() string {
	return "comment"
}

// CommentOnComments 评论的评论（兼容 Flask SFR 数据结构）
type CommentOnComments struct {
	CocID     uint      `json:"coc_id" gorm:"primaryKey;autoIncrement;column:coc_id"`
	CommentID uint      `json:"comment_id" gorm:"not null;index;column:comment_id"`
	ToCocID   *uint     `json:"to_coc_id" gorm:"column:to_coc_id"`
	UserID    uint      `json:"user_id" gorm:"not null;index;column:user_id"`
	Body      string    `json:"body" gorm:"type:text;column:body"`
	IsDeleted bool      `json:"is_deleted" gorm:"default:false;column:is_deleted"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (CommentOnComments) TableName() string {
	return "comment_on_comments"
}
