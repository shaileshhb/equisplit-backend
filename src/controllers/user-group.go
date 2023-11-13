package controllers

import (
	"errors"
	"time"

	"github.com/shaileshhb/equisplit/src/db"
	"github.com/shaileshhb/equisplit/src/models"
	"gorm.io/gorm"
)

type UserGroupController interface {
	AddUserToGroup(userGroup *models.UserGroup) error
	DeleteUserFromGroup(userGroup *models.UserGroup) error
	GetGroupDetails(userGroups *[]models.UserGroup, groupId uint) error
}

type userGroupController struct {
	db *gorm.DB
}

func NewUserGroupController(db *gorm.DB) UserGroupController {
	return &userGroupController{
		db: db,
	}
}

// AddUserToGroup will add specified user to the group.
func (u *userGroupController) AddUserToGroup(userGroup *models.UserGroup) error {
	err := u.doesUserExist(userGroup.UserId)
	if err != nil {
		return err
	}

	err = u.doesGroupExist(userGroup.GroupId)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(u.db)
	defer uow.RollBack()

	var totalCount int64 = 0

	err = uow.DB.Model(&models.UserGroup{}).Where("user_groups.group_id = ?", userGroup.GroupId).Count(&totalCount).Error
	if err != nil {
		return err
	}

	if totalCount >= 10 {
		return errors.New("maximum number of people already added to the group")
	}

	err = uow.DB.Create(&models.UserGroup{
		UserId:  userGroup.UserId,
		GroupId: userGroup.GroupId,
	}).Error
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// DeleteUserFromGroup will delete specified user from the group.
func (u *userGroupController) DeleteUserFromGroup(userGroup *models.UserGroup) error {
	// err := u.doesUserExist(userGroup.UserId)
	// if err != nil {
	// 	return err
	// }

	// err = u.doesGroupExist(userGroup.GroupId)
	// if err != nil {
	// 	return err
	// }

	// err = u.doesUserExistInGroup(userGroup)
	// if err != nil {
	// 	return err
	// }

	err := u.doesUserGroupExist(userGroup.ID)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(u.db)
	defer uow.RollBack()

	err = uow.DB.Where("user_groups.user_id = ? AND user_groups.group_id = ?", userGroup.UserId, userGroup.GroupId).
		Updates(map[string]interface{}{
			"DeletedAt": time.Now(),
		}).Error
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// GetGroupDetails will fetch all user details of specified group.
func (u *userGroupController) GetGroupDetails(userGroups *[]models.UserGroup, groupId uint) error {

	err := u.doesGroupExist(groupId)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(u.db)
	defer uow.RollBack()

	err = uow.DB.Where("user_groups.group_id = ?", groupId).
		Preload("User").Find(userGroups).Error
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// doesUserExist will check if specified user exist or not.
func (u *userGroupController) doesUserExist(userId uint) error {
	err := u.db.Where("users.id = ?", userId).First(&models.User{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found")
		}
		return err
	}
	return nil
}

// doesGroupExist will check if specified group exist or not.
func (u *userGroupController) doesGroupExist(groupId uint) error {
	err := u.db.Where("groups.id = ?", groupId).First(&models.Group{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("group not found")
		}
		return err
	}
	return nil
}

// doesUserGroupExist will check if specified user_group exist or not.
func (u *userGroupController) doesUserGroupExist(userGroupId uint) error {
	err := u.db.Where("user_groups.id = ?", userGroupId).First(&models.UserGroup{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found in group")
		}
		return err
	}
	return nil
}

// doesUserExistInGroup will check if specified user exist or not.
// func (u *userGroupController) doesUserExistInGroup(userGroup *models.UserGroup) error {
// 	tempUserGroup := models.UserGroup{}
// 	err := u.db.Where("user_groups.user_id = ? AND user_groups.group_id = ?", userGroup.UserId, userGroup.GroupId).
// 		First(&tempUserGroup).Error
// 	if err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return errors.New("user not found in this group")
// 		}
// 		return err
// 	}

// 	userGroup.ID = tempUserGroup.ID

// 	return nil
// }
