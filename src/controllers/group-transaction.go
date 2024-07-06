package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shaileshhb/equisplit/src/db"
	"github.com/shaileshhb/equisplit/src/models"
	"gorm.io/gorm"
)

// GroupTransactionController will contain all methods to be implemented by userGroupHistory controller
type GroupTransactionController interface {
	Add(*models.GroupTransaction, ...*db.UnitOfWork) error
	AddMulitple(transaction *[]models.GroupTransaction) error
	MarkTransactionPaid(*models.GroupTransaction, uuid.UUID) error
	GetTransactionDetails(userBalance *[]models.UserBalance, userId, groupId uuid.UUID) error
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
func (g *groupTransactionController) Add(transaction *models.GroupTransaction,
	uows ...*db.UnitOfWork) error {

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

	var uow *db.UnitOfWork

	if len(uows) == 0 {
		uow = db.NewUnitOfWork(g.db)
		defer uow.RollBack()
	} else {
		uow = uows[0]
	}

	err = uow.DB.Create(transaction).Error
	if err != nil {
		return err
	}

	// updates payer incoming amount.
	err = g.setUserIncomingAmount(uow, transaction)
	if err != nil {
		return err
	}

	// updates payees outgoing amount.
	err = g.setUserOutgoingAmount(uow, transaction)
	if err != nil {
		return err
	}

	if len(uows) == 0 {
		uow.Commit()
	}

	return nil
}

func (g *groupTransactionController) setUserIncomingAmount(uow *db.UnitOfWork, transaction *models.GroupTransaction) error {
	fmt.Println("==================setUserIncomingAmount==============================")
	payerAmount := struct {
		Amount  float64
		PayerId uuid.UUID
	}{}

	payerAmount.PayerId = transaction.PayerId

	// set incoming for payer
	err := uow.DB.Select("payer_id, sum(amount) AS amount").Where("payer_id = ? AND is_paid = ? AND is_adjusted = ?",
		transaction.PayerId, false, false).Group("payer_id").Find(&models.GroupTransaction{}).Scan(&payerAmount).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if payerAmount.PayerId != uuid.Nil {
		err = uow.DB.Model(&models.UserGroup{}).Where("user_id = ? AND group_id = ?", transaction.PayerId,
			transaction.GroupId).Updates(map[string]interface{}{
			"IncomingAmount": payerAmount.Amount,
		}).Error
		if err != nil {
			return err
		}
	}

	payerAmount.Amount = 0
	payerAmount.PayerId = transaction.PayeeId

	// set incoming for payee
	err = uow.DB.Select("payer_id, sum(amount) AS amount").Where("payer_id = ? AND is_paid = ? AND is_adjusted = ?",
		transaction.PayeeId, false, false).Group("payer_id").Find(&models.GroupTransaction{}).
		Scan(&payerAmount).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if payerAmount.PayerId != uuid.Nil {
		err = uow.DB.Model(&models.UserGroup{}).Where("user_id = ? AND group_id = ?", transaction.PayeeId,
			transaction.GroupId).Updates(map[string]interface{}{
			"IncomingAmount": payerAmount.Amount,
		}).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *groupTransactionController) setUserOutgoingAmount(uow *db.UnitOfWork, transaction *models.GroupTransaction) error {
	fmt.Println("==================setUserOutgoingAmount==============================")
	payeeAmount := struct {
		Amount  float64
		PayeeId uuid.UUID
	}{}

	payeeAmount.PayeeId = transaction.PayerId

	// set outcoming for payer
	err := uow.DB.Select("payee_id, sum(amount)").Where("payee_id = ? AND is_paid = ? AND is_adjusted = ?",
		transaction.PayerId, false, false).Group("payee_id").Find(&models.GroupTransaction{}).
		Scan(&payeeAmount).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if payeeAmount.PayeeId != uuid.Nil {
		err = uow.DB.Model(&models.UserGroup{}).Where("user_id = ? AND group_id = ?", transaction.PayerId,
			transaction.GroupId).Updates(map[string]interface{}{
			"OutgoingAmount": payeeAmount.Amount,
		}).Error
		if err != nil {
			return err
		}
	}

	payeeAmount.Amount = 0
	payeeAmount.PayeeId = transaction.PayeeId

	// set incoming for payee
	err = uow.DB.Select("payee_id, sum(amount) AS amount").Where("payee_id = ? AND is_paid = ? AND is_adjusted = ?",
		transaction.PayeeId, false, false).Group("payee_id").Find(&models.GroupTransaction{}).
		Scan(&payeeAmount).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if payeeAmount.PayeeId != uuid.Nil {
		err = uow.DB.Model(&models.UserGroup{}).Where("user_id = ? AND group_id = ?", transaction.PayeeId,
			transaction.GroupId).Updates(map[string]interface{}{
			"OutgoingAmount": payeeAmount.Amount,
		}).Error
		if err != nil {
			return err
		}
	}

	return nil
}

// AddMulitple will add new transaction for specified group.
func (g *groupTransactionController) AddMulitple(transaction *[]models.GroupTransaction) error {

	uow := db.NewUnitOfWork(g.db)
	defer uow.RollBack()

	for _, t := range *transaction {
		err := g.Add(&t, uow)
		if err != nil {
			return err
		}
	}

	uow.Commit()
	return nil
}

// MarkTransactionPaid will mark the transaction has paid
func (g *groupTransactionController) MarkTransactionPaid(transaction *models.GroupTransaction, payeeId uuid.UUID) error {

	err := g.doesGroupTransactionExist(transaction.Id)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(g.db)
	defer uow.RollBack()

	err = uow.DB.Where("id = ? AND payee_id = ?", transaction.Id, payeeId).
		First(&models.GroupTransaction{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("only payee can mark transaction as paid")
		}
		return err
	}

	err = uow.DB.Model(&models.GroupTransaction{}).Where("group_transactions.id = ?", transaction.Id).
		Updates(map[string]interface{}{
			"IsPaid": true,
		}).Error
	if err != nil {
		return err
	}

	err = uow.DB.Model(&models.GroupTransaction{}).Where("id = ?", transaction.Id).
		First(&transaction).Error
	if err != nil {
		return err
	}

	// updates payer incoming amount.
	err = g.setUserIncomingAmount(uow, transaction)
	if err != nil {
		return err
	}

	// updates payees outgoing amount.
	err = g.setUserOutgoingAmount(uow, transaction)
	if err != nil {
		return err
	}

	uow.Commit()
	return nil
}

// GetTransactionDetails will fetch amount to be fetched from all users for specified group
func (g *groupTransactionController) GetTransactionDetails(userBalance *[]models.UserBalance, userId, groupId uuid.UUID) error {
	err := g.doesUserExist(userId)
	if err != nil {
		return err
	}

	uow := db.NewUnitOfWork(g.db)
	defer uow.RollBack()

	err = uow.DB.Select("payee_id AS user_id, group_id, sum(amount) AS amount").Table("group_transactions").
		Preload("User").Where("payer_id = ? AND group_id = ? AND is_paid = ? AND is_adjusted = ?",
		userId, groupId, false, false).Group("payee_id, group_id").Order("amount").Find(userBalance).Error
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

	err = uow.DB.Model(&models.GroupTransaction{}).Where("group_transactions.id = ? AND group_transactions.payer_id = ?", transactionId, userId).
		First(&models.GroupTransaction{}).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("only payer can delete a transaction")
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
