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

type GroupRouter interface {
	RegisterRoutes(router fiber.Router)
	createGroup(c *fiber.Ctx) error
	updateGroup(c *fiber.Ctx) error
	deleteGroup(c *fiber.Ctx) error
	getUserGroups(c *fiber.Ctx) error
}

type groupRouter struct {
	con  controllers.GroupController
	auth security.Authentication
	log  zerolog.Logger
}

// NewGroupRouter will create new instance of GroupRouter.
func NewGroupRouter(con controllers.GroupController, auth security.Authentication, log zerolog.Logger) GroupRouter {
	return &groupRouter{
		con:  con,
		auth: auth,
		log:  log,
	}
}

// RegisterRoutes will register routes for group.
func (g *groupRouter) RegisterRoutes(router fiber.Router) {
	router.Get("/:userId<int>/group", g.auth.MandatoryAuthMiddleware, g.getUserGroups)
	router.Post("/:userId<int>/group", g.auth.MandatoryAuthMiddleware, g.createGroup)
	router.Put("/:userId<int>/group/:groupId<int>", g.auth.MandatoryAuthMiddleware, g.updateGroup)
	router.Delete("/:userId<int>/group/:groupId<int>", g.auth.MandatoryAuthMiddleware, g.deleteGroup)

	g.log.Info().Msg("Group routes registered")
}

// createGroup will create new group for specified user.
func (g *groupRouter) createGroup(c *fiber.Ctx) error {
	group := &models.Group{}

	err := c.BodyParser(group)
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id, err := strconv.Atoi(c.Params("userId", "0"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	group.CreatedBy = uint(id)

	err = g.con.CreateGroup(group)
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(nil)
}

// updateGroup will update group for specified user.
func (g *groupRouter) updateGroup(c *fiber.Ctx) error {
	group := &models.Group{}

	err := c.BodyParser(group)
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	id, err := strconv.Atoi(c.Params("userId", "0"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	group.CreatedBy = uint(id)

	id, err = strconv.Atoi(c.Params("groupId", "0"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	group.ID = uint(id)

	err = g.con.CreateGroup(group)
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusAccepted).JSON(nil)
}

// deleteGroup will delete group for specified user.
func (g *groupRouter) deleteGroup(c *fiber.Ctx) error {
	group := &models.Group{}

	id, err := strconv.Atoi(c.Params("userId", "0"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	group.CreatedBy = uint(id)

	id, err = strconv.Atoi(c.Params("groupId", "0"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	group.ID = uint(id)

	err = g.con.DeleteGroup(group)
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusAccepted).JSON(nil)
}

// deleteGroup will delete group for specified user.
func (g *groupRouter) getUserGroups(c *fiber.Ctx) error {
	group := &[]models.Group{}
	parser := util.NewParser(c)

	id, err := strconv.Atoi(c.Params("userId", "0"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var totalCount int64

	err = g.con.GetUserGroups(group, uint(id), &totalCount, parser)
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Response().Header.Add("X-Total-Count", strconv.Itoa(int(totalCount)))

	return c.Status(http.StatusOK).JSON(group)
}
