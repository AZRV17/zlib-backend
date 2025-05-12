package domain

import (
	"time"
)

type Book struct {
	ID                uint            `gorm:"primary_key,autoIncrement" json:"id"`
	Title             string          `json:"title"`
	AuthorID          uint            `json:"author_id"`
	Author            Author          `gorm:"author" json:"author,omitempty"`
	GenreID           uint            `json:"genre_id"`
	Genre             Genre           `json:"genre"`
	Description       string          `json:"description"`
	PublisherID       uint            `json:"publisher_id"`
	Publisher         Publisher       `json:"publisher"`
	ISBN              int             `json:"isbn"`
	YearOfPublication time.Time       `json:"year_of_publication"`
	Picture           string          `json:"picture"`
	Rating            float32         `json:"rating"`
	UniqueCodes       []UniqueCode    `json:"unique_codes,omitempty" gorm:"foreignKey:BookID;references:ID"`
	Favorites         []Favorite      `json:"favorites,omitempty"`
	Reservations      []Reservation   `json:"reservations,omitempty"`
	Reviews           []Review        `json:"reviews,omitempty"`
	EpubFile          string          `json:"epub_file"`
	AudiobookFiles    []AudiobookFile `json:"audiobook_files,omitempty" gorm:"foreignKey:BookID;references:ID"`
	CreatedAt         time.Time       `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time       `gorm:"autoUpdateTime" json:"updated_at"`
}

type AudiobookFile struct {
	ID           uint   `gorm:"primary_key,autoIncrement" json:"id"`
	BookID       uint   `json:"book_id"`
	Book         Book   `gorm:"foreignKey:BookID" json:"book,omitempty"`
	FilePath     string `json:"file_path"`
	ChapterTitle string `json:"chapter_title"`
	Order        int    `json:"order"`
}

type UniqueCode struct {
	ID          uint `gorm:"primary_key,autoIncrement" json:"id"`
	Code        int  `json:"code"`
	BookID      uint `json:"book_id"`
	Book        Book `json:"book" gorm:"foreignKey:BookID;references:ID"`
	IsAvailable bool `json:"is_available"`
}
