package domain

import (
	"time"
)

type Book struct {
	ID                uint          `gorm:"primary_key,autoIncrement" json:"id"`
	Title             string        `json:"title"`
	AuthorID          uint          `json:"author_id"`
	Author            []Author      `gorm:"many2many:author_books;" json:"author,omitempty"`
	GenreID           uint          `json:"genre_id"`
	Genre             Genre         `json:"genre"`
	PublisherID       uint          `json:"publisher_id"`
	Publisher         Publisher     `json:"publisher"`
	ISBN              int           `json:"isbn"`
	YearOfPublication time.Time     `json:"year_of_publication"`
	Picture           string        `json:"picture"`
	Rating            float32       `json:"rating"`
	UniqueCode        int           `gorm:"unique" json:"unique_code"`
	Favorites         []Favorite    `json:"favorites,omitempty"`
	Reservations      []Reservation `json:"reservations,omitempty"`
	Reviews           []Review      `json:"reviews,omitempty"`
	CreatedAt         time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}
