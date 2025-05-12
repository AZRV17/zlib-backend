package domain

import (
	"time"
)

type Log struct {
	ID      uint      `gorm:"primaryKey,autoIncrement" json:"id"`
	UserID  uint      `json:"user_id"`
	User    User      `json:"user"`
	Action  string    `json:"action"`
	Date    time.Time `json:"date"`
	Details string    `json:"details"`
}
