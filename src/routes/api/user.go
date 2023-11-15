package api

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/shaileshhb/equisplit/src/controllers"
	"github.com/shaileshhb/equisplit/src/models"
	"github.com/shaileshhb/equisplit/src/security"
	"github.com/shaileshhb/equisplit/src/util"
)

type UserRouter interface {
	RegisterRoutes(router fiber.Router)
	register(ctx *fiber.Ctx) error
	login(c *fiber.Ctx) error
	logout(c *fiber.Ctx) error
	getUser(c *fiber.Ctx) error
	getUsers(c *fiber.Ctx) error
}

type userRouter struct {
	con  controllers.UserController
	auth security.Authentication
	log  zerolog.Logger
}

// NewUserRouter will create new instance for UserRouter
func NewUserRouter(con controllers.UserController, auth security.Authentication, log zerolog.Logger) UserRouter {
	return &userRouter{
		con:  con,
		auth: auth,
		log:  log,
	}
}

// RegisterRoutes will register user routes.
func (u *userRouter) RegisterRoutes(router fiber.Router) {
	router.Post("/register", u.register)
	router.Post("/login", u.login)
	router.Get("/logout", u.logout)
	router.Get("/user/:userId<int>", u.auth.MandatoryAuthMiddleware, u.getUser)
	router.Get("/user", u.auth.MandatoryAuthMiddleware, u.getUsers)
}

// register will add user.
func (u *userRouter) register(c *fiber.Ctx) error {
	user := &models.User{}

	err := c.BodyParser(user)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = u.con.Register(user)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := security.GenerateJWT(user)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "authorization",
		Value:    token,
		HTTPOnly: false,
		Secure:   true,
	})

	userResponse := map[string]interface{}{
		"userId": user.ID,
		"token":  token,
		"name":   user.Name,
		"email":  user.Email,
	}

	return c.Status(http.StatusCreated).JSON(userResponse)
}

// login will check user details and set the cookie
func (u *userRouter) login(c *fiber.Ctx) error {
	user := &models.User{}

	err := c.BodyParser(user)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = u.con.Login(user)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := security.GenerateJWT(user)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "authorization",
		Value:    token,
		HTTPOnly: false,
		Secure:   true,
	})

	userResponse := map[string]interface{}{
		"userId": user.ID,
		"token":  token,
		"name":   user.Name,
		"email":  user.Email,
	}

	return c.Status(http.StatusOK).JSON(userResponse)
}

// getUser will fetch specified user details.
func (u *userRouter) getUser(c *fiber.Ctx) error {
	user := models.UserDTO{}

	userId, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	user.ID = uint(userId)

	err = u.con.GetUser(&user)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(user)
}

// logout will log user out from the system
func (g *userRouter) logout(c *fiber.Ctx) error {

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "user successfully logged out",
	})
}

// getUsers will fetch specified user details.
func (u *userRouter) getUsers(c *fiber.Ctx) error {
	users := []models.UserDTO{}
	parser := util.NewParser(c)

	err := u.con.GetUsers(&users, parser)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(users)
}
