package models

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

// User db entity
type UserDTO struct {
	BaseDTO
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (*UserDTO) TableName() string {
	return "users"
}
