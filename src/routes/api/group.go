package api

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	router.Get("/user/:userId<uuid>/groups", g.auth.MandatoryAuthMiddleware, g.getUserGroups)
	router.Post("/user/:userId<uuid>/group", g.auth.MandatoryAuthMiddleware, g.createGroup)
	router.Put("/user/:userId<uuid>/group/:groupId<uuid>", g.auth.MandatoryAuthMiddleware, g.updateGroup)
	router.Delete("/user/:userId<uuid>/group/:groupId<uuid>", g.auth.MandatoryAuthMiddleware, g.deleteGroup)

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

	group.CreatedBy, err = uuid.Parse(c.Params("userId"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

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

	group.CreatedBy, err = uuid.Parse(c.Params("userId"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	group.Id, err = uuid.Parse(c.Params("groupId"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

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

	userId, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	group.Id, err = uuid.Parse(c.Params("groupId", "0"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	group.CreatedBy = userId

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
	group := &[]models.GroupDTO{}
	parser := util.NewParser(c)

	userId, err := uuid.Parse(c.Params("userId"))
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	var totalCount int64

	err = g.con.GetUserGroups(group, userId, &totalCount, parser)
	if err != nil {
		g.log.Error().Err(err).Msg("")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	c.Response().Header.Add("X-Total-Count", strconv.Itoa(int(totalCount)))

	return c.Status(http.StatusOK).JSON(group)
}
