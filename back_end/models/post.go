package models

// Post 帖子（社区模块，兼容 Flask SFR 数据结构）
type Post struct {
	PostID    uint   `json:"post_id" gorm:"primaryKey;autoIncrement;column:post_id"`
	UserID    uint   `json:"user_id" gorm:"not null;index;column:user_id"`
	Title     string `json:"title" gorm:"size:255;column:title"`
	Body      string `json:"body" gorm:"size:255;column:body"`
	Tag       string `json:"tag" gorm:"size:50;column:tag"`
	IsDeleted bool   `json:"is_deleted" gorm:"default:false;column:is_deleted"`
}

func (Post) TableName() string {
	return "post"
}
