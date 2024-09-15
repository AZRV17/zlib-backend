package domain

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	UserID uint
	User   User
	BookID uint
	Book   Book
}
