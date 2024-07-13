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

type UserGroupRouter interface {
	RegisterRoutes(router fiber.Router)
	addUserToGroup(c *fiber.Ctx) error
	deleteUserFromGroup(c *fiber.Ctx) error
	getGroupDetails(c *fiber.Ctx) error
	getUserGroups(c *fiber.Ctx) error
}

type userGroupRouter struct {
	con  controllers.UserGroupController
	auth security.Authentication
	log  zerolog.Logger
}

// NewUserGroupRouter will create new instance of UserGroupRouter.
func NewUserGroupRouter(con controllers.UserGroupController, auth security.Authentication, log zerolog.Logger) UserGroupRouter {
	return &userGroupRouter{
		con:  con,
		auth: auth,
		log:  log,
	}
}

// RegisterRoutes will register routes for user-group router.
func (u *userGroupRouter) RegisterRoutes(router fiber.Router) {
	router.Get("/group/:groupId<uuid>", u.auth.MandatoryAuthMiddleware, u.getGroupDetails)
	router.Get("/user/:userId<uuid>/group", u.auth.MandatoryAuthMiddleware, u.getUserGroups)
	router.Post("/group/:groupId<uuid>/user", u.auth.MandatoryAuthMiddleware, u.addUserToGroup)
	router.Get("/group/:groupId<uuid>/users", u.auth.MandatoryAuthMiddleware, u.getGroupUsers)
	router.Delete("/group/:groupId<uuid>/user/:userGroupId<uuid>", u.auth.MandatoryAuthMiddleware, u.deleteUserFromGroup)
	u.log.Info().Msg("UserGroup routes registered")
}

// addUserToGroup will add user to specified group
func (u *userGroupRouter) addUserToGroup(c *fiber.Ctx) error {
	userGroup := models.UserGroup{}

	err := c.BodyParser(&userGroup)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userGroup.GroupId, err = uuid.Parse(c.Params("groupId"))
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = u.con.AddUserToGroup(&userGroup)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(nil)
}

// deleteUserFromGroup will delete specified user from group
func (u *userGroupRouter) deleteUserFromGroup(c *fiber.Ctx) error {
	userGroup := models.UserGroup{}

	id, err := uuid.Parse(c.Params("userGroupId"))
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userGroup.Id = id

	err = u.con.DeleteUserFromGroup(&userGroup)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusAccepted).JSON(nil)
}

// getGroupDetails will fetch all user details from specified group
func (u *userGroupRouter) getGroupDetails(c *fiber.Ctx) error {
	userGroups := []models.UserGroupDTO{}

	groupId, err := uuid.Parse(c.Params("groupId"))
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userInterface := c.Locals("user")
	user := userInterface.(*models.User)

	err = u.con.GetGroupDetails(&userGroups, groupId, user.Id)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(userGroups)
}

// getUserGroups will fetch all groups for specified user
func (u *userGroupRouter) getUserGroups(c *fiber.Ctx) error {
	userGroups := []models.UserGroupDTO{}

	userId, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = u.con.GetUserGroups(&userGroups, userId)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(userGroups)
}

// getGroupUsers will fetch all groups for specified user
func (u *userGroupRouter) getGroupUsers(c *fiber.Ctx) error {
	userGroups := []models.UserGroupDTO{}

	groupId, err := uuid.Parse(c.Params("groupId"))
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = u.con.GetGroupUsers(&userGroups, groupId)
	if err != nil {
		u.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(userGroups)
}
