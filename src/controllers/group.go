package controllers

import (
	"errors"

	"github.com/google/uuid"
	"github.com/shaileshhb/equisplit/src/db"
	"github.com/shaileshhb/equisplit/src/models"
	"github.com/shaileshhb/equisplit/src/util"
	"gorm.io/gorm"
)

type GroupController interface {
	CreateGroup(group *models.Group) error
	UpdateGroup(group *models.Group) error
	DeleteGroup(group *models.Group) error
	GetUserGroups(group *[]models.GroupDTO, userId uuid.UUID, totalCount *int64, parser *util.Parser) error
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

	err := g.doesUserExist(group.CreatedBy)
	if err != nil {
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
		GroupId: group.Id,
	}).Error
	if err != nil {
		return err
	}

	// token, err := security.GenerateInviteJwt(group.Id)
	// if err != nil {
	// 	return err
	// }

	uow.Commit()
	return nil
}

// UpdateGroup will update specified group details.
func (g *groupController) UpdateGroup(group *models.Group) error {
	err := g.doesGroupExist(group.Id)
	if err != nil {
		return err
	}

	err = g.doesUserExist(group.CreatedBy)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(g.db)
	defer uow.RollBack()

	err = uow.DB.Updates(group).Error
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// DeleteGroup will delete specified group.
func (g *groupController) DeleteGroup(group *models.Group) error {
	err := g.doesGroupExist(group.Id)
	if err != nil {
		return err
	}

	err = g.db.
		Where("groups.`created_by` = ? AND groups.id = ?", group.CreatedBy, group.Id).
		First(&models.Group{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("only admin can delete this group")
		}
		return err
	}

	uow := db.NewUnitOfWork(g.db)
	defer uow.RollBack()

	err = uow.DB.Unscoped().Delete(&models.Group{}, group.Id).Error
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// GetUserGroups will fetch all groups for specified userId.
func (g *groupController) GetUserGroups(groups *[]models.GroupDTO, userId uuid.UUID, totalCount *int64, parser *util.Parser) error {

	err := g.doesUserExist(userId)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(g.db)
	defer uow.RollBack()

	whereDB := uow.DB.Where("groups.created_by = ?", userId)

	err = whereDB.Model(&models.Group{}).Count(totalCount).Error
	if err != nil {
		return err
	}

	limit, offset := parser.ParseLimitAndOffset()

	err = whereDB.Limit(limit).Offset(offset).Preload("User").Find(groups).Error
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// getUserGroupCount will fetch count of groups created by a specific user.
func (g *groupController) getUserGroupCount(uow *db.UnitOfWork, userId uuid.UUID, totalCount *int64) error {
	err := uow.DB.Model(&models.Group{}).
		Select("COUNT(groups.id)").
		Where("groups.created_by = ?", userId).
		Count(totalCount).Error
	if err != nil {
		return err
	}

	return nil
}

// doesUserExist will check if specified user exist or not.
func (g *groupController) doesUserExist(userId uuid.UUID) error {
	err := g.db.Where("users.id = ?", userId).First(&models.User{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		return err
	}
	return nil
}

// doesGroupExist will check if specified group exist or not.
func (g *groupController) doesGroupExist(groupId uuid.UUID) error {
	err := g.db.Where("groups.id = ?", groupId).First(&models.Group{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("group not found")
		}
		return err
	}
	return nil
}
