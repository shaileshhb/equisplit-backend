package models

import (
	"errors"

	"github.com/google/uuid"
)

// Payer - Represents the user who has to transfer some amount.
// Payee - Represents the user to whom the amount should be transferred.

// GroupTransaction entity
type GroupTransaction struct {
	Base
	Payer       User      `json:"-" gorm:"foreignKey:PayerId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Payee       User      `json:"-" gorm:"foreignKey:PayeeId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Group       Group     `json:"-" gorm:"foreignKey:GroupId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PayerId     uuid.UUID `json:"payerId" gorm:"index;type:uuid"`
	PayeeId     uuid.UUID `json:"payeeId" gorm:"index;type:uuid"`
	GroupId     uuid.UUID `json:"groupId" gorm:"index;type:uuid"`
	Amount      float64   `json:"amount" gorm:"type:float;default:0"`
	IsPaid      bool      `json:"isPaid" gorm:"default:false"`
	IsAdjusted  bool      `json:"isAdjusted" gorm:"default:false"`
	Description *string   `json:"description" gorm:"type:text"`
}

// TableName specifies name of the table for UserGroupHistory struct.
func (*GroupTransaction) TableName() string {
	return "group_transactions"
}

func (g *GroupTransaction) Validate() error {

	if g.PayerId == uuid.Nil {
		return errors.New("payer must be specified")
	}

	if g.PayeeId == uuid.Nil {
		return errors.New("payee must be specified")
	}

	if g.GroupId == uuid.Nil {
		return errors.New("group must be specified")
	}

	if g.Amount == 0 {
		return errors.New("amount must be greater than zero")
	}
	return nil
}

// UserBalance represents the balance amount to be paid by other users.
type UserBalance struct {
	UserId  uuid.UUID `json:"user_id"`
	User    UserDTO   `json:"user" gorm:"foreignKey:UserId;"`
	GroupId uuid.UUID `json:"group_id"`
	Amount  float64   `json:"amount"`
	// Group   Group     `json:"group" gorm:"foreignKey:GroupId;"`
}
