package domain

import "time"

type Favorite struct {
	ID        uint      `gorm:"primary_key,autoIncrement" json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user"`
	BookID    uint      `json:"book_id"`
	Book      Book      `json:"book"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
