package models

// UserGroup entity
type UserGroup struct {
	Base
	User           User    `json:"-" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Group          Group   `json:"-" gorm:"foreignKey:GroupId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserId         uint    `json:"userId" gorm:"index"`
	GroupId        uint    `json:"groupId" gorm:"index"`
	OutgoingAmount float64 `json:"outgoingAmount" gorm:"type:float;default:0"`
	IncomingAmount float64 `json:"incomingAmount" gorm:"type:float;default:0"`

	// TODO: do I need to have a column to know who added the user in this group?
}

func (*UserGroup) TableName() string {
	return "user_groups"
}

// UserGroup entity
type UserGroupDTO struct {
	BaseDTO
	User           *UserDTO      `json:"user"`
	Group          *Group        `json:"group"`
	UserId         uint          `json:"userId"`
	GroupId        uint          `json:"groupId"`
	OutgoingAmount float64       `json:"outgoingAmount"`
	IncomingAmount float64       `json:"incomingAmount"`
	Summary        *GroupSummary `json:"summary" gorm:"-"`
}

func (*UserGroupDTO) TableName() string {
	return "user_groups"
}

// GroupSummary will contain details of how much a user has outgoing and incoming amount
type GroupSummary struct {
	UserId         uint    `json:"userId"`
	OutgoingAmount float64 `json:"outgoingAmount"`
	IncomingAmount float64 `json:"incomingAmount"`
}
