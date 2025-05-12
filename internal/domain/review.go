package domain

import "time"

type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey,autoIncrement"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	BookID    uint      `json:"book_id"`
	Book      Book      `json:"book" gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
	Rating    float32   `json:"rating"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
