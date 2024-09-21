package controllers

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shaileshhb/equisplit/src/db"
	"github.com/shaileshhb/equisplit/src/models"
	"github.com/shaileshhb/equisplit/src/util"
	"gorm.io/gorm"
)

type UserInvitationController interface {
	Add(invitation *models.UserInvitation) error
	AcceptInvitation(invitation *models.UserInvitation) error
	DeleteInvitation(invitation *models.UserInvitation) error
	GetInvitations(invitations *[]models.UserInvitationDTO, parser *util.Parser) error
	GetGroupInvitation(invitations *[]models.UserInvitation, groupId uuid.UUID) error
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

	if tempInvitation.Id != uuid.Nil {
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

// AcceptInvitation will mark invitation as accepted and add user in the group that they were invited to.
func (ui *userInvitationController) AcceptInvitation(invitation *models.UserInvitation) error {

	err := ui.doesUserInvitationExist(invitation.Id)
	if err != nil {
		return err
	}

	err = ui.doesGroupExist(invitation.GroupId)
	if err != nil {
		return err
	}

	err = ui.doesUserExist(invitation.UserId)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(ui.db)
	defer uow.RollBack()

	err = uow.DB.Updates(map[string]interface{}{
		"IsAccepted": invitation.IsAccepted,
	}).Error
	if err != nil {
		return err
	}

	if invitation.IsAccepted != nil && *invitation.IsAccepted {
		err = uow.DB.Create(&models.UserGroup{
			UserId:  invitation.UserId,
			GroupId: invitation.GroupId,
		}).Error
		if err != nil {
			return err
		}
	}

	uow.Commit()
	return nil
}

// GetInvitations will fetch all invitations.
func (ui *userInvitationController) GetInvitations(invitations *[]models.UserInvitationDTO, parser *util.Parser) error {

	uow := db.NewUnitOfWork(ui.db)
	defer uow.RollBack()

	queryDB := ui.searchQuery(uow, parser)

	err := queryDB.Debug().Preload("User").Preload("Group").Preload("InvitedByUser").Find(&invitations).Error
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// GetGroupInvitation will fetch all invitations of specified group.
func (ui *userInvitationController) GetGroupInvitation(invitations *[]models.UserInvitation, groupId uuid.UUID) error {

	err := ui.doesGroupExist(groupId)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(ui.db)
	defer uow.RollBack()

	err = uow.DB.Where("group_id = ?", groupId).Find(&invitations).Error
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// DeleteInvitation will delete the specified invitation
func (ui *userInvitationController) DeleteInvitation(invitation *models.UserInvitation) error {

	err := ui.doesUserInvitationExist(invitation.Id)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(ui.db)
	defer uow.RollBack()

	err = uow.DB.Delete(invitation).Error
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

func (u *userInvitationController) searchQuery(uow *db.UnitOfWork, parser *util.Parser) *gorm.DB {
	queryDB := uow.DB

	if len(parser.GetQuery("userId")) > 0 {
		queryDB = uow.DB.Where("user_invitations.user_id = ?", parser.GetQuery("userId"))
	}

	return queryDB
}

// doesUserExist will check if specified user exist or not.
func (u *userInvitationController) doesUserExist(userId uuid.UUID) error {
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
func (u *userInvitationController) doesGroupExist(groupId uuid.UUID) error {
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
func (u *userInvitationController) doesUserGroupExist(userId, groupId uuid.UUID) error {
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

// doesUserInvitationExist will check if specified group exist or not.
func (u *userInvitationController) doesUserInvitationExist(invitationId uuid.UUID) error {
	err := u.db.Where("id = ?", invitationId).First(&models.UserInvitation{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("invitation not found")
		}
		return err
	}
	return nil
}
