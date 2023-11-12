package controllers

import (
	"github.com/shaileshhb/equisplit/src/models"
	"gorm.io/gorm"
)

type GroupController interface {
}

type groupController struct {
	db *gorm.DB
}

func NewGroupController(db *gorm.DB) GroupController {
	return &groupController{
		db: db,
	}
}

// CreateGroup will create new group for specified user.
func (g *groupController) CreateGroup(group *models.Group) error {

	return nil
}
