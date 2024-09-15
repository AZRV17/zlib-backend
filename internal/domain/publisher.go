package domain

import "gorm.io/gorm"

type Publisher struct {
	gorm.Model
	Name  string
	Books []Book
}
