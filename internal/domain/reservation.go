package domain

import (
	"time"
)

type Reservation struct {
	ID           uint       `json:"id" gorm:"primaryKey,autoIncrement"`
	UserID       uint       `json:"user_id"`
	User         User       `json:"user" gorm:"foreignKey:UserID"`
	BookID       uint       `json:"book_id"`
	Book         Book       `json:"book" gorm:"foreignKey:BookID"`
	DateOfIssue  time.Time  `json:"date_of_issue"`
	DateOfReturn time.Time  `json:"date_of_return"`
	Status       string     `json:"status"`
	UniqueCodeID uint       `json:"unique_code_id"`
	UniqueCode   UniqueCode `json:"unique_code" gorm:"foreignKey:UniqueCodeID"`
	CreatedAt    time.Time  `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt    time.Time  `json:"updatedAt" gorm:"autoUpdateTime"`
}
