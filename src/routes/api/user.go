package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/shaileshhb/equisplit/src/controllers"
	"github.com/shaileshhb/equisplit/src/models"
	"github.com/shaileshhb/equisplit/src/security"
)

type UserRouter interface {
	RegisterRoutes(router *fiber.Router)
	register(ctx *fiber.Ctx) error
	login(c *fiber.Ctx) error
}

type userRouter struct {
	service controllers.UserController
}

func NewUserRouter(service controllers.UserController) UserRouter {
	return &userRouter{
		service: service,
	}
}

func (u *userRouter) RegisterRoutes(router *fiber.Router) {
	(*router).Post("/register", u.register)
	(*router).Post("/login", u.login)
}

func (u *userRouter) register(c *fiber.Ctx) error {
	user := &models.User{}

	err := c.BodyParser(user)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = u.service.Register(user)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := security.GenerateJWT(user)
	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     "authorization",
		Value:    token,
		HTTPOnly: false,
		Secure:   true,
	})

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User successfully registered",
	})
}

func (u *userRouter) login(c *fiber.Ctx) error {
	user := &models.User{}

	err := c.BodyParser(user)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = u.service.Login(user)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "User successfully logged in",
	})
}
