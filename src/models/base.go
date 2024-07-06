package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base for all models.
type Base struct {
	Id        uuid.UUID      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// BaseDTO for all models.
type BaseDTO struct {
	ID        uuid.UUID      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (u *Base) BeforeCreate(tx *gorm.DB) (err error) {
	u.Id = uuid.New()
	return
}
