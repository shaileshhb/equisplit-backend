package models

import "gorm.io/gorm"

// UserGroup entity
type UserGroup struct {
	gorm.Model
	User           User    `json:"-" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Group          Group   `json:"-" gorm:"foreignKey:GroupId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserId         uint    `json:"userId" gorm:"index"`
	GroupId        uint    `json:"groupId" gorm:"index"`
	OutgoingAmount float64 `json:"outgoingAmount" gorm:"type:float;default:0"`
	IncomingAmount float64 `json:"incomingAmount" gorm:"type:float;default:0"`
}

func (*UserGroup) TableName() string {
	return "user_groups"
}
