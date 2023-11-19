package controllers

import (
	"errors"
	"time"

	"github.com/shaileshhb/equisplit/src/db"
	"github.com/shaileshhb/equisplit/src/models"
	"gorm.io/gorm"
)

type UserInvitationController interface {
	Add(invitation *models.UserInvitation) error
}

type userInvitationController struct {
	db *gorm.DB
}

func NewUserInvitationController(db *gorm.DB) UserInvitationController {
	return &userInvitationController{
		db: db,
	}
}

// Add will add invitation for the specified user in the group.
func (ui *userInvitationController) Add(invitation *models.UserInvitation) error {

	err := ui.doesUserExist(invitation.UserId)
	if err != nil {
		return err
	}

	err = ui.doesGroupExist(invitation.GroupId)
	if err != nil {
		return err
	}

	err = ui.doesUserGroupExist(invitation.UserId, invitation.GroupId)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(ui.db)
	defer uow.RollBack()

	tempInvitation := models.UserInvitation{}

	err = uow.DB.Table("user_invitations").
		Where("group_id = ? AND user_id = ? AND expires_on > NOW() AND is_accepted = ?",
			invitation.GroupId, invitation.UserId, false).
		First(&tempInvitation).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if tempInvitation.ID > 0 {
		return errors.New("user already invited")
	}

	expiry := time.Now().Local().AddDate(0, 0, 30)

	invitation.ExpiresOn = &expiry

	err = uow.DB.Create(invitation).Error
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// doesUserExist will check if specified user exist or not.
func (u *userInvitationController) doesUserExist(userId uint) error {
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
func (u *userInvitationController) doesGroupExist(groupId uint) error {
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
func (u *userInvitationController) doesUserGroupExist(userId, groupId uint) error {
	var totalCount int64 = 0
	err := u.db.Model(&models.UserGroup{}).Where("user_groups.user_id = ? AND user_groups.group_id = ?", userId, groupId).
		Count(&totalCount).Error
	if err != nil {
		return err
	}
	if totalCount > 0 {
		return errors.New("user already exist in group")
	}
	return nil
}