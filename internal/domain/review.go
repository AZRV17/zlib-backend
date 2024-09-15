package domain

import "gorm.io/gorm"

type Review struct {
	gorm.Model
	UserID  uint
	User    User
	BookID  uint
	Book    Book
	Rating  float32
	Message string
}
