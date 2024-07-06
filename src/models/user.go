package models

import (
	"errors"
	"strings"
)

// User db entity
type User struct {
	Base
	Name     string `json:"name" gorm:"type:varchar(80);not null"`
	Email    string `json:"email" gorm:"uniqueIndex;not null"`
	Password string `json:"password" gorm:"not null"`
}

func (*User) TableName() string {
	return "users"
}

func (u *User) ValidateUser() error {
	u.Name = strings.TrimSpace(u.Name)
	u.Email = strings.TrimSpace(u.Email)

	if u.Name == "" {
		return errors.New("name must be specified")
	}

	if u.Email == "" {
		return errors.New("email must be specified")
	}

	return nil
}

// User db entity
type UserDTO struct {
	BaseDTO
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (*UserDTO) TableName() string {
	return "users"
}
