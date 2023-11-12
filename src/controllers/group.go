package controllers

import (
	"errors"

	"github.com/shaileshhb/equisplit/src/db"
	"github.com/shaileshhb/equisplit/src/models"
	"gorm.io/gorm"
)

type GroupController interface {
	CreateGroup(group *models.Group) error
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

	err := g.db.Where("users.id = ?", group.CreatedBy).First(&models.User{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		return err
	}

	uow := db.NewUnitOfWork(g.db)
	defer uow.RollBack()

	var totalCount int64 = 0
	err = g.getUserGroupCount(uow, group.CreatedBy, &totalCount)
	if err != nil {
		return err
	}

	if totalCount > 10 {
		return errors.New("maximum groups already created")
	}

	err = uow.DB.Create(group).Error
	if err != nil {
		return err
	}

	err = uow.DB.Create(&models.UserGroup{
		UserId:  group.CreatedBy,
		GroupId: group.ID,
	}).Error
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// getUserGroupCount will fetch count of groups created by a specific user.
func (g *groupController) getUserGroupCount(uow *db.UnitOfWork, userId uint, totalCount *int64) error {

	err := uow.DB.Model(&models.Group{}).
		Select("COUNT(groups.id)").
		Where("groups.created_by = ?", userId).
		Count(totalCount).Error
	if err != nil {
		return err
	}

	return nil
}
