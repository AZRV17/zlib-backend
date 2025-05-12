package domain

import "time"

type Genre struct {
	ID          uint      `json:"id" gorm:"primary_key,autoIncrement"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Books       []Book    `json:"books" gorm:"foreignKey:GenreID;constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
