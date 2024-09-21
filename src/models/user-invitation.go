package models

import (
	"errors"
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

func (u *UserInvitation) Validate() error {
	if u.UserId == uuid.Nil {
		return errors.New("user must be specified")
	}

	if u.GroupId == uuid.Nil {
		return errors.New("group must be specified")
	}

	if u.InvitedBy == nil || *u.InvitedBy == uuid.Nil {
		return errors.New("invited by must be specified")
	}

	return nil
}

// UserInvitationDTO entity
type UserInvitationDTO struct {
	Base
	User          *User      `json:"user" gorm:"foreignKey:UserId"`
	Group         *Group     `json:"group" gorm:"foreignKey:GroupId"`
	InvitedByUser *User      `json:"invitedByUser" gorm:"foreignKey:InvitedBy"`
	InvitedBy     *uuid.UUID `json:"invitedBy"`
	UserId        uuid.UUID  `json:"userId"`
	GroupId       uuid.UUID  `json:"groupId"`
	ExpiresOn     *time.Time `json:"-"`
	IsAccepted    *bool      `json:"isAccepted"`
}

func (*UserInvitationDTO) TableName() string {
	return "user_invitations"
}
