package models

// Image 图片（兼容 Flask SFR 数据结构）
type Image struct {
	ImageID   uint   `json:"image_id" gorm:"primaryKey;autoIncrement;column:image_id"`
	UserID    uint   `json:"user_id" gorm:"index;column:user_id"`
	ImageURL  string `json:"image_url" gorm:"size:255;column:image_url"`
	PostID    *uint  `json:"post_id" gorm:"column:post_id"`
	CommentID *uint  `json:"comment_id" gorm:"column:comment_id"`
	FishID    *uint  `json:"fish_id" gorm:"column:fish_id"`
	IsAvatar  bool   `json:"is_avatar" gorm:"default:false;column:is_avatar"`
	IsDeleted bool   `json:"is_deleted" gorm:"default:false;column:is_deleted"`
}

func (Image) TableName() string {
	return "image"
}
