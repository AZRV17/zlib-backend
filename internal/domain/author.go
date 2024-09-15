package domain

import (
	"gorm.io/gorm"
	"time"
)

type Author struct {
	gorm.Model
	Name      string
	Lastname  string
	Biography string
	Birthdate time.Time
	Books     []Book `gorm:"many2many:author_books;"`
}
