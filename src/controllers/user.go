package controllers

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/shaileshhb/equisplit/src/db"
	"github.com/shaileshhb/equisplit/src/models"
	"github.com/shaileshhb/equisplit/src/security"
	"github.com/shaileshhb/equisplit/src/util"
	"gorm.io/gorm"
)

type UserController interface {
	Register(user *models.User) error
	Login(user *models.User) error
	GetUser(user *models.UserDTO) error
	GetUsers(users *[]models.UserDTO, parser *util.Parser) error

	// Testing
	Unlimited(ip string) error
}

type userController struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewUserController(db *gorm.DB, rdb *redis.Client) UserController {
	return &userController{
		db:  db,
		rdb: rdb,
	}
}

func (u *userController) Unlimited(ip string) error {
	value, err := u.rdb.Get(db.Ctx, ip).Result()
	if err != nil && err != redis.Nil {
		return err
	}

	// this indicates that no entry exist for the specified IP in cache.
	if err == redis.Nil {
		fmt.Println("setting new value for specified ip")
		err := u.rdb.Set(db.Ctx, ip, os.Getenv("API_QUOTA"), 60*time.Second).Err()
		if err != nil {
			return err
		}
		value = os.Getenv("API_QUOTA")
	}

	fmt.Println("==================value is ->", value, len(value))
	var valueInt int

	if len(value) > 0 {
		valueInt, err = strconv.Atoi(value)
		if err != nil {
			return err
		}
	}

	fmt.Println("==================value int ->", valueInt)

	if valueInt <= 0 {
		limit, err := u.rdb.TTL(db.Ctx, ip).Result()
		if err != nil {
			return err
		}

		fmt.Println("==============limit after calling ttl", limit)
		return errors.New("rate limit exceeded")
	}

	// err = u.rdb.Set(db.Ctx, ip, valueInt-1, 60*time.Second).Err()
	// if err != nil {
	// 	return err
	// }

	return nil
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

	uow.Commit()
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
	user.Name = tempUser.Name

	return nil
}

// GetUser will fetch specified user details
func (u *userController) GetUser(user *models.UserDTO) error {

	err := u.db.First(user).Error
	if err != nil {
		return err
	}

	return nil
}

// GetUsers will fetch all users
func (u *userController) GetUsers(users *[]models.UserDTO, parser *util.Parser) error {

	uow := db.NewUnitOfWork(u.db)
	defer uow.RollBack()

	queryDB := u.searchQuery(uow, parser)

	err := queryDB.Find(users).Error
	if err != nil {
		return err
	}

	uow.Commit()
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

func (u *userController) searchQuery(uow *db.UnitOfWork, parser *util.Parser) *gorm.DB {
	queryDB := uow.DB

	if len(parser.GetQuery("email")) > 0 {
		queryDB = uow.DB.Where("users.email LIKE ?", parser.GetQuery("email")+"%")
	}

	if len(parser.GetQuery("name")) > 0 {
		queryDB = uow.DB.Where("users.name LIKE ?", "%"+parser.GetQuery("name")+"%")
	}

	return queryDB
}
