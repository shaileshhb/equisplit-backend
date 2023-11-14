package models

import (
	"time"

	"gorm.io/gorm"
)

// Base for all models.
type Base struct {
	ID        uint           `json:"id" gorm:"primarykey;autoIncrement"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `json:"-"`
}

// BaseDTO for all models.
type BaseDTO struct {
	ID        uint           `json:"id" gorm:"primarykey;autoIncrement"`
	CreatedAt time.Time      `json:"createdAt"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
