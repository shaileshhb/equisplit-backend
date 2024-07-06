package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base for all models.
type Base struct {
	ID        uuid.UUID      `json:"id" gorm:"primarykey"`
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
