package domain

import (
	"database/sql/driver"
	"fmt"
	"gorm.io/gorm"
)

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
	RoleGuest Role = "guest"
)

func (r *Role) Scan(value interface{}) error {
	*r = Role(fmt.Sprintf("%s", value))
	return nil
}

func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

type User struct {
	gorm.Model
	Login          string
	Password       string
	Role           Role `gorm:"type:role;default:'user'"`
	Email          string
	PhoneNumber    string
	PassportNumber int
	Favorites      []Favorite
	Reservations   []Reservation
	Reviews        []Review
	Notifications  []Notification
	Logs           []Log
}
