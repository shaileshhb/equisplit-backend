package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/shaileshhb/equisplit/src/controllers"
	"github.com/shaileshhb/equisplit/src/models"
	"github.com/shaileshhb/equisplit/src/security"
)

type UserInvitationRouter interface {
	RegisterRoutes(router fiber.Router)
	add(c *fiber.Ctx) error
	acceptInvitation(c *fiber.Ctx) error
}

type userInvitationRouter struct {
	con  controllers.UserInvitationController
	auth security.Authentication
	log  zerolog.Logger
}

// NewUserInvitationRouter will create new instance of UserInvitationRouter.
func NewUserInvitationRouter(con controllers.UserInvitationController, auth security.Authentication, log zerolog.Logger) UserInvitationRouter {
	return &userInvitationRouter{
		con:  con,
		auth: auth,
		log:  log,
	}
}

// RegisterRoutes will register routes for user-group router.
func (u *userInvitationRouter) RegisterRoutes(router fiber.Router) {
	router.Post("/user-invitation", u.auth.MandatoryAuthMiddleware, u.add)
	router.Put("/user-invitation", u.auth.MandatoryAuthMiddleware, u.acceptInvitation)

	u.log.Info().Msg("UserInvitation routes registered")
}

// add will create invitation for specified user in the group
func (u *userInvitationRouter) add(c *fiber.Ctx) error {
	userInvitation := models.UserInvitation{}

	err := c.BodyParser(&userInvitation)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userInterface := c.Locals("user")
	user := userInterface.(*models.User)
	userInvitation.InvitedBy = &user.ID

	err = u.con.Add(&userInvitation)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(nil)
}

// acceptInvitation will mark invitation as accepted and add user in the group that they were invited to.
func (u *userInvitationRouter) acceptInvitation(c *fiber.Ctx) error {
	userInvitation := models.UserInvitation{}

	err := c.BodyParser(&userInvitation)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = u.con.AcceptInvitation(&userInvitation)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusAccepted).JSON(nil)
}
