package api

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/shaileshhb/equisplit/src/controllers"
	"github.com/shaileshhb/equisplit/src/models"
)

type UserGroupRouter interface {
	RegisterRoutes(router *fiber.Router)
	addUserToGroup(c *fiber.Ctx) error
	deleteUserFromGroup(c *fiber.Ctx) error
	getGroupDetails(c *fiber.Ctx) error
	getUserGroups(c *fiber.Ctx) error
}

type userGroupRouter struct {
	con controllers.UserGroupController
}

func NewUserGroupRouter(con controllers.UserGroupController) UserGroupRouter {
	return &userGroupRouter{
		con: con,
	}
}

// RegisterRoutes will register routes for user-group router.
func (u *userGroupRouter) RegisterRoutes(router *fiber.Router) {
	(*router).Get("/group/:groupId<int>/user", u.getGroupDetails)
	(*router).Get("/user/:userId<int>/group", u.getUserGroups)
	(*router).Post("/group/:groupId<int>/user", u.addUserToGroup)
	(*router).Post("/group/:groupId<int>/user/:userGroupId", u.deleteUserFromGroup)
}

// addUserToGroup will add user to specified group
func (u *userGroupRouter) addUserToGroup(c *fiber.Ctx) error {
	userGroup := models.UserGroup{}

	err := c.BodyParser(&userGroup)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id, err := strconv.Atoi(c.Params("groupId"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userGroup.GroupId = uint(id)

	err = u.con.AddUserToGroup(&userGroup)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{})
}

// deleteUserFromGroup will delete specified user from group
func (u *userGroupRouter) deleteUserFromGroup(c *fiber.Ctx) error {
	userGroup := models.UserGroup{}

	id, err := strconv.Atoi(c.Params("userGroupId"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	userGroup.ID = uint(id)

	err = u.con.DeleteUserFromGroup(&userGroup)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusAccepted).JSON(fiber.Map{})
}

// getGroupDetails will fetch all user details from specified group
func (u *userGroupRouter) getGroupDetails(c *fiber.Ctx) error {
	userGroups := []models.UserGroupDTO{}

	groupId, err := strconv.Atoi(c.Params("groupId"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = u.con.GetGroupDetails(&userGroups, uint(groupId))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"data": userGroups,
	})
}

// getUserGroups will fetch all groups for specified user
func (u *userGroupRouter) getUserGroups(c *fiber.Ctx) error {
	userGroups := []models.UserGroupDTO{}

	userId, err := strconv.Atoi(c.Params("userId"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = u.con.GetUserGroups(&userGroups, uint(userId))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"data": userGroups,
	})
}
