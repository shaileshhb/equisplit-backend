package models

import "time"

// UserInvitation entity
type UserInvitation struct {
	Base
	User          User       `json:"-" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Group         Group      `json:"-" gorm:"foreignKey:GroupId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	InvitedByUser User       `json:"-" gorm:"foreignKey:InvitedBy;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserId        uint       `json:"userId" gorm:"index"`
	GroupId       uint       `json:"groupId" gorm:"index;not null"`
	InvitedBy     *uint      `json:"invitedBy" gorm:"index;not null"`
	ExpiresOn     *time.Time `json:"expiresOn" gorm:"not null"`
	IsAccepted    *bool      `json:"isAccepted" gorm:"default:false;not null"`
}

func (*UserInvitation) TableName() string {
	return "user_invitations"
}
