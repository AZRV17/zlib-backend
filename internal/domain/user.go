package domain

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Role string

const (
	RoleUser      Role = "user"
	RoleAdmin     Role = "admin"
	RoleLibrarian Role = "librarian"
)

func (r *Role) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*r = Role(v)
	case string:
		*r = Role(v)
	default:
		return fmt.Errorf("unsupported type: %T", v)
	}
	return nil
}

func (r Role) Value() (driver.Value, error) {
	return string(r), nil
}

type User struct {
	ID                 uint          `json:"id" gorm:"primaryKey,autoIncrement"`
	Login              string        `json:"login" gore:"unique"`
	FullName           string        `json:"full_name"`
	Password           string        `json:"password"`
	Role               Role          `json:"role" gorm:"type:role;default:'user'"`
	Email              string        `json:"email" gorm:"unique"`
	PhoneNumber        string        `json:"phone_number" gorm:"unique"`
	PassportNumber     int           `json:"passport_number" gorm:"unique"`
	Favorites          []Favorite    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Reservations       []Reservation `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Reviews            []Review      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Logs               []Log         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	ResetPasswordToken string        `json:"-" gorm:"type:varchar(100)"`
	ResetTokenExpiry   time.Time     `json:"-"`
	CreatedAt          time.Time     `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt          time.Time     `json:"updatedAt" gorm:"autoUpdateTime"`
}
