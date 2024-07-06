package models

import "github.com/google/uuid"

// UserGroup entity
type UserGroup struct {
	Base
	User           User      `json:"-" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Group          Group     `json:"-" gorm:"foreignKey:GroupId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserId         uuid.UUID `json:"userId" gorm:"index;type:uuid"`
	GroupId        uuid.UUID `json:"groupId" gorm:"index;type:uuid"`
	OutgoingAmount float64   `json:"outgoingAmount" gorm:"type:float;default:0"`
	IncomingAmount float64   `json:"incomingAmount" gorm:"type:float;default:0"`
}

func (*UserGroup) TableName() string {
	return "user_groups"
}

// UserGroup entity
type UserGroupDTO struct {
	BaseDTO
	User           *UserDTO      `json:"user"`
	Group          *Group        `json:"group"`
	UserId         uuid.UUID     `json:"userId"`
	GroupId        uuid.UUID     `json:"groupId"`
	OutgoingAmount float64       `json:"outgoingAmount"`
	IncomingAmount float64       `json:"incomingAmount"`
	Summary        *GroupSummary `json:"summary" gorm:"-"`
}

func (*UserGroupDTO) TableName() string {
	return "user_groups"
}

// GroupSummary will contain details of how much a user has outgoing and incoming amount
type GroupSummary struct {
	UserId         uuid.UUID `json:"userId"`
	OutgoingAmount float64   `json:"outgoingAmount"`
	IncomingAmount float64   `json:"incomingAmount"`
}
