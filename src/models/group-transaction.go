package models

// Payer - Represents the user who has to transfer some amount.
// Payee - Represents the user to whom the amount should be transferred.

// GroupTransaction entity
type GroupTransaction struct {
	Base
	Payer   User    `json:"-" gorm:"foreignKey:PayerId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Payee   User    `json:"-" gorm:"foreignKey:PayeeId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Group   Group   `json:"-" gorm:"foreignKey:GroupId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PayerId uint    `json:"payerId" gorm:"index"`
	PayeeId uint    `json:"payeeId" gorm:"index"`
	GroupId uint    `json:"groupId" gorm:"index"`
	Amount  float64 `json:"amount" gorm:"type:float;default:0"`
	IsPaid  bool    `json:"isPaid" gorm:"default:false"`
}

// TableName specifies name of the table for UserGroupHistory struct.
func (*GroupTransaction) TableName() string {
	return "group_transactions"
}
