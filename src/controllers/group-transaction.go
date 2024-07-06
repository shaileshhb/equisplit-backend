package controllers

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/shaileshhb/equisplit/src/db"
	"github.com/shaileshhb/equisplit/src/models"
	"gorm.io/gorm"
)

// GroupTransactionController will contain all methods to be implemented by userGroupHistory controller
type GroupTransactionController interface {
	Add(transaction *models.GroupTransaction) error
	MarkTransactionPaid(transaction *models.GroupTransaction, payerId uuid.UUID) error
	Delete(userId, transactionId uuid.UUID) error
}

type groupTransactionController struct {
	db *gorm.DB
}

// NewGroupTransactionController will return new instance of GroupTransactionController.
func NewGroupTransactionController(db *gorm.DB) GroupTransactionController {
	return &groupTransactionController{
		db: db,
	}
}

// Add will add new transaction for specified group and user.
func (g *groupTransactionController) Add(transaction *models.GroupTransaction) error {

	err := g.doesUserExist(transaction.PayeeId)
	if err != nil {
		return err
	}

	err = g.doesUserExist(transaction.PayerId)
	if err != nil {
		return err
	}

	err = g.doesGroupExist(transaction.GroupId)
	if err != nil {
		return err
	}

	err = g.doesUserExistInGroup(transaction.PayeeId, transaction.GroupId)
	if err != nil {
		return err
	}

	err = g.doesUserExistInGroup(transaction.PayerId, transaction.GroupId)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(g.db)
	defer uow.RollBack()

	err = uow.DB.Create(transaction).Error
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// MarkTransactionPaid will mark the transaction has paid
func (g *groupTransactionController) MarkTransactionPaid(transaction *models.GroupTransaction, payerId uuid.UUID) error {

	err := g.doesGroupTransactionExist(transaction.ID)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(g.db)
	defer uow.RollBack()

	err = uow.DB.Where("id = ? AND payer_id = ?", transaction.ID, payerId).
		First(&models.GroupTransaction{}).Error
	if err != nil {
		return err
	}

	err = uow.DB.Model(&models.GroupTransaction{}).Where("group_transactions.id = ?", transaction.ID).
		Updates(map[string]interface{}{
			"IsPaid": true,
		}).Error
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// Delete will delete specified transaction
func (g *groupTransactionController) Delete(userId, transactionId uuid.UUID) error {

	err := g.doesUserExist(userId)
	if err != nil {
		return err
	}

	err = g.doesGroupTransactionExist(transactionId)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(g.db)
	defer uow.RollBack()

	err = uow.DB.Model(&models.GroupTransaction{}).Where("group_transactions.id = ? AND group_transactions.payee_id = ?", transactionId, userId).
		First(&models.GroupTransaction{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("only payee can delete a transaction")
		}
		return err
	}

	err = uow.DB.Where("group_transactions.id = ?", transactionId).
		Updates(map[string]interface{}{
			"DeletedAt": time.Now(),
		}).Error
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// doesUserExist will check if specified user exist or not.
func (g *groupTransactionController) doesUserExist(userId uuid.UUID) error {
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
func (g *groupTransactionController) doesGroupExist(groupId uuid.UUID) error {
	err := g.db.Where("groups.id = ?", groupId).First(&models.Group{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("group not found")
		}
		return err
	}
	return nil
}

// doesUserExistInGroup will check if specified user exist or not.
func (g *groupTransactionController) doesUserExistInGroup(userId, groupId uuid.UUID) error {
	err := g.db.Where("user_groups.user_id = ? AND user_groups.group_id = ?", userId, groupId).
		First(&models.UserGroup{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("user not found in this group")
		}
		return err
	}
	return nil
}

// doesGroupTransactionExist will check if specified group exist or not.
func (g *groupTransactionController) doesGroupTransactionExist(transactionId uuid.UUID) error {
	err := g.db.Where("group_transactions.id = ?", transactionId).First(&models.GroupTransaction{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("transaction not found")
		}
		return err
	}
	return nil
}
