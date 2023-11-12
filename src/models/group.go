package models

import "gorm.io/gorm"

// Group entity
type Group struct {
	gorm.Model
	Name       string  `json:"name" gorm:"type:varchar(100);not null;"`
	User       User    `json:"-" gorm:"foreignKey:CreatedBy"` // added to create foregin key. can't create using constraint
	CreatedBy  uint    `json:"createdBy"`
	TotalSpent float64 `json:"totalSpent" gorm:"type:float;default:0"`
	Users      []*User `json:"users" gorm:"many2many:user_groups;"`
}

func (*Group) TableName() string {
	return "groups"
}
