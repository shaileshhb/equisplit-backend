package controllers

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shaileshhb/equisplit/src/db"
	"github.com/shaileshhb/equisplit/src/models"
	"gorm.io/gorm"
)

type UserGroupController interface {
	AddUserToGroup(userGroup *models.UserGroup) error
	DeleteUserFromGroup(userGroup *models.UserGroup) error
	GetGroupDetails(userGroups *[]models.UserGroupDTO, groupId, userId uuid.UUID) error
	GetUserGroups(userGroups *[]models.UserGroupDTO, userId uuid.UUID) error
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

	err = u.doesUserExistInGroup(userGroup.UserId, userGroup.GroupId)
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

	err := u.doesUserGroupExist(userGroup.Id)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(u.db)
	defer uow.RollBack()

	err = uow.DB.Model(&models.UserGroup{}).Where("user_groups.user_id = ? AND user_groups.group_id = ?", userGroup.UserId, userGroup.GroupId).
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
func (u *userGroupController) GetGroupDetails(userGroups *[]models.UserGroupDTO, groupId, userId uuid.UUID) error {

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

	for index := range *userGroups {
		(*userGroups)[index].Summary = &models.GroupSummary{
			UserId: (*userGroups)[index].UserId,
		}
		if userId == (*userGroups)[index].UserId {
			continue
		}

		err = uow.DB.Select("SUM(amount) AS incoming_amount").Table("group_transactions").
			Where("group_id = ? AND payer_id = ? AND payee_id = ? AND is_paid = ?",
				(*userGroups)[index].GroupId, (*userGroups)[index].UserId, userId, false).
			Scan(&(*userGroups)[index].Summary).Error
		if err != nil {
			return err
		}

		err = uow.DB.Select("SUM(amount) AS outgoing_amount").Table("group_transactions").
			Where("group_id = ? AND payee_id = ? AND payer_id = ? AND is_paid = ?",
				(*userGroups)[index].GroupId, (*userGroups)[index].UserId, userId, false).
			Scan(&(*userGroups)[index].Summary).Error
		if err != nil {
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetUserGroups will fetch all groups for specific user.
func (u *userGroupController) GetUserGroups(userGroups *[]models.UserGroupDTO, userId uuid.UUID) error {

	err := u.doesUserExist(userId)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(u.db)
	defer uow.RollBack()

	err = uow.DB.Where("user_groups.user_id = ?", userId).Preload("Group").Find(userGroups).Error
	if err != nil {
		return err
	}

	for index := range *userGroups {
		err = uow.DB.Select("SUM(amount) AS outgoing_amount").Table("group_transactions").
			Where("group_transactions.group_id = ? AND group_transactions.payer_id = ? AND is_paid = ?",
				(*userGroups)[index].GroupId, userId, false).
			Scan(&(*userGroups)[index].Summary).Error
		if err != nil {
			return err
		}

		err = uow.DB.Select("SUM(amount) AS incoming_amount").Table("group_transactions").
			Where("group_transactions.group_id = ? AND group_transactions.payee_id = ? AND is_paid = ?",
				(*userGroups)[index].GroupId, userId, false).
			Scan(&(*userGroups)[index].Summary).Error
		if err != nil {
			return err
		}
	}

	uow.Commit()
	return nil
}

// doesUserExist will check if specified user exist or not.
func (u *userGroupController) doesUserExist(userId uuid.UUID) error {
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
func (u *userGroupController) doesGroupExist(groupId uuid.UUID) error {
	err := u.db.Where("groups.id = ?", groupId).First(&models.Group{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("group not found")
		}
		return err
	}
	return nil
}

// doesUserExistInGroup will check if specified user exist in group or not.
func (u *userGroupController) doesUserExistInGroup(userId, groupId uuid.UUID) error {
	err := u.db.Where("user_groups.user_id = ? AND user_groups.group_id = ?", userId, groupId).First(&models.UserGroup{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		return err
	}
	return errors.New("user already exists in specified group")
}

// doesUserGroupExist will check if specified user_group exist or not.
func (u *userGroupController) doesUserGroupExist(userGroupId uuid.UUID) error {
	err := u.db.Where("user_groups.id = ?", userGroupId).First(&models.UserGroup{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found in group")
		}
		return err
	}
	return nil
}
