package domain

import (
	"gorm.io/gorm"
	"time"
)

type Book struct {
	gorm.Model
	Title             string
	AuthorID          uint
	Author            []Author `gorm:"many2many:author_books;"`
	GenreID           uint
	Genre             Genre
	PublisherID       uint
	Publisher         Publisher
	ISBN              int
	YearOfPublication time.Time
	Picture           string
	Rating            float32
	UniqueCode        int
	Favorites         []Favorite
	Reservations      []Reservation
	Reviews           []Review
}
