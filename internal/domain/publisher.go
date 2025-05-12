package domain

import "time"

type Publisher struct {
	ID        uint      `json:"id" gorm:"primaryKey,autoIncrement"`
	Name      string    `json:"name"`
	Books     []Book    `json:"books" gorm:"foreignKey:PublisherID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
}
