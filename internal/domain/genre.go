package domain

import "gorm.io/gorm"

type Genre struct {
	gorm.Model
	Name        string
	Description string
	Books       []Book
}
