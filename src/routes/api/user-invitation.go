package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/shaileshhb/equisplit/src/controllers"
	"github.com/shaileshhb/equisplit/src/models"
	"github.com/shaileshhb/equisplit/src/security"
)

type UserInvitationRouter interface {
	RegisterRoutes(router fiber.Router)
	add(c *fiber.Ctx) error
	updateInvitation(c *fiber.Ctx) error
	deleteInvitation(c *fiber.Ctx) error
	getGroupInvitation(c *fiber.Ctx) error
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
	router.Put("/user-invitation/:userInvitationId<uint>", u.auth.MandatoryAuthMiddleware, u.updateInvitation)
	router.Delete("/user-invitation/:userInvitationId<uint>", u.auth.MandatoryAuthMiddleware, u.deleteInvitation)
	router.Get("/group/:groupId<uint>/user-invitation", u.auth.MandatoryAuthMiddleware, u.getGroupInvitation)

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
	userInvitation.InvitedBy = &user.Id

	err = u.con.Add(&userInvitation)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(nil)
}

// updateInvitation will mark invitation as accepted and add user in the group that they were invited to.
func (u *userInvitationRouter) updateInvitation(c *fiber.Ctx) error {
	userInvitation := models.UserInvitation{}

	err := c.BodyParser(&userInvitation)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userInvitation.Id, err = uuid.Parse(c.Params("userInvitationId"))
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = u.con.UpdateInvitation(&userInvitation)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusAccepted).JSON(nil)
}

// deleteInvitation will delete the specified invitation
func (u *userInvitationRouter) deleteInvitation(c *fiber.Ctx) error {
	userInvitation := models.UserInvitation{}

	err := c.BodyParser(&userInvitation)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userInvitation.Id, err = uuid.Parse(c.Params("userInvitationId"))
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = u.con.DeleteInvitation(&userInvitation)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusAccepted).JSON(nil)
}

// getGroupInvitation will fetch all invitations of specified group.
func (u *userInvitationRouter) getGroupInvitation(c *fiber.Ctx) error {
	userInvitation := []models.UserInvitation{}

	groupId, err := uuid.Parse(c.Params("groupId"))
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = u.con.GetGroupInvitation(&userInvitation, groupId)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(nil)
}
