package models

import "gorm.io/gorm"

// User db entity
type User struct {
	// ID       uint   `gorm:"primaryKey;autoIncrement"`
	gorm.Model
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

func (*User) TableName() string {
	return "users"
}
