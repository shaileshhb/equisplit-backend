package models

import "github.com/google/uuid"

// Group entity
type Group struct {
	Base
	Name       string    `json:"name" gorm:"type:varchar(100);not null;"`
	User       User      `json:"-" gorm:"foreignKey:CreatedBy"` // added to create foregin key. can't create using constraint
	CreatedBy  uuid.UUID `json:"createdBy" gorm:"index;type:uuid"`
	TotalSpent float64   `json:"totalSpent" gorm:"type:float;default:0"`
	Tag        *string   `json:"tag" gorm:"type:varchar(50)"`
}

func (*Group) TableName() string {
	return "groups"
}
