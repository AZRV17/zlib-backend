package domain

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	UserID  uint
	User    User
	Message string
	IsRead  bool
}
