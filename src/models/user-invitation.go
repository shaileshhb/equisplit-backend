package models

import (
	"time"

	"github.com/google/uuid"
)

// UserInvitation entity
type UserInvitation struct {
	Base
	User          User       `json:"-" gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Group         Group      `json:"-" gorm:"foreignKey:GroupId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	InvitedByUser User       `json:"-" gorm:"foreignKey:InvitedBy;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserId        uuid.UUID  `json:"userId" gorm:"index;type:uuid"`
	GroupId       uuid.UUID  `json:"groupId" gorm:"index;not null;type:uuid"`
	InvitedBy     *uuid.UUID `json:"invitedBy" gorm:"index;not null;type:uuid"`
	ExpiresOn     *time.Time `json:"expiresOn" gorm:"not null"`
	IsAccepted    *bool      `json:"isAccepted" gorm:"default:false;not null"`
}

func (*UserInvitation) TableName() string {
	return "user_invitations"
}
