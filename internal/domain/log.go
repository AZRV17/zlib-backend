package domain

import (
	"gorm.io/gorm"
	"time"
)

type Log struct {
	gorm.Model
	UserID  uint
	User    User
	Action  string
	Date    time.Time
	Details string
}
