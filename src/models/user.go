package models

import "gorm.io/gorm"

// User db entity
type User struct {
	// ID       uint   `gorm:"primaryKey;autoIncrement"`
	gorm.Model
	Name     string `json:"name" gorm:"type:varchar(80);not null"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Password string `json:"password" gorm:"not null"`
}

func (*User) TableName() string {
	return "users"
}
