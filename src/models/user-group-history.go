package models

import "gorm.io/gorm"

// Payer - Represents the user who has to transfer some amount.
// Payee - Represents the user to whom the amount should be transferred.

// UserGroupHistory entity
type UserGroupHistory struct {
	gorm.Model
	Payer   User    `json:"-" gorm:"foreignKey:PayerId"`
	Payee   User    `json:"-" gorm:"foreignKey:PayeeId"`
	Group   Group   `json:"-" gorm:"foreignKey:GroupId"`
	PayerId uint    `json:"payerId" gorm:"index"`
	PayeeId uint    `json:"payeeId" gorm:"index"`
	GroupId uint    `json:"groupId" gorm:"index"`
	Amount  float64 `json:"amount" gorm:"type:float;default:0"`
}

func (*UserGroupHistory) TableName() string {
	return "user_group_history"
}
