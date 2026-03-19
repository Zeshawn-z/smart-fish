package models

import "time"

type User struct {
	BaseModel
	Username     string         `json:"username" gorm:"uniqueIndex;size:50;not null"`
	PasswordHash string         `json:"-" gorm:"size:255;not null"`
	Role         string         `json:"role" gorm:"size:20;default:user;not null"` // user, staff, admin
	Phone        string         `json:"phone" gorm:"size:20"`
	Email        string         `json:"email" gorm:"size:100"`
	RegisterTime time.Time      `json:"register_time" gorm:"autoCreateTime"`
	Favorites    []*FishingSpot `json:"favorites,omitempty" gorm:"many2many:user_favorites;"`
}

// UserResponse 用户信息响应（不含密码）
type UserResponse struct {
	ID           uint      `json:"id"`
	Username     string    `json:"username"`
	Role         string    `json:"role"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	RegisterTime time.Time `json:"register_time"`
}

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:           u.ID,
		Username:     u.Username,
		Role:         u.Role,
		Phone:        u.Phone,
		Email:        u.Email,
		RegisterTime: u.RegisterTime,
	}
}
