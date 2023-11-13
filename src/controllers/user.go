package controllers

import (
	"errors"

	"github.com/shaileshhb/equisplit/src/db"
	"github.com/shaileshhb/equisplit/src/models"
	"github.com/shaileshhb/equisplit/src/security"
	"gorm.io/gorm"
)

type UserController interface {
	Register(user *models.User) error
	Login(user *models.User) error
	GetUser(user *models.User) error
}

type userController struct {
	db *gorm.DB
}

func NewUserController(db *gorm.DB) UserController {
	return &userController{
		db: db,
	}
}

// Register will register new user in the system.
func (u *userController) Register(user *models.User) error {

	err := u.validateUser(user)
	if err != nil {
		return err
	}

	password, err := security.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = string(password)

	uow := db.NewUnitOfWork(u.db)
	defer uow.RollBack()

	err = uow.DB.Create(user).Error
	if err != nil {
		return err
	}

	// uow.Commit()
	return nil
}

// Login user.
func (u *userController) Login(user *models.User) error {

	uow := db.NewUnitOfWork(u.db)
	defer uow.RollBack()

	tempUser := &models.User{}
	err := uow.DB.Where("users.email = ?", user.Email).First(tempUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New("email not registered")
		}
		return err
	}

	err = security.ComparePassword(tempUser.Password, user.Password)
	if err != nil {
		return errors.New("email or password did not match")
	}

	user.ID = tempUser.ID

	return nil
}

// GetUser will fetch specified user details
func (u *userController) GetUser(user *models.User) error {

	err := u.db.First(user).Error
	if err != nil {
		return err
	}

	return nil
}

// validateUer will check if it is unique user.
func (u *userController) validateUser(user *models.User) error {

	var count int64 = 0
	err := u.db.Model(&models.User{}).
		Select("COUNT(DISTINCT(id))").
		Where("users.id != ? AND users.email = ?", user.ID, user.Email).
		Unscoped().
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("email already exist")
	}

	return nil
}
