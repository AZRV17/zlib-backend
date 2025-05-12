package domain

import (
	"time"
)

type Message struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ChatID     uint       `json:"chat_id" binding:"required"`
	SenderID   uint       `json:"sender_id" binding:"required"`
	SenderRole string     `json:"sender_role" binding:"required,oneof=user librarian"`
	SenderName string     `json:"sender_name" binding:"required,min=1,max=100"`
	Content    string     `json:"content" binding:"required,min=1,max=1000"`
	ReadAt     *time.Time `json:"read_at" gorm:"default:null"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"created_at"`
}

type Chat struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID       uint      `json:"user_id" binding:"required"`
	LibrarianID  *uint     `json:"librarian_id" gorm:"default:null"`
	Status       string    `json:"status" binding:"required,oneof=active closed waiting"`
	Title        string    `json:"title" binding:"required,min=1,max=100"`
	LastActivity time.Time `gorm:"autoUpdateTime" json:"last_activity"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
