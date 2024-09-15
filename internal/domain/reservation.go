package domain

import (
	"gorm.io/gorm"
	"time"
)

type Reservation struct {
	gorm.Model
	UserID       uint
	User         User
	BookID       uint
	Book         Book
	DateOfIssue  time.Time
	DateOfReturn time.Time
	Status       string
}
