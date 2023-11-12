package models

import "gorm.io/gorm"

// Group entity
type Group struct {
	gorm.Model
	Name       string  `json:"name" gorm:"type:varchar(100);not null;"`
	User       User    `json:"-" gorm:"foreignKey:CreatedBy"` // added to create foregin key. can't create using constraint
	CreatedBy  uint    `json:"createdBy" gorm:"index"`
	TotalSpent float64 `json:"totalSpent" gorm:"type:float;default:0"`
}

func (*Group) TableName() string {
	return "groups"
}
