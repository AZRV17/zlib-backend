package domain

import (
	"time"
)

type Author struct {
	ID        uint      `gorm:"primaryKey,autoIncrement" json:"id"`
	Name      string    `json:"name"`
	Lastname  string    `json:"lastname"`
	Biography string    `json:"biography"`
	Birthdate time.Time `json:"birthdate"`
	Books     []Book    `gorm:"many2many:author_books;" json:"books"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
