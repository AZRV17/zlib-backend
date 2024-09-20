package domain

import "time"

type Notification struct {
	ID        uint      `json:"id" gorm:"primaryKey,autoIncrement"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read" gorm:"default:false"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
