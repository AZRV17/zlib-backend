package domain

import (
	"time"
)

type Book struct {
	ID                uint          `gorm:"primary_key,autoIncrement" json:"id"`
	Title             string        `json:"title"`
	AuthorID          uint          `json:"author_id"`
	Author            Author        `gorm:"author" json:"author,omitempty"`
	GenreID           uint          `json:"genre_id"`
	Genre             Genre         `json:"genre"`
	Description       string        `json:"description"`
	PublisherID       uint          `json:"publisher_id"`
	Publisher         Publisher     `json:"publisher"`
	ISBN              int           `json:"isbn"`
	YearOfPublication time.Time     `json:"year_of_publication"`
	Picture           string        `json:"picture"`
	Rating            float32       `json:"rating"`
	UniqueCodes       []UniqueCode  `json:"unique_codes,omitempty" gorm:"foreignKey:BookID;references:ID"`
	Favorites         []Favorite    `json:"favorites,omitempty"`
	Reservations      []Reservation `json:"reservations,omitempty"`
	Reviews           []Review      `json:"reviews,omitempty"`
	CreatedAt         time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

type AggregatedBook struct {
	IDs               []uint    `json:"ids"`
	UniqueCodes       []int     `json:"unique_codes"`
	Title             string    `json:"title"`
	Author            Author    `json:"author"`
	Genre             Genre     `json:"genre"`
	Description       string    `json:"description"`
	Publisher         Publisher `json:"publisher"`
	ISBN              int       `json:"isbn"`
	YearOfPublication time.Time `json:"year_of_publication"`
	Picture           string    `json:"picture"`
	Rating            float32   `json:"rating"`
	IsAvailable       bool      `json:"is_available"`
}

type UniqueCode struct {
	ID          uint `gorm:"primary_key,autoIncrement" json:"id"`
	Code        int  `json:"code"`
	BookID      uint `json:"book_id"`
	Book        Book `json:"book" gorm:"foreignKey:BookID;references:ID"`
	IsAvailable bool `json:"is_available"`
}
